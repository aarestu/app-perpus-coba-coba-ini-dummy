package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app-perpus/internal/database"
	"app-perpus/internal/middleware"
	"app-perpus/internal/models"
	"app-perpus/internal/services"
)

func setupTestRouter(t *testing.T) (*http.ServeMux, *services.AuthService) {
	t.Helper()
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Gagal inisialisasi database testing: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	authService := services.NewAuthService(db, "test-secret-key", 1*time.Hour)
	authHandler := NewAuthHandler(authService)
	authMiddleware := middleware.AuthMiddleware(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.Handle("/api/auth/me", authMiddleware(http.HandlerFunc(authHandler.Me)))

	return mux, authService
}

func TestAuthHandler_RegisterAndLogin(t *testing.T) {
	mux, _ := setupTestRouter(t)

	// 1. Uji Registrasi Sukses (HTTP 201)
	regPayload := map[string]any{
		"name":     "Rudi Tabuti",
		"email":    "rudi@perpus.local",
		"password": "securepassword123",
		"height":   172,
	}
	body, _ := json.Marshal(regPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Status registrasi diharapkan 201, didapat %d. Body: %s", rec.Code, rec.Body.String())
	}

	var regResponse struct {
		Success bool `json:"success"`
		Data    struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &regResponse); err != nil {
		t.Fatalf("Gagal unmarshal response registrasi: %v", err)
	}
	if !regResponse.Success {
		t.Errorf("Expected success: true")
	}
	if regResponse.Data.User.Height != 172 {
		t.Errorf("Expected user height 172, got %d", regResponse.Data.User.Height)
	}

	// 2. Uji Registrasi Gagal - Validasi Input (HTTP 400)
	invalidRegPayload := map[string]string{
		"name":     "",
		"email":    "invalid-email",
		"password": "123",
	}
	body, _ = json.Marshal(invalidRegPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Status diharapkan 400 untuk validasi gagal, didapat %d", rec.Code)
	}

	// 3. Uji Registrasi Duplikat (HTTP 409)
	body, _ = json.Marshal(regPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("Status diharapkan 409 untuk email duplikat, didapat %d", rec.Code)
	}

	// 4. Uji Login Sukses (HTTP 200)
	loginPayload := map[string]string{
		"email":    "rudi@perpus.local",
		"password": "securepassword123",
	}
	body, _ = json.Marshal(loginPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Status login diharapkan 200, didapat %d. Body: %s", rec.Code, rec.Body.String())
	}

	var loginResponse struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &loginResponse); err != nil {
		t.Fatalf("Gagal unmarshal response login: %v", err)
	}
	token := loginResponse.Data.Token
	if token == "" {
		t.Fatalf("Token tidak ditemukan di response login")
	}

	// 5. Uji Login Gagal - Password Salah (HTTP 401)
	wrongLoginPayload := map[string]string{
		"email":    "rudi@perpus.local",
		"password": "wrongpassword",
	}
	body, _ = json.Marshal(wrongLoginPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Status login kredensial salah diharapkan 401, didapat %d", rec.Code)
	}

	// 6. Uji Akses Rute Terproteksi /api/auth/me dengan Token Valid (HTTP 200)
	req = httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Status /api/auth/me diharapkan 200, didapat %d. Body: %s", rec.Code, rec.Body.String())
	}

	var meResponse struct {
		Success bool        `json:"success"`
		Data    models.User `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &meResponse); err != nil {
		t.Fatalf("Gagal unmarshal response /api/auth/me: %v", err)
	}
	if meResponse.Data.Height != 172 {
		t.Errorf("Expected user height 172 in /api/auth/me, got %d", meResponse.Data.Height)
	}

	// 7. Uji Akses Rute Terproteksi /api/auth/me tanpa Token (HTTP 401)
	req = httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Status diharapkan 401 tanpa token, didapat %d", rec.Code)
	}

	// 8. Uji Akses Rute Terproteksi /api/auth/me dengan Token Palsu (HTTP 401)
	req = httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer invalid-tampered-token")
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	// 9. Uji Registrasi dengan alias 'tinggi_badan' (HTTP 201)
	aliasRegPayload := map[string]any{
		"name":         "Siti Nurbaya",
		"email":        "siti@perpus.local",
		"password":     "securepassword123",
		"tinggi_badan": 160,
	}
	body, _ = json.Marshal(aliasRegPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Status registrasi diharapkan 201, didapat %d", rec.Code)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &regResponse); err != nil {
		t.Fatalf("Gagal unmarshal response registrasi alias: %v", err)
	}
	if regResponse.Data.User.Height != 160 {
		t.Errorf("Expected user height 160 from tinggi_badan alias, got %d", regResponse.Data.User.Height)
	}
}

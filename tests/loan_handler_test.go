package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app-perpus/internal/database"
	"app-perpus/internal/handlers"
	"app-perpus/internal/middleware"
	"app-perpus/internal/models"
	"app-perpus/internal/services"
)

func setupLoanRouter(t *testing.T) (*http.ServeMux, *services.AuthService, *services.LoanService, string, int64) {
	t.Helper()
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Gagal inisialisasi database testing: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	authService := services.NewAuthService(db, "test-secret-key-loans", 1*time.Hour)
	authHandler := handlers.NewAuthHandler(authService)
	loanService := services.NewLoanService(db)
	loanHandler := handlers.NewLoanHandler(loanService)
	authMiddleware := middleware.AuthMiddleware(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)
	mux.Handle("/api/auth/me", authMiddleware(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/api/loans", authMiddleware(http.HandlerFunc(loanHandler.HandleLoans)))

	// Daftarkan user contoh untuk test token
	regPayload := map[string]any{
		"name":     "Gita Gutawa",
		"email":    "gita@perpus.local",
		"password": "password123",
		"height":   165,
	}
	body, _ := json.Marshal(regPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var resp struct {
		Data struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		} `json:"data"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)

	return mux, authService, loanService, resp.Data.Token, resp.Data.User.ID
}

func TestLoanHandler_CreateAndGetLoans(t *testing.T) {
	mux, _, _, token, userID := setupLoanRouter(t)

	// 1. Uji Peminjaman Buku Sukses (HTTP 201)
	loanPayload := map[string]any{
		"book_title":    "Belajar Pemrograman Go Dasar",
		"duration_days": 10,
	}
	body, _ := json.Marshal(loanPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/loans", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Status peminjaman diharapkan 201, didapat %d. Body: %s", rec.Code, rec.Body.String())
	}

	var loanResp struct {
		Success bool        `json:"success"`
		Data    models.Loan `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &loanResp); err != nil {
		t.Fatalf("Gagal unmarshal response peminjaman: %v", err)
	}
	if !loanResp.Success {
		t.Errorf("Expected success: true")
	}
	if loanResp.Data.UserID != userID {
		t.Errorf("Expected UserID %d, got %d", userID, loanResp.Data.UserID)
	}
	if loanResp.Data.BookTitle != "Belajar Pemrograman Go Dasar" {
		t.Errorf("Expected BookTitle 'Belajar Pemrograman Go Dasar', got '%s'", loanResp.Data.BookTitle)
	}
	if loanResp.Data.Status != "borrowed" {
		t.Errorf("Expected Status 'borrowed', got '%s'", loanResp.Data.Status)
	}

	// 2. Uji Peminjaman Buku dengan Alias 'judul_buku' (HTTP 201)
	aliasPayload := map[string]any{
		"judul_buku": "Arsitektur Microservices Modern",
	}
	body, _ = json.Marshal(aliasPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/loans", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Status peminjaman dengan alias diharapkan 201, didapat %d", rec.Code)
	}

	// 3. Uji Peminjaman Buku Gagal - Tanpa Token / Belum Login (HTTP 401)
	req = httptest.NewRequest(http.MethodPost, "/api/loans", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Status diharapkan 401 untuk request tanpa token, didapat %d", rec.Code)
	}

	// 4. Uji Peminjaman Buku Gagal - Validasi Input (Judul Kosong) (HTTP 400)
	invalidPayload := map[string]any{
		"book_title": "",
	}
	body, _ = json.Marshal(invalidPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/loans", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Status diharapkan 400 untuk validasi judul kosong, didapat %d", rec.Code)
	}

	// 5. Uji Peminjaman Buku Gagal - JSON Rusak (HTTP 400)
	req = httptest.NewRequest(http.MethodPost, "/api/loans", bytes.NewReader([]byte("{bad-json}")))
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Status diharapkan 400 untuk body json rusak, didapat %d", rec.Code)
	}

	// 6. Uji Ambil Daftar Peminjaman Pengguna (GET /api/loans) (HTTP 200)
	req = httptest.NewRequest(http.MethodGet, "/api/loans", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Status GET /api/loans diharapkan 200, didapat %d. Body: %s", rec.Code, rec.Body.String())
	}

	var listResp struct {
		Success bool          `json:"success"`
		Data    []models.Loan `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &listResp); err != nil {
		t.Fatalf("Gagal unmarshal response daftar peminjaman: %v", err)
	}
	if len(listResp.Data) != 2 {
		t.Errorf("Diharapkan 2 peminjaman pada riwayat, didapat %d", len(listResp.Data))
	}

	// 7. Uji Metode HTTP Tidak Diizinkan (DELETE /api/loans) (HTTP 405)
	req = httptest.NewRequest(http.MethodDelete, "/api/loans", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Status diharapkan 405 untuk method DELETE, didapat %d", rec.Code)
	}
}

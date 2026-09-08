package services

import (
	"context"
	"testing"
	"time"

	"app-perpus/internal/database"
	"app-perpus/internal/models"
)

func setupTestDB(t *testing.T) *AuthService {
	t.Helper()
	// Menggunakan in-memory SQLite database terisolasi untuk unit testing
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Gagal inisialisasi database testing in-memory: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	return NewAuthService(db, "test-jwt-secret-key-12345", 1*time.Hour)
}

func TestPasswordHashing(t *testing.T) {
	service := setupTestDB(t)

	password := "RahasiaSuper123!"
	hash, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword gagal: %v", err)
	}

	if hash == password {
		t.Errorf("Hash tidak boleh sama persis dengan plaintext password")
	}

	if !service.CheckPassword(password, hash) {
		t.Errorf("CheckPassword harus mengembalikan true untuk password yang cocok")
	}

	if service.CheckPassword("PasswordSalah123", hash) {
		t.Errorf("CheckPassword harus mengembalikan false untuk password yang salah")
	}
}

func TestJWTTokenGenerationAndValidation(t *testing.T) {
	service := setupTestDB(t)

	user := &models.User{
		ID:    42,
		Name:  "Budi Santoso",
		Email: "budi@perpus.local",
		Role:  "admin",
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken gagal: %v", err)
	}

	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken gagal: %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("UserID mismatch: didapat %d, diharapkan %d", claims.UserID, user.ID)
	}
	if claims.Email != user.Email {
		t.Errorf("Email mismatch: didapat %s, diharapkan %s", claims.Email, user.Email)
	}
	if claims.Role != user.Role {
		t.Errorf("Role mismatch: didapat %s, diharapkan %s", claims.Role, user.Role)
	}

	// Validasi token yang rusak / invalid
	_, err = service.ValidateToken("invalid.token.string")
	if err == nil {
		t.Errorf("ValidateToken harus gagal untuk format token acak")
	}
}

func TestRegisterAndLoginTableDriven(t *testing.T) {
	service := setupTestDB(t)
	ctx := context.Background()

	// 1. Sukses registrasi pertama
	user1, token1, err := service.Register(ctx, models.RegisterRequest{
		Name:     "Dewi Sartika",
		Email:    "dewi@perpus.local",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Registrasi user pertama gagal: %v", err)
	}
	if user1.ID == 0 || token1 == "" {
		t.Errorf("ID user dan token tidak boleh kosong")
	}

	// Table driven tests untuk skenario login & registrasi lanjutan
	tests := []struct {
		name          string
		action        string // "register" atau "login"
		regReq        models.RegisterRequest
		loginReq      models.LoginRequest
		expectErr     bool
		expectedError error
	}{
		{
			name:   "Registrasi gagal karena email duplikat",
			action: "register",
			regReq: models.RegisterRequest{
				Name:     "Dewi Duplikat",
				Email:    "dewi@perpus.local",
				Password: "passwordBaru123",
			},
			expectErr:     true,
			expectedError: ErrUserAlreadyExists,
		},
		{
			name:   "Login sukses dengan kredensial benar",
			action: "login",
			loginReq: models.LoginRequest{
				Email:    "dewi@perpus.local",
				Password: "password123",
			},
			expectErr: false,
		},
		{
			name:   "Login gagal dengan password salah",
			action: "login",
			loginReq: models.LoginRequest{
				Email:    "dewi@perpus.local",
				Password: "passwordSalahBanget",
			},
			expectErr:     true,
			expectedError: ErrInvalidCredentials,
		},
		{
			name:   "Login gagal dengan email tidak terdaftar",
			action: "login",
			loginReq: models.LoginRequest{
				Email:    "tidakada@perpus.local",
				Password: "password123",
			},
			expectErr:     true,
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.action == "register" {
				_, _, err := service.Register(ctx, tc.regReq)
				if tc.expectErr && err == nil {
					t.Errorf("Diharapkan error, tapi mendapatkan nil")
				}
				if tc.expectedError != nil && err != tc.expectedError {
					t.Errorf("Expected error %v, got %v", tc.expectedError, err)
				}
			} else if tc.action == "login" {
				user, token, err := service.Login(ctx, tc.loginReq)
				if tc.expectErr {
					if err == nil {
						t.Errorf("Diharapkan error login, tapi mendapatkan nil")
					}
					if tc.expectedError != nil && err != tc.expectedError {
						t.Errorf("Expected error %v, got %v", tc.expectedError, err)
					}
				} else {
					if err != nil {
						t.Fatalf("Login gagal tak terduga: %v", err)
					}
					if user.Email != tc.loginReq.Email {
						t.Errorf("Expected email %s, got %s", tc.loginReq.Email, user.Email)
					}
					if token == "" {
						t.Errorf("Token tidak boleh kosong")
					}
				}
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	service := setupTestDB(t)
	ctx := context.Background()

	user, _, err := service.Register(ctx, models.RegisterRequest{
		Name:     "Andi",
		Email:    "andi@perpus.local",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Register gagal: %v", err)
	}

	foundUser, err := service.GetUserByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetUserByID gagal: %v", err)
	}
	if foundUser.Name != "Andi" {
		t.Errorf("Expected name 'Andi', got '%s'", foundUser.Name)
	}

	_, err = service.GetUserByID(ctx, 99999)
	if err != ErrUserNotFound {
		t.Errorf("Diharapkan ErrUserNotFound untuk ID tidak dikenal, didapat: %v", err)
	}
}

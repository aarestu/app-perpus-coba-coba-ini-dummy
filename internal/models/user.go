package models

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User merepresentasikan entitas pengguna dalam sistem
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	Age          int       `json:"age"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RegisterRequest merepresentasikan payload body registrasi user
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

// Validate memvalidasi kelayakan field RegisterRequest
func (r *RegisterRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Name) == "" {
		errs["name"] = "Nama wajib diisi"
	}
	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = "Email wajib diisi"
	} else if !strings.Contains(r.Email, "@") || !strings.Contains(r.Email, ".") {
		errs["email"] = "Format email tidak valid"
	}
	if len(r.Password) < 6 {
		errs["password"] = "Password minimal terdiri dari 6 karakter"
	}
	if r.Age < 0 {
		errs["age"] = "Umur tidak boleh bernilai negatif"
	}
	return errs
}

// LoginRequest merepresentasikan payload body login user
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate memvalidasi kelayakan field LoginRequest
func (r *LoginRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.Email) == "" {
		errs["email"] = "Email wajib diisi"
	}
	if strings.TrimSpace(r.Password) == "" {
		errs["password"] = "Password wajib diisi"
	}
	return errs
}

// AuthResponse merepresentasikan payload balasan sukses setelah login/register
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// JWTClaims merepresentasikan klaim payload dalam token JWT
type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

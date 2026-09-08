package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"app-perpus/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("email sudah terdaftar dalam sistem")
	ErrInvalidCredentials = errors.New("email atau password yang dimasukkan salah")
	ErrUserNotFound      = errors.New("pengguna tidak ditemukan")
	ErrInvalidToken      = errors.New("token autentikasi tidak valid atau telah kedaluwarsa")
)

type AuthService struct {
	db            *sql.DB
	jwtSecret     []byte
	tokenDuration time.Duration
}

// NewAuthService membuat instance baru AuthService
func NewAuthService(db *sql.DB, jwtSecret string, tokenDuration time.Duration) *AuthService {
	if tokenDuration <= 0 {
		tokenDuration = 24 * time.Hour
	}
	return &AuthService{
		db:            db,
		jwtSecret:     []byte(jwtSecret),
		tokenDuration: tokenDuration,
	}
}

// HashPassword melakukan enkripsi password menggunakan bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("gagal mengenkripsi password: %w", err)
	}
	return string(bytes), nil
}

// CheckPassword membandingkan plain password dengan bcrypt hash
func (s *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken membuat JWT token yang ditandatangani untuk pengguna
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(s.tokenDuration)
	claims := &models.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("gagal menerbitkan token JWT: %w", err)
	}

	return tokenString, nil
}

// ValidateToken memvalidasi string JWT token dan mengembalikan klaimnya
func (s *AuthService) ValidateToken(tokenStr string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// Register menambahkan pengguna baru ke basis data SQLite
func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, string, error) {
	emailNormalized := strings.ToLower(strings.TrimSpace(req.Email))

	// Cek apakah email sudah terdaftar
	var existingID int64
	err := s.db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", emailNormalized).Scan(&existingID)
	if err == nil {
		return nil, "", ErrUserAlreadyExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, "", fmt.Errorf("gagal mengecek ketersediaan email: %w", err)
	}

	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	now := time.Now()
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO users (name, email, password_hash, role, height, created_at, updated_at) 
		 VALUES (?, ?, ?, 'user', ?, ?, ?)`,
		strings.TrimSpace(req.Name), emailNormalized, hashedPassword, req.Height, now, now,
	)
	if err != nil {
		return nil, "", fmt.Errorf("gagal menyimpan data user: %w", err)
	}

	userID, err := res.LastInsertId()
	if err != nil {
		return nil, "", fmt.Errorf("gagal mendapatkan ID user: %w", err)
	}

	user := &models.User{
		ID:        userID,
		Name:      strings.TrimSpace(req.Name),
		Email:     emailNormalized,
		Role:      "user",
		Height:    req.Height,
		CreatedAt: now,
		UpdatedAt: now,
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login memverifikasi email & password dan menerbitkan token jika valid
func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.User, string, error) {
	emailNormalized := strings.ToLower(strings.TrimSpace(req.Email))

	user := &models.User{}
	query := `SELECT id, name, email, password_hash, role, height, created_at, updated_at FROM users WHERE email = ?`
	err := s.db.QueryRowContext(ctx, query, emailNormalized).Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.Height, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", ErrInvalidCredentials
	} else if err != nil {
		return nil, "", fmt.Errorf("gagal memproses query login: %w", err)
	}

	if !s.CheckPassword(req.Password, user.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUserByID mengambil data user berdasarkan ID
func (s *AuthService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, name, email, role, height, created_at, updated_at FROM users WHERE id = ?`
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.Height, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("gagal mengambil data user: %w", err)
	}

	return user, nil
}

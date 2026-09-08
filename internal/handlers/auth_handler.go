package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"app-perpus/internal/middleware"
	"app-perpus/internal/models"
	"app-perpus/internal/services"
	"app-perpus/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register menangani pembuatan akun pengguna baru (POST /api/auth/register)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Metode HTTP tidak diizinkan")
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Format JSON tidak valid atau rusak")
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		utils.JSONError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Data registrasi tidak memenuhi validasi", validationErrors)
		return
	}

	user, token, err := h.authService.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			utils.JSONError(w, http.StatusConflict, "USER_ALREADY_EXISTS", err.Error())
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Gagal memproses pendaftaran pengguna")
		return
	}

	utils.JSONSuccess(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Login memvalidasi kredensial pengguna dan mengembalikan token JWT (POST /api/auth/login)
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Metode HTTP tidak diizinkan")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Format JSON tidak valid atau rusak")
		return
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		utils.JSONError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Data login tidak memenuhi validasi", validationErrors)
		return
	}

	user, token, err := h.authService.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.JSONError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Gagal memproses autentikasi")
		return
	}

	utils.JSONSuccess(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

// Me mengembalikan informasi profil pengguna yang sedang login (GET /api/auth/me)
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Metode HTTP tidak diizinkan")
		return
	}

	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok || claims == nil {
		utils.JSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Sesi tidak valid")
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			utils.JSONError(w, http.StatusNotFound, "USER_NOT_FOUND", err.Error())
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Gagal mengambil data profil pengguna")
		return
	}

	utils.JSONSuccess(w, http.StatusOK, user)
}

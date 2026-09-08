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

type LoanHandler struct {
	loanService *services.LoanService
}

// NewLoanHandler membuat instance baru LoanHandler
func NewLoanHandler(loanService *services.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

// HandleLoans mengarahkan request /api/loans berdasarkan metode HTTP
func (h *LoanHandler) HandleLoans(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateLoan(w, r)
	case http.MethodGet:
		h.GetUserLoans(w, r)
	default:
		utils.JSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Metode HTTP tidak diizinkan")
	}
}

// CreateLoan menangani pembuatan peminjaman buku baru (POST /api/loans)
func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.JSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Metode HTTP tidak diizinkan")
		return
	}

	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok || claims == nil {
		utils.JSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Sesi tidak valid atau belum terautentikasi")
		return
	}

	var req models.CreateLoanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "Format JSON tidak valid atau rusak")
		return
	}

	if valErrs := req.Validate(); len(valErrs) > 0 {
		utils.JSONError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Data peminjaman tidak memenuhi validasi", valErrs)
		return
	}

	loan, err := h.loanService.CreateLoan(r.Context(), claims.UserID, req)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			utils.JSONError(w, http.StatusNotFound, "USER_NOT_FOUND", err.Error())
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Gagal memproses peminjaman buku")
		return
	}

	utils.JSONSuccess(w, http.StatusCreated, loan)
}

// GetUserLoans mengembalikan riwayat peminjaman buku untuk pengguna yang sedang login (GET /api/loans)
func (h *LoanHandler) GetUserLoans(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.JSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Metode HTTP tidak diizinkan")
		return
	}

	claims, ok := middleware.GetClaimsFromContext(r.Context())
	if !ok || claims == nil {
		utils.JSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Sesi tidak valid atau belum terautentikasi")
		return
	}

	loans, err := h.loanService.GetUserLoans(r.Context(), claims.UserID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Gagal mengambil riwayat peminjaman buku")
		return
	}

	utils.JSONSuccess(w, http.StatusOK, loans)
}

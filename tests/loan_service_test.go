package tests

import (
	"context"
	"testing"
	"time"

	"app-perpus/internal/database"
	"app-perpus/internal/models"
	"app-perpus/internal/services"
)

func setupLoanTestDB(t *testing.T) (*services.LoanService, *services.AuthService) {
	t.Helper()
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Gagal inisialisasi database testing in-memory: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	authService := services.NewAuthService(db, "test-jwt-secret-key-12345", 1*time.Hour)
	loanService := services.NewLoanService(db)

	return loanService, authService
}

func TestLoanService_CreateLoan(t *testing.T) {
	loanService, authService := setupLoanTestDB(t)
	ctx := context.Background()

	// Buat pengguna terlebih dahulu
	user, _, err := authService.Register(ctx, models.RegisterRequest{
		Name:     "Citra Lestari",
		Email:    "citra@perpus.local",
		Password: "password123",
		Height:   168,
	})
	if err != nil {
		t.Fatalf("Gagal membuat user untuk testing: %v", err)
	}

	t.Run("Peminjaman berhasil dengan durasi default (7 hari)", func(t *testing.T) {
		req := models.CreateLoanRequest{
			BookTitle: "Pemrograman Go untuk Pemula",
		}
		loan, err := loanService.CreateLoan(ctx, user.ID, req)
		if err != nil {
			t.Fatalf("CreateLoan gagal: %v", err)
		}

		if loan.ID == 0 {
			t.Errorf("ID pinjaman tidak boleh 0")
		}
		if loan.UserID != user.ID {
			t.Errorf("UserID mismatch: didapat %d, diharapkan %d", loan.UserID, user.ID)
		}
		if loan.BookTitle != req.BookTitle {
			t.Errorf("BookTitle mismatch: didapat %s, diharapkan %s", loan.BookTitle, req.BookTitle)
		}
		if loan.Status != "borrowed" {
			t.Errorf("Status mismatch: didapat %s, diharapkan 'borrowed'", loan.Status)
		}
		if loan.User == nil || loan.User.Email != user.Email {
			t.Errorf("User relasi tidak terisi dengan benar")
		}

		diffDays := int(loan.DueDate.Sub(loan.LoanDate).Hours() / 24)
		if diffDays != 7 {
			t.Errorf("Durasi default diharapkan 7 hari, didapat %d hari", diffDays)
		}
	})

	t.Run("Peminjaman berhasil dengan durasi kustom (14 hari)", func(t *testing.T) {
		req := models.CreateLoanRequest{
			BookTitle:    "Struktur Data & Algoritma",
			DurationDays: 14,
		}
		loan, err := loanService.CreateLoan(ctx, user.ID, req)
		if err != nil {
			t.Fatalf("CreateLoan gagal: %v", err)
		}

		diffDays := int(loan.DueDate.Sub(loan.LoanDate).Hours() / 24)
		if diffDays != 14 {
			t.Errorf("Durasi kustom diharapkan 14 hari, didapat %d hari", diffDays)
		}
	})

	t.Run("Peminjaman gagal untuk user yang tidak terdaftar", func(t *testing.T) {
		req := models.CreateLoanRequest{
			BookTitle: "Buku Tidak Ditemukan",
		}
		_, err := loanService.CreateLoan(ctx, 99999, req)
		if err == nil {
			t.Fatalf("Diharapkan error user tidak ditemukan, tapi berhasil")
		}
		if err != services.ErrUserNotFound {
			t.Errorf("Diharapkan ErrUserNotFound, didapat: %v", err)
		}
	})
}

func TestLoanService_GetUserLoans(t *testing.T) {
	loanService, authService := setupLoanTestDB(t)
	ctx := context.Background()

	user1, _, err := authService.Register(ctx, models.RegisterRequest{
		Name:     "Doni",
		Email:    "doni@perpus.local",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Gagal registrasi user1: %v", err)
	}

	user2, _, err := authService.Register(ctx, models.RegisterRequest{
		Name:     "Eka",
		Email:    "eka@perpus.local",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Gagal registrasi user2: %v", err)
	}

	// Buat 2 peminjaman untuk user1
	_, err = loanService.CreateLoan(ctx, user1.ID, models.CreateLoanRequest{BookTitle: "Buku A"})
	if err != nil {
		t.Fatalf("Gagal membuat pinjaman 1: %v", err)
	}
	_, err = loanService.CreateLoan(ctx, user1.ID, models.CreateLoanRequest{BookTitle: "Buku B"})
	if err != nil {
		t.Fatalf("Gagal membuat pinjaman 2: %v", err)
	}

	// Ambil pinjaman user1
	loans1, err := loanService.GetUserLoans(ctx, user1.ID)
	if err != nil {
		t.Fatalf("GetUserLoans user1 gagal: %v", err)
	}
	if len(loans1) != 2 {
		t.Errorf("Expected 2 loans for user1, got %d", len(loans1))
	}

	// Ambil pinjaman user2 (belum ada peminjaman)
	loans2, err := loanService.GetUserLoans(ctx, user2.ID)
	if err != nil {
		t.Fatalf("GetUserLoans user2 gagal: %v", err)
	}
	if len(loans2) != 0 {
		t.Errorf("Expected 0 loans for user2, got %d", len(loans2))
	}
}

func TestLoanService_GetLoanByID(t *testing.T) {
	loanService, authService := setupLoanTestDB(t)
	ctx := context.Background()

	user, _, err := authService.Register(ctx, models.RegisterRequest{
		Name:     "Fajar",
		Email:    "fajar@perpus.local",
		Password: "password123",
		Height:   170,
	})
	if err != nil {
		t.Fatalf("Gagal registrasi user: %v", err)
	}

	created, err := loanService.CreateLoan(ctx, user.ID, models.CreateLoanRequest{BookTitle: "Clean Architecture"})
	if err != nil {
		t.Fatalf("Gagal create loan: %v", err)
	}

	loan, err := loanService.GetLoanByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetLoanByID gagal: %v", err)
	}
	if loan.BookTitle != "Clean Architecture" {
		t.Errorf("Expected book title 'Clean Architecture', got '%s'", loan.BookTitle)
	}
	if loan.User == nil || loan.User.Name != "Fajar" {
		t.Errorf("User relation tidak terpanggil dengan benar")
	}

	_, err = loanService.GetLoanByID(ctx, 88888)
	if err != services.ErrLoanNotFound {
		t.Errorf("Diharapkan ErrLoanNotFound untuk ID pinjaman tidak valid, didapat: %v", err)
	}
}

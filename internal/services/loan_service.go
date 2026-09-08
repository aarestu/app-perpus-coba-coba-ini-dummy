package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"app-perpus/internal/models"
)

var (
	ErrLoanNotFound = errors.New("data peminjaman tidak ditemukan")
)

type LoanService struct {
	db *sql.DB
}

// NewLoanService membuat instance baru LoanService
func NewLoanService(db *sql.DB) *LoanService {
	return &LoanService{
		db: db,
	}
}

// CreateLoan memvalidasi peminjam dan mencatat peminjaman buku baru
func (s *LoanService) CreateLoan(ctx context.Context, userID int64, req models.CreateLoanRequest) (*models.Loan, error) {
	// Pastikan user ada dan valid
	user := &models.User{}
	userQuery := `SELECT id, name, email, role, height, created_at, updated_at FROM users WHERE id = ?`
	err := s.db.QueryRowContext(ctx, userQuery, userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Role, &user.Height, &user.CreatedAt, &user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, fmt.Errorf("gagal memvalidasi data pengguna: %w", err)
	}

	durationDays := req.DurationDays
	if durationDays <= 0 {
		durationDays = 7 // Durasi peminjaman standar 7 hari
	}

	now := time.Now().Truncate(time.Second)
	dueDate := now.AddDate(0, 0, durationDays)
	bookTitle := strings.TrimSpace(req.BookTitle)

	insertQuery := `
		INSERT INTO loans (user_id, book_title, loan_date, due_date, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, 'borrowed', ?, ?)
	`
	res, err := s.db.ExecContext(ctx, insertQuery, userID, bookTitle, now, dueDate, now, now)
	if err != nil {
		return nil, fmt.Errorf("gagal menyimpan data peminjaman: %w", err)
	}

	loanID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil ID peminjaman: %w", err)
	}

	loan := &models.Loan{
		ID:        loanID,
		UserID:    userID,
		BookTitle: bookTitle,
		LoanDate:  now,
		DueDate:   dueDate,
		Status:    "borrowed",
		CreatedAt: now,
		UpdatedAt: now,
		User:      user,
	}

	return loan, nil
}

// GetUserLoans mengambil seluruh riwayat peminjaman buku milik pengguna
func (s *LoanService) GetUserLoans(ctx context.Context, userID int64) ([]models.Loan, error) {
	query := `
		SELECT id, user_id, book_title, loan_date, due_date, return_date, status, created_at, updated_at
		FROM loans
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar peminjaman: %w", err)
	}
	defer rows.Close()

	loans := make([]models.Loan, 0)
	for rows.Next() {
		var l models.Loan
		var returnDate sql.NullTime
		if err := rows.Scan(&l.ID, &l.UserID, &l.BookTitle, &l.LoanDate, &l.DueDate, &returnDate, &l.Status, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, fmt.Errorf("gagal membaca data baris peminjaman: %w", err)
		}
		if returnDate.Valid {
			l.ReturnDate = &returnDate.Time
		}
		loans = append(loans, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error pada iterasi data peminjaman: %w", err)
	}

	return loans, nil
}

// GetLoanByID mengambil detail peminjaman buku beserta informasi user terkait
func (s *LoanService) GetLoanByID(ctx context.Context, loanID int64) (*models.Loan, error) {
	query := `
		SELECT l.id, l.user_id, l.book_title, l.loan_date, l.due_date, l.return_date, l.status, l.created_at, l.updated_at,
		       u.id, u.name, u.email, u.role, u.height, u.created_at, u.updated_at
		FROM loans l
		INNER JOIN users u ON l.user_id = u.id
		WHERE l.id = ?
	`
	var l models.Loan
	var u models.User
	var returnDate sql.NullTime
	err := s.db.QueryRowContext(ctx, query, loanID).Scan(
		&l.ID, &l.UserID, &l.BookTitle, &l.LoanDate, &l.DueDate, &returnDate, &l.Status, &l.CreatedAt, &l.UpdatedAt,
		&u.ID, &u.Name, &u.Email, &u.Role, &u.Height, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrLoanNotFound
	} else if err != nil {
		return nil, fmt.Errorf("gagal mengambil data peminjaman: %w", err)
	}

	if returnDate.Valid {
		l.ReturnDate = &returnDate.Time
	}
	l.User = &u

	return &l, nil
}

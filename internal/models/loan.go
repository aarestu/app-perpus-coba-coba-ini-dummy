package models

import (
	"encoding/json"
	"strings"
	"time"
)

// Loan merepresentasikan entitas peminjaman buku dalam sistem
type Loan struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	BookTitle  string     `json:"book_title"`
	LoanDate   time.Time  `json:"loan_date"`
	DueDate    time.Time  `json:"due_date"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	User       *User      `json:"user,omitempty"`
}

// CreateLoanRequest merepresentasikan payload body permohonan peminjaman buku
type CreateLoanRequest struct {
	BookTitle    string `json:"book_title"`
	DurationDays int    `json:"duration_days,omitempty"`
}

// UnmarshalJSON mendukung mapping dari field 'book_title' maupun 'judul_buku'
func (r *CreateLoanRequest) UnmarshalJSON(data []byte) error {
	type Alias CreateLoanRequest
	aux := &struct {
		JudulBuku *string `json:"judul_buku"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.JudulBuku != nil && strings.TrimSpace(r.BookTitle) == "" {
		r.BookTitle = *aux.JudulBuku
	}
	return nil
}

// Validate memvalidasi kelayakan input CreateLoanRequest
func (r *CreateLoanRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.BookTitle) == "" {
		errs["book_title"] = "Judul buku wajib diisi"
	}
	if r.DurationDays < 0 {
		errs["duration_days"] = "Durasi peminjaman tidak boleh bernilai negatif"
	}
	return errs
}

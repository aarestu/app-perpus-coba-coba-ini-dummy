package database

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaSQL string

// InitDB menginisialisasi koneksi SQLite dan menjalankan migrasi skema tabel
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal ping database: %w", err)
	}

	// Aktifkan foreign keys dan journal mode WAL pada SQLite
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, fmt.Errorf("gagal mengaktifkan PRAGMA foreign_keys: %w", err)
	}

	// Jalankan migrasi skema
	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, fmt.Errorf("gagal mengeksekusi skema database: %w", err)
	}

	return db, nil
}

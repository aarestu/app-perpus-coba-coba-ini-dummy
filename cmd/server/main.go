package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app-perpus/internal/database"
	"app-perpus/internal/handlers"
	"app-perpus/internal/middleware"
	"app-perpus/internal/services"
	"app-perpus/internal/utils"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "app.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecret-perpus-app-key-replace-in-prod"
	}

	// 1. Inisialisasi Database SQLite
	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Fatal: Gagal menginisialisasi database: %v", err)
	}
	defer db.Close()
	log.Printf("Terkoneksi ke basis data SQLite: %s", dbPath)

	// 2. Inisialisasi Service & Handler
	authService := services.NewAuthService(db, jwtSecret, 24*time.Hour)
	authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := middleware.AuthMiddleware(authService)

	loanService := services.NewLoanService(db)
	loanHandler := handlers.NewLoanHandler(loanService)

	// 3. Router & Rute API
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		utils.JSONSuccess(w, http.StatusOK, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Public auth routes
	mux.HandleFunc("/api/auth/register", authHandler.Register)
	mux.HandleFunc("/api/auth/login", authHandler.Login)

	// Protected routes
	mux.Handle("/api/auth/me", authMiddleware(http.HandlerFunc(authHandler.Me)))
	mux.Handle("/api/loans", authMiddleware(http.HandlerFunc(loanHandler.HandleLoans)))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel untuk menangani graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server perpustakaan berjalan pada port %s (http://localhost:%s)", port, port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error saat menjalankan server: %v", err)
		}
	}()

	<-stopChan
	log.Println("Menerima sinyal terminasi, mematikan server secara graceful...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server gagal shutdown secara graceful: %v", err)
	}
	log.Println("Server berhasil dimatikan dengan aman.")
}

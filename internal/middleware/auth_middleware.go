package middleware

import (
	"context"
	"net/http"
	"strings"

	"app-perpus/internal/models"
	"app-perpus/internal/services"
	"app-perpus/internal/utils"
)

type contextKey string

const (
	UserClaimsContextKey contextKey = "user_claims"
)

// AuthMiddleware memeriksa keberadaan dan keabsahan token JWT pada header Authorization
func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.JSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Header Authorization tidak ditemukan")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				utils.JSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Format token tidak valid. Gunakan format 'Bearer <token>'")
				return
			}

			claims, err := authService.ValidateToken(parts[1])
			if err != nil {
				utils.JSONError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Token autentikasi tidak valid atau telah kedaluwarsa")
				return
			}

			ctx := context.WithValue(r.Context(), UserClaimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaimsFromContext mengambil data klaim JWT dari request context
func GetClaimsFromContext(ctx context.Context) (*models.JWTClaims, bool) {
	claims, ok := ctx.Value(UserClaimsContextKey).(*models.JWTClaims)
	return claims, ok
}

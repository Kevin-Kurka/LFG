package middleware

import (
	"net/http"
	"strings"

	"lfg/shared/auth"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

// Authenticate validates JWT token and injects user ID into headers
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Check Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			respondError(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			if err == auth.ErrExpiredToken {
				respondError(w, "Token has expired", http.StatusUnauthorized)
				return
			}
			respondError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Inject user ID and email into request headers for downstream services
		r.Header.Set("X-User-ID", claims.UserID.String())
		r.Header.Set("X-User-Email", claims.Email)
		if claims.Role != "" {
			r.Header.Set("X-User-Role", claims.Role)
		}

		// Forward to next handler
		next.ServeHTTP(w, r)
	})
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error":"` + message + `"}`))
}

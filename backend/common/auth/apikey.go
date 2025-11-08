package auth

import (
	"log"
	"net/http"
	"os"
)

// GetInternalAPIKey returns the internal API key from environment
// This is used for service-to-service authentication
func GetInternalAPIKey() string {
	apiKey := os.Getenv("INTERNAL_API_KEY")
	if apiKey == "" {
		log.Fatal("SECURITY ERROR: INTERNAL_API_KEY environment variable is not set. Application cannot start without a secure internal API key.")
	}
	if len(apiKey) < 32 {
		log.Fatalf("SECURITY ERROR: INTERNAL_API_KEY must be at least 32 characters long. Current length: %d", len(apiKey))
	}
	return apiKey
}

// InternalAPIKeyMiddleware validates the internal API key for service-to-service calls
func InternalAPIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Internal-API-Key")
		if apiKey == "" {
			http.Error(w, `{"error":"unauthorized","message":"missing internal API key"}`, http.StatusUnauthorized)
			return
		}

		expectedKey := GetInternalAPIKey()
		if apiKey != expectedKey {
			http.Error(w, `{"error":"unauthorized","message":"invalid internal API key"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetAdminAPIKey returns the admin API key from environment
// This is used for admin operations
func GetAdminAPIKey() string {
	apiKey := os.Getenv("ADMIN_API_KEY")
	if apiKey == "" {
		log.Fatal("SECURITY ERROR: ADMIN_API_KEY environment variable is not set. Application cannot start without a secure admin API key.")
	}
	if len(apiKey) < 32 {
		log.Fatalf("SECURITY ERROR: ADMIN_API_KEY must be at least 32 characters long. Current length: %d", len(apiKey))
	}
	return apiKey
}

// AdminAPIKeyMiddleware validates the admin API key for admin operations
func AdminAPIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Admin-API-Key")
		if apiKey == "" {
			http.Error(w, `{"error":"unauthorized","message":"missing admin API key"}`, http.StatusUnauthorized)
			return
		}

		expectedKey := GetAdminAPIKey()
		if apiKey != expectedKey {
			http.Error(w, `{"error":"unauthorized","message":"invalid admin API key"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

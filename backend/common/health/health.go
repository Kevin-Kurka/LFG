package health

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

// Status represents the health status of the service
type Status struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Service   string            `json:"service"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// LivenessHandler returns a simple liveness check (always returns OK if service is running)
func LivenessHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := Status{
			Status:    "alive",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Service:   serviceName,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}
}

// ReadinessHandler returns a readiness check that verifies database connectivity
func ReadinessHandler(serviceName string, db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checks := make(map[string]string)
		healthy := true

		// Check database connection if provided
		if db != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			if err := db.PingContext(ctx); err != nil {
				checks["database"] = "unhealthy: " + err.Error()
				healthy = false
			} else {
				checks["database"] = "healthy"
			}
		}

		status := Status{
			Service:   serviceName,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Checks:    checks,
		}

		w.Header().Set("Content-Type", "application/json")

		if healthy {
			status.Status = "ready"
			w.WriteHeader(http.StatusOK)
		} else {
			status.Status = "not ready"
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(status)
	}
}

// SimpleReadinessHandler returns a simple readiness check for services without database
func SimpleReadinessHandler(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := Status{
			Status:    "ready",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Service:   serviceName,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}
}

// CheckDatabaseHealth checks if database is healthy
func CheckDatabaseHealth(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}

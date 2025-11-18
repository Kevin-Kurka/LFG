package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Status represents health check status
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusDegraded  Status = "degraded"
)

// CheckResult represents a health check result
type CheckResult struct {
	Status    Status                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]CheckDetail `json:"checks"`
	Version   string                 `json:"version,omitempty"`
}

// CheckDetail represents details of a specific health check
type CheckDetail struct {
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"`
}

// Checker performs health checks
type Checker struct {
	db      *pgxpool.Pool
	redis   *redis.Client
	version string
	mu      sync.RWMutex
}

// NewChecker creates a new health checker
func NewChecker(db *pgxpool.Pool, redis *redis.Client, version string) *Checker {
	return &Checker{
		db:      db,
		redis:   redis,
		version: version,
	}
}

// Health returns HTTP handler for health check endpoint
func (c *Checker) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := c.Check(r.Context())

		w.Header().Set("Content-Type", "application/json")

		if result.Status == StatusHealthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(result)
	}
}

// Ready returns HTTP handler for readiness check endpoint
func (c *Checker) Ready() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := c.CheckReadiness(r.Context())

		w.Header().Set("Content-Type", "application/json")

		if result.Status == StatusHealthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		json.NewEncoder(w).Encode(result)
	}
}

// Check performs comprehensive health check
func (c *Checker) Check(ctx context.Context) CheckResult {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := CheckResult{
		Timestamp: time.Now(),
		Checks:    make(map[string]CheckDetail),
		Version:   c.version,
	}

	// Check database
	if c.db != nil {
		dbCheck := c.checkDatabase(ctx)
		result.Checks["database"] = dbCheck
	}

	// Check Redis
	if c.redis != nil {
		redisCheck := c.checkRedis(ctx)
		result.Checks["redis"] = redisCheck
	}

	// Determine overall status
	result.Status = c.calculateOverallStatus(result.Checks)

	return result
}

// CheckReadiness performs readiness check (lighter than health check)
func (c *Checker) CheckReadiness(ctx context.Context) CheckResult {
	result := CheckResult{
		Timestamp: time.Now(),
		Checks:    make(map[string]CheckDetail),
	}

	// Quick database ping
	if c.db != nil {
		start := time.Now()
		err := c.db.Ping(ctx)
		latency := time.Since(start)

		if err == nil {
			result.Checks["database"] = CheckDetail{
				Status:  StatusHealthy,
				Latency: latency.String(),
			}
		} else {
			result.Checks["database"] = CheckDetail{
				Status:  StatusUnhealthy,
				Message: "database not reachable",
			}
		}
	}

	result.Status = c.calculateOverallStatus(result.Checks)

	return result
}

// checkDatabase performs database health check
func (c *Checker) checkDatabase(ctx context.Context) CheckDetail {
	start := time.Now()

	// Ping database
	err := c.db.Ping(ctx)
	latency := time.Since(start)

	if err != nil {
		return CheckDetail{
			Status:  StatusUnhealthy,
			Message: err.Error(),
			Latency: latency.String(),
		}
	}

	// Check connection pool
	stats := c.db.Stat()
	if stats.TotalConns() > 0 && stats.IdleConns() == 0 {
		return CheckDetail{
			Status:  StatusDegraded,
			Message: "no idle connections available",
			Latency: latency.String(),
		}
	}

	return CheckDetail{
		Status:  StatusHealthy,
		Latency: latency.String(),
	}
}

// checkRedis performs Redis health check
func (c *Checker) checkRedis(ctx context.Context) CheckDetail {
	start := time.Now()

	err := c.redis.Ping(ctx).Err()
	latency := time.Since(start)

	if err != nil {
		return CheckDetail{
			Status:  StatusUnhealthy,
			Message: err.Error(),
			Latency: latency.String(),
		}
	}

	return CheckDetail{
		Status:  StatusHealthy,
		Latency: latency.String(),
	}
}

// calculateOverallStatus calculates overall health status
func (c *Checker) calculateOverallStatus(checks map[string]CheckDetail) Status {
	if len(checks) == 0 {
		return StatusHealthy
	}

	hasUnhealthy := false
	hasDegraded := false

	for _, check := range checks {
		switch check.Status {
		case StatusUnhealthy:
			hasUnhealthy = true
		case StatusDegraded:
			hasDegraded = true
		}
	}

	if hasUnhealthy {
		return StatusUnhealthy
	}

	if hasDegraded {
		return StatusDegraded
	}

	return StatusHealthy
}

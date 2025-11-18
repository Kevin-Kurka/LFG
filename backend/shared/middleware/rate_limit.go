package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"lfg/shared/errors"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
// ratePerSecond: requests per second
// burst: maximum burst size
func NewRateLimiter(ratePerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(ratePerSecond),
		burst:    burst,
	}
}

// getLimiter returns limiter for a key (IP or user ID)
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check after acquiring write lock
		limiter, exists = rl.limiters[key]
		if !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.limiters[key] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

// Allow checks if request is allowed
func (rl *RateLimiter) Allow(key string) bool {
	limiter := rl.getLimiter(key)
	return limiter.Allow()
}

// Cleanup removes old limiters (run periodically)
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Simple cleanup: clear all if too many entries
	if len(rl.limiters) > 10000 {
		rl.limiters = make(map[string]*rate.Limiter)
	}
}

// RateLimitMiddleware creates rate limiting middleware
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	// Start cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			limiter.Cleanup()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get key (prefer user ID, fallback to IP)
			key := r.Header.Get("X-User-ID")
			if key == "" {
				key = getIP(r)
			}

			if !limiter.Allow(key) {
				errors.RespondWithError(w, errors.ErrRateLimited())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getIP extracts IP from request
func getIP(r *http.Request) string {
	// Check X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr
}

package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	requests int
	window   time.Duration
	clients  map[string]*clientBucket
	mu       sync.RWMutex
	cleanup  *time.Ticker
}

type clientBucket struct {
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: requests,
		window:   window,
		clients:  make(map[string]*clientBucket),
		cleanup:  time.NewTicker(time.Minute),
	}

	// Start cleanup goroutine
	go rl.cleanupExpired()

	return rl
}

// Limit applies rate limiting to requests
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get client identifier (IP address)
		clientIP := getClientIP(r)

		// Check rate limit
		if !rl.allow(clientIP) {
			respondError(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Forward to next handler
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(clientIP string) bool {
	rl.mu.Lock()
	bucket, exists := rl.clients[clientIP]
	if !exists {
		bucket = &clientBucket{
			tokens:     rl.requests,
			lastRefill: time.Now(),
		}
		rl.clients[clientIP] = bucket
	}
	rl.mu.Unlock()

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)
	if elapsed >= rl.window {
		bucket.tokens = rl.requests
		bucket.lastRefill = now
	}

	// Check if request is allowed
	if bucket.tokens > 0 {
		bucket.tokens--
		return true
	}

	return false
}

func (rl *RateLimiter) cleanupExpired() {
	for range rl.cleanup.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, bucket := range rl.clients {
			bucket.mu.Lock()
			if now.Sub(bucket.lastRefill) > rl.window*2 {
				delete(rl.clients, ip)
			}
			bucket.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (proxy/load balancer)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr
}

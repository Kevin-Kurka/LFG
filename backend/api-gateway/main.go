package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/middleware"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

func main() {
	log.Println("Starting API Gateway...")

	// Define the target URLs for the microservices
	userServiceURL, _ := url.Parse("http://localhost:8080")
	walletServiceURL, _ := url.Parse("http://localhost:8081")
	orderServiceURL, _ := url.Parse("http://localhost:8082")
	marketServiceURL, _ := url.Parse("http://localhost:8083")
	creditExchangeURL, _ := url.Parse("http://localhost:8084")
	notificationServiceURL, _ := url.Parse("http://localhost:8085")
	sportsbookServiceURL, _ := url.Parse("http://localhost:8086")
	cryptoServiceURL, _ := url.Parse("http://localhost:8087")

	// Create reverse proxies for each service
	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	walletProxy := httputil.NewSingleHostReverseProxy(walletServiceURL)
	orderProxy := httputil.NewSingleHostReverseProxy(orderServiceURL)
	marketProxy := httputil.NewSingleHostReverseProxy(marketServiceURL)
	creditExchangeProxy := httputil.NewSingleHostReverseProxy(creditExchangeURL)
	notificationProxy := httputil.NewSingleHostReverseProxy(notificationServiceURL)
	sportsbookProxy := httputil.NewSingleHostReverseProxy(sportsbookServiceURL)
	cryptoProxy := httputil.NewSingleHostReverseProxy(cryptoServiceURL)

	// Create rate limiter (10 requests per second with burst of 20)
	rateLimiter := middleware.NewIPRateLimiter(rate.Limit(10), 20)

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Public routes (no auth required, but rate limited)
	mux.Handle("/register", applyPublicMiddleware(rateLimiter, userProxy))
	mux.Handle("/login", applyPublicMiddleware(rateLimiter, userProxy))
	mux.Handle("/markets", applyPublicMiddleware(rateLimiter, marketProxy))
	mux.Handle("/markets/", applyPublicMiddleware(rateLimiter, marketProxy))

	// Protected routes (require JWT auth + rate limiting)
	mux.Handle("/profile", applyAuthMiddleware(rateLimiter, userProxy))
	mux.Handle("/balance", applyAuthMiddleware(rateLimiter, walletProxy))
	mux.Handle("/transactions", applyAuthMiddleware(rateLimiter, walletProxy))
	mux.Handle("/orders/", applyAuthMiddleware(rateLimiter, orderProxy))
	mux.Handle("/exchange/", applyAuthMiddleware(rateLimiter, creditExchangeProxy))
	mux.Handle("/sportsbook/", applyAuthMiddleware(rateLimiter, sportsbookProxy))
	mux.Handle("/crypto/", applyAuthMiddleware(rateLimiter, cryptoProxy))

	// WebSocket (rate limited only)
	mux.Handle("/ws", applyPublicMiddleware(rateLimiter, notificationProxy))

	// CORS configuration
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	origins := []string{"http://localhost:3000"} // Default for development
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr:         ":8000",
		Handler:      corsHandler.Handler(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Println("API Gateway listening on port 8000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// applyPublicMiddleware applies rate limiting to public endpoints
func applyPublicMiddleware(rateLimiter *middleware.IPRateLimiter, next http.Handler) http.Handler {
	return middleware.RateLimitMiddleware(rateLimiter)(next)
}

// applyAuthMiddleware applies JWT authentication and rate limiting
func applyAuthMiddleware(rateLimiter *middleware.IPRateLimiter, next http.Handler) http.Handler {
	return middleware.RateLimitMiddleware(rateLimiter)(auth.AuthMiddleware(next))
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/middleware"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

func main() {
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

	fmt.Println("API Gateway listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", corsHandler.Handler(mux)))
}

// applyPublicMiddleware applies rate limiting to public endpoints
func applyPublicMiddleware(rateLimiter *middleware.IPRateLimiter, next http.Handler) http.Handler {
	return middleware.RateLimitMiddleware(rateLimiter)(next)
}

// applyAuthMiddleware applies JWT authentication and rate limiting
func applyAuthMiddleware(rateLimiter *middleware.IPRateLimiter, next http.Handler) http.Handler {
	return middleware.RateLimitMiddleware(rateLimiter)(auth.AuthMiddleware(next))
}

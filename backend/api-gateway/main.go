package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lfg/shared/auth"
	"lfg/shared/config"
	"lfg/api-gateway/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRequests, cfg.RateLimitWindow)
	corsMiddleware := middleware.NewCORSMiddleware(cfg.CORSAllowedOrigins)

	// Parse service URLs
	userServiceURL, err := url.Parse(cfg.UserServiceURL)
	if err != nil {
		log.Fatalf("Invalid user service URL: %v", err)
	}

	walletServiceURL, err := url.Parse(cfg.WalletServiceURL)
	if err != nil {
		log.Fatalf("Invalid wallet service URL: %v", err)
	}

	orderServiceURL, err := url.Parse(cfg.OrderServiceURL)
	if err != nil {
		log.Fatalf("Invalid order service URL: %v", err)
	}

	marketServiceURL, err := url.Parse(cfg.MarketServiceURL)
	if err != nil {
		log.Fatalf("Invalid market service URL: %v", err)
	}

	creditExchangeURL, err := url.Parse(cfg.CreditExchangeURL)
	if err != nil {
		log.Fatalf("Invalid credit exchange URL: %v", err)
	}

	notificationServiceURL, err := url.Parse(cfg.NotificationServiceURL)
	if err != nil {
		log.Fatalf("Invalid notification service URL: %v", err)
	}

	// Create reverse proxies
	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	walletProxy := httputil.NewSingleHostReverseProxy(walletServiceURL)
	orderProxy := httputil.NewSingleHostReverseProxy(orderServiceURL)
	marketProxy := httputil.NewSingleHostReverseProxy(marketServiceURL)
	creditExchangeProxy := httputil.NewSingleHostReverseProxy(creditExchangeURL)
	notificationProxy := httputil.NewSingleHostReverseProxy(notificationServiceURL)

	// Setup routes
	mux := http.NewServeMux()

	// Health check (no auth required)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Public endpoints (no auth required, but rate limited)
	mux.Handle("/register", applyMiddleware(userProxy, rateLimiter))
	mux.Handle("/login", applyMiddleware(userProxy, rateLimiter))

	// Protected endpoints (auth + rate limiting)
	mux.Handle("/profile", applyMiddleware(userProxy, rateLimiter, authMiddleware))
	mux.Handle("/balance", applyMiddleware(walletProxy, rateLimiter, authMiddleware))
	mux.Handle("/transactions", applyMiddleware(walletProxy, rateLimiter, authMiddleware))
	mux.Handle("/orders/", applyMiddleware(orderProxy, rateLimiter, authMiddleware))
	mux.Handle("/exchange/", applyMiddleware(creditExchangeProxy, rateLimiter, authMiddleware))

	// Public market endpoints (rate limited, no auth)
	mux.Handle("/markets", applyMiddleware(marketProxy, rateLimiter))
	mux.Handle("/markets/", applyMiddleware(marketProxy, rateLimiter))

	// WebSocket endpoint (auth required)
	mux.Handle("/ws", applyMiddleware(notificationProxy, authMiddleware))

	// Apply CORS to all routes
	handler := corsMiddleware.Handle(mux)

	// Create server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("API Gateway listening on port %s...\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("API Gateway exited")
}

// applyMiddleware chains middleware and the final handler
func applyMiddleware(handler http.Handler, middlewares ...interface{}) http.Handler {
	result := handler

	// Apply middleware in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		switch mw := middlewares[i].(type) {
		case *middleware.AuthMiddleware:
			result = mw.Authenticate(result)
		case *middleware.RateLimiter:
			result = mw.Limit(result)
		}
	}

	return result
}

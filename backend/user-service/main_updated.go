package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"lfg/shared/auth"
	"lfg/shared/cache"
	"lfg/shared/config"
	"lfg/shared/db"
	"lfg/shared/health"
	"lfg/shared/logging"
	"lfg/shared/middleware"
	"lfg/user-service/handlers"
	"lfg/user-service/repository"
)

const serviceName = "user-service"
const version = "1.0.0"

func main() {
	// Initialize logger
	logger := logging.NewLogger(logging.Config{
		Service: serviceName,
		Level:   getEnv("LOG_LEVEL", "info"),
		Pretty:  getEnv("ENVIRONMENT", "production") == "development",
	})

	logger.Info("Starting user service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Validate JWT secret in production
	if cfg.Env == "production" && len(cfg.JWTSecret) < 32 {
		logger.Fatal("JWT secret must be at least 32 characters in production")
	}

	ctx := context.Background()

	// Initialize database connection
	dbCfg := db.Config{
		Host:            cfg.DBHost,
		Port:            cfg.DBPort,
		User:            cfg.DBUser,
		Password:        cfg.DBPassword,
		Database:        cfg.DBName,
		SSLMode:         cfg.DBSSLMode,
		MaxConns:        cfg.DBMaxConns,
		MinConns:        cfg.DBMinConns,
		MaxConnLifetime: 1 * time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}

	pool, err := db.NewPool(ctx, dbCfg)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close(pool)

	logger.Info("Connected to database successfully")

	// Initialize Redis cache
	redisCache, err := cache.NewCache(cache.Config{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	if err != nil {
		logger.Warnf("Failed to connect to Redis: %v (continuing without cache)", err)
		redisCache = nil
	} else {
		defer redisCache.Close()
		logger.Info("Connected to Redis successfully")
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)

	// Initialize repository
	userRepo := repository.NewUserRepository(pool)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo, jwtManager, redisCache, logger)

	// Initialize health checker
	healthChecker := health.NewChecker(pool, redisCache.Client(), version)

	// Create rate limiter (100 requests per minute per user/IP)
	rateLimiter := middleware.NewRateLimiter(100.0/60.0, 10)

	// Setup HTTP routes
	mux := http.NewServeMux()

	// Health and metrics endpoints (no auth required)
	mux.HandleFunc("/health", healthChecker.Health())
	mux.HandleFunc("/ready", healthChecker.Ready())
	mux.Handle("/metrics", promhttp.Handler())

	// Public endpoints
	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)

	// Protected endpoints (require auth)
	mux.HandleFunc("/profile", userHandler.Profile)
	mux.HandleFunc("/refresh", userHandler.RefreshToken)

	// Apply middleware chain
	handler := middleware.SecurityHeaders(mux)
	handler = middleware.CORS(cfg.CORSAllowedOrigins)(handler)
	handler = middleware.RateLimitMiddleware(rateLimiter)(handler)
	handler = middleware.RequestLogging(logger, serviceName)(handler)

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Infof("%s listening on port %s...", serviceName, cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with 30 second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Errorf(err, "Server forced to shutdown")
	}

	logger.Info("Server exited gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Environment string
	LogLevel    string

	// Server
	Port string

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBMaxConns int32
	DBMinConns int32

	// JWT
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	// NATS
	NATSURL string

	// Service URLs
	UserServiceURL         string
	WalletServiceURL       string
	OrderServiceURL        string
	MarketServiceURL       string
	CreditExchangeURL      string
	NotificationServiceURL string
	MatchingEngineGRPC     string

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// CORS
	CORSAllowedOrigins []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),

		Port: getEnv("PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "lfg"),
		DBPassword: getEnv("DB_PASSWORD", "lfg_dev_password"),
		DBName:     getEnv("DB_NAME", "lfg"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		DBMaxConns: int32(getEnvAsInt("DB_MAX_CONNS", 25)),
		DBMinConns: int32(getEnvAsInt("DB_MIN_CONNS", 5)),

		JWTSecret:     getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
		JWTAccessTTL:  getEnvAsDuration("JWT_ACCESS_TTL", 15*time.Minute),
		JWTRefreshTTL: getEnvAsDuration("JWT_REFRESH_TTL", 7*24*time.Hour),

		NATSURL: getEnv("NATS_URL", "nats://localhost:4222"),

		UserServiceURL:         getEnv("USER_SERVICE_URL", "http://localhost:8080"),
		WalletServiceURL:       getEnv("WALLET_SERVICE_URL", "http://localhost:8081"),
		OrderServiceURL:        getEnv("ORDER_SERVICE_URL", "http://localhost:8082"),
		MarketServiceURL:       getEnv("MARKET_SERVICE_URL", "http://localhost:8083"),
		CreditExchangeURL:      getEnv("CREDIT_EXCHANGE_URL", "http://localhost:8084"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8085"),
		MatchingEngineGRPC:     getEnv("MATCHING_ENGINE_GRPC", "localhost:50051"),

		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsDuration("RATE_LIMIT_WINDOW", 1*time.Minute),

		CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
	}

	// Validate required fields
	if cfg.JWTSecret == "dev-secret-key-change-in-production" && cfg.Environment == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	// Simple split by comma
	var result []string
	for i, j := 0, 0; j <= len(valueStr); j++ {
		if j == len(valueStr) || valueStr[j] == ',' {
			if j > i {
				result = append(result, valueStr[i:j])
			}
			i = j + 1
		}
	}
	return result
}

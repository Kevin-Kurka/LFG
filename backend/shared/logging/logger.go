package logging

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
)

// Logger wraps zerolog with additional context
type Logger struct {
	logger zerolog.Logger
}

// Config holds logger configuration
type Config struct {
	Service string
	Level   string
	Pretty  bool
}

// NewLogger creates a new structured logger
func NewLogger(cfg Config) *Logger {
	// Set log level
	level := zerolog.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	}

	zerolog.SetGlobalLevel(level)

	// Configure output
	var output zerolog.ConsoleWriter
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	// Create base logger
	var baseLogger zerolog.Logger
	if cfg.Pretty {
		baseLogger = log.Output(output)
	} else {
		baseLogger = zerolog.New(os.Stdout)
	}

	// Add default fields
	logger := baseLogger.With().
		Timestamp().
		Str("service", cfg.Service).
		Logger()

	return &Logger{logger: logger}
}

// WithRequestID adds request ID to logger context
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		logger: l.logger.With().Str("request_id", requestID).Logger(),
	}
}

// WithUserID adds user ID to logger context
func (l *Logger) WithUserID(userID string) *Logger {
	return &Logger{
		logger: l.logger.With().Str("user_id", userID).Logger(),
	}
}

// WithField adds a custom field to logger context
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		logger: l.logger.With().Interface(key, value).Logger(),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debug().Msgf(format, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Info().Msgf(format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warn().Msgf(format, args...)
}

// Error logs an error message
func (l *Logger) Error(err error, msg string) {
	l.logger.Error().Err(err).Msg(msg)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(err error, format string, args ...interface{}) {
	l.logger.Error().Err(err).Msgf(format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal().Msgf(format, args...)
}

// ContextWithRequestID adds request ID to context
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestIDFromContext extracts request ID from context
func RequestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// ContextWithUserID adds user ID to context
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// UserIDFromContext extracts user ID from context
func UserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

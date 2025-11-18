package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	// Client errors (4xx)
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeConflict       ErrorCode = "CONFLICT"
	ErrCodeBadRequest     ErrorCode = "BAD_REQUEST"
	ErrCodeRateLimited    ErrorCode = "RATE_LIMITED"

	// Server errors (5xx)
	ErrCodeInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrCodeDatabaseError  ErrorCode = "DATABASE_ERROR"
)

// AppError represents a structured application error
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// ToJSON converts error to JSON bytes
func (e *AppError) ToJSON() []byte {
	data, _ := json.Marshal(e)
	return data
}

// New creates a new AppError
func New(code ErrorCode, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// Predefined errors
var (
	ErrValidation = func(message string) *AppError {
		return New(ErrCodeValidation, message, http.StatusBadRequest)
	}

	ErrUnauthorized = func(message string) *AppError {
		return New(ErrCodeUnauthorized, message, http.StatusUnauthorized)
	}

	ErrForbidden = func(message string) *AppError {
		return New(ErrCodeForbidden, message, http.StatusForbidden)
	}

	ErrNotFound = func(resource string) *AppError {
		return New(ErrCodeNotFound, fmt.Sprintf("%s not found", resource), http.StatusNotFound)
	}

	ErrConflict = func(message string) *AppError {
		return New(ErrCodeConflict, message, http.StatusConflict)
	}

	ErrBadRequest = func(message string) *AppError {
		return New(ErrCodeBadRequest, message, http.StatusBadRequest)
	}

	ErrRateLimited = func() *AppError {
		return New(ErrCodeRateLimited, "Rate limit exceeded", http.StatusTooManyRequests)
	}

	ErrInternalServer = func(message string) *AppError {
		return New(ErrCodeInternalServer, message, http.StatusInternalServerError)
	}

	ErrServiceUnavailable = func(service string) *AppError {
		return New(ErrCodeServiceUnavailable, fmt.Sprintf("%s service unavailable", service), http.StatusServiceUnavailable)
	}

	ErrDatabase = func(operation string) *AppError {
		return New(ErrCodeDatabaseError, fmt.Sprintf("Database %s failed", operation), http.StatusInternalServerError)
	}
)

// RespondWithError writes an error response
func RespondWithError(w http.ResponseWriter, err *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)
	w.Write(err.ToJSON())
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// AsAppError converts an error to AppError
func AsAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return ErrInternalServer(err.Error())
}

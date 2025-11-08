package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with HTTP status code
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// NewAppError creates a new AppError
func NewAppError(code int, message string, err error) *AppError {
	appErr := &AppError{
		Code:    code,
		Message: message,
	}
	if err != nil {
		appErr.Error = err.Error()
	}
	return appErr
}

// BadRequest creates a 400 error
func BadRequest(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

// Unauthorized creates a 401 error
func Unauthorized(message string, err error) *AppError {
	return NewAppError(http.StatusUnauthorized, message, err)
}

// Forbidden creates a 403 error
func Forbidden(message string, err error) *AppError {
	return NewAppError(http.StatusForbidden, message, err)
}

// NotFound creates a 404 error
func NotFound(message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, message, err)
}

// Conflict creates a 409 error
func Conflict(message string, err error) *AppError {
	return NewAppError(http.StatusConflict, message, err)
}

// InternalServerError creates a 500 error
func InternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}

// ValidationError creates a validation error (400)
func ValidationError(field string, err error) *AppError {
	message := fmt.Sprintf("validation failed for field '%s'", field)
	if err != nil {
		message = fmt.Sprintf("validation failed for field '%s': %s", field, err.Error())
	}
	return BadRequest(message, err)
}

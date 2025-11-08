package response

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kevin-Kurka/LFG/backend/common/errors"
)

// JSONResponse represents a standard JSON response
type JSONResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
}

// ErrorData represents error information in the response
type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// JSON sends a JSON response
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := JSONResponse{
		Success: statusCode >= 200 && statusCode < 300,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// Error sends an error JSON response
func Error(w http.ResponseWriter, appErr *errors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)

	response := JSONResponse{
		Success: false,
		Error: &ErrorData{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Error,
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}

// Success sends a success JSON response with data
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

// Created sends a 201 Created response
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}

// NoContent sends a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(w http.ResponseWriter, message string, err error) {
	Error(w, errors.BadRequest(message, err))
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(w http.ResponseWriter, message string, err error) {
	Error(w, errors.Unauthorized(message, err))
}

// NotFound sends a 404 Not Found response
func NotFound(w http.ResponseWriter, message string, err error) {
	Error(w, errors.NotFound(message, err))
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(w http.ResponseWriter, message string, err error) {
	Error(w, errors.InternalServerError(message, err))
}

// Conflict sends a 409 Conflict response
func Conflict(w http.ResponseWriter, message string, err error) {
	Error(w, errors.Conflict(message, err))
}

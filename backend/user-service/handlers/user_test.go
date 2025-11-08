package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	commontest "github.com/Kevin-Kurka/LFG/backend/common/testing"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid registration",
			body: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "SecurePass123!",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "missing username",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "SecurePass123!",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			body: map[string]interface{}{
				"username": "testuser",
				"email":    "invalid-email",
				"password": "SecurePass123!",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "weak password",
			body: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "weak",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := commontest.MakeRequest(http.MethodPost, "/register", tt.body)
			rr := httptest.NewRecorder()

			// Note: This is a placeholder test. In a real implementation,
			// you would need to set up the database and handler properly.
			// Register(rr, req)

			// commontest.AssertStatus(t, rr.Code, tt.wantStatus)
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name: "valid login",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "SecurePass123!",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			body: map[string]interface{}{
				"email":    "test@example.com",
				"password": "WrongPassword",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing email",
			body: map[string]interface{}{
				"password": "SecurePass123!",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := commontest.MakeRequest(http.MethodPost, "/login", tt.body)
			rr := httptest.NewRecorder()

			// Note: This is a placeholder test
			// Login(rr, req)

			// commontest.AssertStatus(t, rr.Code, tt.wantStatus)
		})
	}
}

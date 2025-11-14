package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"lfg/shared/auth"
	"lfg/shared/models"
	"lfg/user-service/repository"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	repo       *repository.UserRepository
	jwtManager *auth.JWTManager
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo *repository.UserRepository, jwtManager *auth.JWTManager) *UserHandler {
	return &UserHandler{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.UserRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate password strength
	if err := auth.ValidatePasswordStrength(req.Password); err != nil {
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Status:       models.UserStatusActive,
	}

	if err := h.repo.Create(r.Context(), user); err != nil {
		if err == repository.ErrEmailAlreadyExists {
			respondError(w, "Email already registered", http.StatusConflict)
			return
		}
		respondError(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Generate JWT tokens
	accessToken, expiresAt, err := h.jwtManager.GenerateToken(user.ID, user.Email, "user")
	if err != nil {
		respondError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, _, err := h.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		respondError(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.UserLoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresAt:    expiresAt,
	}

	respondJSON(w, response, http.StatusCreated)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user by email
	user, err := h.repo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			respondError(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}
		respondError(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	// Check if user is active
	if user.Status != models.UserStatusActive {
		respondError(w, "Account is not active", http.StatusForbidden)
		return
	}

	// Compare password
	if err := auth.ComparePassword(user.PasswordHash, req.Password); err != nil {
		respondError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT tokens
	accessToken, expiresAt, err := h.jwtManager.GenerateToken(user.ID, user.Email, "user")
	if err != nil {
		respondError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, _, err := h.jwtManager.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		respondError(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.UserLoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresAt:    expiresAt,
	}

	respondJSON(w, response, http.StatusOK)
}

// Profile handles getting user profile
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from context (set by auth middleware)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get user
	user, err := h.repo.GetByID(r.Context(), userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			respondError(w, "User not found", http.StatusNotFound)
			return
		}
		respondError(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	respondJSON(w, user, http.StatusOK)
}

// Health check handler
func Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{"status": "healthy"}, http.StatusOK)
}

// Helper functions
func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	respondJSON(w, map[string]string{"error": message}, statusCode)
}

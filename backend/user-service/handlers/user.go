package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/common/response"
	"github.com/Kevin-Kurka/LFG/backend/user-service/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Email         string `json:"email,omitempty"`
	WalletAddress string `json:"wallet_address,omitempty"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Register handles user registration
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	// Validate input
	if req.Email == "" {
		response.BadRequest(w, "email is required", nil)
		return
	}
	if req.Password == "" {
		response.BadRequest(w, "password is required", nil)
		return
	}
	if len(req.Password) < 8 {
		response.BadRequest(w, "password must be at least 8 characters", nil)
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		response.InternalServerError(w, "failed to hash password", err)
		return
	}

	// Create user
	userID := uuid.New()
	now := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start transaction
	tx, err := database.GetDB().Begin(ctx)
	if err != nil {
		response.InternalServerError(w, "failed to start transaction", err)
		return
	}
	defer tx.Rollback(ctx)

	// Insert user
	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, req.Email, hashedPassword, now, now)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		response.Conflict(w, "email already exists", err)
		return
	}

	// Create wallet for user
	walletID := uuid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO wallets (id, user_id, balance_credits, created_at)
		VALUES ($1, $2, 0, $3)
	`, walletID, userID, now)
	if err != nil {
		log.Printf("Failed to create wallet: %v", err)
		response.InternalServerError(w, "failed to create wallet", err)
		return
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		response.InternalServerError(w, "failed to commit transaction", err)
		return
	}

	user := &models.User{
		ID:        userID,
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	response.Created(w, user)
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		response.BadRequest(w, "email and password are required", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get user from database
	var user models.User
	err := database.GetDB().QueryRow(ctx, `
		SELECT id, email, password_hash, wallet_address, created_at, updated_at
		FROM users
		WHERE email = $1
	`, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.WalletAddress, &user.CreatedAt, &user.UpdatedAt)

	if err == pgx.ErrNoRows {
		response.Unauthorized(w, "invalid email or password", nil)
		return
	}
	if err != nil {
		log.Printf("Failed to query user: %v", err)
		response.InternalServerError(w, "failed to query user", err)
		return
	}

	// Verify password
	if err := auth.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		response.Unauthorized(w, "invalid email or password", nil)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		response.InternalServerError(w, "failed to generate token", err)
		return
	}

	loginResp := LoginResponse{
		Token: token,
		User:  &user,
	}

	response.Success(w, loginResp)
}

// GetProfile handles retrieving user profile (requires authentication)
func GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	// Get user from context
	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get user from database
	var user models.User
	err = database.GetDB().QueryRow(ctx, `
		SELECT id, email, wallet_address, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userCtx.UserID).Scan(&user.ID, &user.Email, &user.WalletAddress, &user.CreatedAt, &user.UpdatedAt)

	if err == pgx.ErrNoRows {
		response.NotFound(w, "user not found", nil)
		return
	}
	if err != nil {
		log.Printf("Failed to query user: %v", err)
		response.InternalServerError(w, "failed to query user", err)
		return
	}

	response.Success(w, user)
}

// UpdateProfile handles updating user profile (requires authentication)
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	// Get user from context
	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Build dynamic update query
	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.WalletAddress != "" {
		updates["wallet_address"] = req.WalletAddress
	}
	updates["updated_at"] = time.Now()

	if len(updates) == 1 { // Only updated_at
		response.BadRequest(w, "no fields to update", nil)
		return
	}

	// Update user
	_, err = database.GetDB().Exec(ctx, `
		UPDATE users
		SET email = COALESCE(NULLIF($1, ''), email),
		    wallet_address = COALESCE(NULLIF($2, ''), wallet_address),
		    updated_at = $3
		WHERE id = $4
	`, req.Email, req.WalletAddress, updates["updated_at"], userCtx.UserID)

	if err != nil {
		log.Printf("Failed to update user: %v", err)
		response.InternalServerError(w, "failed to update user", err)
		return
	}

	// Get updated user
	var user models.User
	err = database.GetDB().QueryRow(ctx, `
		SELECT id, email, wallet_address, created_at, updated_at
		FROM users
		WHERE id = $1
	`, userCtx.UserID).Scan(&user.ID, &user.Email, &user.WalletAddress, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.Printf("Failed to query updated user: %v", err)
		response.InternalServerError(w, "failed to query updated user", err)
		return
	}

	response.Success(w, user)
}

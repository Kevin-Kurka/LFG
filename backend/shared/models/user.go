package models

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusBanned    UserStatus = "BANNED"
)

// User represents the user model corresponding to the "users" table
type User struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Email         string     `json:"email" db:"email" validate:"required,email"`
	PasswordHash  string     `json:"-" db:"password_hash" validate:"required,min=60"`
	WalletAddress *string    `json:"wallet_address,omitempty" db:"wallet_address"`
	Status        UserStatus `json:"status" db:"status" validate:"required,oneof=ACTIVE SUSPENDED BANNED"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// UserRegistrationRequest represents the request payload for user registration
type UserRegistrationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// UserLoginRequest represents the request payload for user login
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserLoginResponse represents the response payload for user login
type UserLoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	User         *User     `json:"user"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// UserUpdateRequest represents the request payload for updating user profile
type UserUpdateRequest struct {
	WalletAddress *string `json:"wallet_address,omitempty" validate:"omitempty,len=42"`
}

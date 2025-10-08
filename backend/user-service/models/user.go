package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user model corresponding to the "users" table.
type User struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"` // Do not expose password hash in JSON responses
	WalletAddress string    `json:"wallet_address,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

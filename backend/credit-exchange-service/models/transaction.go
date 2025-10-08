package models

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a credit exchange transaction.
type Transaction struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Type          string    `json:"type"` // "PURCHASE" or "SALE"
	CryptoType    string    `json:"crypto_type"` // e.g., "BTC"
	CryptoAmount  float64   `json:"crypto_amount"`
	CreditAmount  float64   `json:"credit_amount"`
	Status        string    `json:"status"` // "PENDING", "COMPLETED", "FAILED"
	CreatedAt     time.Time `json:"created_at"`
}

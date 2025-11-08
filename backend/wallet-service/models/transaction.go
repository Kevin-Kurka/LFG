package models

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a wallet transaction
type Transaction struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Type        string    `json:"type"` // DEPOSIT, WITHDRAWAL, ORDER_LOCK, ORDER_UNLOCK, TRADE, CREDIT_BUY, CREDIT_SELL
	Amount      float64   `json:"amount"`
	Balance     float64   `json:"balance"` // Balance after transaction
	Reference   string    `json:"reference,omitempty"`
	ReferenceID *uuid.UUID `json:"reference_id,omitempty"` // Order ID, Trade ID, etc.
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

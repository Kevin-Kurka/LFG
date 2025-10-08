package models

import (
	"time"

	"github.com/google/uuid"
)

// Wallet represents the wallet model corresponding to the "wallets" table.
type Wallet struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	BalanceCredits float64   `json:"balance_credits"`
	CreatedAt      time.Time `json:"created_at"`
}

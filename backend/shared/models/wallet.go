package models

import (
	"time"

	"github.com/google/uuid"
)

// Wallet represents the wallet model corresponding to the "wallets" table
type Wallet struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id" validate:"required"`
	BalanceCredits float64   `json:"balance_credits" db:"balance_credits" validate:"min=0"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// WalletBalanceResponse represents the wallet balance response
type WalletBalanceResponse struct {
	Balance          float64 `json:"balance"`
	AvailableBalance float64 `json:"available_balance"`
	LockedBalance    float64 `json:"locked_balance"`
}

// TransactionType represents the type of a wallet transaction
type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "DEPOSIT"
	TransactionTypeWithdrawal TransactionType = "WITHDRAWAL"
	TransactionTypeTrade      TransactionType = "TRADE"
	TransactionTypeRefund     TransactionType = "REFUND"
)

// WalletTransaction represents a wallet balance change
type WalletTransaction struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	WalletID        uuid.UUID       `json:"wallet_id" db:"wallet_id"`
	Type            TransactionType `json:"type" db:"type"`
	Amount          float64         `json:"amount" db:"amount"`
	BalanceBefore   float64         `json:"balance_before" db:"balance_before"`
	BalanceAfter    float64         `json:"balance_after" db:"balance_after"`
	ReferenceID     *uuid.UUID      `json:"reference_id,omitempty" db:"reference_id"`
	ReferenceType   *string         `json:"reference_type,omitempty" db:"reference_type"`
	Description     string          `json:"description" db:"description"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
}

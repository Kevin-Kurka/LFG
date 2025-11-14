package models

import (
	"time"

	"github.com/google/uuid"
)

// CreditTransactionType represents the type of credit transaction
type CreditTransactionType string

const (
	CreditTransactionTypePurchase CreditTransactionType = "PURCHASE"
	CreditTransactionTypeSale     CreditTransactionType = "SALE"
)

// CreditTransactionStatus represents the status of a credit transaction
type CreditTransactionStatus string

const (
	CreditTransactionStatusPending   CreditTransactionStatus = "PENDING"
	CreditTransactionStatusCompleted CreditTransactionStatus = "COMPLETED"
	CreditTransactionStatusFailed    CreditTransactionStatus = "FAILED"
)

// CreditTransaction represents a credit exchange transaction
type CreditTransaction struct {
	ID           uuid.UUID               `json:"id" db:"id"`
	UserID       uuid.UUID               `json:"user_id" db:"user_id" validate:"required"`
	Type         CreditTransactionType   `json:"type" db:"type" validate:"required"`
	CryptoType   string                  `json:"crypto_type" db:"crypto_type" validate:"required"`
	CryptoAmount float64                 `json:"crypto_amount" db:"crypto_amount" validate:"required,gt=0"`
	CreditAmount float64                 `json:"credit_amount" db:"credit_amount" validate:"required,gt=0"`
	Status       CreditTransactionStatus `json:"status" db:"status" validate:"required"`
	CreatedAt    time.Time               `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at" db:"updated_at"`
}

// BuyCreditsRequest represents the request to buy credits
type BuyCreditsRequest struct {
	CryptoType   string  `json:"crypto_type" validate:"required,oneof=BTC ETH USDC"`
	CryptoAmount float64 `json:"crypto_amount" validate:"required,gt=0"`
}

// SellCreditsRequest represents the request to sell credits
type SellCreditsRequest struct {
	CryptoType   string  `json:"crypto_type" validate:"required,oneof=BTC ETH USDC"`
	CreditAmount float64 `json:"credit_amount" validate:"required,gt=0"`
}

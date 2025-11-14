package models

import (
	"time"

	"github.com/google/uuid"
)

// Trade represents a trade execution
type Trade struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ContractID    uuid.UUID `json:"contract_id" db:"contract_id" validate:"required"`
	MakerOrderID  uuid.UUID `json:"maker_order_id" db:"maker_order_id" validate:"required"`
	TakerOrderID  uuid.UUID `json:"taker_order_id" db:"taker_order_id" validate:"required"`
	Quantity      int       `json:"quantity" db:"quantity" validate:"required,min=1"`
	PriceCredits  float64   `json:"price_credits" db:"price_credits" validate:"required,gt=0,lte=1"`
	ExecutedAt    time.Time `json:"executed_at" db:"executed_at"`
}

// TradeExecutedEvent represents an event published when a trade is executed
type TradeExecutedEvent struct {
	TradeID       uuid.UUID `json:"trade_id"`
	ContractID    uuid.UUID `json:"contract_id"`
	MakerOrderID  uuid.UUID `json:"maker_order_id"`
	TakerOrderID  uuid.UUID `json:"taker_order_id"`
	MakerUserID   uuid.UUID `json:"maker_user_id"`
	TakerUserID   uuid.UUID `json:"taker_user_id"`
	Quantity      int       `json:"quantity"`
	PriceCredits  float64   `json:"price_credits"`
	ExecutedAt    time.Time `json:"executed_at"`
}

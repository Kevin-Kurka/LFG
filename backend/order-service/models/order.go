package models

import (
	"time"

	"github.com/google/uuid"
)

// Order represents the order model corresponding to the "orders" table.
// This model is used within the order-service to manage order state.
type Order struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	ContractID        uuid.UUID `json:"contract_id"`
	Type              string    `json:"type"`              // MARKET, LIMIT, STOP, STOP_LIMIT
	Status            string    `json:"status"`            // PENDING, ACTIVE, FILLED, CANCELLED
	Quantity          int       `json:"quantity"`
	QuantityFilled    int       `json:"quantity_filled"`
	LimitPriceCredits float64   `json:"limit_price_credits,omitempty"`
	StopPriceCredits  float64   `json:"stop_price_credits,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

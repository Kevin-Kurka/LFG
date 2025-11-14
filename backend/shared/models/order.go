package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderType represents the type of an order
type OrderType string

const (
	OrderTypeMarket    OrderType = "MARKET"
	OrderTypeLimit     OrderType = "LIMIT"
	OrderTypeStop      OrderType = "STOP"
	OrderTypeStopLimit OrderType = "STOP_LIMIT"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending         OrderStatus = "PENDING"
	OrderStatusActive          OrderStatus = "ACTIVE"
	OrderStatusFilled          OrderStatus = "FILLED"
	OrderStatusPartiallyFilled OrderStatus = "PARTIALLY_FILLED"
	OrderStatusCancelled       OrderStatus = "CANCELLED"
	OrderStatusRejected        OrderStatus = "REJECTED"
)

// Order represents the order model corresponding to the "orders" table
type Order struct {
	ID                uuid.UUID   `json:"id" db:"id"`
	UserID            uuid.UUID   `json:"user_id" db:"user_id" validate:"required"`
	ContractID        uuid.UUID   `json:"contract_id" db:"contract_id" validate:"required"`
	Type              OrderType   `json:"type" db:"type" validate:"required"`
	Status            OrderStatus `json:"status" db:"status" validate:"required"`
	Quantity          int         `json:"quantity" db:"quantity" validate:"required,min=1"`
	QuantityFilled    int         `json:"quantity_filled" db:"quantity_filled" validate:"min=0"`
	LimitPriceCredits *float64    `json:"limit_price_credits,omitempty" db:"limit_price_credits" validate:"omitempty,gt=0,lte=1"`
	StopPriceCredits  *float64    `json:"stop_price_credits,omitempty" db:"stop_price_credits" validate:"omitempty,gt=0,lte=1"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`
}

// OrderPlaceRequest represents the request to place a new order
type OrderPlaceRequest struct {
	ContractID        uuid.UUID `json:"contract_id" validate:"required"`
	Type              OrderType `json:"type" validate:"required,oneof=MARKET LIMIT STOP STOP_LIMIT"`
	Quantity          int       `json:"quantity" validate:"required,min=1,max=10000"`
	LimitPriceCredits *float64  `json:"limit_price_credits,omitempty" validate:"omitempty,gt=0,lte=1"`
	StopPriceCredits  *float64  `json:"stop_price_credits,omitempty" validate:"omitempty,gt=0,lte=1"`
}

// OrderPlaceResponse represents the response after placing an order
type OrderPlaceResponse struct {
	OrderID        uuid.UUID   `json:"order_id"`
	Status         OrderStatus `json:"status"`
	QuantityFilled int         `json:"quantity_filled"`
	AveragePrice   float64     `json:"average_price,omitempty"`
}

// OrderCancelRequest represents the request to cancel an order
type OrderCancelRequest struct {
	OrderID uuid.UUID `json:"order_id" validate:"required"`
}

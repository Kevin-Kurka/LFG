package models

import (
	"time"

	"github.com/google/uuid"
)

// MarketStatus represents the status of a market
type MarketStatus string

const (
	MarketStatusUpcoming  MarketStatus = "UPCOMING"
	MarketStatusOpen      MarketStatus = "OPEN"
	MarketStatusClosed    MarketStatus = "CLOSED"
	MarketStatusResolved  MarketStatus = "RESOLVED"
	MarketStatusCancelled MarketStatus = "CANCELLED"
)

// MarketOutcome represents the outcome of a resolved market
type MarketOutcome string

const (
	MarketOutcomeYes       MarketOutcome = "YES"
	MarketOutcomeNo        MarketOutcome = "NO"
	MarketOutcomeCancelled MarketOutcome = "CANCELLED"
)

// Market represents the market model corresponding to the "markets" table
type Market struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	Ticker           string         `json:"ticker" db:"ticker" validate:"required,uppercase,max=50"`
	Question         string         `json:"question" db:"question" validate:"required,min=10,max=500"`
	Rules            string         `json:"rules" db:"rules" validate:"required"`
	ResolutionSource string         `json:"resolution_source" db:"resolution_source" validate:"required,max=255"`
	Status           MarketStatus   `json:"status" db:"status" validate:"required"`
	ExpiresAt        time.Time      `json:"expires_at" db:"expires_at" validate:"required"`
	ResolvedAt       *time.Time     `json:"resolved_at,omitempty" db:"resolved_at"`
	Outcome          *MarketOutcome `json:"outcome,omitempty" db:"outcome"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

// ContractSide represents the side of a contract (YES or NO)
type ContractSide string

const (
	ContractSideYes ContractSide = "YES"
	ContractSideNo  ContractSide = "NO"
)

// Contract represents a contract (e.g., a "Yes" or "No" share) for a given market
type Contract struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	MarketID  uuid.UUID    `json:"market_id" db:"market_id" validate:"required"`
	Side      ContractSide `json:"side" db:"side" validate:"required,oneof=YES NO"`
	Ticker    string       `json:"ticker" db:"ticker" validate:"required,max=60"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
}

// MarketCreateRequest represents the request to create a new market
type MarketCreateRequest struct {
	Ticker           string    `json:"ticker" validate:"required,uppercase,max=50"`
	Question         string    `json:"question" validate:"required,min=10,max=500"`
	Rules            string    `json:"rules" validate:"required"`
	ResolutionSource string    `json:"resolution_source" validate:"required,max=255"`
	ExpiresAt        time.Time `json:"expires_at" validate:"required"`
}

// MarketResolveRequest represents the request to resolve a market
type MarketResolveRequest struct {
	Outcome MarketOutcome `json:"outcome" validate:"required,oneof=YES NO CANCELLED"`
}

// MarketListResponse represents the response for listing markets
type MarketListResponse struct {
	Markets    []*Market `json:"markets"`
	TotalCount int       `json:"total_count"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
}

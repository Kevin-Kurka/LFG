package models

import (
	"time"

	"github.com/google/uuid"
)

// Market represents the market model corresponding to the "markets" table.
type Market struct {
	ID                uuid.UUID `json:"id"`
	Ticker            string    `json:"ticker"`
	Question          string    `json:"question"`
	Rules             string    `json:"rules"`
	ResolutionSource  string    `json:"resolution_source"`
	Status            string    `json:"status"` // UPCOMING, OPEN, CLOSED, RESOLVED
	ExpiresAt         time.Time `json:"expires_at"`
	ResolvedAt        time.Time `json:"resolved_at,omitempty"`
	Outcome           string    `json:"outcome,omitempty"` // YES or NO
}

// Contract represents a contract (e.g., a "Yes" or "No" share) for a given market.
// This is defined here as it's closely related to the Market model.
type Contract struct {
	ID       uuid.UUID `json:"id"`
	MarketID uuid.UUID `json:"market_id"`
	Side     string    `json:"side"` // YES or NO
	Ticker   string    `json:"ticker"`
}

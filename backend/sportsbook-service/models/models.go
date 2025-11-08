package models

import (
	"time"

	"github.com/google/uuid"
)

type SportsbookAccount struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             uuid.UUID  `json:"user_id"`
	SportsbookID       uuid.UUID  `json:"sportsbook_id"`
	SportsbookName     string     `json:"sportsbook_name"`
	UsernameEncrypted  string     `json:"-"` // Hidden
	CredentialsEncrypted string   `json:"-"` // Hidden
	IsConnected        bool       `json:"is_connected"`
	LastSyncAt         *time.Time `json:"last_sync_at,omitempty"`
	Balance            *float64   `json:"balance,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type SportsEvent struct {
	ID         uuid.UUID  `json:"id"`
	SportID    uuid.UUID  `json:"sport_id"`
	SportName  string     `json:"sport_name"`
	HomeTeam   string     `json:"home_team"`
	AwayTeam   string     `json:"away_team"`
	EventTime  time.Time  `json:"event_time"`
	League     string     `json:"league"`
	Status     string     `json:"status"`
	HomeScore  *int       `json:"home_score,omitempty"`
	AwayScore  *int       `json:"away_score,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Odds       []OddsData `json:"odds,omitempty"`
}

type OddsData struct {
	ID             uuid.UUID `json:"id"`
	EventID        uuid.UUID `json:"event_id"`
	SportsbookID   uuid.UUID `json:"sportsbook_id"`
	SportsbookName string    `json:"sportsbook_name"`
	MarketType     string    `json:"market_type"` // MONEYLINE, SPREAD, TOTAL
	BetType        string    `json:"bet_type"`
	Selection      string    `json:"selection"`
	OddsAmerican   *int      `json:"odds_american,omitempty"`
	OddsDecimal    *float64  `json:"odds_decimal,omitempty"`
	Line           *float64  `json:"line,omitempty"`
	IsActive       bool      `json:"is_active"`
	LastUpdated    time.Time `json:"last_updated"`
}

type ArbitrageOpportunity struct {
	ID                  uuid.UUID              `json:"id"`
	EventID             uuid.UUID              `json:"event_id"`
	Event               *SportsEvent           `json:"event,omitempty"`
	MarketType          string                 `json:"market_type"`
	BetDescription      string                 `json:"bet_description"`
	ProfitPercentage    float64                `json:"profit_percentage"`
	TotalStakeRequired  float64                `json:"total_stake_required"`
	EstimatedProfit     float64                `json:"estimated_profit"`
	OddsData            map[string]interface{} `json:"odds_data"`
	IsActive            bool                   `json:"is_active"`
	DetectedAt          time.Time              `json:"detected_at"`
	ExpiresAt           time.Time              `json:"expires_at"`
}

type HedgeOpportunity struct {
	ID                 uuid.UUID              `json:"id"`
	UserID             uuid.UUID              `json:"user_id"`
	OriginalBetID      *uuid.UUID             `json:"original_bet_id,omitempty"`
	EventID            uuid.UUID              `json:"event_id"`
	Event              *SportsEvent           `json:"event,omitempty"`
	HedgeType          string                 `json:"hedge_type"`
	OriginalStake      float64                `json:"original_stake"`
	HedgeStakeRequired float64                `json:"hedge_stake_required"`
	GuaranteedProfit   float64                `json:"guaranteed_profit"`
	HedgeOddsData      map[string]interface{} `json:"hedge_odds_data"`
	IsActive           bool                   `json:"is_active"`
	DetectedAt         time.Time              `json:"detected_at"`
}

type UserBet struct {
	ID                   uuid.UUID  `json:"id"`
	UserID               uuid.UUID  `json:"user_id"`
	SportsbookAccountID  uuid.UUID  `json:"sportsbook_account_id"`
	SportsbookName       string     `json:"sportsbook_name,omitempty"`
	EventID              *uuid.UUID `json:"event_id,omitempty"`
	Event                *SportsEvent `json:"event,omitempty"`
	BetType              string     `json:"bet_type"`
	Selection            string     `json:"selection"`
	Stake                float64    `json:"stake"`
	OddsAmerican         *int       `json:"odds_american,omitempty"`
	OddsDecimal          *float64   `json:"odds_decimal,omitempty"`
	PotentialPayout      float64    `json:"potential_payout"`
	Status               string     `json:"status"` // PENDING, WON, LOST, CANCELLED
	Result               *string    `json:"result,omitempty"`
	ProfitLoss           *float64   `json:"profit_loss,omitempty"`
	PlacedAt             time.Time  `json:"placed_at"`
	SettledAt            *time.Time `json:"settled_at,omitempty"`
	Notes                *string    `json:"notes,omitempty"`
}

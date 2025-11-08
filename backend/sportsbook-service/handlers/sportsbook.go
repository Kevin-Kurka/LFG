package handlers

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/crypto"
	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/common/response"
	"github.com/Kevin-Kurka/LFG/backend/sportsbook-service/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type LinkAccountRequest struct {
	SportsbookID string `json:"sportsbook_id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

type TrackBetRequest struct {
	SportsbookAccountID string   `json:"sportsbook_account_id"`
	EventID             *string  `json:"event_id,omitempty"`
	BetType             string   `json:"bet_type"`
	Selection           string   `json:"selection"`
	Stake               float64  `json:"stake"`
	OddsAmerican        *int     `json:"odds_american,omitempty"`
	OddsDecimal         *float64 `json:"odds_decimal,omitempty"`
	PotentialPayout     float64  `json:"potential_payout"`
	Notes               *string  `json:"notes,omitempty"`
}

// LinkAccount links a sportsbook account
func LinkAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	var req LinkAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	sportsbookID, err := uuid.Parse(req.SportsbookID)
	if err != nil {
		response.BadRequest(w, "invalid sportsbook_id", err)
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.Unauthorized(w, "invalid user", err)
		return
	}

	// Encrypt credentials
	usernameEnc, err := crypto.Encrypt(req.Username)
	if err != nil {
		response.InternalServerError(w, "failed to encrypt username", err)
		return
	}

	passwordEnc, err := crypto.Encrypt(req.Password)
	if err != nil {
		response.InternalServerError(w, "failed to encrypt password", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accountID := uuid.New()
	now := time.Now()

	_, err = database.GetDB().Exec(ctx, `
		INSERT INTO user_sportsbook_accounts (id, user_id, sportsbook_id, username_encrypted, credentials_encrypted, is_connected, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, true, $6, $7)
	`, accountID, userID, sportsbookID, usernameEnc, passwordEnc, now, now)

	if err != nil {
		log.Printf("Failed to link account: %v", err)
		response.Conflict(w, "account already linked or database error", err)
		return
	}

	response.Created(w, map[string]string{
		"id":      accountID.String(),
		"message": "sportsbook account linked successfully",
	})
}

// GetLinkedAccounts returns user's linked sportsbook accounts
func GetLinkedAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.Unauthorized(w, "invalid user", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := database.GetDB().Query(ctx, `
		SELECT usa.id, usa.user_id, usa.sportsbook_id, sp.display_name, usa.is_connected, usa.last_sync_at, usa.balance, usa.created_at, usa.updated_at
		FROM user_sportsbook_accounts usa
		JOIN sportsbook_providers sp ON usa.sportsbook_id = sp.id
		WHERE usa.user_id = $1
	`, userID)
	if err != nil {
		log.Printf("Failed to query accounts: %v", err)
		response.InternalServerError(w, "failed to query accounts", err)
		return
	}
	defer rows.Close()

	accounts := make([]models.SportsbookAccount, 0)
	for rows.Next() {
		var account models.SportsbookAccount
		err := rows.Scan(&account.ID, &account.UserID, &account.SportsbookID, &account.SportsbookName, &account.IsConnected, &account.LastSyncAt, &account.Balance, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			continue
		}
		accounts = append(accounts, account)
	}

	response.Success(w, map[string]interface{}{
		"accounts": accounts,
		"count":    len(accounts),
	})
}

// DeleteLinkedAccount removes a linked account
func DeleteLinkedAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	vars := mux.Vars(r)
	accountID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid account_id", err)
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.Unauthorized(w, "invalid user", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := database.GetDB().Exec(ctx, `
		DELETE FROM user_sportsbook_accounts
		WHERE id = $1 AND user_id = $2
	`, accountID, userID)

	if err != nil {
		log.Printf("Failed to delete account: %v", err)
		response.InternalServerError(w, "failed to delete account", err)
		return
	}

	if result.RowsAffected() == 0 {
		response.NotFound(w, "account not found", nil)
		return
	}

	response.Success(w, map[string]string{"message": "account deleted successfully"})
}

// GetSportsEvents returns sports events with odds from all books
func GetSportsEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	sport := r.URL.Query().Get("sport")
	league := r.URL.Query().Get("league")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		SELECT se.id, se.sport_id, s.display_name, se.home_team, se.away_team, se.event_time, se.league, se.status, se.created_at, se.updated_at
		FROM sports_events se
		JOIN sports s ON se.sport_id = s.id
		WHERE se.status = 'UPCOMING'
	`
	args := []interface{}{}
	argIdx := 1

	if sport != "" {
		query += " AND s.name = $1"
		args = append(args, sport)
		argIdx++
	}

	if league != "" {
		query += " AND se.league = $" + string(rune(argIdx+'0'))
		args = append(args, league)
	}

	query += " ORDER BY se.event_time ASC LIMIT 50"

	rows, err := database.GetDB().Query(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to query events: %v", err)
		response.InternalServerError(w, "failed to query events", err)
		return
	}
	defer rows.Close()

	events := make([]models.SportsEvent, 0)
	for rows.Next() {
		var event models.SportsEvent
		err := rows.Scan(&event.ID, &event.SportID, &event.SportName, &event.HomeTeam, &event.AwayTeam, &event.EventTime, &event.League, &event.Status, &event.CreatedAt, &event.UpdatedAt)
		if err != nil {
			continue
		}

		// Get odds for this event
		oddsRows, err := database.GetDB().Query(ctx, `
			SELECT o.id, o.event_id, o.sportsbook_id, sp.display_name, o.market_type, o.bet_type, o.selection, o.odds_american, o.odds_decimal, o.line, o.is_active, o.last_updated
			FROM odds o
			JOIN sportsbook_providers sp ON o.sportsbook_id = sp.id
			WHERE o.event_id = $1 AND o.is_active = true
			ORDER BY o.market_type, o.bet_type
		`, event.ID)
		if err == nil {
			defer oddsRows.Close()
			odds := make([]models.OddsData, 0)
			for oddsRows.Next() {
				var odd models.OddsData
				err := oddsRows.Scan(&odd.ID, &odd.EventID, &odd.SportsbookID, &odd.SportsbookName, &odd.MarketType, &odd.BetType, &odd.Selection, &odd.OddsAmerican, &odd.OddsDecimal, &odd.Line, &odd.IsActive, &odd.LastUpdated)
				if err == nil {
					odds = append(odds, odd)
				}
			}
			event.Odds = odds
		}

		events = append(events, event)
	}

	response.Success(w, map[string]interface{}{
		"events": events,
		"count":  len(events),
	})
}

// GetEventDetails returns detailed information for a single event
func GetEventDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	vars := mux.Vars(r)
	eventID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid event_id", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var event models.SportsEvent
	err = database.GetDB().QueryRow(ctx, `
		SELECT se.id, se.sport_id, s.display_name, se.home_team, se.away_team, se.event_time, se.league, se.status, se.home_score, se.away_score, se.created_at, se.updated_at
		FROM sports_events se
		JOIN sports s ON se.sport_id = s.id
		WHERE se.id = $1
	`, eventID).Scan(&event.ID, &event.SportID, &event.SportName, &event.HomeTeam, &event.AwayTeam, &event.EventTime, &event.League, &event.Status, &event.HomeScore, &event.AwayScore, &event.CreatedAt, &event.UpdatedAt)

	if err == pgx.ErrNoRows {
		response.NotFound(w, "event not found", nil)
		return
	}
	if err != nil {
		log.Printf("Failed to query event: %v", err)
		response.InternalServerError(w, "failed to query event", err)
		return
	}

	// Get all odds for this event grouped by market type
	oddsRows, err := database.GetDB().Query(ctx, `
		SELECT o.id, o.event_id, o.sportsbook_id, sp.display_name, o.market_type, o.bet_type, o.selection, o.odds_american, o.odds_decimal, o.line, o.is_active, o.last_updated
		FROM odds o
		JOIN sportsbook_providers sp ON o.sportsbook_id = sp.id
		WHERE o.event_id = $1 AND o.is_active = true
		ORDER BY o.market_type, o.sportsbook_id
	`, eventID)
	if err != nil {
		log.Printf("Failed to query odds: %v", err)
	} else {
		defer oddsRows.Close()
		odds := make([]models.OddsData, 0)
		for oddsRows.Next() {
			var odd models.OddsData
			err := oddsRows.Scan(&odd.ID, &odd.EventID, &odd.SportsbookID, &odd.SportsbookName, &odd.MarketType, &odd.BetType, &odd.Selection, &odd.OddsAmerican, &odd.OddsDecimal, &odd.Line, &odd.IsActive, &odd.LastUpdated)
			if err == nil {
				odds = append(odds, odd)
			}
		}
		event.Odds = odds
	}

	response.Success(w, event)
}

// GetArbitrageOpportunities finds arbitrage opportunities
func GetArbitrageOpportunities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all active events
	rows, err := database.GetDB().Query(ctx, `
		SELECT id FROM sports_events
		WHERE status = 'UPCOMING' AND event_time > NOW()
		LIMIT 100
	`)
	if err != nil {
		response.InternalServerError(w, "failed to query events", err)
		return
	}
	defer rows.Close()

	arbitrages := make([]models.ArbitrageOpportunity, 0)

	for rows.Next() {
		var eventID uuid.UUID
		rows.Scan(&eventID)

		// Get moneyline odds for this event
		oddsRows, err := database.GetDB().Query(ctx, `
			SELECT selection, sportsbook_id, odds_decimal
			FROM odds
			WHERE event_id = $1 AND market_type = 'MONEYLINE' AND is_active = true AND odds_decimal IS NOT NULL
		`, eventID)
		if err != nil {
			continue
		}

		type OddsEntry struct {
			Selection    string
			SportsbookID uuid.UUID
			OddsDecimal  float64
		}

		oddsMap := make(map[string][]OddsEntry)
		for oddsRows.Next() {
			var entry OddsEntry
			oddsRows.Scan(&entry.Selection, &entry.SportsbookID, &entry.OddsDecimal)
			oddsMap[entry.Selection] = append(oddsMap[entry.Selection], entry)
		}
		oddsRows.Close()

		// Check for arbitrage on 2-way markets
		if len(oddsMap) >= 2 {
			bestOdds := make(map[string]float64)
			for selection, entries := range oddsMap {
				best := 0.0
				for _, entry := range entries {
					if entry.OddsDecimal > best {
						best = entry.OddsDecimal
					}
				}
				bestOdds[selection] = best
			}

			// Calculate implied probabilities
			totalImplied := 0.0
			for _, odds := range bestOdds {
				totalImplied += 1.0 / odds
			}

			// If total implied probability < 1.0, there's an arbitrage opportunity
			if totalImplied < 1.0 {
				profitPercentage := (1.0/totalImplied - 1.0) * 100
				totalStake := 1000.0 // Example stake
				profit := totalStake * (1.0/totalImplied - 1.0)

				arb := models.ArbitrageOpportunity{
					ID:                 uuid.New(),
					EventID:            eventID,
					MarketType:         "MONEYLINE",
					BetDescription:     "Moneyline arbitrage",
					ProfitPercentage:   profitPercentage,
					TotalStakeRequired: totalStake,
					EstimatedProfit:    profit,
					OddsData:           map[string]interface{}{"best_odds": bestOdds},
					IsActive:           true,
					DetectedAt:         time.Now(),
					ExpiresAt:          time.Now().Add(1 * time.Hour),
				}
				arbitrages = append(arbitrages, arb)
			}
		}
	}

	response.Success(w, map[string]interface{}{
		"arbitrages": arbitrages,
		"count":      len(arbitrages),
	})
}

// GetHedgeOpportunities finds hedge opportunities for a user
func GetHedgeOpportunities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.Unauthorized(w, "invalid user", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user's active bets
	rows, err := database.GetDB().Query(ctx, `
		SELECT id, event_id, bet_type, selection, stake, odds_decimal, potential_payout
		FROM user_bets
		WHERE user_id = $1 AND status = 'PENDING' AND event_id IS NOT NULL
	`, userID)
	if err != nil {
		response.InternalServerError(w, "failed to query bets", err)
		return
	}
	defer rows.Close()

	hedges := make([]models.HedgeOpportunity, 0)

	for rows.Next() {
		var betID, eventID uuid.UUID
		var betType, selection string
		var stake, oddsDecimal, potentialPayout float64
		rows.Scan(&betID, &eventID, &betType, &selection, &stake, &oddsDecimal, &potentialPayout)

		// Find opposite side odds
		var oppositeOdds float64
		var oppSelection string

		if selection == "home" || selection == "HOME" {
			oppSelection = "AWAY"
		} else {
			oppSelection = "HOME"
		}

		err := database.GetDB().QueryRow(ctx, `
			SELECT MAX(odds_decimal)
			FROM odds
			WHERE event_id = $1 AND market_type = 'MONEYLINE' AND selection ILIKE $2 AND is_active = true AND odds_decimal IS NOT NULL
		`, eventID, "%"+oppSelection+"%").Scan(&oppositeOdds)

		if err == nil && oppositeOdds > 0 {
			// Calculate hedge stake to guarantee profit
			hedgeStake := potentialPayout / oppositeOdds
			profit := potentialPayout - stake - hedgeStake

			if profit > 0 {
				hedge := models.HedgeOpportunity{
					ID:                 uuid.New(),
					UserID:             userID,
					OriginalBetID:      &betID,
					EventID:            eventID,
					HedgeType:          "FULL_HEDGE",
					OriginalStake:      stake,
					HedgeStakeRequired: hedgeStake,
					GuaranteedProfit:   profit,
					HedgeOddsData:      map[string]interface{}{"hedge_odds": oppositeOdds, "selection": oppSelection},
					IsActive:           true,
					DetectedAt:         time.Now(),
				}
				hedges = append(hedges, hedge)
			}
		}
	}

	response.Success(w, map[string]interface{}{
		"hedges": hedges,
		"count":  len(hedges),
	})
}

// TrackBet tracks a bet placed externally
func TrackBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	var req TrackBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	accountID, err := uuid.Parse(req.SportsbookAccountID)
	if err != nil {
		response.BadRequest(w, "invalid sportsbook_account_id", err)
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.Unauthorized(w, "invalid user", err)
		return
	}

	var eventID *uuid.UUID
	if req.EventID != nil {
		eid, err := uuid.Parse(*req.EventID)
		if err == nil {
			eventID = &eid
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	betID := uuid.New()
	now := time.Now()

	_, err = database.GetDB().Exec(ctx, `
		INSERT INTO user_bets (id, user_id, sportsbook_account_id, event_id, bet_type, selection, stake, odds_american, odds_decimal, potential_payout, status, placed_at, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'PENDING', $11, $12)
	`, betID, userID, accountID, eventID, req.BetType, req.Selection, req.Stake, req.OddsAmerican, req.OddsDecimal, req.PotentialPayout, now, req.Notes)

	if err != nil {
		log.Printf("Failed to track bet: %v", err)
		response.InternalServerError(w, "failed to track bet", err)
		return
	}

	response.Created(w, map[string]string{
		"id":      betID.String(),
		"message": "bet tracked successfully",
	})
}

// GetUserBets returns user's bet history
func GetUserBets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.Unauthorized(w, "invalid user", err)
		return
	}

	status := r.URL.Query().Get("status")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
		SELECT ub.id, ub.user_id, ub.sportsbook_account_id, sp.display_name, ub.event_id, ub.bet_type, ub.selection, ub.stake, ub.odds_american, ub.odds_decimal, ub.potential_payout, ub.status, ub.result, ub.profit_loss, ub.placed_at, ub.settled_at, ub.notes
		FROM user_bets ub
		JOIN user_sportsbook_accounts usa ON ub.sportsbook_account_id = usa.id
		JOIN sportsbook_providers sp ON usa.sportsbook_id = sp.id
		WHERE ub.user_id = $1
	`
	args := []interface{}{userID}

	if status != "" {
		query += " AND ub.status = $2"
		args = append(args, status)
	}

	query += " ORDER BY ub.placed_at DESC LIMIT 100"

	rows, err := database.GetDB().Query(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to query bets: %v", err)
		response.InternalServerError(w, "failed to query bets", err)
		return
	}
	defer rows.Close()

	bets := make([]models.UserBet, 0)
	for rows.Next() {
		var bet models.UserBet
		err := rows.Scan(&bet.ID, &bet.UserID, &bet.SportsbookAccountID, &bet.SportsbookName, &bet.EventID, &bet.BetType, &bet.Selection, &bet.Stake, &bet.OddsAmerican, &bet.OddsDecimal, &bet.PotentialPayout, &bet.Status, &bet.Result, &bet.ProfitLoss, &bet.PlacedAt, &bet.SettledAt, &bet.Notes)
		if err != nil {
			continue
		}
		bets = append(bets, bet)
	}

	response.Success(w, map[string]interface{}{
		"bets":  bets,
		"count": len(bets),
	})
}

// Helper function to convert American odds to decimal
func americanToDecimal(american int) float64 {
	if american > 0 {
		return float64(american)/100.0 + 1.0
	}
	return 100.0/math.Abs(float64(american)) + 1.0
}

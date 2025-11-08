package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/common/response"
	"github.com/Kevin-Kurka/LFG/backend/market-service/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type CreateMarketRequest struct {
	Ticker           string    `json:"ticker"`
	Question         string    `json:"question"`
	Rules            string    `json:"rules"`
	ResolutionSource string    `json:"resolution_source"`
	ExpiresAt        time.Time `json:"expires_at"`
}

type UpdateMarketRequest struct {
	Question         string    `json:"question,omitempty"`
	Rules            string    `json:"rules,omitempty"`
	ResolutionSource string    `json:"resolution_source,omitempty"`
	Status           string    `json:"status,omitempty"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
}

type ResolveMarketRequest struct {
	Outcome string `json:"outcome"` // YES or NO
}

type OrderBookEntry struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type OrderBook struct {
	ContractID uuid.UUID        `json:"contract_id"`
	Side       string           `json:"side"`
	Bids       []OrderBookEntry `json:"bids"`
	Asks       []OrderBookEntry `json:"asks"`
}

// ListMarkets returns all markets with optional filtering
func ListMarkets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	// Parse query parameters
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build query
	query := `SELECT id, ticker, question, rules, resolution_source, status, expires_at, resolved_at, outcome FROM markets WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if status != "" {
		query += " AND status = $" + string(rune(argIdx+'0'))
		args = append(args, status)
		argIdx++
	}

	if search != "" {
		query += " AND (question ILIKE $" + string(rune(argIdx+'0')) + " OR ticker ILIKE $" + string(rune(argIdx+'0')) + ")"
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += " ORDER BY created_at DESC"

	rows, err := database.GetDB().Query(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to query markets: %v", err)
		response.InternalServerError(w, "failed to query markets", err)
		return
	}
	defer rows.Close()

	markets := make([]models.Market, 0)
	for rows.Next() {
		var market models.Market
		var resolvedAt *time.Time
		var outcome *string

		err := rows.Scan(
			&market.ID,
			&market.Ticker,
			&market.Question,
			&market.Rules,
			&market.ResolutionSource,
			&market.Status,
			&market.ExpiresAt,
			&resolvedAt,
			&outcome,
		)
		if err != nil {
			log.Printf("Failed to scan market: %v", err)
			continue
		}

		if resolvedAt != nil {
			market.ResolvedAt = *resolvedAt
		}
		if outcome != nil {
			market.Outcome = *outcome
		}

		markets = append(markets, market)
	}

	response.Success(w, map[string]interface{}{
		"markets": markets,
		"count":   len(markets),
	})
}

// GetMarket returns a specific market by ID
func GetMarket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	vars := mux.Vars(r)
	marketID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid market ID", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var market models.Market
	var resolvedAt *time.Time
	var outcome *string

	err = database.GetDB().QueryRow(ctx, `
		SELECT id, ticker, question, rules, resolution_source, status, expires_at, resolved_at, outcome
		FROM markets
		WHERE id = $1
	`, marketID).Scan(
		&market.ID,
		&market.Ticker,
		&market.Question,
		&market.Rules,
		&market.ResolutionSource,
		&market.Status,
		&market.ExpiresAt,
		&resolvedAt,
		&outcome,
	)

	if err == pgx.ErrNoRows {
		response.NotFound(w, "market not found", nil)
		return
	}
	if err != nil {
		log.Printf("Failed to query market: %v", err)
		response.InternalServerError(w, "failed to query market", err)
		return
	}

	if resolvedAt != nil {
		market.ResolvedAt = *resolvedAt
	}
	if outcome != nil {
		market.Outcome = *outcome
	}

	// Get contracts for this market
	rows, err := database.GetDB().Query(ctx, `
		SELECT id, market_id, side, ticker
		FROM contracts
		WHERE market_id = $1
	`, marketID)
	if err != nil {
		log.Printf("Failed to query contracts: %v", err)
	} else {
		defer rows.Close()
		contracts := make([]models.Contract, 0)
		for rows.Next() {
			var contract models.Contract
			if err := rows.Scan(&contract.ID, &contract.MarketID, &contract.Side, &contract.Ticker); err == nil {
				contracts = append(contracts, contract)
			}
		}

		responseData := map[string]interface{}{
			"market":    market,
			"contracts": contracts,
		}
		response.Success(w, responseData)
		return
	}

	response.Success(w, market)
}

// CreateMarket creates a new prediction market
func CreateMarket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	var req CreateMarketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	// Validate input
	if req.Ticker == "" || req.Question == "" || req.Rules == "" || req.ResolutionSource == "" {
		response.BadRequest(w, "missing required fields", nil)
		return
	}

	if req.ExpiresAt.Before(time.Now()) {
		response.BadRequest(w, "expires_at must be in the future", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start transaction
	tx, err := database.GetDB().Begin(ctx)
	if err != nil {
		response.InternalServerError(w, "failed to start transaction", err)
		return
	}
	defer tx.Rollback(ctx)

	// Create market
	marketID := uuid.New()
	_, err = tx.Exec(ctx, `
		INSERT INTO markets (id, ticker, question, rules, resolution_source, status, expires_at)
		VALUES ($1, $2, $3, $4, $5, 'OPEN', $6)
	`, marketID, req.Ticker, req.Question, req.Rules, req.ResolutionSource, req.ExpiresAt)

	if err != nil {
		log.Printf("Failed to create market: %v", err)
		response.Conflict(w, "ticker already exists or database error", err)
		return
	}

	// Create YES and NO contracts
	yesContractID := uuid.New()
	noContractID := uuid.New()

	_, err = tx.Exec(ctx, `
		INSERT INTO contracts (id, market_id, side, ticker)
		VALUES ($1, $2, 'YES', $3), ($4, $5, 'NO', $6)
	`, yesContractID, marketID, req.Ticker+"-YES", noContractID, marketID, req.Ticker+"-NO")

	if err != nil {
		log.Printf("Failed to create contracts: %v", err)
		response.InternalServerError(w, "failed to create contracts", err)
		return
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		response.InternalServerError(w, "failed to commit transaction", err)
		return
	}

	market := &models.Market{
		ID:               marketID,
		Ticker:           req.Ticker,
		Question:         req.Question,
		Rules:            req.Rules,
		ResolutionSource: req.ResolutionSource,
		Status:           "OPEN",
		ExpiresAt:        req.ExpiresAt,
	}

	response.Created(w, market)
}

// UpdateMarket updates a market
func UpdateMarket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	vars := mux.Vars(r)
	marketID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid market ID", err)
		return
	}

	var req UpdateMarketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Build dynamic update query
	query := "UPDATE markets SET "
	args := []interface{}{}
	argIdx := 1
	updates := []string{}

	if req.Question != "" {
		updates = append(updates, "question = $"+string(rune(argIdx+'0')))
		args = append(args, req.Question)
		argIdx++
	}
	if req.Rules != "" {
		updates = append(updates, "rules = $"+string(rune(argIdx+'0')))
		args = append(args, req.Rules)
		argIdx++
	}
	if req.ResolutionSource != "" {
		updates = append(updates, "resolution_source = $"+string(rune(argIdx+'0')))
		args = append(args, req.ResolutionSource)
		argIdx++
	}
	if req.Status != "" {
		updates = append(updates, "status = $"+string(rune(argIdx+'0')))
		args = append(args, req.Status)
		argIdx++
	}
	if req.ExpiresAt != nil {
		updates = append(updates, "expires_at = $"+string(rune(argIdx+'0')))
		args = append(args, *req.ExpiresAt)
		argIdx++
	}

	if len(updates) == 0 {
		response.BadRequest(w, "no fields to update", nil)
		return
	}

	query += updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = $" + string(rune(argIdx+'0'))
	args = append(args, marketID)

	_, err = database.GetDB().Exec(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to update market: %v", err)
		response.InternalServerError(w, "failed to update market", err)
		return
	}

	response.Success(w, map[string]string{"message": "market updated successfully"})
}

// ResolveMarket resolves a market with an outcome
func ResolveMarket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	vars := mux.Vars(r)
	marketID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid market ID", err)
		return
	}

	var req ResolveMarketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	if req.Outcome != "YES" && req.Outcome != "NO" {
		response.BadRequest(w, "outcome must be YES or NO", nil)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	_, err = database.GetDB().Exec(ctx, `
		UPDATE markets
		SET status = 'RESOLVED', outcome = $1, resolved_at = $2
		WHERE id = $3 AND status != 'RESOLVED'
	`, req.Outcome, now, marketID)

	if err != nil {
		log.Printf("Failed to resolve market: %v", err)
		response.InternalServerError(w, "failed to resolve market", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"message":     "market resolved successfully",
		"outcome":     req.Outcome,
		"resolved_at": now,
	})
}

// GetOrderBook returns the order book for a market
func GetOrderBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	vars := mux.Vars(r)
	marketID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid market ID", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get contracts for this market
	rows, err := database.GetDB().Query(ctx, `
		SELECT id, side
		FROM contracts
		WHERE market_id = $1
	`, marketID)
	if err != nil {
		log.Printf("Failed to query contracts: %v", err)
		response.InternalServerError(w, "failed to query contracts", err)
		return
	}
	defer rows.Close()

	orderBooks := make([]OrderBook, 0)
	for rows.Next() {
		var contractID uuid.UUID
		var side string
		if err := rows.Scan(&contractID, &side); err != nil {
			continue
		}

		// Get active orders for this contract (this is a simplified version)
		// In production, the matching engine would provide this data
		bids := []OrderBookEntry{}
		asks := []OrderBookEntry{}

		orderBook := OrderBook{
			ContractID: contractID,
			Side:       side,
			Bids:       bids,
			Asks:       asks,
		}
		orderBooks = append(orderBooks, orderBook)
	}

	response.Success(w, map[string]interface{}{
		"market_id":   marketID,
		"order_books": orderBooks,
	})
}

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/common/response"
	"github.com/Kevin-Kurka/LFG/backend/order-service/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type PlaceOrderRequest struct {
	ContractID string  `json:"contract_id"`
	Type       string  `json:"type"`  // LIMIT or MARKET
	Side       string  `json:"side"`  // BUY or SELL
	Quantity   int     `json:"quantity"`
	LimitPrice *float64 `json:"limit_price,omitempty"`
}

type CancelOrderRequest struct {
	OrderID string `json:"order_id"`
}

// PlaceOrder places a new order
func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	var req PlaceOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	// Validate input
	contractID, err := uuid.Parse(req.ContractID)
	if err != nil {
		response.BadRequest(w, "invalid contract_id", err)
		return
	}

	if req.Type != "LIMIT" && req.Type != "MARKET" {
		response.BadRequest(w, "type must be LIMIT or MARKET", nil)
		return
	}

	if req.Side != "BUY" && req.Side != "SELL" {
		response.BadRequest(w, "side must be BUY or SELL", nil)
		return
	}

	if req.Quantity <= 0 {
		response.BadRequest(w, "quantity must be positive", nil)
		return
	}

	if req.Type == "LIMIT" && (req.LimitPrice == nil || *req.LimitPrice <= 0) {
		response.BadRequest(w, "limit_price is required for LIMIT orders", nil)
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.Unauthorized(w, "invalid user", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create order in database
	orderID := uuid.New()
	now := time.Now()

	var limitPrice interface{}
	if req.LimitPrice != nil {
		limitPrice = *req.LimitPrice
	}

	_, err = database.GetDB().Exec(ctx, `
		INSERT INTO orders (id, user_id, contract_id, type, status, quantity, quantity_filled, limit_price_credits, created_at)
		VALUES ($1, $2, $3, $4, 'PENDING', $5, 0, $6, $7)
	`, orderID, userID, contractID, req.Type, req.Quantity, limitPrice, now)

	if err != nil {
		log.Printf("Failed to create order: %v", err)
		response.InternalServerError(w, "failed to create order", err)
		return
	}

	// Submit to matching engine
	matchingEngineURL := os.Getenv("MATCHING_ENGINE_URL")
	if matchingEngineURL == "" {
		matchingEngineURL = "http://localhost:8084"
	}

	price := 0.0
	if req.LimitPrice != nil {
		price = *req.LimitPrice
	} else {
		// For market orders, use a very high/low price to ensure immediate matching
		if req.Side == "BUY" {
			price = 1.0 // Max price for prediction markets
		} else {
			price = 0.0 // Min price
		}
	}

	submitReq := map[string]interface{}{
		"order_id":    orderID.String(),
		"user_id":     userID.String(),
		"contract_id": contractID.String(),
		"side":        req.Side,
		"price":       price,
		"quantity":    req.Quantity,
	}

	submitBody, err := json.Marshal(submitReq)
	if err != nil {
		log.Printf("Failed to marshal submit request: %v", err)
	} else {
		// Create HTTP client with timeout
		client := &http.Client{Timeout: 10 * time.Second}
		httpReq, err := http.NewRequest("POST", matchingEngineURL+"/submit", bytes.NewBuffer(submitBody))
		if err != nil {
			log.Printf("Failed to create request: %v", err)
		} else {
			httpReq.Header.Set("Content-Type", "application/json")
			httpReq.Header.Set("X-Internal-API-Key", os.Getenv("INTERNAL_API_KEY"))

			resp, err := client.Do(httpReq)
			if err != nil {
				log.Printf("Failed to submit to matching engine: %v", err)
				// Order is still created in DB, will be processed eventually
			} else {
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)
				log.Printf("Matching engine response: %s", string(body))
			}
		}
	}

	var limitPriceCredits float64
	if req.LimitPrice != nil {
		limitPriceCredits = *req.LimitPrice
	}

	order := &models.Order{
		ID:                orderID,
		UserID:            userID,
		ContractID:        contractID,
		Type:              req.Type,
		Status:            "PENDING",
		Quantity:          req.Quantity,
		QuantityFilled:    0,
		LimitPriceCredits: limitPriceCredits,
		CreatedAt:         now,
	}

	response.Created(w, order)
}

// CancelOrder cancels an existing order
func CancelOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	var req CancelOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		response.BadRequest(w, "invalid order_id", err)
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.Unauthorized(w, "invalid user", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Atomically update order status with validation to prevent race conditions
	// This ensures we only cancel orders that are in a valid state and owned by the user
	var contractID uuid.UUID
	var updatedStatus string
	err = database.GetDB().QueryRow(ctx, `
		UPDATE orders
		SET status = 'CANCELLED'
		WHERE id = $1 AND user_id = $2 AND status NOT IN ('FILLED', 'CANCELLED')
		RETURNING contract_id, status
	`, orderID, userID).Scan(&contractID, &updatedStatus)

	if err == pgx.ErrNoRows {
		// Check why it failed - either not found, not authorized, or already in final state
		var checkStatus string
		var checkUserID uuid.UUID
		checkErr := database.GetDB().QueryRow(ctx, `
			SELECT user_id, status FROM orders WHERE id = $1
		`, orderID).Scan(&checkUserID, &checkStatus)

		if checkErr == pgx.ErrNoRows {
			response.NotFound(w, "order not found", nil)
			return
		}
		if checkUserID != userID {
			response.Unauthorized(w, "not authorized to cancel this order", nil)
			return
		}
		if checkStatus == "FILLED" || checkStatus == "CANCELLED" {
			response.BadRequest(w, "order cannot be cancelled", nil)
			return
		}
		response.InternalServerError(w, "failed to cancel order", nil)
		return
	}
	if err != nil {
		log.Printf("Failed to cancel order: %v", err)
		response.InternalServerError(w, "failed to cancel order", err)
		return
	}

	// Notify matching engine
	matchingEngineURL := os.Getenv("MATCHING_ENGINE_URL")
	if matchingEngineURL == "" {
		matchingEngineURL = "http://localhost:8084"
	}

	cancelReq := map[string]string{
		"order_id":    orderID.String(),
		"contract_id": contractID.String(),
	}

	cancelBody, err := json.Marshal(cancelReq)
	if err != nil {
		log.Printf("Failed to marshal cancel request: %v", err)
	} else {
		// Create HTTP client with timeout
		client := &http.Client{Timeout: 10 * time.Second}
		httpReq, err := http.NewRequest("POST", matchingEngineURL+"/cancel", bytes.NewBuffer(cancelBody))
		if err != nil {
			log.Printf("Failed to create request: %v", err)
		} else {
			httpReq.Header.Set("Content-Type", "application/json")
			httpReq.Header.Set("X-Internal-API-Key", os.Getenv("INTERNAL_API_KEY"))

			resp, err := client.Do(httpReq)
			if err != nil {
				log.Printf("Failed to notify matching engine: %v", err)
			} else {
				defer resp.Body.Close()
			}
		}
	}

	response.Success(w, map[string]string{"message": "order cancelled successfully"})
}

// GetOrders returns user's orders
func GetOrders(w http.ResponseWriter, r *http.Request) {
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
		SELECT id, user_id, contract_id, type, status, quantity, quantity_filled, limit_price_credits, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`, userID)
	if err != nil {
		log.Printf("Failed to query orders: %v", err)
		response.InternalServerError(w, "failed to query orders", err)
		return
	}
	defer rows.Close()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var order models.Order
		var contractID uuid.UUID
		err := rows.Scan(&order.ID, &order.UserID, &contractID, &order.Type, &order.Status, &order.Quantity, &order.QuantityFilled, &order.LimitPriceCredits, &order.CreatedAt)
		if err != nil {
			log.Printf("Failed to scan order: %v", err)
			continue
		}
		order.ContractID = contractID
		orders = append(orders, order)
	}

	response.Success(w, map[string]interface{}{
		"orders": orders,
		"count":  len(orders),
	})
}

// GetOrder returns a specific order
func GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	vars := mux.Vars(r)
	orderID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.BadRequest(w, "invalid order_id", err)
		return
	}

	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.Unauthorized(w, "invalid user", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order models.Order
	var contractID uuid.UUID
	err = database.GetDB().QueryRow(ctx, `
		SELECT id, user_id, contract_id, type, status, quantity, quantity_filled, limit_price_credits, created_at
		FROM orders
		WHERE id = $1
	`, orderID).Scan(&order.ID, &order.UserID, &contractID, &order.Type, &order.Status, &order.Quantity, &order.QuantityFilled, &order.LimitPriceCredits, &order.CreatedAt)

	if err == pgx.ErrNoRows {
		response.NotFound(w, "order not found", nil)
		return
	}
	if err != nil {
		log.Printf("Failed to query order: %v", err)
		response.InternalServerError(w, "failed to query order", err)
		return
	}

	if order.UserID != userID {
		response.Unauthorized(w, "not authorized to view this order", nil)
		return
	}

	order.ContractID = contractID
	response.Success(w, order)
}

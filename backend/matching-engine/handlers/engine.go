package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Kevin-Kurka/LFG/backend/common/response"
	"github.com/Kevin-Kurka/LFG/backend/matching-engine/engine"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SubmitOrderRequest struct {
	OrderID    string  `json:"order_id"`
	UserID     string  `json:"user_id"`
	ContractID string  `json:"contract_id"`
	Side       string  `json:"side"`
	Price      float64 `json:"price"`
	Quantity   int     `json:"quantity"`
}

type CancelOrderRequest struct {
	OrderID    string `json:"order_id"`
	ContractID string `json:"contract_id"`
}

// SubmitOrder handles order submission to the matching engine
func SubmitOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	var req SubmitOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		response.BadRequest(w, "invalid order_id", err)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.BadRequest(w, "invalid user_id", err)
		return
	}

	contractID, err := uuid.Parse(req.ContractID)
	if err != nil {
		response.BadRequest(w, "invalid contract_id", err)
		return
	}

	if req.Side != "BUY" && req.Side != "SELL" {
		response.BadRequest(w, "side must be BUY or SELL", nil)
		return
	}

	order := &engine.Order{
		ID:         orderID,
		UserID:     userID,
		ContractID: contractID,
		Side:       req.Side,
		Price:      req.Price,
		Quantity:   req.Quantity,
		Filled:     0,
		Timestamp:  time.Now(),
	}

	trades, err := engine.GlobalEngine.SubmitOrder(order)
	if err != nil {
		response.InternalServerError(w, "failed to submit order", err)
		return
	}

	response.Success(w, map[string]interface{}{
		"order":  order,
		"trades": trades,
	})
}

// CancelOrder handles order cancellation
func CancelOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
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

	contractID, err := uuid.Parse(req.ContractID)
	if err != nil {
		response.BadRequest(w, "invalid contract_id", err)
		return
	}

	success := engine.GlobalEngine.CancelOrder(contractID, orderID)
	if !success {
		response.NotFound(w, "order not found", nil)
		return
	}

	response.Success(w, map[string]string{"message": "order cancelled successfully"})
}

// GetOrderBook returns the order book for a contract
func GetOrderBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	vars := mux.Vars(r)
	contractID, err := uuid.Parse(vars["contractId"])
	if err != nil {
		response.BadRequest(w, "invalid contract_id", err)
		return
	}

	bids, asks, exists := engine.GlobalEngine.GetOrderBook(contractID)
	if !exists {
		response.Success(w, map[string]interface{}{
			"contract_id": contractID,
			"bids":        []engine.Order{},
			"asks":        []engine.Order{},
		})
		return
	}

	response.Success(w, map[string]interface{}{
		"contract_id": contractID,
		"bids":        bids,
		"asks":        asks,
	})
}

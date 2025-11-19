package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"lfg/shared/models"
	"lfg/order-service/repository"
	pb "lfg/matching-engine/proto"
)

// OrderHandler handles HTTP requests for order operations
type OrderHandler struct {
	repo               *repository.OrderRepository
	matchingEngineAddr string
	walletServiceURL   string
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(repo *repository.OrderRepository, matchingEngineAddr string, walletServiceURL string) *OrderHandler {
	return &OrderHandler{
		repo:               repo,
		matchingEngineAddr: matchingEngineAddr,
		walletServiceURL:   walletServiceURL,
	}
}

// PlaceOrder handles order placement
func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from header (set by API gateway)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse request
	var req models.OrderPlaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Quantity <= 0 {
		respondError(w, "Quantity must be positive", http.StatusBadRequest)
		return
	}

	if req.Type == models.OrderTypeLimit && (req.LimitPriceCredits == nil || *req.LimitPriceCredits <= 0) {
		respondError(w, "Limit price required for limit orders", http.StatusBadRequest)
		return
	}

	// Check wallet balance before placing order
	if req.Type == models.OrderTypeLimit && req.LimitPriceCredits != nil {
		requiredBalance := float64(req.Quantity) * (*req.LimitPriceCredits)
		balance, err := h.checkWalletBalance(userIDStr)
		if err != nil {
			respondError(w, "Failed to check wallet balance", http.StatusInternalServerError)
			return
		}

		if balance < requiredBalance {
			errorMsg := fmt.Sprintf("Insufficient balance. Required: %.2f credits, Available: %.2f credits", requiredBalance, balance)
			respondError(w, errorMsg, http.StatusBadRequest)
			return
		}
	}

	// Create order in database
	order := &models.Order{
		ID:                uuid.New(),
		UserID:            userID,
		ContractID:        req.ContractID,
		Type:              req.Type,
		Status:            models.OrderStatusPending,
		Quantity:          req.Quantity,
		QuantityFilled:    0,
		LimitPriceCredits: req.LimitPriceCredits,
	}

	if err := h.repo.Create(r.Context(), order); err != nil {
		respondError(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	// Submit to matching engine via gRPC
	conn, err := grpc.NewClient(h.matchingEngineAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		respondError(w, "Failed to connect to matching engine", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewMatchingEngineClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Determine order side from contract
	var contractSide string
	err = h.repo.GetPool().QueryRow(r.Context(), `
		SELECT side FROM contracts WHERE id = $1
	`, req.ContractID).Scan(&contractSide)

	if err != nil {
		respondError(w, "Failed to fetch contract details", http.StatusInternalServerError)
		return
	}

	orderSide := pb.OrderSide_BUY
	if contractSide == "NO" {
		orderSide = pb.OrderSide_SELL
	}

	// Determine order type
	orderType := pb.OrderType_MARKET
	if req.Type == models.OrderTypeLimit {
		orderType = pb.OrderType_LIMIT
	}

	limitPrice := 0.0
	if req.LimitPriceCredits != nil {
		limitPrice = *req.LimitPriceCredits
	}

	grpcReq := &pb.PlaceOrderRequest{
		OrderId:    order.ID.String(),
		UserId:     userID.String(),
		ContractId: req.ContractID.String(),
		Type:       orderType,
		Side:       orderSide,
		Quantity:   int32(req.Quantity),
		LimitPrice: limitPrice,
	}

	resp, err := client.PlaceOrder(ctx, grpcReq)
	if err != nil {
		// Update order status to rejected
		h.repo.UpdateStatus(r.Context(), order.ID, models.OrderStatusRejected, 0)
		respondError(w, "Failed to place order", http.StatusInternalServerError)
		return
	}

	// Update order based on matching engine response
	status := models.OrderStatusActive
	if resp.Status == "FILLED" {
		status = models.OrderStatusFilled
	} else if resp.Status == "PARTIALLY_FILLED" {
		status = models.OrderStatusPartiallyFilled
	}

	h.repo.UpdateStatus(r.Context(), order.ID, status, int(resp.QuantityFilled))

	// Return response
	response := models.OrderPlaceResponse{
		OrderID:        order.ID,
		Status:         status,
		QuantityFilled: int(resp.QuantityFilled),
		AveragePrice:   resp.AveragePrice,
	}

	respondJSON(w, response, http.StatusCreated)
}

// CancelOrder handles order cancellation
func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse request
	var req models.OrderCancelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get order
	order, err := h.repo.GetByID(r.Context(), req.OrderID)
	if err != nil {
		if err == repository.ErrOrderNotFound {
			respondError(w, "Order not found", http.StatusNotFound)
			return
		}
		respondError(w, "Failed to get order", http.StatusInternalServerError)
		return
	}

	// Verify ownership
	if order.UserID != userID {
		respondError(w, "Not authorized to cancel this order", http.StatusForbidden)
		return
	}

	// Check if order can be cancelled
	if order.Status != models.OrderStatusActive && order.Status != models.OrderStatusPartiallyFilled {
		respondError(w, "Order cannot be cancelled", http.StatusBadRequest)
		return
	}

	// Cancel in matching engine
	conn, err := grpc.NewClient(h.matchingEngineAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		respondError(w, "Failed to connect to matching engine", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewMatchingEngineClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.CancelOrderRequest{
		OrderId:    order.ID.String(),
		ContractId: order.ContractID.String(),
	}

	_, err = client.CancelOrder(ctx, grpcReq)
	if err != nil {
		respondError(w, "Failed to cancel order", http.StatusInternalServerError)
		return
	}

	// Update order status in database
	if err := h.repo.Cancel(r.Context(), order.ID); err != nil {
		respondError(w, "Failed to update order status", http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"status": "cancelled"}, http.StatusOK)
}

// GetOrderStatus handles order status retrieval
func (h *OrderHandler) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		respondError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get order ID from query
	orderIDStr := r.URL.Query().Get("id")
	if orderIDStr == "" {
		// Return all user orders
		status := r.URL.Query().Get("status")
		orders, err := h.repo.GetByUserID(r.Context(), userID, status, 100)
		if err != nil {
			respondError(w, "Failed to get orders", http.StatusInternalServerError)
			return
		}
		respondJSON(w, map[string]interface{}{"orders": orders}, http.StatusOK)
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		respondError(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Get specific order
	order, err := h.repo.GetByID(r.Context(), orderID)
	if err != nil {
		if err == repository.ErrOrderNotFound {
			respondError(w, "Order not found", http.StatusNotFound)
			return
		}
		respondError(w, "Failed to get order", http.StatusInternalServerError)
		return
	}

	// Verify ownership
	if order.UserID != userID {
		respondError(w, "Not authorized to view this order", http.StatusForbidden)
		return
	}

	respondJSON(w, order, http.StatusOK)
}

// Health check handler
func Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{"status": "healthy"}, http.StatusOK)
}

// checkWalletBalance checks the user's wallet balance via HTTP call to wallet service
func (h *OrderHandler) checkWalletBalance(userID string) (float64, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Create request to wallet service
	req, err := http.NewRequest(http.MethodGet, h.walletServiceURL+"/balance", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set user ID header
	req.Header.Set("X-User-ID", userID)

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call wallet service: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("wallet service returned status %d", resp.StatusCode)
	}

	// Parse response
	var balanceResp struct {
		Balance float64 `json:"balance"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&balanceResp); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return balanceResp.Balance, nil
}

// Helper functions
func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	respondJSON(w, map[string]string{"error": message}, statusCode)
}

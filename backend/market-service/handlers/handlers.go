package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "lfg/matching-engine/proto"
	"lfg/shared/models"
	"lfg/market-service/repository"
)

// MarketHandler handles HTTP requests for market operations
type MarketHandler struct {
	repo               *repository.MarketRepository
	matchingEngineAddr string
}

// NewMarketHandler creates a new market handler
func NewMarketHandler(repo *repository.MarketRepository, matchingEngineAddr string) *MarketHandler {
	return &MarketHandler{
		repo:               repo,
		matchingEngineAddr: matchingEngineAddr,
	}
}

// ListMarkets handles listing all available markets
func (h *MarketHandler) ListMarkets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters
	status := r.URL.Query().Get("status")
	search := r.URL.Query().Get("search")
	page := parseInt(r.URL.Query().Get("page"), 1)
	pageSize := parseInt(r.URL.Query().Get("page_size"), 20)

	// Validate page size
	if pageSize > 100 {
		pageSize = 100
	}

	// Get markets
	markets, totalCount, err := h.repo.List(r.Context(), status, search, page, pageSize)
	if err != nil {
		respondError(w, "Failed to list markets", http.StatusInternalServerError)
		return
	}

	// Return response
	response := models.MarketListResponse{
		Markets:    markets,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}

	respondJSON(w, response, http.StatusOK)
}

// MarketDetail handles retrieving detailed information for a single market
func (h *MarketHandler) MarketDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get market ID from query parameter
	marketIDStr := r.URL.Query().Get("id")
	if marketIDStr == "" {
		respondError(w, "Market ID is required", http.StatusBadRequest)
		return
	}

	marketID, err := uuid.Parse(marketIDStr)
	if err != nil {
		respondError(w, "Invalid market ID", http.StatusBadRequest)
		return
	}

	// Get market
	market, err := h.repo.GetByID(r.Context(), marketID)
	if err != nil {
		if err == repository.ErrMarketNotFound {
			respondError(w, "Market not found", http.StatusNotFound)
			return
		}
		respondError(w, "Failed to get market", http.StatusInternalServerError)
		return
	}

	// Get contracts
	contracts, err := h.repo.GetContractsByMarketID(r.Context(), marketID)
	if err != nil {
		respondError(w, "Failed to get contracts", http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]interface{}{
		"market":    market,
		"contracts": contracts,
	}

	respondJSON(w, response, http.StatusOK)
}

// OrderBook handles retrieving the current order book for a market
func (h *MarketHandler) OrderBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get contract ID from query parameter
	contractIDStr := r.URL.Query().Get("contract_id")
	if contractIDStr == "" {
		respondError(w, "Contract ID is required", http.StatusBadRequest)
		return
	}

	contractID, err := uuid.Parse(contractIDStr)
	if err != nil {
		respondError(w, "Invalid contract ID", http.StatusBadRequest)
		return
	}

	// Fetch order book from matching engine via gRPC
	conn, err := grpc.NewClient(h.matchingEngineAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		respondError(w, "Failed to connect to matching engine", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewMatchingEngineClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.GetOrderBookRequest{
		ContractId: contractID.String(),
		Depth:      20,
	}

	resp, err := client.GetOrderBook(ctx, grpcReq)
	if err != nil {
		// Return empty order book if matching engine call fails
		respondJSON(w, map[string]interface{}{
			"contract_id": contractID,
			"bids":        []interface{}{},
			"asks":        []interface{}{},
		}, http.StatusOK)
		return
	}

	// Transform protobuf response to JSON
	bids := make([]map[string]interface{}, len(resp.Bids))
	for i, bid := range resp.Bids {
		bids[i] = map[string]interface{}{
			"price":       bid.Price,
			"quantity":    bid.Quantity,
			"order_count": bid.OrderCount,
		}
	}

	asks := make([]map[string]interface{}, len(resp.Asks))
	for i, ask := range resp.Asks {
		asks[i] = map[string]interface{}{
			"price":       ask.Price,
			"quantity":    ask.Quantity,
			"order_count": ask.OrderCount,
		}
	}

	respondJSON(w, map[string]interface{}{
		"contract_id": contractID,
		"bids":        bids,
		"asks":        asks,
	}, http.StatusOK)
}

// Health check handler
func Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{"status": "healthy"}, http.StatusOK)
}

// Helper functions
func parseInt(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return v
}

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	respondJSON(w, map[string]string{"error": message}, statusCode)
}

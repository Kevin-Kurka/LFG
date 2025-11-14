package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"lfg/shared/models"
	"lfg/market-service/repository"
)

// MarketHandler handles HTTP requests for market operations
type MarketHandler struct {
	repo *repository.MarketRepository
}

// NewMarketHandler creates a new market handler
func NewMarketHandler(repo *repository.MarketRepository) *MarketHandler {
	return &MarketHandler{repo: repo}
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

	_, err := uuid.Parse(contractIDStr)
	if err != nil {
		respondError(w, "Invalid contract ID", http.StatusBadRequest)
		return
	}

	// TODO: Fetch order book from matching engine via gRPC
	// For now, return empty order book
	orderBook := map[string]interface{}{
		"bids": []interface{}{},
		"asks": []interface{}{},
	}

	respondJSON(w, orderBook, http.StatusOK)
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

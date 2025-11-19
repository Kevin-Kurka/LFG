package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"lfg/shared/models"
	"lfg/wallet-service/repository"
)

// WalletHandler handles HTTP requests for wallet operations
type WalletHandler struct {
	repo *repository.WalletRepository
}

// NewWalletHandler creates a new wallet handler
func NewWalletHandler(repo *repository.WalletRepository) *WalletHandler {
	return &WalletHandler{repo: repo}
}

// Balance handles getting wallet balance
func (h *WalletHandler) Balance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from header (set by API gateway auth middleware)
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

	// Get wallet
	wallet, err := h.repo.GetByUserID(r.Context(), userID)
	if err != nil {
		if err == repository.ErrWalletNotFound {
			respondError(w, "Wallet not found", http.StatusNotFound)
			return
		}
		respondError(w, "Failed to get wallet", http.StatusInternalServerError)
		return
	}

	// Calculate locked balance from active orders
	lockedBalance, err := h.repo.GetLockedBalance(r.Context(), userID)
	if err != nil {
		respondError(w, "Failed to calculate locked balance", http.StatusInternalServerError)
		return
	}

	// Calculate available balance
	availableBalance := wallet.BalanceCredits - lockedBalance
	if availableBalance < 0 {
		availableBalance = 0
	}

	// Return balance response
	response := models.WalletBalanceResponse{
		Balance:          wallet.BalanceCredits,
		AvailableBalance: availableBalance,
		LockedBalance:    lockedBalance,
	}

	respondJSON(w, response, http.StatusOK)
}

// Transactions handles getting transaction history
func (h *WalletHandler) Transactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from header
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

	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get wallet
	wallet, err := h.repo.GetByUserID(r.Context(), userID)
	if err != nil {
		if err == repository.ErrWalletNotFound {
			respondError(w, "Wallet not found", http.StatusNotFound)
			return
		}
		respondError(w, "Failed to get wallet", http.StatusInternalServerError)
		return
	}

	// Get transaction history from credit_transactions table
	transactions, err := h.repo.GetTransactionHistory(r.Context(), userID, limit, offset)
	if err != nil {
		respondError(w, "Failed to fetch transaction history", http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]interface{}{
		"wallet_id":    wallet.ID,
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
	}, http.StatusOK)
}

// Credit handles adding credits to a wallet (internal use)
func (h *WalletHandler) Credit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID      string  `json:"user_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		respondError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		respondError(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	// Credit the wallet
	if err := h.repo.Credit(r.Context(), userID, req.Amount, req.Description); err != nil {
		respondError(w, "Failed to credit wallet", http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"status": "success"}, http.StatusOK)
}

// Debit handles removing credits from a wallet (internal use)
func (h *WalletHandler) Debit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID      string  `json:"user_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		respondError(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		respondError(w, "Amount must be positive", http.StatusBadRequest)
		return
	}

	// Debit the wallet
	if err := h.repo.Debit(r.Context(), userID, req.Amount, req.Description); err != nil {
		if err == repository.ErrInsufficientBalance {
			respondError(w, "Insufficient balance", http.StatusBadRequest)
			return
		}
		respondError(w, "Failed to debit wallet", http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"status": "success"}, http.StatusOK)
}

// Health check handler
func Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{"status": "healthy"}, http.StatusOK)
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

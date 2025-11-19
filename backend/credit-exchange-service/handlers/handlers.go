package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"lfg/shared/models"
	"lfg/credit-exchange-service/repository"
)

// ExchangeHandler handles credit exchange operations
type ExchangeHandler struct {
	txRepo           *repository.CreditTransactionRepository
	walletServiceURL string
}

// NewExchangeHandler creates a new exchange handler
func NewExchangeHandler(txRepo *repository.CreditTransactionRepository, walletServiceURL string) *ExchangeHandler {
	return &ExchangeHandler{
		txRepo:           txRepo,
		walletServiceURL: walletServiceURL,
	}
}

// BuyCredits handles credit purchase with cryptocurrency
func (h *ExchangeHandler) BuyCredits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	// Parse request
	var req models.BuyCreditsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate crypto type
	if req.CryptoType != "BTC" && req.CryptoType != "ETH" && req.CryptoType != "USDC" {
		respondError(w, "Invalid crypto type. Supported: BTC, ETH, USDC", http.StatusBadRequest)
		return
	}

	if req.CryptoAmount <= 0 {
		respondError(w, "Crypto amount must be positive", http.StatusBadRequest)
		return
	}

	// Calculate credit amount (mock exchange rates)
	// In production, fetch real-time rates from API
	creditAmount := calculateCreditsFromCrypto(req.CryptoType, req.CryptoAmount)

	// Create transaction record
	tx := &models.CreditTransaction{
		ID:           uuid.New(),
		UserID:       userID,
		Type:         models.CreditTransactionTypePurchase,
		CryptoType:   req.CryptoType,
		CryptoAmount: req.CryptoAmount,
		CreditAmount: creditAmount,
		Status:       models.CreditTransactionStatusPending,
	}

	if err := h.txRepo.Create(r.Context(), tx); err != nil {
		respondError(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	// TODO: In production, integrate with actual crypto payment gateway
	// For demo purposes, immediately mark as completed
	h.txRepo.UpdateStatus(r.Context(), tx.ID, models.CreditTransactionStatusCompleted)

	// Credit user's wallet via internal API call
	if err := h.creditWallet(userIDStr, creditAmount, tx.ID.String(), "PURCHASE"); err != nil {
		log.Printf("Failed to credit wallet: %v", err)
		respondError(w, "Transaction recorded but wallet credit failed", http.StatusInternalServerError)
		return
	}

	// Return response
	response := map[string]interface{}{
		"transaction_id": tx.ID,
		"status":         "COMPLETED",
		"crypto_type":    req.CryptoType,
		"crypto_amount":  req.CryptoAmount,
		"credits":        creditAmount,
		"message":        fmt.Sprintf("Successfully purchased %.2f credits", creditAmount),
	}

	respondJSON(w, response, http.StatusCreated)
}

// SellCredits handles credit sale for cryptocurrency
func (h *ExchangeHandler) SellCredits(w http.ResponseWriter, r *http.Request) {
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
	var req models.SellCreditsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate crypto type
	if req.CryptoType != "BTC" && req.CryptoType != "ETH" && req.CryptoType != "USDC" {
		respondError(w, "Invalid crypto type. Supported: BTC, ETH, USDC", http.StatusBadRequest)
		return
	}

	if req.CreditAmount <= 0 {
		respondError(w, "Credit amount must be positive", http.StatusBadRequest)
		return
	}

	// Check user has sufficient credits in wallet
	balance, err := h.getWalletBalance(userIDStr)
	if err != nil {
		respondError(w, "Failed to check wallet balance", http.StatusInternalServerError)
		return
	}

	if balance < req.CreditAmount {
		errorMsg := fmt.Sprintf("Insufficient balance. Required: %.2f credits, Available: %.2f credits", req.CreditAmount, balance)
		respondError(w, errorMsg, http.StatusBadRequest)
		return
	}

	// Calculate crypto amount (mock exchange rates)
	cryptoAmount := calculateCryptoFromCredits(req.CryptoType, req.CreditAmount)

	// Create transaction record
	tx := &models.CreditTransaction{
		ID:           uuid.New(),
		UserID:       userID,
		Type:         models.CreditTransactionTypeSale,
		CryptoType:   req.CryptoType,
		CryptoAmount: cryptoAmount,
		CreditAmount: req.CreditAmount,
		Status:       models.CreditTransactionStatusPending,
	}

	if err := h.txRepo.Create(r.Context(), tx); err != nil {
		respondError(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	// Debit user's wallet
	if err := h.debitWallet(userIDStr, req.CreditAmount, tx.ID.String(), "SALE"); err != nil {
		log.Printf("Failed to debit wallet: %v", err)
		respondError(w, "Transaction recorded but wallet debit failed", http.StatusInternalServerError)
		return
	}

	// TODO: Process crypto payout
	// For demo, mark as completed
	h.txRepo.UpdateStatus(r.Context(), tx.ID, models.CreditTransactionStatusCompleted)

	// Return response
	response := map[string]interface{}{
		"transaction_id": tx.ID,
		"status":         "COMPLETED",
		"crypto_type":    req.CryptoType,
		"crypto_amount":  cryptoAmount,
		"credits":        req.CreditAmount,
		"message":        fmt.Sprintf("Successfully sold %.2f credits for %.8f %s", req.CreditAmount, cryptoAmount, req.CryptoType),
	}

	respondJSON(w, response, http.StatusCreated)
}

// ExchangeHistory handles retrieving transaction history
func (h *ExchangeHandler) ExchangeHistory(w http.ResponseWriter, r *http.Request) {
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

	// Get transaction history
	transactions, err := h.txRepo.GetByUserID(r.Context(), userID, 100)
	if err != nil {
		respondError(w, "Failed to get transaction history", http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]interface{}{
		"transactions": transactions,
		"count":        len(transactions),
	}, http.StatusOK)
}

// Health check
func Health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, map[string]string{"status": "healthy"}, http.StatusOK)
}

// Helper functions
func calculateCreditsFromCrypto(cryptoType string, cryptoAmount float64) float64 {
	// Mock exchange rates (1 credit = $1 USD)
	rates := map[string]float64{
		"BTC":  50000.0, // 1 BTC = 50,000 credits
		"ETH":  3000.0,  // 1 ETH = 3,000 credits
		"USDC": 1.0,     // 1 USDC = 1 credit
	}

	return cryptoAmount * rates[cryptoType]
}

func calculateCryptoFromCredits(cryptoType string, creditAmount float64) float64 {
	// Mock exchange rates (1 credit = $1 USD)
	rates := map[string]float64{
		"BTC":  50000.0,
		"ETH":  3000.0,
		"USDC": 1.0,
	}

	return creditAmount / rates[cryptoType]
}

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	respondJSON(w, map[string]string{"error": message}, statusCode)
}

// getWalletBalance fetches the user's wallet balance
func (h *ExchangeHandler) getWalletBalance(userID string) (float64, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", h.walletServiceURL+"/balance", nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-User-ID", userID)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("wallet service returned status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	balance, ok := result["available_balance"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid balance response from wallet service")
	}

	return balance, nil
}

// creditWallet credits the user's wallet
func (h *ExchangeHandler) creditWallet(userID string, amount float64, transactionID string, transactionType string) error {
	client := &http.Client{Timeout: 5 * time.Second}

	payload := map[string]interface{}{
		"amount":           amount,
		"transaction_id":   transactionID,
		"transaction_type": transactionType,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", h.walletServiceURL+"/credit", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("wallet service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// debitWallet debits the user's wallet
func (h *ExchangeHandler) debitWallet(userID string, amount float64, transactionID string, transactionType string) error {
	client := &http.Client{Timeout: 5 * time.Second}

	payload := map[string]interface{}{
		"amount":           amount,
		"transaction_id":   transactionID,
		"transaction_type": transactionType,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", h.walletServiceURL+"/debit", bytes.NewBuffer(payloadJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("wallet service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

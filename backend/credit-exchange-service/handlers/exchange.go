package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/response"
	"github.com/Kevin-Kurka/LFG/backend/common/validation"
)

type BuyCreditsRequest struct {
	Amount          float64 `json:"amount"`
	PaymentMethod   string  `json:"payment_method"`
	PaymentDetails  string  `json:"payment_details,omitempty"`
}

type SellCreditsRequest struct {
	Amount float64 `json:"amount"`
}

func BuyCredits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	var req BuyCreditsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	// Validate amount (min $1, max $10,000)
	if err := validation.ValidateAmount(req.Amount, 0, 10000); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	// Validate payment method
	if err := validation.ValidatePaymentMethod(req.PaymentMethod); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	// Sanitize payment details
	if req.PaymentDetails != "" {
		req.PaymentDetails = validation.SanitizeString(req.PaymentDetails)
	}

	// Mock payment processing (in production, integrate with Stripe/PayPal)
	log.Printf("Processing payment of $%.2f for user %s", req.Amount, userCtx.UserID)

	// Create wallet transaction via wallet service
	walletServiceURL := os.Getenv("WALLET_SERVICE_URL")
	if walletServiceURL == "" {
		walletServiceURL = "http://localhost:8081"
	}

	txReq := map[string]interface{}{
		"user_id":     userCtx.UserID,
		"type":        "CREDIT_BUY",
		"amount":      req.Amount,
		"reference":   req.PaymentMethod,
		"description": "Credits purchased via " + req.PaymentMethod,
	}

	txBody, err := json.Marshal(txReq)
	if err != nil {
		log.Printf("Failed to marshal transaction request: %v", err)
		response.InternalServerError(w, "failed to process request", nil)
		return
	}

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", walletServiceURL+"/internal/transactions", bytes.NewBuffer(txBody))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		response.InternalServerError(w, "failed to process request", nil)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-API-Key", os.Getenv("INTERNAL_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to create wallet transaction: %v", err)
		response.InternalServerError(w, "failed to process credit purchase", nil)
		return
	}
	defer resp.Body.Close()

	response.Success(w, map[string]interface{}{
		"message": "credits purchased successfully",
		"amount":  req.Amount,
	})
}

func SellCredits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	var req SellCreditsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	// Validate amount (min $1, max $10,000)
	if err := validation.ValidateAmount(req.Amount, 0, 10000); err != nil {
		response.BadRequest(w, err.Error(), nil)
		return
	}

	// Create wallet transaction (negative amount for withdrawal)
	walletServiceURL := os.Getenv("WALLET_SERVICE_URL")
	if walletServiceURL == "" {
		walletServiceURL = "http://localhost:8081"
	}

	txReq := map[string]interface{}{
		"user_id":     userCtx.UserID,
		"type":        "CREDIT_SELL",
		"amount":      -req.Amount, // Negative for withdrawal
		"description": "Credits sold/withdrawn",
	}

	txBody, err := json.Marshal(txReq)
	if err != nil {
		log.Printf("Failed to marshal transaction request: %v", err)
		response.InternalServerError(w, "failed to process request", nil)
		return
	}

	// Create HTTP client with timeout
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", walletServiceURL+"/internal/transactions", bytes.NewBuffer(txBody))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		response.InternalServerError(w, "failed to process request", nil)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-API-Key", os.Getenv("INTERNAL_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to create wallet transaction: %v", err)
		response.InternalServerError(w, "failed to process credit sale", nil)
		return
	}
	defer resp.Body.Close()

	response.Success(w, map[string]interface{}{
		"message": "credits sold successfully",
		"amount":  req.Amount,
	})
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	_, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	// Forward to wallet service
	walletServiceURL := os.Getenv("WALLET_SERVICE_URL")
	if walletServiceURL == "" {
		walletServiceURL = "http://localhost:8081"
	}

	req, _ := http.NewRequest("GET", walletServiceURL+"/transactions?type=CREDIT_BUY&type=CREDIT_SELL", nil)
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		response.InternalServerError(w, "failed to get history", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	response.Success(w, result)
}

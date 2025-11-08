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

	if req.Amount <= 0 {
		response.BadRequest(w, "amount must be positive", nil)
		return
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

	txBody, _ := json.Marshal(txReq)
	resp, err := http.Post(walletServiceURL+"/internal/transactions", "application/json", bytes.NewBuffer(txBody))
	if err != nil {
		log.Printf("Failed to create wallet transaction: %v", err)
		response.InternalServerError(w, "failed to process credit purchase", err)
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

	if req.Amount <= 0 {
		response.BadRequest(w, "amount must be positive", nil)
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

	txBody, _ := json.Marshal(txReq)
	resp, err := http.Post(walletServiceURL+"/internal/transactions", "application/json", bytes.NewBuffer(txBody))
	if err != nil {
		log.Printf("Failed to create wallet transaction: %v", err)
		response.InternalServerError(w, "failed to process credit sale", err)
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

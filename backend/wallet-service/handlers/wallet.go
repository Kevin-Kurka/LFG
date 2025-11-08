package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/common/response"
	"github.com/Kevin-Kurka/LFG/backend/wallet-service/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetBalance returns the user's wallet balance
func GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	// Get user from context
	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get wallet from database
	var wallet models.Wallet
	err = database.GetDB().QueryRow(ctx, `
		SELECT id, user_id, balance_credits, created_at
		FROM wallets
		WHERE user_id = $1
	`, userCtx.UserID).Scan(&wallet.ID, &wallet.UserID, &wallet.BalanceCredits, &wallet.CreatedAt)

	if err == pgx.ErrNoRows {
		response.NotFound(w, "wallet not found", nil)
		return
	}
	if err != nil {
		log.Printf("Failed to query wallet: %v", err)
		response.InternalServerError(w, "failed to query wallet", err)
		return
	}

	response.Success(w, wallet)
}

// GetTransactions returns the user's transaction history
func GetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	// Get user from context
	userCtx, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		response.Unauthorized(w, "unauthorized", err)
		return
	}

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	txType := r.URL.Query().Get("type")

	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build query
	query := `
		SELECT id, user_id, type, amount, balance, reference, reference_id, description, created_at
		FROM wallet_transactions
		WHERE user_id = $1
	`
	args := []interface{}{userCtx.UserID}

	if txType != "" {
		query += " AND type = $2"
		args = append(args, txType)
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := database.GetDB().Query(ctx, query, args...)
	if err != nil {
		log.Printf("Failed to query transactions: %v", err)
		response.InternalServerError(w, "failed to query transactions", err)
		return
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		var tx models.Transaction
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Type,
			&tx.Amount,
			&tx.Balance,
			&tx.Reference,
			&tx.ReferenceID,
			&tx.Description,
			&tx.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan transaction: %v", err)
			continue
		}
		transactions = append(transactions, tx)
	}

	responseData := map[string]interface{}{
		"transactions": transactions,
		"limit":        limit,
		"offset":       offset,
		"count":        len(transactions),
	}

	response.Success(w, responseData)
}

type CreateTransactionRequest struct {
	UserID      string  `json:"user_id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Reference   string  `json:"reference,omitempty"`
	ReferenceID string  `json:"reference_id,omitempty"`
	Description string  `json:"description"`
}

// CreateTransaction creates a new transaction (internal use only)
func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.BadRequest(w, "method not allowed", nil)
		return
	}

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body", err)
		return
	}

	// Validate input
	if req.UserID == "" {
		response.BadRequest(w, "user_id is required", nil)
		return
	}
	if req.Type == "" {
		response.BadRequest(w, "type is required", nil)
		return
	}
	if req.Amount == 0 {
		response.BadRequest(w, "amount must be non-zero", nil)
		return
	}
	if req.Description == "" {
		response.BadRequest(w, "description is required", nil)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.BadRequest(w, "invalid user_id", err)
		return
	}

	var referenceID *uuid.UUID
	if req.ReferenceID != "" {
		refID, err := uuid.Parse(req.ReferenceID)
		if err != nil {
			response.BadRequest(w, "invalid reference_id", err)
			return
		}
		referenceID = &refID
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

	// Update wallet balance
	var newBalance float64
	err = tx.QueryRow(ctx, `
		UPDATE wallets
		SET balance_credits = balance_credits + $1
		WHERE user_id = $2
		RETURNING balance_credits
	`, req.Amount, userID).Scan(&newBalance)

	if err != nil {
		log.Printf("Failed to update wallet balance: %v", err)
		response.InternalServerError(w, "failed to update wallet balance", err)
		return
	}

	// Create transaction record
	transactionID := uuid.New()
	now := time.Now()

	_, err = tx.Exec(ctx, `
		INSERT INTO wallet_transactions (id, user_id, type, amount, balance, reference, reference_id, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`, transactionID, userID, req.Type, req.Amount, newBalance, req.Reference, referenceID, req.Description, now)

	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		response.InternalServerError(w, "failed to create transaction", err)
		return
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		response.InternalServerError(w, "failed to commit transaction", err)
		return
	}

	transaction := &models.Transaction{
		ID:          transactionID,
		UserID:      userID,
		Type:        req.Type,
		Amount:      req.Amount,
		Balance:     newBalance,
		Reference:   req.Reference,
		ReferenceID: referenceID,
		Description: req.Description,
		CreatedAt:   now,
	}

	response.Created(w, transaction)
}

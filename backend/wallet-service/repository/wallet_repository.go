package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lfg/shared/models"
)

var (
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

// WalletRepository handles wallet database operations
type WalletRepository struct {
	pool *pgxpool.Pool
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(pool *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{pool: pool}
}

// Create creates a new wallet
func (r *WalletRepository) Create(ctx context.Context, wallet *models.Wallet) error {
	query := `
		INSERT INTO wallets (id, user_id, balance_credits, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`

	_, err := r.pool.Exec(ctx, query,
		wallet.ID,
		wallet.UserID,
		wallet.BalanceCredits,
	)

	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

// GetByUserID retrieves a wallet by user ID
func (r *WalletRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, balance_credits, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
	`

	var wallet models.Wallet
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.BalanceCredits,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

// GetByID retrieves a wallet by ID
func (r *WalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	query := `
		SELECT id, user_id, balance_credits, created_at, updated_at
		FROM wallets
		WHERE id = $1
	`

	var wallet models.Wallet
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.BalanceCredits,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

// UpdateBalance updates wallet balance in a transaction-safe way
func (r *WalletRepository) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount float64, description string) error {
	// Start transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get current balance with row lock
	var currentBalance float64
	err = tx.QueryRow(ctx, `
		SELECT balance_credits FROM wallets WHERE id = $1 FOR UPDATE
	`, walletID).Scan(&currentBalance)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrWalletNotFound
		}
		return fmt.Errorf("failed to lock wallet: %w", err)
	}

	// Check if balance would go negative
	newBalance := currentBalance + amount
	if newBalance < 0 {
		return ErrInsufficientBalance
	}

	// Update balance
	_, err = tx.Exec(ctx, `
		UPDATE wallets SET balance_credits = $1, updated_at = NOW() WHERE id = $2
	`, newBalance, walletID)

	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Credit adds credits to a wallet (atomic)
func (r *WalletRepository) Credit(ctx context.Context, userID uuid.UUID, amount float64, description string) error {
	wallet, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return r.UpdateBalance(ctx, wallet.ID, amount, description)
}

// Debit removes credits from a wallet (atomic)
func (r *WalletRepository) Debit(ctx context.Context, userID uuid.UUID, amount float64, description string) error {
	wallet, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	return r.UpdateBalance(ctx, wallet.ID, -amount, description)
}

// Transfer transfers credits between wallets (atomic)
func (r *WalletRepository) Transfer(ctx context.Context, fromUserID, toUserID uuid.UUID, amount float64, description string) error {
	// Start transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get both wallets with locks
	var fromWalletID, toWalletID uuid.UUID
	var fromBalance, toBalance float64

	err = tx.QueryRow(ctx, `
		SELECT id, balance_credits FROM wallets WHERE user_id = $1 FOR UPDATE
	`, fromUserID).Scan(&fromWalletID, &fromBalance)

	if err != nil {
		return fmt.Errorf("failed to lock from wallet: %w", err)
	}

	err = tx.QueryRow(ctx, `
		SELECT id, balance_credits FROM wallets WHERE user_id = $1 FOR UPDATE
	`, toUserID).Scan(&toWalletID, &toBalance)

	if err != nil {
		return fmt.Errorf("failed to lock to wallet: %w", err)
	}

	// Check sufficient balance
	if fromBalance < amount {
		return ErrInsufficientBalance
	}

	// Update both balances
	_, err = tx.Exec(ctx, `
		UPDATE wallets SET balance_credits = balance_credits - $1, updated_at = NOW() WHERE id = $2
	`, amount, fromWalletID)

	if err != nil {
		return fmt.Errorf("failed to debit from wallet: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE wallets SET balance_credits = balance_credits + $1, updated_at = NOW() WHERE id = $2
	`, amount, toWalletID)

	if err != nil {
		return fmt.Errorf("failed to credit to wallet: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetLockedBalance calculates the locked balance from active orders
func (r *WalletRepository) GetLockedBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(SUM(quantity * limit_price_credits), 0) as locked_balance
		FROM orders
		WHERE user_id = $1 AND status = 'ACTIVE'
	`

	var lockedBalance float64
	err := r.pool.QueryRow(ctx, query, userID).Scan(&lockedBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate locked balance: %w", err)
	}

	return lockedBalance, nil
}

// GetTransactionHistory retrieves transaction history from credit_transactions table
func (r *WalletRepository) GetTransactionHistory(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.CreditTransaction, error) {
	query := `
		SELECT id, user_id, type, crypto_type, crypto_amount, credit_amount, status, created_at, updated_at
		FROM credit_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query transaction history: %w", err)
	}
	defer rows.Close()

	transactions := []models.CreditTransaction{}
	for rows.Next() {
		var t models.CreditTransaction
		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Type,
			&t.CryptoType,
			&t.CryptoAmount,
			&t.CreditAmount,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			continue
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}

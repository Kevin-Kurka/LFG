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
	ErrTransactionNotFound = errors.New("transaction not found")
)

// CreditTransactionRepository handles credit transaction database operations
type CreditTransactionRepository struct {
	pool *pgxpool.Pool
}

// NewCreditTransactionRepository creates a new repository
func NewCreditTransactionRepository(pool *pgxpool.Pool) *CreditTransactionRepository {
	return &CreditTransactionRepository{pool: pool}
}

// Create creates a new credit transaction
func (r *CreditTransactionRepository) Create(ctx context.Context, tx *models.CreditTransaction) error {
	query := `
		INSERT INTO credit_transactions (id, user_id, type, crypto_type, crypto_amount, credit_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`

	_, err := r.pool.Exec(ctx, query,
		tx.ID,
		tx.UserID,
		tx.Type,
		tx.CryptoType,
		tx.CryptoAmount,
		tx.CreditAmount,
		tx.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a transaction by ID
func (r *CreditTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CreditTransaction, error) {
	query := `
		SELECT id, user_id, type, crypto_type, crypto_amount, credit_amount, status, created_at, updated_at
		FROM credit_transactions
		WHERE id = $1
	`

	var tx models.CreditTransaction
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&tx.ID,
		&tx.UserID,
		&tx.Type,
		&tx.CryptoType,
		&tx.CryptoAmount,
		&tx.CreditAmount,
		&tx.Status,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &tx, nil
}

// GetByUserID retrieves transactions for a user
func (r *CreditTransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*models.CreditTransaction, error) {
	query := `
		SELECT id, user_id, type, crypto_type, crypto_amount, credit_amount, status, created_at, updated_at
		FROM credit_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	transactions := []*models.CreditTransaction{}
	for rows.Next() {
		var tx models.CreditTransaction
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.Type,
			&tx.CryptoType,
			&tx.CryptoAmount,
			&tx.CreditAmount,
			&tx.Status,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, &tx)
	}

	return transactions, nil
}

// UpdateStatus updates transaction status
func (r *CreditTransactionRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.CreditTransactionStatus) error {
	query := `
		UPDATE credit_transactions
		SET status = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update transaction status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

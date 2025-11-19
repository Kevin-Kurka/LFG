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
	ErrOrderNotFound = errors.New("order not found")
)

// OrderRepository handles order database operations
type OrderRepository struct {
	pool *pgxpool.Pool
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

// Create creates a new order
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (id, user_id, contract_id, type, status, quantity, quantity_filled, limit_price_credits, stop_price_credits, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`

	_, err := r.pool.Exec(ctx, query,
		order.ID,
		order.UserID,
		order.ContractID,
		order.Type,
		order.Status,
		order.Quantity,
		order.QuantityFilled,
		order.LimitPriceCredits,
		order.StopPriceCredits,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	query := `
		SELECT id, user_id, contract_id, type, status, quantity, quantity_filled, limit_price_credits, stop_price_credits, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var order models.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.ContractID,
		&order.Type,
		&order.Status,
		&order.Quantity,
		&order.QuantityFilled,
		&order.LimitPriceCredits,
		&order.StopPriceCredits,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// GetByUserID retrieves all orders for a user
func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID, status string, limit int) ([]*models.Order, error) {
	query := `
		SELECT id, user_id, contract_id, type, status, quantity, quantity_filled, limit_price_credits, stop_price_credits, created_at, updated_at
		FROM orders
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argCount := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argCount)
	args = append(args, limit)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	orders := []*models.Order{}
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.ContractID,
			&order.Type,
			&order.Status,
			&order.Quantity,
			&order.QuantityFilled,
			&order.LimitPriceCredits,
			&order.StopPriceCredits,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

// UpdateStatus updates order status and filled quantity
func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus, quantityFilled int) error {
	query := `
		UPDATE orders
		SET status = $2, quantity_filled = $3, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query, orderID, status, quantityFilled)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrOrderNotFound
	}

	return nil
}

// Cancel cancels an order
func (r *OrderRepository) Cancel(ctx context.Context, orderID uuid.UUID) error {
	return r.UpdateStatus(ctx, orderID, models.OrderStatusCancelled, 0)
}

// GetPool returns the underlying database connection pool
func (r *OrderRepository) GetPool() *pgxpool.Pool {
	return r.pool
}

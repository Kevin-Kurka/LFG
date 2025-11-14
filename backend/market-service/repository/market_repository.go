package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"lfg/shared/models"
)

var (
	ErrMarketNotFound = errors.New("market not found")
)

// MarketRepository handles market database operations
type MarketRepository struct {
	pool *pgxpool.Pool
}

// NewMarketRepository creates a new market repository
func NewMarketRepository(pool *pgxpool.Pool) *MarketRepository {
	return &MarketRepository{pool: pool}
}

// List retrieves markets with filtering and pagination
func (r *MarketRepository) List(ctx context.Context, status string, search string, page, pageSize int) ([]*models.Market, int, error) {
	// Build query
	query := `
		SELECT id, ticker, question, rules, resolution_source, status, expires_at, resolved_at, outcome, created_at, updated_at
		FROM markets
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(*) FROM markets WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	// Add status filter
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		countQuery += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, status)
		argCount++
	}

	// Add search filter
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query += fmt.Sprintf(" AND (LOWER(ticker) LIKE $%d OR LOWER(question) LIKE $%d)", argCount, argCount)
		countQuery += fmt.Sprintf(" AND (LOWER(ticker) LIKE $%d OR LOWER(question) LIKE $%d)", argCount, argCount)
		args = append(args, searchPattern)
		argCount++
	}

	// Get total count
	var totalCount int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count markets: %w", err)
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)

	// Execute query
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query markets: %w", err)
	}
	defer rows.Close()

	// Parse results
	markets := []*models.Market{}
	for rows.Next() {
		var market models.Market
		err := rows.Scan(
			&market.ID,
			&market.Ticker,
			&market.Question,
			&market.Rules,
			&market.ResolutionSource,
			&market.Status,
			&market.ExpiresAt,
			&market.ResolvedAt,
			&market.Outcome,
			&market.CreatedAt,
			&market.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan market: %w", err)
		}
		markets = append(markets, &market)
	}

	return markets, totalCount, nil
}

// GetByID retrieves a market by ID
func (r *MarketRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Market, error) {
	query := `
		SELECT id, ticker, question, rules, resolution_source, status, expires_at, resolved_at, outcome, created_at, updated_at
		FROM markets
		WHERE id = $1
	`

	var market models.Market
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&market.ID,
		&market.Ticker,
		&market.Question,
		&market.Rules,
		&market.ResolutionSource,
		&market.Status,
		&market.ExpiresAt,
		&market.ResolvedAt,
		&market.Outcome,
		&market.CreatedAt,
		&market.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMarketNotFound
		}
		return nil, fmt.Errorf("failed to get market: %w", err)
	}

	return &market, nil
}

// GetByTicker retrieves a market by ticker
func (r *MarketRepository) GetByTicker(ctx context.Context, ticker string) (*models.Market, error) {
	query := `
		SELECT id, ticker, question, rules, resolution_source, status, expires_at, resolved_at, outcome, created_at, updated_at
		FROM markets
		WHERE ticker = $1
	`

	var market models.Market
	err := r.pool.QueryRow(ctx, query, ticker).Scan(
		&market.ID,
		&market.Ticker,
		&market.Question,
		&market.Rules,
		&market.ResolutionSource,
		&market.Status,
		&market.ExpiresAt,
		&market.ResolvedAt,
		&market.Outcome,
		&market.CreatedAt,
		&market.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMarketNotFound
		}
		return nil, fmt.Errorf("failed to get market: %w", err)
	}

	return &market, nil
}

// GetContractsByMarketID retrieves contracts for a market
func (r *MarketRepository) GetContractsByMarketID(ctx context.Context, marketID uuid.UUID) ([]*models.Contract, error) {
	query := `
		SELECT id, market_id, side, ticker, created_at
		FROM contracts
		WHERE market_id = $1
		ORDER BY side
	`

	rows, err := r.pool.Query(ctx, query, marketID)
	if err != nil {
		return nil, fmt.Errorf("failed to query contracts: %w", err)
	}
	defer rows.Close()

	contracts := []*models.Contract{}
	for rows.Next() {
		var contract models.Contract
		err := rows.Scan(
			&contract.ID,
			&contract.MarketID,
			&contract.Side,
			&contract.Ticker,
			&contract.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contract: %w", err)
		}
		contracts = append(contracts, &contract)
	}

	return contracts, nil
}

// Create creates a new market
func (r *MarketRepository) Create(ctx context.Context, market *models.Market) error {
	query := `
		INSERT INTO markets (id, ticker, question, rules, resolution_source, status, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`

	_, err := r.pool.Exec(ctx, query,
		market.ID,
		market.Ticker,
		market.Question,
		market.Rules,
		market.ResolutionSource,
		market.Status,
		market.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create market: %w", err)
	}

	return nil
}

// CreateContract creates a new contract
func (r *MarketRepository) CreateContract(ctx context.Context, contract *models.Contract) error {
	query := `
		INSERT INTO contracts (id, market_id, side, ticker, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := r.pool.Exec(ctx, query,
		contract.ID,
		contract.MarketID,
		contract.Side,
		contract.Ticker,
	)

	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}

	return nil
}

package testing

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDatabase represents a test database container
type TestDatabase struct {
	Container testcontainers.Container
	Pool      *pgxpool.Pool
	DSN       string
}

// SetupTestDB creates a PostgreSQL test container and returns a connection pool
func SetupTestDB(t *testing.T) *TestDatabase {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Create connection pool
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := runTestMigrations(ctx, pool); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return &TestDatabase{
		Container: container,
		Pool:      pool,
		DSN:       dsn,
	}
}

// Cleanup cleans up the test database
func (db *TestDatabase) Cleanup(t *testing.T) {
	ctx := context.Background()
	if db.Pool != nil {
		db.Pool.Close()
	}
	if db.Container != nil {
		if err := db.Container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}
}

// runTestMigrations runs database migrations for tests
func runTestMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	migrations := []string{
		`CREATE TYPE user_status AS ENUM ('ACTIVE', 'SUSPENDED', 'BANNED')`,
		`CREATE TYPE market_status AS ENUM ('UPCOMING', 'OPEN', 'CLOSED', 'RESOLVED', 'CANCELLED')`,
		`CREATE TYPE order_type AS ENUM ('MARKET', 'LIMIT', 'STOP', 'STOP_LIMIT')`,
		`CREATE TYPE order_status AS ENUM ('PENDING', 'ACTIVE', 'FILLED', 'PARTIALLY_FILLED', 'CANCELLED', 'REJECTED')`,
		`CREATE TYPE contract_side AS ENUM ('YES', 'NO')`,
		`CREATE TYPE market_outcome AS ENUM ('YES', 'NO', 'CANCELLED')`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			wallet_address VARCHAR(255) NULL,
			status user_status NOT NULL DEFAULT 'ACTIVE',
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS wallets (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID UNIQUE NOT NULL,
			balance_credits DECIMAL(18, 8) NOT NULL DEFAULT 0,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			CONSTRAINT balance_non_negative CHECK (balance_credits >= 0),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS markets (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			ticker VARCHAR(50) UNIQUE NOT NULL,
			question TEXT NOT NULL,
			rules TEXT NOT NULL,
			resolution_source VARCHAR(255) NOT NULL,
			status market_status NOT NULL DEFAULT 'UPCOMING',
			expires_at TIMESTAMPTZ NOT NULL,
			resolved_at TIMESTAMPTZ NULL,
			outcome market_outcome NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS contracts (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			market_id UUID NOT NULL,
			side contract_side NOT NULL,
			ticker VARCHAR(60) UNIQUE NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			FOREIGN KEY (market_id) REFERENCES markets(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL,
			contract_id UUID NOT NULL,
			type order_type NOT NULL,
			status order_status NOT NULL DEFAULT 'PENDING',
			quantity INTEGER NOT NULL,
			quantity_filled INTEGER NOT NULL DEFAULT 0,
			limit_price_credits DECIMAL(10, 8) NULL,
			stop_price_credits DECIMAL(10, 8) NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (contract_id) REFERENCES contracts(id) ON DELETE CASCADE
		)`,
	}

	for _, migration := range migrations {
		if _, err := pool.Exec(ctx, migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// TruncateTables truncates all tables for clean test state
func (db *TestDatabase) TruncateTables(t *testing.T) {
	ctx := context.Background()
	tables := []string{"orders", "contracts", "markets", "wallets", "users"}

	for _, table := range tables {
		_, err := db.Pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}

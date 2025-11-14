-- Enhanced schema with indexes, constraints, and optimizations
-- Migration: 001_enhanced_schema

-- Create enum types for better data integrity
CREATE TYPE user_status AS ENUM ('ACTIVE', 'SUSPENDED', 'BANNED');
CREATE TYPE market_status AS ENUM ('UPCOMING', 'OPEN', 'CLOSED', 'RESOLVED', 'CANCELLED');
CREATE TYPE order_type AS ENUM ('MARKET', 'LIMIT', 'STOP', 'STOP_LIMIT');
CREATE TYPE order_status AS ENUM ('PENDING', 'ACTIVE', 'FILLED', 'PARTIALLY_FILLED', 'CANCELLED', 'REJECTED');
CREATE TYPE contract_side AS ENUM ('YES', 'NO');
CREATE TYPE market_outcome AS ENUM ('YES', 'NO', 'CANCELLED');

-- Users Table with enhanced constraints
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    wallet_address VARCHAR(255) NULL,
    status user_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT password_hash_length CHECK (length(password_hash) >= 60)
);

-- Index for email lookups (login)
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- Wallets Table with balance constraints
CREATE TABLE IF NOT EXISTS wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL,
    balance_credits DECIMAL(18, 8) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT balance_non_negative CHECK (balance_credits >= 0),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Index for user_id lookups
CREATE INDEX idx_wallets_user_id ON wallets(user_id);

-- Markets Table with enhanced constraints
CREATE TABLE IF NOT EXISTS markets (
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
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT ticker_format CHECK (ticker ~* '^[A-Z0-9_-]+$'),
    CONSTRAINT expires_future CHECK (expires_at > created_at),
    CONSTRAINT resolution_logic CHECK (
        (status = 'RESOLVED' AND resolved_at IS NOT NULL AND outcome IS NOT NULL) OR
        (status != 'RESOLVED' AND resolved_at IS NULL)
    )
);

-- Indexes for market queries
CREATE INDEX idx_markets_ticker ON markets(ticker);
CREATE INDEX idx_markets_status ON markets(status);
CREATE INDEX idx_markets_expires_at ON markets(expires_at);
CREATE INDEX idx_markets_created_at ON markets(created_at DESC);

-- Contracts Table
CREATE TABLE IF NOT EXISTS contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    market_id UUID NOT NULL,
    side contract_side NOT NULL,
    ticker VARCHAR(60) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    FOREIGN KEY (market_id) REFERENCES markets(id) ON DELETE CASCADE,
    CONSTRAINT unique_market_side UNIQUE(market_id, side)
);

-- Indexes for contract lookups
CREATE INDEX idx_contracts_market_id ON contracts(market_id);
CREATE INDEX idx_contracts_ticker ON contracts(ticker);

-- Orders Table with comprehensive constraints
CREATE TABLE IF NOT EXISTS orders (
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

    CONSTRAINT quantity_positive CHECK (quantity > 0),
    CONSTRAINT quantity_filled_valid CHECK (quantity_filled >= 0 AND quantity_filled <= quantity),
    CONSTRAINT limit_price_valid CHECK (
        (type IN ('LIMIT', 'STOP_LIMIT') AND limit_price_credits > 0 AND limit_price_credits <= 1) OR
        (type NOT IN ('LIMIT', 'STOP_LIMIT') AND limit_price_credits IS NULL)
    ),
    CONSTRAINT stop_price_valid CHECK (
        (type IN ('STOP', 'STOP_LIMIT') AND stop_price_credits > 0 AND stop_price_credits <= 1) OR
        (type NOT IN ('STOP', 'STOP_LIMIT') AND stop_price_credits IS NULL)
    ),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (contract_id) REFERENCES contracts(id) ON DELETE CASCADE
);

-- Indexes for order queries
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_contract_id ON orders(contract_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX idx_orders_type ON orders(type);
-- Composite index for active orders by contract
CREATE INDEX idx_orders_active_contract ON orders(contract_id, status) WHERE status IN ('ACTIVE', 'PARTIALLY_FILLED');

-- Trades Table
CREATE TABLE IF NOT EXISTS trades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contract_id UUID NOT NULL,
    maker_order_id UUID NOT NULL,
    taker_order_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    price_credits DECIMAL(10, 8) NOT NULL,
    executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT trade_quantity_positive CHECK (quantity > 0),
    CONSTRAINT trade_price_valid CHECK (price_credits > 0 AND price_credits <= 1),
    FOREIGN KEY (contract_id) REFERENCES contracts(id) ON DELETE CASCADE,
    FOREIGN KEY (maker_order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (taker_order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- Indexes for trade queries
CREATE INDEX idx_trades_contract_id ON trades(contract_id);
CREATE INDEX idx_trades_maker_order_id ON trades(maker_order_id);
CREATE INDEX idx_trades_taker_order_id ON trades(taker_order_id);
CREATE INDEX idx_trades_executed_at ON trades(executed_at DESC);

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for automatic updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_wallets_updated_at BEFORE UPDATE ON wallets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_markets_updated_at BEFORE UPDATE ON markets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Credit exchange transactions table
CREATE TABLE IF NOT EXISTS credit_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL,
    crypto_type VARCHAR(10) NOT NULL,
    crypto_amount DECIMAL(18, 8) NOT NULL,
    credit_amount DECIMAL(18, 8) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT credit_amount_positive CHECK (credit_amount > 0),
    CONSTRAINT crypto_amount_positive CHECK (crypto_amount > 0),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_credit_transactions_user_id ON credit_transactions(user_id);
CREATE INDEX idx_credit_transactions_status ON credit_transactions(status);
CREATE INDEX idx_credit_transactions_created_at ON credit_transactions(created_at DESC);

CREATE TRIGGER update_credit_transactions_updated_at BEFORE UPDATE ON credit_transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

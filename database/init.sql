-- Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    wallet_address VARCHAR(255) NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- Wallets Table
CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    user_id UUID UNIQUE NOT NULL,
    balance_credits DECIMAL(18, 8) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Markets Table
CREATE TABLE markets (
    id UUID PRIMARY KEY,
    ticker VARCHAR(50) UNIQUE NOT NULL,
    question TEXT NOT NULL,
    rules TEXT NOT NULL,
    resolution_source VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'UPCOMING',
    expires_at TIMESTAMPTZ NOT NULL,
    resolved_at TIMESTAMPTZ NULL,
    outcome VARCHAR(10) NULL
);

-- Contracts Table
CREATE TABLE contracts (
    id UUID PRIMARY KEY,
    market_id UUID NOT NULL,
    side VARCHAR(10) NOT NULL,
    ticker VARCHAR(60) UNIQUE NOT NULL,
    FOREIGN KEY (market_id) REFERENCES markets(id)
);

-- Orders Table
CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    contract_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    quantity INTEGER NOT NULL,
    quantity_filled INTEGER NOT NULL DEFAULT 0,
    limit_price_credits DECIMAL(10, 8) NULL,
    stop_price_credits DECIMAL(10, 8) NULL,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (contract_id) REFERENCES contracts(id)
);

-- Trades Table
CREATE TABLE trades (
    id UUID PRIMARY KEY,
    contract_id UUID NOT NULL,
    maker_order_id UUID NOT NULL,
    taker_order_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    price_credits DECIMAL(10, 8) NOT NULL,
    executed_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (contract_id) REFERENCES contracts(id),
    FOREIGN KEY (maker_order_id) REFERENCES orders(id),
    FOREIGN KEY (taker_order_id) REFERENCES orders(id)
);

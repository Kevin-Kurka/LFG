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
    FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT positive_balance CHECK (balance_credits >= 0)
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
    outcome VARCHAR(10) NULL,
    CONSTRAINT valid_market_status CHECK (status IN ('UPCOMING', 'OPEN', 'CLOSED', 'RESOLVED', 'CANCELLED')),
    CONSTRAINT valid_market_outcome CHECK (outcome IS NULL OR outcome IN ('YES', 'NO'))
);

-- Contracts Table
CREATE TABLE contracts (
    id UUID PRIMARY KEY,
    market_id UUID NOT NULL,
    side VARCHAR(10) NOT NULL,
    ticker VARCHAR(60) UNIQUE NOT NULL,
    FOREIGN KEY (market_id) REFERENCES markets(id),
    CONSTRAINT valid_contract_side CHECK (side IN ('YES', 'NO'))
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
    FOREIGN KEY (contract_id) REFERENCES contracts(id),
    CONSTRAINT valid_order_type CHECK (type IN ('LIMIT', 'MARKET', 'STOP')),
    CONSTRAINT valid_order_status CHECK (status IN ('PENDING', 'OPEN', 'PARTIALLY_FILLED', 'FILLED', 'CANCELLED')),
    CONSTRAINT positive_quantity CHECK (quantity > 0),
    CONSTRAINT valid_filled_quantity CHECK (quantity_filled >= 0 AND quantity_filled <= quantity),
    CONSTRAINT valid_limit_price CHECK (limit_price_credits IS NULL OR (limit_price_credits >= 0 AND limit_price_credits <= 1))
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
    FOREIGN KEY (taker_order_id) REFERENCES orders(id),
    CONSTRAINT positive_trade_quantity CHECK (quantity > 0),
    CONSTRAINT valid_trade_price CHECK (price_credits >= 0 AND price_credits <= 1)
);

-- Sportsbook Providers Table
CREATE TABLE sportsbook_providers (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    logo_url VARCHAR(500) NULL,
    website_url VARCHAR(500) NOT NULL,
    supports_api BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL
);

-- User Sportsbook Accounts Table (encrypted credentials)
CREATE TABLE user_sportsbook_accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    sportsbook_id UUID NOT NULL,
    username_encrypted TEXT NOT NULL,
    credentials_encrypted TEXT NOT NULL,
    is_connected BOOLEAN DEFAULT false,
    last_sync_at TIMESTAMPTZ NULL,
    balance DECIMAL(18, 2) NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sportsbook_id) REFERENCES sportsbook_providers(id),
    UNIQUE(user_id, sportsbook_id)
);

-- Sports Table
CREATE TABLE sports (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true
);

-- Sports Events Table
CREATE TABLE sports_events (
    id UUID PRIMARY KEY,
    sport_id UUID NOT NULL,
    home_team VARCHAR(255) NOT NULL,
    away_team VARCHAR(255) NOT NULL,
    event_time TIMESTAMPTZ NOT NULL,
    league VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'UPCOMING',
    home_score INTEGER NULL,
    away_score INTEGER NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (sport_id) REFERENCES sports(id)
);

-- Odds Table (stores odds from different sportsbooks)
CREATE TABLE odds (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL,
    sportsbook_id UUID NOT NULL,
    market_type VARCHAR(50) NOT NULL,
    bet_type VARCHAR(100) NOT NULL,
    selection VARCHAR(255) NOT NULL,
    odds_american INTEGER NULL,
    odds_decimal DECIMAL(10, 3) NULL,
    odds_fractional VARCHAR(20) NULL,
    line DECIMAL(10, 2) NULL,
    is_active BOOLEAN DEFAULT true,
    last_updated TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (event_id) REFERENCES sports_events(id) ON DELETE CASCADE,
    FOREIGN KEY (sportsbook_id) REFERENCES sportsbook_providers(id)
);

-- Arbitrage Opportunities Table
CREATE TABLE arbitrage_opportunities (
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL,
    market_type VARCHAR(50) NOT NULL,
    bet_description TEXT NOT NULL,
    profit_percentage DECIMAL(10, 4) NOT NULL,
    total_stake_required DECIMAL(18, 2) NOT NULL,
    estimated_profit DECIMAL(18, 2) NOT NULL,
    odds_data JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    detected_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (event_id) REFERENCES sports_events(id) ON DELETE CASCADE
);

-- Hedge Opportunities Table
CREATE TABLE hedge_opportunities (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    original_bet_id UUID NULL,
    event_id UUID NOT NULL,
    hedge_type VARCHAR(50) NOT NULL,
    original_stake DECIMAL(18, 2) NOT NULL,
    hedge_stake_required DECIMAL(18, 2) NOT NULL,
    guaranteed_profit DECIMAL(18, 2) NOT NULL,
    hedge_odds_data JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    detected_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES sports_events(id) ON DELETE CASCADE
);

-- User Bets Table (track bets across all sportsbooks)
CREATE TABLE user_bets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    sportsbook_account_id UUID NOT NULL,
    event_id UUID NULL,
    bet_type VARCHAR(100) NOT NULL,
    selection VARCHAR(255) NOT NULL,
    stake DECIMAL(18, 2) NOT NULL,
    odds_american INTEGER NULL,
    odds_decimal DECIMAL(10, 3) NULL,
    potential_payout DECIMAL(18, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'PENDING',
    result VARCHAR(50) NULL,
    profit_loss DECIMAL(18, 2) NULL,
    placed_at TIMESTAMPTZ NOT NULL,
    settled_at TIMESTAMPTZ NULL,
    notes TEXT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sportsbook_account_id) REFERENCES user_sportsbook_accounts(id),
    FOREIGN KEY (event_id) REFERENCES sports_events(id) ON DELETE SET NULL,
    CONSTRAINT positive_stake CHECK (stake > 0),
    CONSTRAINT positive_payout CHECK (potential_payout > 0),
    CONSTRAINT valid_bet_status CHECK (status IN ('PENDING', 'ACTIVE', 'WON', 'LOST', 'VOID', 'CANCELLED'))
);

-- Wallet Transactions Table
CREATE TABLE IF NOT EXISTS wallet_transactions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    amount DECIMAL(18, 8) NOT NULL,
    balance DECIMAL(18, 8) NOT NULL,
    reference VARCHAR(255) NULL,
    reference_id UUID NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for performance
-- User and wallet indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_wallets_user ON wallets(user_id);

-- Market and contract indexes
CREATE INDEX idx_markets_status ON markets(status, expires_at);
CREATE INDEX idx_markets_ticker ON markets(ticker);
CREATE INDEX idx_contracts_market ON contracts(market_id);

-- Order and trade indexes
CREATE INDEX idx_orders_user ON orders(user_id, status, created_at DESC);
CREATE INDEX idx_orders_contract ON orders(contract_id, status);
CREATE INDEX idx_orders_status ON orders(status) WHERE status IN ('PENDING', 'OPEN');
CREATE INDEX idx_trades_contract ON trades(contract_id, executed_at DESC);
CREATE INDEX idx_trades_orders ON trades(maker_order_id, taker_order_id);

-- Wallet transaction indexes
CREATE INDEX idx_wallet_transactions_user ON wallet_transactions(user_id, created_at DESC);
CREATE INDEX idx_wallet_transactions_type ON wallet_transactions(type);
CREATE INDEX idx_wallet_transactions_reference ON wallet_transactions(reference_id) WHERE reference_id IS NOT NULL;

-- Sportsbook indexes
CREATE INDEX idx_odds_event_sportsbook ON odds(event_id, sportsbook_id);
CREATE INDEX idx_odds_market_type ON odds(market_type, is_active);
CREATE INDEX idx_events_time ON sports_events(event_time);
CREATE INDEX idx_events_sport ON sports_events(sport_id, status);
CREATE INDEX idx_user_accounts_user ON user_sportsbook_accounts(user_id);
CREATE INDEX idx_user_bets_user ON user_bets(user_id, status);
CREATE INDEX idx_arbitrage_active ON arbitrage_opportunities(is_active, expires_at);
CREATE INDEX idx_hedge_active ON hedge_opportunities(is_active, user_id);

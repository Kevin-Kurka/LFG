-- Performance Optimization Indexes for LFG Platform
-- Run these after initial migrations for optimal query performance

-- Orders table optimizations
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_user_status
    ON orders(user_id, status)
    WHERE status IN ('ACTIVE', 'PARTIALLY_FILLED');

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_contract_status
    ON orders(contract_id, status)
    WHERE status IN ('ACTIVE', 'PARTIALLY_FILLED');

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_created_at_desc
    ON orders(created_at DESC);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_orders_composite
    ON orders(contract_id, status, limit_price_credits, created_at);

-- Trades table optimizations (if exists)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_trades_market_created
    ON trades(market_id, created_at DESC)
    WHERE created_at > NOW() - INTERVAL '30 days';

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_trades_user_created
    ON trades(user_id, created_at DESC)
    WHERE created_at > NOW() - INTERVAL '90 days';

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_trades_contract
    ON trades(contract_id);

-- Wallets table optimizations
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_wallets_balance
    ON wallets(balance_credits)
    WHERE balance_credits > 0;

-- Markets table optimizations
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_markets_status_expires
    ON markets(status, expires_at)
    WHERE status = 'OPEN';

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_markets_ticker_lower
    ON markets(LOWER(ticker));

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_markets_question_trgm
    ON markets USING gin(to_tsvector('english', question));

-- Contracts table optimizations
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_contracts_market_side
    ON contracts(market_id, side);

-- Users table optimizations
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email_lower
    ON users(LOWER(email));

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_status_created
    ON users(status, created_at DESC);

-- Credit transactions table optimizations (if exists)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_credit_txn_user_created
    ON credit_transactions(user_id, created_at DESC);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_credit_txn_status
    ON credit_transactions(status);

-- Analyze tables for query planner
ANALYZE users;
ANALYZE wallets;
ANALYZE markets;
ANALYZE contracts;
ANALYZE orders;

-- Vacuum tables
VACUUM ANALYZE users;
VACUUM ANALYZE wallets;
VACUUM ANALYZE markets;
VACUUM ANALYZE orders;

-- Create statistics for better query planning
CREATE STATISTICS IF NOT EXISTS orders_price_qty_stats
    ON limit_price_credits, quantity
    FROM orders;

CREATE STATISTICS IF NOT EXISTS markets_status_time_stats
    ON status, created_at, expires_at
    FROM markets;

-- Set autovacuum settings for high-traffic tables
ALTER TABLE orders SET (
    autovacuum_vacuum_scale_factor = 0.05,
    autovacuum_analyze_scale_factor = 0.02
);

ALTER TABLE trades SET (
    autovacuum_vacuum_scale_factor = 0.05,
    autovacuum_analyze_scale_factor = 0.02
);

-- Performance monitoring views
CREATE OR REPLACE VIEW slow_queries AS
SELECT
    query,
    calls,
    total_time,
    mean_time,
    max_time,
    stddev_time
FROM pg_stat_statements
WHERE mean_time > 100  -- Queries averaging > 100ms
ORDER BY mean_time DESC
LIMIT 50;

-- Index usage statistics view
CREATE OR REPLACE VIEW index_usage AS
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans,
    idx_tup_read as tuples_read,
    idx_tup_fetch as tuples_fetched
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- Table bloat monitoring
CREATE OR REPLACE VIEW table_bloat AS
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as total_size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) as table_size,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) -
                   pg_relation_size(schemaname||'.'||tablename)) as index_size,
    n_dead_tup as dead_tuples,
    n_live_tup as live_tuples,
    ROUND(100.0 * n_dead_tup / NULLIF(n_live_tup + n_dead_tup, 0), 2) as dead_tuple_percent
FROM pg_stat_user_tables
WHERE n_live_tup > 0
ORDER BY n_dead_tup DESC;

-- Success message
DO $$ BEGIN
    RAISE NOTICE 'Performance optimizations applied successfully';
    RAISE NOTICE 'Created % indexes', (SELECT count(*) FROM pg_indexes WHERE schemaname = 'public');
    RAISE NOTICE 'Run EXPLAIN ANALYZE on your queries to verify index usage';
END $$;

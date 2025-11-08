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

CREATE INDEX idx_wallet_transactions_user ON wallet_transactions(user_id, created_at DESC);
CREATE INDEX idx_wallet_transactions_type ON wallet_transactions(type);

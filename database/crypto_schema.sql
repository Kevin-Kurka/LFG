-- Cryptocurrency System Schema for LFG Platform
-- Supports BTC, ETH, USDC, USDT, LTC deposits and withdrawals

-- Crypto Wallets Table
-- Stores user crypto wallet addresses and encrypted private keys
CREATE TABLE IF NOT EXISTS crypto_wallets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    currency VARCHAR(10) NOT NULL, -- BTC, ETH, USDC, USDT, LTC
    address VARCHAR(255) NOT NULL UNIQUE,
    private_key_encrypted TEXT NOT NULL, -- Encrypted with AES-256
    balance_crypto DECIMAL(28, 18) NOT NULL DEFAULT 0, -- Track crypto balance
    derivation_path VARCHAR(100) NULL, -- HD wallet derivation path (m/44'/0'/0'/0/0)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, currency)
);

CREATE INDEX idx_crypto_wallets_user ON crypto_wallets(user_id);
CREATE INDEX idx_crypto_wallets_currency ON crypto_wallets(currency);
CREATE INDEX idx_crypto_wallets_address ON crypto_wallets(address);

-- Crypto Deposits Table
-- Tracks all incoming cryptocurrency deposits
CREATE TABLE IF NOT EXISTS crypto_deposits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    wallet_id UUID NOT NULL,
    currency VARCHAR(10) NOT NULL,
    amount_crypto DECIMAL(28, 18) NOT NULL, -- Amount in crypto
    amount_credits DECIMAL(18, 8) NOT NULL, -- Amount credited in platform credits
    exchange_rate DECIMAL(18, 8) NOT NULL, -- Rate used for conversion
    address VARCHAR(255) NOT NULL, -- Deposit address
    tx_hash VARCHAR(255) NOT NULL UNIQUE, -- Blockchain transaction hash
    confirmations INTEGER NOT NULL DEFAULT 0,
    required_confirmations INTEGER NOT NULL, -- BTC: 3, ETH/USDC/USDT: 12, LTC: 6
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, CONFIRMING, CONFIRMED, CREDITED, FAILED
    block_number BIGINT NULL,
    detected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMPTZ NULL,
    credited_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (wallet_id) REFERENCES crypto_wallets(id) ON DELETE CASCADE
);

CREATE INDEX idx_crypto_deposits_user ON crypto_deposits(user_id, created_at DESC);
CREATE INDEX idx_crypto_deposits_status ON crypto_deposits(status, currency);
CREATE INDEX idx_crypto_deposits_tx_hash ON crypto_deposits(tx_hash);
CREATE INDEX idx_crypto_deposits_wallet ON crypto_deposits(wallet_id);

-- Crypto Withdrawals Table
-- Tracks all outgoing cryptocurrency withdrawals
CREATE TABLE IF NOT EXISTS crypto_withdrawals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    wallet_id UUID NOT NULL,
    currency VARCHAR(10) NOT NULL,
    amount_credits DECIMAL(18, 8) NOT NULL, -- Amount debited in platform credits
    amount_crypto DECIMAL(28, 18) NOT NULL, -- Amount sent in crypto
    exchange_rate DECIMAL(18, 8) NOT NULL, -- Rate used for conversion
    network_fee DECIMAL(28, 18) NOT NULL DEFAULT 0, -- Network transaction fee
    platform_fee DECIMAL(18, 8) NOT NULL DEFAULT 0, -- Platform fee (1%)
    total_fee_credits DECIMAL(18, 8) NOT NULL, -- Total fee in credits
    to_address VARCHAR(255) NOT NULL, -- Recipient address
    tx_hash VARCHAR(255) NULL UNIQUE, -- Blockchain transaction hash (null until broadcast)
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- PENDING, PROCESSING, SENT, CONFIRMED, FAILED, CANCELLED
    confirmations INTEGER NOT NULL DEFAULT 0,
    block_number BIGINT NULL,
    requires_2fa BOOLEAN NOT NULL DEFAULT true,
    verified_2fa BOOLEAN NOT NULL DEFAULT false,
    error_message TEXT NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ NULL,
    confirmed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (wallet_id) REFERENCES crypto_wallets(id) ON DELETE CASCADE
);

CREATE INDEX idx_crypto_withdrawals_user ON crypto_withdrawals(user_id, created_at DESC);
CREATE INDEX idx_crypto_withdrawals_status ON crypto_withdrawals(status, currency);
CREATE INDEX idx_crypto_withdrawals_tx_hash ON crypto_withdrawals(tx_hash);
CREATE INDEX idx_crypto_withdrawals_wallet ON crypto_withdrawals(wallet_id);

-- Exchange Rates Table
-- Stores current exchange rates for all supported cryptocurrencies
CREATE TABLE IF NOT EXISTS exchange_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    currency VARCHAR(10) NOT NULL UNIQUE,
    usd_rate DECIMAL(18, 8) NOT NULL, -- Current rate in USD
    provider VARCHAR(50) NOT NULL, -- CoinGecko, CryptoCompare, etc.
    last_updated TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_exchange_rates_currency ON exchange_rates(currency);
CREATE INDEX idx_exchange_rates_updated ON exchange_rates(last_updated DESC);

-- Insert initial exchange rate records
INSERT INTO exchange_rates (currency, usd_rate, provider)
VALUES
    ('BTC', 0, 'CoinGecko'),
    ('ETH', 0, 'CoinGecko'),
    ('USDC', 1.00, 'CoinGecko'),
    ('USDT', 1.00, 'CoinGecko'),
    ('LTC', 0, 'CoinGecko')
ON CONFLICT (currency) DO NOTHING;

-- Crypto Transactions Table
-- Comprehensive audit trail for all crypto-related transactions
CREATE TABLE IF NOT EXISTS crypto_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    wallet_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL, -- DEPOSIT, WITHDRAWAL, FEE
    currency VARCHAR(10) NOT NULL,
    amount_crypto DECIMAL(28, 18) NOT NULL,
    amount_credits DECIMAL(18, 8) NULL,
    exchange_rate DECIMAL(18, 8) NULL,
    tx_hash VARCHAR(255) NULL,
    status VARCHAR(20) NOT NULL,
    reference_id UUID NULL, -- Links to deposit or withdrawal ID
    reference_type VARCHAR(20) NULL, -- DEPOSIT, WITHDRAWAL
    metadata JSONB NULL, -- Additional transaction data
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (wallet_id) REFERENCES crypto_wallets(id) ON DELETE CASCADE
);

CREATE INDEX idx_crypto_transactions_user ON crypto_transactions(user_id, created_at DESC);
CREATE INDEX idx_crypto_transactions_wallet ON crypto_transactions(wallet_id);
CREATE INDEX idx_crypto_transactions_type ON crypto_transactions(type, status);
CREATE INDEX idx_crypto_transactions_tx_hash ON crypto_transactions(tx_hash);

-- Crypto Wallet Monitoring Jobs Table
-- Tracks monitoring jobs for deposit detection
CREATE TABLE IF NOT EXISTS crypto_monitoring_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    currency VARCHAR(10) NOT NULL,
    last_scanned_block BIGINT NOT NULL DEFAULT 0,
    last_scan_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_running BOOLEAN NOT NULL DEFAULT false,
    error_count INTEGER NOT NULL DEFAULT 0,
    last_error TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(currency)
);

-- Insert initial monitoring jobs
INSERT INTO crypto_monitoring_jobs (currency, last_scanned_block)
VALUES
    ('BTC', 0),
    ('ETH', 0),
    ('USDC', 0),
    ('USDT', 0),
    ('LTC', 0)
ON CONFLICT (currency) DO NOTHING;

-- Crypto Configuration Table
-- Stores system-wide crypto configuration
CREATE TABLE IF NOT EXISTS crypto_config (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    description TEXT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default configuration
INSERT INTO crypto_config (key, value, description) VALUES
    ('btc_min_deposit', '0.001', 'Minimum BTC deposit amount'),
    ('eth_min_deposit', '0.01', 'Minimum ETH deposit amount'),
    ('usdc_min_deposit', '10', 'Minimum USDC deposit amount'),
    ('usdt_min_deposit', '10', 'Minimum USDT deposit amount'),
    ('ltc_min_deposit', '0.1', 'Minimum LTC deposit amount'),
    ('btc_min_withdrawal', '0.001', 'Minimum BTC withdrawal amount'),
    ('eth_min_withdrawal', '0.01', 'Minimum ETH withdrawal amount'),
    ('usdc_min_withdrawal', '10', 'Minimum USDC withdrawal amount'),
    ('usdt_min_withdrawal', '10', 'Minimum USDT withdrawal amount'),
    ('ltc_min_withdrawal', '0.1', 'Minimum LTC withdrawal amount'),
    ('platform_fee_percentage', '1.0', 'Platform fee percentage for withdrawals'),
    ('btc_confirmations_required', '3', 'Required confirmations for BTC'),
    ('eth_confirmations_required', '12', 'Required confirmations for ETH'),
    ('usdc_confirmations_required', '12', 'Required confirmations for USDC'),
    ('usdt_confirmations_required', '12', 'Required confirmations for USDT'),
    ('ltc_confirmations_required', '6', 'Required confirmations for LTC'),
    ('deposit_scan_interval_seconds', '30', 'How often to scan for new deposits'),
    ('withdrawal_processing_enabled', 'true', 'Enable/disable withdrawal processing'),
    ('btc_network_fee_satoshis', '10000', 'Estimated BTC network fee in satoshis'),
    ('eth_gas_price_gwei', '20', 'Estimated ETH gas price in gwei'),
    ('ltc_network_fee', '0.001', 'Estimated LTC network fee')
ON CONFLICT (key) DO NOTHING;

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_crypto_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers for updated_at
CREATE TRIGGER trigger_crypto_wallets_updated_at
    BEFORE UPDATE ON crypto_wallets
    FOR EACH ROW
    EXECUTE FUNCTION update_crypto_updated_at();

CREATE TRIGGER trigger_crypto_monitoring_jobs_updated_at
    BEFORE UPDATE ON crypto_monitoring_jobs
    FOR EACH ROW
    EXECUTE FUNCTION update_crypto_updated_at();

CREATE TRIGGER trigger_crypto_config_updated_at
    BEFORE UPDATE ON crypto_config
    FOR EACH ROW
    EXECUTE FUNCTION update_crypto_updated_at();

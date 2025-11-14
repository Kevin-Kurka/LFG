-- Seed data for LFG Platform
-- This file populates the database with test data for development

-- Test users (passwords are all "password123")
-- Password hash generated with bcrypt cost 12
INSERT INTO users (id, email, password_hash, wallet_address, status, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'alice@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzpLpVF3jO', '0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1', 'ACTIVE', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440001', 'bob@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzpLpVF3jO', '0x5A0b54D5dc17e0AadC383d2db43B0a0D3E029c4c', 'ACTIVE', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002', 'charlie@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzpLpVF3jO', NULL, 'ACTIVE', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003', 'diana@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzpLpVF3jO', NULL, 'ACTIVE', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440004', 'eve@example.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzpLpVF3jO', NULL, 'SUSPENDED', NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440005', 'admin@lfg.com', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYzpLpVF3jO', NULL, 'ACTIVE', NOW(), NOW());

-- Wallets for test users
INSERT INTO wallets (id, user_id, balance_credits, created_at, updated_at) VALUES
('660e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 1000.00000000, NOW(), NOW()),
('660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 500.00000000, NOW(), NOW()),
('660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 250.00000000, NOW(), NOW()),
('660e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440003', 100.00000000, NOW(), NOW()),
('660e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440004', 0.00000000, NOW(), NOW()),
('660e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440005', 10000.00000000, NOW(), NOW());

-- Test markets
INSERT INTO markets (id, ticker, question, rules, resolution_source, status, expires_at, created_at, updated_at) VALUES
('770e8400-e29b-41d4-a716-446655440000', 'BTC100K2025', 'Will Bitcoin reach $100,000 by end of 2025?', 'Market resolves YES if BTC/USD reaches or exceeds $100,000 on any major exchange by Dec 31, 2025 23:59:59 UTC. Otherwise resolves NO.', 'CoinMarketCap', 'OPEN', '2025-12-31 23:59:59+00', NOW(), NOW()),
('770e8400-e29b-41d4-a716-446655440001', 'ETH5K2025', 'Will Ethereum reach $5,000 by end of 2025?', 'Market resolves YES if ETH/USD reaches or exceeds $5,000 on any major exchange by Dec 31, 2025 23:59:59 UTC. Otherwise resolves NO.', 'CoinMarketCap', 'OPEN', '2025-12-31 23:59:59+00', NOW(), NOW()),
('770e8400-e29b-41d4-a716-446655440002', 'SUPERBOWL2026', 'Will the Kansas City Chiefs win Super Bowl LX?', 'Market resolves YES if Chiefs win Super Bowl LX. Resolves NO if any other team wins. Resolves CANCELLED if Super Bowl is not played.', 'NFL Official', 'UPCOMING', '2026-02-08 23:59:59+00', NOW(), NOW()),
('770e8400-e29b-41d4-a716-446655440003', 'AI_AGI_2026', 'Will AGI be achieved by end of 2026?', 'Market resolves YES if a credible AI lab (OpenAI, DeepMind, Anthropic, Meta) publicly announces AGI achievement by Dec 31, 2026. Expert panel determines credibility.', 'AI Research Community', 'OPEN', '2026-12-31 23:59:59+00', NOW(), NOW()),
('770e8400-e29b-41d4-a716-446655440004', 'MARS2030', 'Will humans land on Mars by 2030?', 'Market resolves YES if humans successfully land on Mars and survive for at least 24 hours by Dec 31, 2030 23:59:59 UTC.', 'NASA Official', 'OPEN', '2030-12-31 23:59:59+00', NOW(), NOW());

-- Contracts for markets (YES/NO for each market)
INSERT INTO contracts (id, market_id, side, ticker, created_at) VALUES
-- BTC100K2025
('880e8400-e29b-41d4-a716-446655440000', '770e8400-e29b-41d4-a716-446655440000', 'YES', 'BTC100K2025-YES', NOW()),
('880e8400-e29b-41d4-a716-446655440001', '770e8400-e29b-41d4-a716-446655440000', 'NO', 'BTC100K2025-NO', NOW()),
-- ETH5K2025
('880e8400-e29b-41d4-a716-446655440002', '770e8400-e29b-41d4-a716-446655440001', 'YES', 'ETH5K2025-YES', NOW()),
('880e8400-e29b-41d4-a716-446655440003', '770e8400-e29b-41d4-a716-446655440001', 'NO', 'ETH5K2025-NO', NOW()),
-- SUPERBOWL2026
('880e8400-e29b-41d4-a716-446655440004', '770e8400-e29b-41d4-a716-446655440002', 'YES', 'SUPERBOWL2026-YES', NOW()),
('880e8400-e29b-41d4-a716-446655440005', '770e8400-e29b-41d4-a716-446655440002', 'NO', 'SUPERBOWL2026-NO', NOW()),
-- AI_AGI_2026
('880e8400-e29b-41d4-a716-446655440006', '770e8400-e29b-41d4-a716-446655440003', 'YES', 'AI_AGI_2026-YES', NOW()),
('880e8400-e29b-41d4-a716-446655440007', '770e8400-e29b-41d4-a716-446655440003', 'NO', 'AI_AGI_2026-NO', NOW()),
-- MARS2030
('880e8400-e29b-41d4-a716-446655440008', '770e8400-e29b-41d4-a716-446655440004', 'YES', 'MARS2030-YES', NOW()),
('880e8400-e29b-41d4-a716-446655440009', '770e8400-e29b-41d4-a716-446655440004', 'NO', 'MARS2030-NO', NOW());

-- Sample orders
INSERT INTO orders (id, user_id, contract_id, type, status, quantity, quantity_filled, limit_price_credits, created_at, updated_at) VALUES
-- Alice's orders
('990e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', '880e8400-e29b-41d4-a716-446655440000', 'LIMIT', 'ACTIVE', 100, 0, 0.65, NOW(), NOW()),
('990e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440000', '880e8400-e29b-41d4-a716-446655440002', 'LIMIT', 'ACTIVE', 50, 0, 0.55, NOW(), NOW()),
-- Bob's orders
('990e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440001', 'LIMIT', 'ACTIVE', 75, 0, 0.40, NOW(), NOW()),
('990e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440003', 'LIMIT', 'ACTIVE', 100, 0, 0.50, NOW(), NOW()),
-- Charlie's orders
('990e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440006', 'LIMIT', 'ACTIVE', 200, 0, 0.20, NOW(), NOW());

-- Sample completed trades
INSERT INTO trades (id, contract_id, maker_order_id, taker_order_id, quantity, price_credits, executed_at) VALUES
('aa0e8400-e29b-41d4-a716-446655440000', '880e8400-e29b-41d4-a716-446655440000', '990e8400-e29b-41d4-a716-446655440000', '990e8400-e29b-41d4-a716-446655440002', 25, 0.60, NOW() - INTERVAL '1 hour'),
('aa0e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440002', '990e8400-e29b-41d4-a716-446655440001', '990e8400-e29b-41d4-a716-446655440003', 10, 0.52, NOW() - INTERVAL '30 minutes');

-- Credit transactions
INSERT INTO credit_transactions (id, user_id, type, crypto_type, crypto_amount, credit_amount, status, created_at, updated_at) VALUES
('bb0e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 'PURCHASE', 'USDC', 1000.00000000, 1000.00000000, 'COMPLETED', NOW() - INTERVAL '1 day', NOW()),
('bb0e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'PURCHASE', 'ETH', 0.25000000, 500.00000000, 'COMPLETED', NOW() - INTERVAL '12 hours', NOW()),
('bb0e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'PURCHASE', 'BTC', 0.01000000, 250.00000000, 'COMPLETED', NOW() - INTERVAL '6 hours', NOW());

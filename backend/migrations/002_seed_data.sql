-- Seed data for LFG prediction market platform

-- Insert test users (passwords are bcrypt hashed "password123")
INSERT INTO users (id, email, username, password_hash, email_verified, is_active, created_at, updated_at)
VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'alice@example.com', 'alice', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', true, true, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440002', 'bob@example.com', 'bob', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', true, true, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440003', 'charlie@example.com', 'charlie', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', true, true, NOW(), NOW()),
  ('550e8400-e29b-41d4-a716-446655440004', 'diana@example.com', 'diana', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', true, true, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert user wallets with initial balances
INSERT INTO wallets (id, user_id, balance, created_at, updated_at)
VALUES
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440001', 10000.00, NOW(), NOW()),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440002', 5000.00, NOW(), NOW()),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440003', 7500.00, NOW(), NOW()),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440004', 3000.00, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Insert markets (prediction market categories)
INSERT INTO markets (id, name, description, category, status, created_at, updated_at)
VALUES
  ('660e8400-e29b-41d4-a716-446655440001', 'US Politics 2025', 'Political events and elections in the United States', 'politics', 'active', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440002', 'Tech & AI 2025', 'Technology and artificial intelligence developments', 'technology', 'active', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440003', 'Sports 2025', 'Major sporting events and championships', 'sports', 'active', NOW(), NOW()),
  ('660e8400-e29b-41d4-a716-446655440004', 'Crypto Markets', 'Cryptocurrency price predictions', 'crypto', 'active', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert contracts (specific prediction questions)
INSERT INTO contracts (id, market_id, question, description, side, resolution_date, status, created_at, updated_at)
VALUES
  -- US Politics contracts
  ('770e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001',
   'Will Democrats control the Senate after 2026 midterms?',
   'Resolves YES if Democrats have 50+ seats, NO otherwise',
   'YES', '2026-11-30', 'active', NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440001',
   'Will Democrats control the Senate after 2026 midterms?',
   'Resolves YES if Democrats have 50+ seats, NO otherwise',
   'NO', '2026-11-30', 'active', NOW(), NOW()),

  -- Tech & AI contracts
  ('770e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440002',
   'Will OpenAI release GPT-5 in 2025?',
   'Resolves YES if GPT-5 is officially released by Dec 31, 2025',
   'YES', '2025-12-31', 'active', NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440002',
   'Will OpenAI release GPT-5 in 2025?',
   'Resolves YES if GPT-5 is officially released by Dec 31, 2025',
   'NO', '2025-12-31', 'active', NOW(), NOW()),

  -- Sports contracts
  ('770e8400-e29b-41d4-a716-446655440005', '660e8400-e29b-41d4-a716-446655440003',
   'Will the Lakers win the 2025 NBA Championship?',
   'Resolves YES if Los Angeles Lakers win the 2025 NBA Finals',
   'YES', '2025-06-30', 'active', NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440006', '660e8400-e29b-41d4-a716-446655440003',
   'Will the Lakers win the 2025 NBA Championship?',
   'Resolves YES if Los Angeles Lakers win the 2025 NBA Finals',
   'NO', '2025-06-30', 'active', NOW(), NOW()),

  -- Crypto contracts
  ('770e8400-e29b-41d4-a716-446655440007', '660e8400-e29b-41d4-a716-446655440004',
   'Will Bitcoin reach $150,000 in 2025?',
   'Resolves YES if BTC/USD reaches $150k at any point in 2025',
   'YES', '2025-12-31', 'active', NOW(), NOW()),
  ('770e8400-e29b-41d4-a716-446655440008', '660e8400-e29b-41d4-a716-446655440004',
   'Will Bitcoin reach $150,000 in 2025?',
   'Resolves YES if BTC/USD reaches $150k at any point in 2025',
   'NO', '2025-12-31', 'active', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Insert some sample orders
INSERT INTO orders (id, user_id, contract_id, type, side, quantity, limit_price_credits, status, created_at, updated_at)
VALUES
  -- Alice places buy orders for YES contracts
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440001', '770e8400-e29b-41d4-a716-446655440001',
   'limit', 'buy', 100, 0.65, 'active', NOW(), NOW()),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440001', '770e8400-e29b-41d4-a716-446655440003',
   'limit', 'buy', 50, 0.72, 'active', NOW(), NOW()),

  -- Bob places sell orders for YES contracts
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440002', '770e8400-e29b-41d4-a716-446655440001',
   'limit', 'sell', 75, 0.70, 'active', NOW(), NOW()),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440002', '770e8400-e29b-41d4-a716-446655440005',
   'limit', 'sell', 200, 0.45, 'active', NOW(), NOW()),

  -- Charlie places buy orders for NO contracts
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440003', '770e8400-e29b-41d4-a716-446655440002',
   'limit', 'buy', 150, 0.35, 'active', NOW(), NOW()),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440003', '770e8400-e29b-41d4-a716-446655440008',
   'limit', 'buy', 80, 0.55, 'active', NOW(), NOW()),

  -- Diana places mixed orders
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440004', '770e8400-e29b-41d4-a716-446655440007',
   'limit', 'buy', 120, 0.48, 'active', NOW(), NOW()),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440004', '770e8400-e29b-41d4-a716-446655440004',
   'limit', 'sell', 90, 0.62, 'active', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Insert sample credit transactions
INSERT INTO credit_transactions (id, user_id, transaction_type, crypto_type, crypto_amount, credit_amount, status, created_at, updated_at)
VALUES
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440001', 'purchase', 'USDC', 10000.00, 10000.00, 'completed', NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days'),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440002', 'purchase', 'ETH', 1.5, 4500.00, 'completed', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days'),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440002', 'purchase', 'USDC', 500.00, 500.00, 'completed', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440003', 'purchase', 'BTC', 0.15, 7500.00, 'completed', NOW() - INTERVAL '10 days', NOW() - INTERVAL '10 days'),
  (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440004', 'purchase', 'USDC', 3000.00, 3000.00, 'completed', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days')
ON CONFLICT DO NOTHING;

-- Summary statistics
DO $$
DECLARE
  user_count INTEGER;
  market_count INTEGER;
  contract_count INTEGER;
  order_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO user_count FROM users;
  SELECT COUNT(*) INTO market_count FROM markets;
  SELECT COUNT(*) INTO contract_count FROM contracts;
  SELECT COUNT(*) INTO order_count FROM orders;

  RAISE NOTICE 'Database seeded successfully!';
  RAISE NOTICE 'Users: %', user_count;
  RAISE NOTICE 'Markets: %', market_count;
  RAISE NOTICE 'Contracts: %', contract_count;
  RAISE NOTICE 'Orders: %', order_count;
END $$;

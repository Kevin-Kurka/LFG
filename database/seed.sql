-- Seed Sportsbook Providers
INSERT INTO sportsbook_providers (id, name, display_name, logo_url, website_url, supports_api, is_active, created_at) VALUES
('11111111-1111-1111-1111-111111111111', 'fanduel', 'FanDuel', 'https://www.fanduel.com/favicon.ico', 'https://www.fanduel.com', false, true, NOW()),
('22222222-2222-2222-2222-222222222222', 'draftkings', 'DraftKings', 'https://www.draftkings.com/favicon.ico', 'https://www.draftkings.com', false, true, NOW()),
('33333333-3333-3333-3333-333333333333', 'betmgm', 'BetMGM', 'https://sports.betmgm.com/favicon.ico', 'https://sports.betmgm.com', false, true, NOW()),
('44444444-4444-4444-4444-444444444444', 'caesars', 'Caesars Sportsbook', 'https://www.caesars.com/favicon.ico', 'https://www.caesars.com/sportsbook', false, true, NOW()),
('55555555-5555-5555-5555-555555555555', 'espnbet', 'ESPN BET', 'https://espnbet.com/favicon.ico', 'https://espnbet.com', false, true, NOW()),
('66666666-6666-6666-6666-666666666666', 'fanatics', 'Fanatics Sportsbook', 'https://www.fanaticssportsbook.com/favicon.ico', 'https://www.fanaticssportsbook.com', false, true, NOW());

-- Seed Sports
INSERT INTO sports (id, name, display_name, is_active) VALUES
('a0000000-0000-0000-0000-000000000001', 'nfl', 'NFL', true),
('a0000000-0000-0000-0000-000000000002', 'nba', 'NBA', true),
('a0000000-0000-0000-0000-000000000003', 'mlb', 'MLB', true),
('a0000000-0000-0000-0000-000000000004', 'nhl', 'NHL', true),
('a0000000-0000-0000-0000-000000000005', 'ncaaf', 'NCAA Football', true),
('a0000000-0000-0000-0000-000000000006', 'ncaab', 'NCAA Basketball', true),
('a0000000-0000-0000-0000-000000000007', 'soccer', 'Soccer', true),
('a0000000-0000-0000-0000-000000000008', 'mma', 'MMA/UFC', true);

-- Sample Sports Events (upcoming NFL games)
INSERT INTO sports_events (id, sport_id, home_team, away_team, event_time, league, status, created_at, updated_at) VALUES
('e0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'Kansas City Chiefs', 'Buffalo Bills', NOW() + INTERVAL '2 days', 'NFL', 'UPCOMING', NOW(), NOW()),
('e0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'San Francisco 49ers', 'Dallas Cowboys', NOW() + INTERVAL '3 days', 'NFL', 'UPCOMING', NOW(), NOW()),
('e0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'Philadelphia Eagles', 'Green Bay Packers', NOW() + INTERVAL '4 days', 'NFL', 'UPCOMING', NOW(), NOW()),
('e0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000002', 'Los Angeles Lakers', 'Boston Celtics', NOW() + INTERVAL '1 day', 'NBA', 'UPCOMING', NOW(), NOW()),
('e0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000002', 'Golden State Warriors', 'Phoenix Suns', NOW() + INTERVAL '2 days', 'NBA', 'UPCOMING', NOW(), NOW());

-- Sample Odds for Chiefs vs Bills (Moneyline)
INSERT INTO odds (id, event_id, sportsbook_id, market_type, bet_type, selection, odds_american, odds_decimal, line, is_active, last_updated) VALUES
-- FanDuel
('o1111111-1111-1111-1111-111111111111', 'e0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'MONEYLINE', 'MONEYLINE', 'Kansas City Chiefs', -140, 1.714, NULL, true, NOW()),
('o1111111-1111-1111-1111-111111111112', 'e0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'MONEYLINE', 'MONEYLINE', 'Buffalo Bills', 120, 2.20, NULL, true, NOW()),
-- DraftKings
('o2222222-2222-2222-2222-222222222221', 'e0000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'MONEYLINE', 'MONEYLINE', 'Kansas City Chiefs', -145, 1.690, NULL, true, NOW()),
('o2222222-2222-2222-2222-222222222222', 'e0000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'MONEYLINE', 'MONEYLINE', 'Buffalo Bills', 125, 2.25, NULL, true, NOW()),
-- BetMGM
('o3333333-3333-3333-3333-333333333331', 'e0000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'MONEYLINE', 'MONEYLINE', 'Kansas City Chiefs', -135, 1.741, NULL, true, NOW()),
('o3333333-3333-3333-3333-333333333332', 'e0000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'MONEYLINE', 'MONEYLINE', 'Buffalo Bills', 115, 2.15, NULL, true, NOW()),
-- Caesars
('o4444444-4444-4444-4444-444444444441', 'e0000000-0000-0000-0000-000000000001', '44444444-4444-4444-4444-444444444444', 'MONEYLINE', 'MONEYLINE', 'Kansas City Chiefs', -138, 1.725, NULL, true, NOW()),
('o4444444-4444-4444-4444-444444444442', 'e0000000-0000-0000-0000-000000000001', '44444444-4444-4444-4444-444444444444', 'MONEYLINE', 'MONEYLINE', 'Buffalo Bills', 118, 2.18, NULL, true, NOW()),
-- ESPN BET
('o5555555-5555-5555-5555-555555555551', 'e0000000-0000-0000-0000-000000000001', '55555555-5555-5555-5555-555555555555', 'MONEYLINE', 'MONEYLINE', 'Kansas City Chiefs', -142, 1.704, NULL, true, NOW()),
('o5555555-5555-5555-5555-555555555552', 'e0000000-0000-0000-0000-000000000001', '55555555-5555-5555-5555-555555555555', 'MONEYLINE', 'MONEYLINE', 'Buffalo Bills', 122, 2.22, NULL, true, NOW()),
-- Fanatics
('o6666666-6666-6666-6666-666666666661', 'e0000000-0000-0000-0000-000000000001', '66666666-6666-6666-6666-666666666666', 'MONEYLINE', 'MONEYLINE', 'Kansas City Chiefs', -150, 1.667, NULL, true, NOW()),
('o6666666-6666-6666-6666-666666666662', 'e0000000-0000-0000-0000-000000000001', '66666666-6666-6666-6666-666666666666', 'MONEYLINE', 'MONEYLINE', 'Buffalo Bills', 130, 2.30, NULL, true, NOW());

-- Sample Odds for Chiefs vs Bills (Spread)
INSERT INTO odds (id, event_id, sportsbook_id, market_type, bet_type, selection, odds_american, odds_decimal, line, is_active, last_updated) VALUES
-- FanDuel
('os111111-1111-1111-1111-111111111111', 'e0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'SPREAD', 'SPREAD', 'Kansas City Chiefs', -110, 1.909, -2.5, true, NOW()),
('os111111-1111-1111-1111-111111111112', 'e0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'SPREAD', 'SPREAD', 'Buffalo Bills', -110, 1.909, 2.5, true, NOW()),
-- DraftKings
('os222222-2222-2222-2222-222222222221', 'e0000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'SPREAD', 'SPREAD', 'Kansas City Chiefs', -112, 1.893, -2.5, true, NOW()),
('os222222-2222-2222-2222-222222222222', 'e0000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'SPREAD', 'SPREAD', 'Buffalo Bills', -108, 1.926, 2.5, true, NOW()),
-- BetMGM
('os333333-3333-3333-3333-333333333331', 'e0000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'SPREAD', 'SPREAD', 'Kansas City Chiefs', -110, 1.909, -3.0, true, NOW()),
('os333333-3333-3333-3333-333333333332', 'e0000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'SPREAD', 'SPREAD', 'Buffalo Bills', -110, 1.909, 3.0, true, NOW());

-- Sample Odds for Chiefs vs Bills (Total)
INSERT INTO odds (id, event_id, sportsbook_id, market_type, bet_type, selection, odds_american, odds_decimal, line, is_active, last_updated) VALUES
-- FanDuel
('ot111111-1111-1111-1111-111111111111', 'e0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'TOTAL', 'OVER', 'Over', -115, 1.870, 48.5, true, NOW()),
('ot111111-1111-1111-1111-111111111112', 'e0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'TOTAL', 'UNDER', 'Under', -105, 1.952, 48.5, true, NOW()),
-- DraftKings
('ot222222-2222-2222-2222-222222222221', 'e0000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'TOTAL', 'OVER', 'Over', -110, 1.909, 49.0, true, NOW()),
('ot222222-2222-2222-2222-222222222222', 'e0000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'TOTAL', 'UNDER', 'Under', -110, 1.909, 49.0, true, NOW());

-- Sample Prediction Market
INSERT INTO markets (id, ticker, question, rules, resolution_source, status, expires_at) VALUES
('m0000000-0000-0000-0000-000000000001', 'CHIEFS-SB-2025', 'Will the Kansas City Chiefs win Super Bowl LIX in 2025?', 'Resolves YES if Chiefs win Super Bowl LIX. Resolves NO otherwise.', 'Official NFL records', 'OPEN', NOW() + INTERVAL '6 months');

-- Contracts for prediction market
INSERT INTO contracts (id, market_id, side, ticker) VALUES
('c0000000-0000-0000-0000-000000000001', 'm0000000-0000-0000-0000-000000000001', 'YES', 'CHIEFS-SB-2025-YES'),
('c0000000-0000-0000-0000-000000000002', 'm0000000-0000-0000-0000-000000000001', 'NO', 'CHIEFS-SB-2025-NO');

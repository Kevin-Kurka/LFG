export interface User {
  id: string;
  email: string;
  username: string;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface WalletBalance {
  user_id: string;
  balance: number;
  updated_at: string;
}

export interface Transaction {
  id: string;
  user_id: string;
  type: string;
  amount: number;
  balance_after: number;
  description: string;
  reference_id?: string;
  created_at: string;
}

export interface Market {
  id: string;
  title: string;
  description: string;
  category: string;
  end_time: string;
  resolution_time?: string;
  status: 'open' | 'closed' | 'resolved';
  outcome?: string;
  created_at: string;
  contracts?: Contract[];
}

export interface Contract {
  id: string;
  market_id: string;
  outcome: string;
  current_price: number;
  total_volume: number;
}

export interface Order {
  id: string;
  user_id: string;
  contract_id: string;
  order_type: 'buy' | 'sell';
  price: number;
  quantity: number;
  filled_quantity: number;
  status: 'open' | 'filled' | 'partially_filled' | 'cancelled';
  created_at: string;
  updated_at: string;
}

export interface OrderBookEntry {
  price: number;
  quantity: number;
  total: number;
}

export interface OrderBook {
  contract_id: string;
  bids: OrderBookEntry[];
  asks: OrderBookEntry[];
}

export interface Trade {
  id: string;
  contract_id: string;
  buyer_id: string;
  seller_id: string;
  price: number;
  quantity: number;
  created_at: string;
}

export interface SportsbookProvider {
  id: string;
  name: string;
  logo_url?: string;
  supported: boolean;
}

export interface SportsbookAccount {
  id: string;
  user_id: string;
  provider_id: string;
  provider_name: string;
  username: string;
  balance: number;
  status: 'active' | 'inactive' | 'error';
  last_synced?: string;
  created_at: string;
}

export interface Sport {
  id: string;
  name: string;
  slug: string;
}

export interface SportsEvent {
  id: string;
  sport_id: string;
  sport_name: string;
  home_team: string;
  away_team: string;
  start_time: string;
  status: 'upcoming' | 'live' | 'finished';
  odds?: EventOdds[];
}

export interface EventOdds {
  id: string;
  event_id: string;
  provider_id: string;
  provider_name: string;
  market_type: string;
  outcome: string;
  odds_american: number;
  odds_decimal: number;
  odds_fractional: string;
  updated_at: string;
}

export interface ArbitrageOpportunity {
  id: string;
  event_id: string;
  event_description: string;
  start_time: string;
  market_type: string;
  legs: ArbitrageLeg[];
  total_implied_probability: number;
  profit_percentage: number;
  created_at: string;
}

export interface ArbitrageLeg {
  provider_id: string;
  provider_name: string;
  outcome: string;
  odds_decimal: number;
  implied_probability: number;
  stake_percentage: number;
}

export interface HedgeOpportunity {
  id: string;
  original_bet_id: string;
  original_event: string;
  original_outcome: string;
  original_stake: number;
  original_odds: number;
  potential_payout: number;
  hedge_provider_id: string;
  hedge_provider_name: string;
  hedge_outcome: string;
  hedge_odds: number;
  hedge_stake: number;
  guaranteed_profit: number;
  profit_percentage: number;
  created_at: string;
}

export interface Bet {
  id: string;
  user_id: string;
  event_id: string;
  event_description: string;
  provider_id: string;
  provider_name: string;
  market_type: string;
  outcome: string;
  stake: number;
  odds_decimal: number;
  potential_payout: number;
  status: 'pending' | 'won' | 'lost' | 'void' | 'cashed_out';
  placed_at: string;
  settled_at?: string;
  payout?: number;
}

export type OddsFormat = 'american' | 'decimal' | 'fractional';

export interface NotificationMessage {
  type: 'trade' | 'order_filled' | 'market_update' | 'arbitrage' | 'hedge';
  data: any;
  timestamp: string;
}

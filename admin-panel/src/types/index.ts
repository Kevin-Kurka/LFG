// User types
export interface User {
  id: string;
  email: string;
  username: string;
  wallet_address?: string;
  created_at: string;
  updated_at: string;
  is_admin?: boolean;
}

export interface UserWallet {
  user_id: string;
  balance: number;
  currency: string;
  locked_balance: number;
}

export interface UserActivity {
  id: string;
  user_id: string;
  action: string;
  details: any;
  timestamp: string;
}

// Market types
export interface Market {
  id: string;
  title: string;
  description: string;
  category: string;
  creator_id: string;
  resolution_source: string;
  end_date: string;
  status: 'ACTIVE' | 'CLOSED' | 'RESOLVED' | 'CANCELLED';
  created_at: string;
  updated_at: string;
  outcomes?: Outcome[];
  total_volume?: number;
}

export interface Outcome {
  id: string;
  market_id: string;
  name: string;
  probability: number;
  total_volume: number;
  is_winning?: boolean;
}

export interface CreateMarketInput {
  title: string;
  description: string;
  category: string;
  resolution_source: string;
  end_date: string;
  outcomes: string[];
}

// Order types
export interface Order {
  id: string;
  user_id: string;
  market_id: string;
  outcome_id: string;
  order_type: 'BUY' | 'SELL';
  quantity: number;
  price: number;
  status: 'OPEN' | 'FILLED' | 'CANCELLED' | 'PARTIALLY_FILLED';
  created_at: string;
  updated_at: string;
  filled_quantity?: number;
  market?: Market;
  user?: User;
}

export interface OrderBook {
  market_id: string;
  outcome_id: string;
  bids: OrderBookEntry[];
  asks: OrderBookEntry[];
}

export interface OrderBookEntry {
  price: number;
  quantity: number;
  orders_count: number;
}

// Trade types
export interface Trade {
  id: string;
  market_id: string;
  outcome_id: string;
  buyer_id: string;
  seller_id: string;
  price: number;
  quantity: number;
  timestamp: string;
  market?: Market;
}

// Sportsbook types
export interface SportsbookProvider {
  id: string;
  name: string;
  api_key?: string;
  enabled: boolean;
  last_sync?: string;
  created_at: string;
  updated_at: string;
}

export interface Sport {
  id: string;
  name: string;
  active: boolean;
}

export interface SportsEvent {
  id: string;
  provider_id: string;
  sport_id: string;
  home_team: string;
  away_team: string;
  start_time: string;
  status: 'SCHEDULED' | 'LIVE' | 'FINISHED' | 'CANCELLED';
  created_at: string;
  updated_at: string;
  odds?: Odds[];
}

export interface Odds {
  id: string;
  event_id: string;
  market_type: string;
  selection: string;
  odds: number;
  last_updated: string;
}

// Bet types
export interface Bet {
  id: string;
  user_id: string;
  event_id: string;
  market_type: string;
  selection: string;
  stake: number;
  odds: number;
  potential_payout: number;
  status: 'PENDING' | 'WON' | 'LOST' | 'CANCELLED' | 'VOID';
  placed_at: string;
  settled_at?: string;
  user?: User;
  event?: SportsEvent;
}

// Arbitrage types
export interface ArbitrageOpportunity {
  id: string;
  event_id: string;
  sport: string;
  selections: ArbitrageSelection[];
  profit_percentage: number;
  total_stake: number;
  guaranteed_profit: number;
  detected_at: string;
  expires_at?: string;
  status: 'ACTIVE' | 'EXPIRED' | 'EXECUTED';
}

export interface ArbitrageSelection {
  provider_id: string;
  provider_name: string;
  selection: string;
  odds: number;
  stake: number;
}

// Hedge types
export interface HedgeOpportunity {
  id: string;
  original_bet_id: string;
  event_id: string;
  hedge_selection: string;
  hedge_stake: number;
  hedge_odds: number;
  guaranteed_profit: number;
  detected_at: string;
  status: 'ACTIVE' | 'EXPIRED' | 'EXECUTED';
}

// Dashboard stats
export interface DashboardStats {
  total_users: number;
  total_markets: number;
  total_bets: number;
  total_volume: number;
  active_markets: number;
  active_users_24h: number;
  total_trades: number;
  arbitrage_opportunities: number;
}

// API Response types
export interface ApiResponse<T> {
  data: T;
  message?: string;
  error?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

// Auth types
export interface LoginCredentials {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

// Filter types
export interface MarketFilters {
  status?: string;
  category?: string;
  search?: string;
}

export interface OrderFilters {
  market_id?: string;
  user_id?: string;
  status?: string;
  order_type?: string;
}

export interface BetFilters {
  user_id?: string;
  event_id?: string;
  status?: string;
  date_from?: string;
  date_to?: string;
}

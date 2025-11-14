export interface Market {
  id: string;
  question: string;
  description: string;
  category: string;
  status: 'OPEN' | 'CLOSED' | 'RESOLVED';
  resolution_date: string;
  created_at: string;
  outcome?: string;
}

export interface Contract {
  id: string;
  market_id: string;
  side: 'YES' | 'NO';
  current_price_credits: number;
  volume: number;
}

export interface User {
  id: string;
  email: string;
  display_name: string;
  status: 'ACTIVE' | 'SUSPENDED' | 'BANNED';
  created_at: string;
}

export interface Order {
  id: string;
  user_id: string;
  contract_id: string;
  type: 'MARKET' | 'LIMIT';
  status: 'PENDING' | 'ACTIVE' | 'FILLED' | 'CANCELLED';
  quantity: number;
  quantity_filled: number;
  limit_price_credits?: number;
  created_at: string;
}

export interface DashboardStats {
  total_markets: number;
  active_markets: number;
  total_users: number;
  total_volume: number;
  total_orders: number;
  active_orders: number;
}

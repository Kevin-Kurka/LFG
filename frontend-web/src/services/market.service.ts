import { api } from './api';
import { Market, OrderBook } from '../types';

const MARKET_BASE_URL = import.meta.env.VITE_MARKET_SERVICE_URL || 'http://localhost:8082';

export const marketService = {
  async getMarkets(status?: string, search?: string): Promise<Market[]> {
    const params = new URLSearchParams();
    if (status) params.append('status', status);
    if (search) params.append('search', search);

    const response = await api.get<Market[]>(
      `${MARKET_BASE_URL}/markets${params.toString() ? `?${params.toString()}` : ''}`
    );
    return response.data;
  },

  async getMarket(id: string): Promise<Market> {
    const response = await api.get<Market>(`${MARKET_BASE_URL}/markets/${id}`);
    return response.data;
  },

  async getOrderBook(marketId: string, contractId: string): Promise<OrderBook> {
    const response = await api.get<OrderBook>(
      `${MARKET_BASE_URL}/markets/${marketId}/orderbook?contract_id=${contractId}`
    );
    return response.data;
  },

  async createMarket(
    title: string,
    description: string,
    category: string,
    endTime: string
  ): Promise<Market> {
    const response = await api.post<Market>(`${MARKET_BASE_URL}/markets`, {
      title,
      description,
      category,
      end_time: endTime,
    });
    return response.data;
  },
};

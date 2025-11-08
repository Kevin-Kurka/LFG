import { marketServiceApi } from './api';
import { Market, CreateMarketInput, Outcome, PaginatedResponse, MarketFilters } from '../types';

export const marketService = {
  async getMarkets(filters?: MarketFilters, page: number = 1, pageSize: number = 20): Promise<PaginatedResponse<Market>> {
    const params: any = { page, page_size: pageSize };
    if (filters?.status) params.status = filters.status;
    if (filters?.category) params.category = filters.category;
    if (filters?.search) params.search = filters.search;

    return marketServiceApi.get<PaginatedResponse<Market>>('/markets', { params });
  },

  async getMarketById(id: string): Promise<Market> {
    return marketServiceApi.get<Market>(`/markets/${id}`);
  },

  async createMarket(data: CreateMarketInput): Promise<Market> {
    return marketServiceApi.post<Market>('/markets', data);
  },

  async updateMarket(id: string, data: Partial<CreateMarketInput>): Promise<Market> {
    return marketServiceApi.put<Market>(`/markets/${id}`, data);
  },

  async resolveMarket(id: string, winningOutcomeId: string): Promise<Market> {
    return marketServiceApi.post<Market>(`/markets/${id}/resolve`, {
      winning_outcome_id: winningOutcomeId,
    });
  },

  async cancelMarket(id: string): Promise<Market> {
    return marketServiceApi.post<Market>(`/markets/${id}/cancel`);
  },

  async getMarketOutcomes(marketId: string): Promise<Outcome[]> {
    return marketServiceApi.get<Outcome[]>(`/markets/${marketId}/outcomes`);
  },

  async getMarketStats(): Promise<any> {
    return marketServiceApi.get<any>('/markets/stats');
  },
};

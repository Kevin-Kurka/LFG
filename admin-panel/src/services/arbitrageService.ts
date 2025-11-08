import { arbitrageServiceApi } from './api';
import { ArbitrageOpportunity, HedgeOpportunity, PaginatedResponse } from '../types';

export const arbitrageService = {
  async getArbitrageOpportunities(
    status?: string,
    page: number = 1,
    pageSize: number = 20
  ): Promise<PaginatedResponse<ArbitrageOpportunity>> {
    const params: any = { page, page_size: pageSize };
    if (status) params.status = status;

    return arbitrageServiceApi.get<PaginatedResponse<ArbitrageOpportunity>>('/arbitrage', { params });
  },

  async getArbitrageById(id: string): Promise<ArbitrageOpportunity> {
    return arbitrageServiceApi.get<ArbitrageOpportunity>(`/arbitrage/${id}`);
  },

  async executeArbitrage(id: string): Promise<ArbitrageOpportunity> {
    return arbitrageServiceApi.post<ArbitrageOpportunity>(`/arbitrage/${id}/execute`);
  },

  async getHedgeOpportunities(
    status?: string,
    page: number = 1,
    pageSize: number = 20
  ): Promise<PaginatedResponse<HedgeOpportunity>> {
    const params: any = { page, page_size: pageSize };
    if (status) params.status = status;

    return arbitrageServiceApi.get<PaginatedResponse<HedgeOpportunity>>('/hedges', { params });
  },

  async getHedgeById(id: string): Promise<HedgeOpportunity> {
    return arbitrageServiceApi.get<HedgeOpportunity>(`/hedges/${id}`);
  },

  async executeHedge(id: string): Promise<HedgeOpportunity> {
    return arbitrageServiceApi.post<HedgeOpportunity>(`/hedges/${id}/execute`);
  },

  async getArbitrageStats(): Promise<any> {
    return arbitrageServiceApi.get<any>('/arbitrage/stats');
  },
};

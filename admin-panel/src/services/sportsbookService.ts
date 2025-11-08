import { sportsbookServiceApi } from './api';
import { SportsbookProvider, Sport, SportsEvent, Odds, Bet, PaginatedResponse, BetFilters } from '../types';

export const sportsbookService = {
  // Providers
  async getProviders(): Promise<SportsbookProvider[]> {
    return sportsbookServiceApi.get<SportsbookProvider[]>('/providers');
  },

  async updateProvider(id: string, data: Partial<SportsbookProvider>): Promise<SportsbookProvider> {
    return sportsbookServiceApi.put<SportsbookProvider>(`/providers/${id}`, data);
  },

  async toggleProvider(id: string, enabled: boolean): Promise<SportsbookProvider> {
    return sportsbookServiceApi.patch<SportsbookProvider>(`/providers/${id}`, { enabled });
  },

  async syncProvider(id: string): Promise<void> {
    return sportsbookServiceApi.post<void>(`/providers/${id}/sync`);
  },

  // Sports
  async getSports(): Promise<Sport[]> {
    return sportsbookServiceApi.get<Sport[]>('/sports');
  },

  // Events
  async getEvents(sportId?: string, status?: string, page: number = 1): Promise<PaginatedResponse<SportsEvent>> {
    const params: any = { page, page_size: 20 };
    if (sportId) params.sport_id = sportId;
    if (status) params.status = status;

    return sportsbookServiceApi.get<PaginatedResponse<SportsEvent>>('/events', { params });
  },

  async getEventById(id: string): Promise<SportsEvent> {
    return sportsbookServiceApi.get<SportsEvent>(`/events/${id}`);
  },

  async getEventOdds(eventId: string): Promise<Odds[]> {
    return sportsbookServiceApi.get<Odds[]>(`/events/${eventId}/odds`);
  },

  // Bets
  async getBets(filters?: BetFilters, page: number = 1, pageSize: number = 20): Promise<PaginatedResponse<Bet>> {
    const params: any = { page, page_size: pageSize };
    if (filters?.user_id) params.user_id = filters.user_id;
    if (filters?.event_id) params.event_id = filters.event_id;
    if (filters?.status) params.status = filters.status;
    if (filters?.date_from) params.date_from = filters.date_from;
    if (filters?.date_to) params.date_to = filters.date_to;

    return sportsbookServiceApi.get<PaginatedResponse<Bet>>('/bets', { params });
  },

  async getBetById(id: string): Promise<Bet> {
    return sportsbookServiceApi.get<Bet>(`/bets/${id}`);
  },

  async settleBet(id: string, status: 'WON' | 'LOST' | 'VOID'): Promise<Bet> {
    return sportsbookServiceApi.post<Bet>(`/bets/${id}/settle`, { status });
  },

  async getBetStats(): Promise<any> {
    return sportsbookServiceApi.get<any>('/bets/stats');
  },
};

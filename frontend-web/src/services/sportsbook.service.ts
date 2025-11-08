import { api } from './api';
import {
  SportsbookAccount,
  SportsEvent,
  ArbitrageOpportunity,
  HedgeOpportunity,
  Bet,
} from '../types';

const SPORTSBOOK_BASE_URL = import.meta.env.VITE_SPORTSBOOK_SERVICE_URL || 'http://localhost:8088';

export const sportsbookService = {
  async linkAccount(
    providerId: string,
    username: string,
    password: string
  ): Promise<SportsbookAccount> {
    const response = await api.post<SportsbookAccount>(
      `${SPORTSBOOK_BASE_URL}/sportsbook/accounts`,
      {
        provider_id: providerId,
        username,
        password,
      }
    );
    return response.data;
  },

  async getAccounts(): Promise<SportsbookAccount[]> {
    const response = await api.get<SportsbookAccount[]>(
      `${SPORTSBOOK_BASE_URL}/sportsbook/accounts`
    );
    return response.data;
  },

  async deleteAccount(accountId: string): Promise<void> {
    await api.delete(`${SPORTSBOOK_BASE_URL}/sportsbook/accounts/${accountId}`);
  },

  async getEvents(sportId?: string, status?: string): Promise<SportsEvent[]> {
    const params = new URLSearchParams();
    if (sportId) params.append('sport_id', sportId);
    if (status) params.append('status', status);

    const response = await api.get<SportsEvent[]>(
      `${SPORTSBOOK_BASE_URL}/sportsbook/events${params.toString() ? `?${params.toString()}` : ''}`
    );
    return response.data;
  },

  async getEvent(eventId: string): Promise<SportsEvent> {
    const response = await api.get<SportsEvent>(
      `${SPORTSBOOK_BASE_URL}/sportsbook/events/${eventId}`
    );
    return response.data;
  },

  async getArbitrageOpportunities(
    minProfit?: number,
    sportId?: string
  ): Promise<ArbitrageOpportunity[]> {
    const params = new URLSearchParams();
    if (minProfit !== undefined) params.append('min_profit', minProfit.toString());
    if (sportId) params.append('sport_id', sportId);

    const response = await api.get<ArbitrageOpportunity[]>(
      `${SPORTSBOOK_BASE_URL}/sportsbook/arbitrage${params.toString() ? `?${params.toString()}` : ''}`
    );
    return response.data;
  },

  async getHedgeOpportunities(): Promise<HedgeOpportunity[]> {
    const response = await api.get<HedgeOpportunity[]>(
      `${SPORTSBOOK_BASE_URL}/sportsbook/hedges`
    );
    return response.data;
  },

  async trackBet(
    eventId: string,
    providerId: string,
    marketType: string,
    outcome: string,
    stake: number,
    oddsDecimal: number
  ): Promise<Bet> {
    const response = await api.post<Bet>(`${SPORTSBOOK_BASE_URL}/sportsbook/bets`, {
      event_id: eventId,
      provider_id: providerId,
      market_type: marketType,
      outcome,
      stake,
      odds_decimal: oddsDecimal,
    });
    return response.data;
  },

  async getBets(status?: string, providerId?: string): Promise<Bet[]> {
    const params = new URLSearchParams();
    if (status) params.append('status', status);
    if (providerId) params.append('provider_id', providerId);

    const response = await api.get<Bet[]>(
      `${SPORTSBOOK_BASE_URL}/sportsbook/bets${params.toString() ? `?${params.toString()}` : ''}`
    );
    return response.data;
  },
};

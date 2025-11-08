import { marketServiceApi, userServiceApi, orderServiceApi, sportsbookServiceApi, arbitrageServiceApi } from './api';
import { DashboardStats } from '../types';

export const dashboardService = {
  async getDashboardStats(): Promise<DashboardStats> {
    try {
      // Make parallel requests to all services
      const [marketStats, userStats, orderStats, betStats, arbStats] = await Promise.all([
        marketServiceApi.get<any>('/markets/stats').catch(() => ({ total: 0, active: 0, volume: 0 })),
        userServiceApi.get<any>('/users/stats').catch(() => ({ total: 0, active_24h: 0 })),
        orderServiceApi.get<any>('/orders/stats').catch(() => ({ total_trades: 0 })),
        sportsbookServiceApi.get<any>('/bets/stats').catch(() => ({ total: 0 })),
        arbitrageServiceApi.get<any>('/arbitrage/stats').catch(() => ({ opportunities: 0 })),
      ]);

      return {
        total_users: userStats.total || 0,
        total_markets: marketStats.total || 0,
        total_bets: betStats.total || 0,
        total_volume: marketStats.volume || 0,
        active_markets: marketStats.active || 0,
        active_users_24h: userStats.active_24h || 0,
        total_trades: orderStats.total_trades || 0,
        arbitrage_opportunities: arbStats.opportunities || 0,
      };
    } catch (error) {
      console.error('Error fetching dashboard stats:', error);
      return {
        total_users: 0,
        total_markets: 0,
        total_bets: 0,
        total_volume: 0,
        active_markets: 0,
        active_users_24h: 0,
        total_trades: 0,
        arbitrage_opportunities: 0,
      };
    }
  },
};

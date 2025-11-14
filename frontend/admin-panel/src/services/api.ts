const API_BASE_URL = '/api';

export class ApiService {
  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Request failed' }));
      throw new Error(error.error || 'Request failed');
    }

    return response.json();
  }

  // Markets
  async getMarkets(params?: { status?: string; search?: string }) {
    const query = new URLSearchParams(params as Record<string, string>).toString();
    return this.request<{ markets: any[] }>(`/markets?${query}`);
  }

  async getMarket(id: string) {
    return this.request<any>(`/markets/${id}`);
  }

  // Dashboard stats (mock for now)
  async getDashboardStats() {
    // In a real implementation, this would call a backend endpoint
    const markets = await this.getMarkets();
    return {
      total_markets: markets.markets.length,
      active_markets: markets.markets.filter((m: any) => m.status === 'OPEN').length,
      total_users: 0, // Would come from user service
      total_volume: 0, // Would come from analytics service
      total_orders: 0, // Would come from order service
      active_orders: 0, // Would come from order service
    };
  }
}

export const api = new ApiService();

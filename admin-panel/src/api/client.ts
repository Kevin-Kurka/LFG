const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8000';

export interface ApiResponse<T> {
  data?: T;
  error?: string;
}

class ApiClient {
  private baseUrl: string;
  private token: string | null = null;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
    this.token = localStorage.getItem('auth_token');
  }

  setToken(token: string) {
    this.token = token;
    localStorage.setItem('auth_token', token);
  }

  clearToken() {
    this.token = null;
    localStorage.removeItem('auth_token');
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...((options.headers as Record<string, string>) || {}),
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    try {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        ...options,
        headers,
      });

      if (!response.ok) {
        const error = await response.text();
        return { error: error || `HTTP ${response.status}` };
      }

      const data = await response.json();
      return { data };
    } catch (error) {
      return { error: error instanceof Error ? error.message : 'Network error' };
    }
  }

  async get<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, body?: any): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async put<T>(endpoint: string, body?: any): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }

  // Health check
  async checkHealth() {
    return this.get<{ status: string }>('/health');
  }

  // Auth endpoints
  async login(email: string, password: string) {
    return this.post<{ token: string; user: any }>('/login', {
      email,
      password,
    });
  }

  async register(email: string, password: string, username: string) {
    return this.post<{ token: string; user: any }>('/register', {
      email,
      password,
      username,
    });
  }

  // Market endpoints
  async getMarkets(page: number = 1, limit: number = 20) {
    return this.get<{ markets: any[]; total: number; page: number }>(
      `/markets?page=${page}&limit=${limit}`
    );
  }

  async getMarket(id: string) {
    return this.get<any>(`/markets/${id}`);
  }

  async createMarket(market: any) {
    return this.post<any>('/markets', market);
  }

  // User endpoints
  async getUsers(page: number = 1, limit: number = 20) {
    return this.get<{ users: any[]; total: number }>(
      `/users?page=${page}&limit=${limit}`
    );
  }

  async getUser(id: string) {
    return this.get<any>(`/users/${id}`);
  }

  // Wallet endpoints
  async getWallet(userId: string) {
    return this.get<{ balance: number; currency: string }>(
      `/balance`
    );
  }

  // Order endpoints
  async getOrders(userId?: string, contractId?: string) {
    const params = new URLSearchParams();
    if (userId) params.append('user_id', userId);
    if (contractId) params.append('contract_id', contractId);
    return this.get<{ orders: any[] }>(
      `/orders${params.toString() ? '?' + params.toString() : ''}`
    );
  }
}

export const apiClient = new ApiClient();

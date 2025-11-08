import { api } from './api';
import { AuthResponse, User } from '../types';

export const authService = {
  async register(email: string, username: string, password: string): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/register', {
      email,
      username,
      password,
    });
    return response.data;
  },

  async login(email: string, password: string): Promise<AuthResponse> {
    const response = await api.post<AuthResponse>('/login', {
      email,
      password,
    });
    return response.data;
  },

  async getProfile(): Promise<User> {
    const response = await api.get<User>('/profile');
    return response.data;
  },

  async updateProfile(username: string): Promise<User> {
    const response = await api.put<User>('/profile', { username });
    return response.data;
  },

  logout(): void {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },
};

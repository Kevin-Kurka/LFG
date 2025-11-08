import { userServiceApi } from './api';
import { User, UserWallet, UserActivity, PaginatedResponse } from '../types';

export const userService = {
  async getUsers(page: number = 1, pageSize: number = 20, search?: string): Promise<PaginatedResponse<User>> {
    const params: any = { page, page_size: pageSize };
    if (search) {
      params.search = search;
    }
    return userServiceApi.get<PaginatedResponse<User>>('/users', { params });
  },

  async getUserById(id: string): Promise<User> {
    return userServiceApi.get<User>(`/users/${id}`);
  },

  async getUserWallet(userId: string): Promise<UserWallet> {
    return userServiceApi.get<UserWallet>(`/users/${userId}/wallet`);
  },

  async getUserActivity(userId: string, page: number = 1): Promise<PaginatedResponse<UserActivity>> {
    return userServiceApi.get<PaginatedResponse<UserActivity>>(`/users/${userId}/activity`, {
      params: { page, page_size: 50 },
    });
  },

  async updateUserBalance(userId: string, amount: number): Promise<UserWallet> {
    return userServiceApi.post<UserWallet>(`/users/${userId}/wallet/adjust`, { amount });
  },

  async getUserStats(): Promise<any> {
    return userServiceApi.get<any>('/users/stats');
  },
};

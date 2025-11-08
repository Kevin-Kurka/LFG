import { api } from './api';
import { WalletBalance, Transaction } from '../types';

const WALLET_BASE_URL = import.meta.env.VITE_WALLET_SERVICE_URL || 'http://localhost:8081';

export const walletService = {
  async getBalance(): Promise<WalletBalance> {
    const response = await api.get<WalletBalance>(`${WALLET_BASE_URL}/balance`);
    return response.data;
  },

  async getTransactions(
    type?: string,
    limit: number = 50,
    offset: number = 0
  ): Promise<Transaction[]> {
    const params = new URLSearchParams({
      limit: limit.toString(),
      offset: offset.toString(),
    });
    if (type) params.append('type', type);

    const response = await api.get<Transaction[]>(
      `${WALLET_BASE_URL}/transactions?${params.toString()}`
    );
    return response.data;
  },
};

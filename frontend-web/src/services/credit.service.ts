import { api } from './api';
import { Transaction } from '../types';

const CREDIT_BASE_URL = import.meta.env.VITE_CREDIT_SERVICE_URL || 'http://localhost:8086';

export const creditService = {
  async buyCredits(amount: number, paymentMethod: string): Promise<Transaction> {
    const response = await api.post<Transaction>(`${CREDIT_BASE_URL}/exchange/buy`, {
      amount,
      payment_method: paymentMethod,
    });
    return response.data;
  },

  async sellCredits(amount: number, withdrawMethod: string): Promise<Transaction> {
    const response = await api.post<Transaction>(`${CREDIT_BASE_URL}/exchange/sell`, {
      amount,
      withdraw_method: withdrawMethod,
    });
    return response.data;
  },

  async getHistory(): Promise<Transaction[]> {
    const response = await api.get<Transaction[]>(`${CREDIT_BASE_URL}/exchange/history`);
    return response.data;
  },
};

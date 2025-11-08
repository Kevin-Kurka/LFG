import { api } from './api';
import { Order } from '../types';

const ORDER_BASE_URL = import.meta.env.VITE_ORDER_SERVICE_URL || 'http://localhost:8085';

export const orderService = {
  async placeOrder(
    contractId: string,
    orderType: 'buy' | 'sell',
    price: number,
    quantity: number
  ): Promise<Order> {
    const response = await api.post<Order>(`${ORDER_BASE_URL}/orders/place`, {
      contract_id: contractId,
      order_type: orderType,
      price,
      quantity,
    });
    return response.data;
  },

  async cancelOrder(orderId: string): Promise<void> {
    await api.post(`${ORDER_BASE_URL}/orders/cancel`, {
      order_id: orderId,
    });
  },

  async getOrders(status?: string): Promise<Order[]> {
    const params = status ? `?status=${status}` : '';
    const response = await api.get<Order[]>(`${ORDER_BASE_URL}/orders${params}`);
    return response.data;
  },

  async getOrder(orderId: string): Promise<Order> {
    const response = await api.get<Order>(`${ORDER_BASE_URL}/orders/${orderId}`);
    return response.data;
  },
};

import { orderServiceApi } from './api';
import { Order, OrderBook, Trade, PaginatedResponse, OrderFilters } from '../types';

export const orderService = {
  async getOrders(filters?: OrderFilters, page: number = 1, pageSize: number = 20): Promise<PaginatedResponse<Order>> {
    const params: any = { page, page_size: pageSize };
    if (filters?.market_id) params.market_id = filters.market_id;
    if (filters?.user_id) params.user_id = filters.user_id;
    if (filters?.status) params.status = filters.status;
    if (filters?.order_type) params.order_type = filters.order_type;

    return orderServiceApi.get<PaginatedResponse<Order>>('/orders', { params });
  },

  async getOrderById(id: string): Promise<Order> {
    return orderServiceApi.get<Order>(`/orders/${id}`);
  },

  async getOrderBook(marketId: string, outcomeId: string): Promise<OrderBook> {
    return orderServiceApi.get<OrderBook>(`/orderbook/${marketId}/${outcomeId}`);
  },

  async getTrades(marketId?: string, page: number = 1, pageSize: number = 50): Promise<PaginatedResponse<Trade>> {
    const params: any = { page, page_size: pageSize };
    if (marketId) params.market_id = marketId;

    return orderServiceApi.get<PaginatedResponse<Trade>>('/trades', { params });
  },

  async cancelOrder(id: string): Promise<Order> {
    return orderServiceApi.post<Order>(`/orders/${id}/cancel`);
  },

  async getOrderStats(): Promise<any> {
    return orderServiceApi.get<any>('/orders/stats');
  },
};

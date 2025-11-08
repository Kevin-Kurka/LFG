import api from './api';

export interface AdminCryptoStats {
  total_deposits: number;
  total_withdrawals: number;
  pending_deposits: number;
  pending_withdrawals: number;
  total_volume_usd: number;
  deposits_24h: number;
  withdrawals_24h: number;
}

export interface CryptoDeposit {
  id: string;
  user_id: string;
  wallet_id: string;
  currency: string;
  amount_crypto: number;
  amount_credits: number;
  exchange_rate: number;
  address: string;
  tx_hash: string;
  confirmations: number;
  required_confirmations: number;
  status: string;
  block_number?: number;
  detected_at: string;
  confirmed_at?: string;
  credited_at?: string;
  created_at: string;
}

export interface CryptoWithdrawal {
  id: string;
  user_id: string;
  wallet_id: string;
  currency: string;
  amount_credits: number;
  amount_crypto: number;
  exchange_rate: number;
  network_fee: number;
  platform_fee: number;
  total_fee_credits: number;
  to_address: string;
  tx_hash?: string;
  status: string;
  confirmations: number;
  block_number?: number;
  error_message?: string;
  requested_at: string;
  processed_at?: string;
  confirmed_at?: string;
  created_at: string;
}

export interface ExchangeRate {
  currency: string;
  usd_rate: number;
  provider: string;
  last_updated: string;
}

export interface CryptoMonitoringJob {
  id: string;
  currency: string;
  last_scanned_block: number;
  last_scan_at: string;
  is_running: boolean;
  error_count: number;
  last_error?: string;
}

class CryptoService {
  async getStats(): Promise<AdminCryptoStats> {
    const response = await api.get('/admin/crypto/stats');
    return response.data;
  }

  async getAllDeposits(limit: number = 100): Promise<CryptoDeposit[]> {
    const response = await api.get('/admin/crypto/deposits', { params: { limit } });
    return response.data.deposits || [];
  }

  async getPendingDeposits(): Promise<CryptoDeposit[]> {
    const response = await api.get('/admin/crypto/deposits/pending');
    return response.data.deposits || [];
  }

  async getAllWithdrawals(limit: number = 100): Promise<CryptoWithdrawal[]> {
    const response = await api.get('/admin/crypto/withdrawals', { params: { limit } });
    return response.data.withdrawals || [];
  }

  async getPendingWithdrawals(): Promise<CryptoWithdrawal[]> {
    const response = await api.get('/admin/crypto/withdrawals/pending');
    return response.data.withdrawals || [];
  }

  async approveWithdrawal(withdrawalId: string): Promise<void> {
    await api.post(`/admin/crypto/withdrawals/${withdrawalId}/approve`);
  }

  async rejectWithdrawal(withdrawalId: string, reason: string): Promise<void> {
    await api.post(`/admin/crypto/withdrawals/${withdrawalId}/reject`, { reason });
  }

  async getExchangeRates(): Promise<ExchangeRate[]> {
    const response = await api.get('/admin/crypto/rates');
    return response.data.rates || [];
  }

  async updateExchangeRate(currency: string, rate: number): Promise<void> {
    await api.post('/admin/crypto/rates/update', { currency, rate });
  }

  async getMonitoringJobs(): Promise<CryptoMonitoringJob[]> {
    const response = await api.get('/admin/crypto/monitoring');
    return response.data.jobs || [];
  }

  async restartMonitoring(currency: string): Promise<void> {
    await api.post('/admin/crypto/monitoring/restart', { currency });
  }
}

export default new CryptoService();

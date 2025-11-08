import api from './api';

export interface CryptoWallet {
  id: string;
  user_id: string;
  currency: string;
  address: string;
  balance_crypto: number;
  created_at: string;
  updated_at: string;
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
  id: string;
  currency: string;
  usd_rate: number;
  provider: string;
  last_updated: string;
}

export interface ConversionResult {
  amount: number;
  result: number;
  exchange_rate: number;
  currency: string;
  direction: string;
}

class CryptoService {
  // Get all crypto wallets for the user
  async getWallets(): Promise<CryptoWallet[]> {
    const response = await api.get('/crypto/wallets');
    return response.data.wallets || [];
  }

  // Create a new wallet for a currency
  async createWallet(currency: string): Promise<{ wallet: CryptoWallet; mnemonic?: string }> {
    const response = await api.post(`/crypto/wallets/${currency}`);
    return response.data;
  }

  // Get deposit history
  async getDeposits(): Promise<CryptoDeposit[]> {
    const response = await api.get('/crypto/deposits');
    return response.data.deposits || [];
  }

  // Get pending deposits
  async getPendingDeposits(): Promise<CryptoDeposit[]> {
    const response = await api.get('/crypto/deposits/pending');
    return response.data.deposits || [];
  }

  // Simulate a deposit (for testing)
  async simulateDeposit(currency: string, amountCrypto: number, txHash?: string): Promise<CryptoDeposit> {
    const response = await api.post('/crypto/deposits/simulate', {
      currency,
      amount_crypto: amountCrypto,
      tx_hash: txHash,
    });
    return response.data.deposit;
  }

  // Request a withdrawal
  async requestWithdrawal(currency: string, amountCredits: number, toAddress: string): Promise<CryptoWithdrawal> {
    const response = await api.post('/crypto/withdraw', {
      currency,
      amount_credits: amountCredits,
      to_address: toAddress,
    });
    return response.data.withdrawal;
  }

  // Get withdrawal history
  async getWithdrawals(): Promise<CryptoWithdrawal[]> {
    const response = await api.get('/crypto/withdrawals');
    return response.data.withdrawals || [];
  }

  // Get current exchange rates
  async getExchangeRates(): Promise<Record<string, ExchangeRate>> {
    const response = await api.get('/crypto/rates');
    return response.data.rates || {};
  }

  // Convert between crypto and credits
  async convertAmount(
    currency: string,
    amount: number,
    direction: 'crypto_to_credits' | 'credits_to_crypto'
  ): Promise<ConversionResult> {
    const response = await api.get('/crypto/convert', {
      params: { currency, amount, direction },
    });
    return response.data;
  }
}

export default new CryptoService();

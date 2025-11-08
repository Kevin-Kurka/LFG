import React, { useState, useEffect } from 'react';
import { walletService } from '../services/wallet.service';
import { creditService } from '../services/credit.service';
import { WalletBalance, Transaction } from '../types';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { ErrorMessage } from '../components/ErrorMessage';
import { formatCurrency, formatDate } from '../utils/format';
import { handleApiError } from '../services/api';

export const Wallet: React.FC = () => {
  const [balance, setBalance] = useState<WalletBalance | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showBuyForm, setShowBuyForm] = useState(false);
  const [buyAmount, setBuyAmount] = useState('');
  const [buyLoading, setBuyLoading] = useState(false);
  const [buyError, setBuyError] = useState('');
  const [buySuccess, setBuySuccess] = useState('');

  useEffect(() => {
    loadWalletData();
  }, []);

  const loadWalletData = async () => {
    try {
      setLoading(true);
      setError('');
      const [balanceData, transactionsData] = await Promise.all([
        walletService.getBalance(),
        walletService.getTransactions(),
      ]);
      setBalance(balanceData);
      setTransactions(transactionsData);
    } catch (err: any) {
      setError(err.message || 'Failed to load wallet data');
    } finally {
      setLoading(false);
    }
  };

  const handleBuyCredits = async (e: React.FormEvent) => {
    e.preventDefault();
    setBuyError('');
    setBuySuccess('');
    setBuyLoading(true);

    try {
      await creditService.buyCredits(parseFloat(buyAmount), 'card');
      setBuySuccess('Credits purchased successfully!');
      setBuyAmount('');
      setShowBuyForm(false);
      loadWalletData();
    } catch (err) {
      setBuyError(handleApiError(err));
    } finally {
      setBuyLoading(false);
    }
  };

  const getTransactionTypeColor = (type: string) => {
    switch (type) {
      case 'deposit':
      case 'credit':
      case 'win':
        return 'text-green-600 dark:text-green-400';
      case 'withdraw':
      case 'debit':
      case 'bet':
        return 'text-red-600 dark:text-red-400';
      default:
        return 'text-gray-600 dark:text-gray-400';
    }
  };

  const getTransactionSign = (type: string) => {
    return ['deposit', 'credit', 'win'].includes(type) ? '+' : '-';
  };

  if (loading) {
    return <LoadingSpinner size="lg" text="Loading wallet..." />;
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">Wallet</h1>
          <p className="text-gray-600 dark:text-gray-400">
            Manage your balance and view transaction history
          </p>
        </div>

        {error && <ErrorMessage message={error} onRetry={loadWalletData} />}

        <div className="bg-gradient-to-br from-primary-600 to-primary-700 rounded-lg shadow-lg p-8 mb-6 text-white">
          <p className="text-primary-100 mb-2">Available Balance</p>
          <h2 className="text-5xl font-bold mb-6">
            {balance ? formatCurrency(balance.balance) : '$0.00'}
          </h2>
          <div className="flex space-x-4">
            <button
              onClick={() => setShowBuyForm(!showBuyForm)}
              className="px-6 py-3 bg-white text-primary-600 font-semibold rounded-lg hover:bg-primary-50 transition-colors"
            >
              Buy Credits
            </button>
            <button className="px-6 py-3 bg-primary-500 hover:bg-primary-400 text-white font-semibold rounded-lg transition-colors">
              Withdraw
            </button>
          </div>
        </div>

        {showBuyForm && (
          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-6 mb-6">
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
              Buy Credits
            </h3>

            {buySuccess && (
              <div className="mb-4 p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg text-green-800 dark:text-green-400">
                {buySuccess}
              </div>
            )}

            {buyError && <ErrorMessage message={buyError} />}

            <form onSubmit={handleBuyCredits} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Amount (USD)
                </label>
                <input
                  type="number"
                  step="0.01"
                  min="1"
                  required
                  value={buyAmount}
                  onChange={(e) => setBuyAmount(e.target.value)}
                  className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  placeholder="100.00"
                />
              </div>

              <div className="flex space-x-4">
                <button
                  type="submit"
                  disabled={buyLoading}
                  className="flex-1 px-4 py-3 bg-primary-600 hover:bg-primary-700 disabled:bg-primary-400 text-white font-semibold rounded-lg transition-colors"
                >
                  {buyLoading ? 'Processing...' : 'Purchase'}
                </button>
                <button
                  type="button"
                  onClick={() => setShowBuyForm(false)}
                  className="px-6 py-3 bg-gray-200 dark:bg-dark-700 text-gray-700 dark:text-gray-300 font-semibold rounded-lg hover:bg-gray-300 dark:hover:bg-dark-600 transition-colors"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        )}

        <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg overflow-hidden">
          <div className="p-6 border-b border-gray-200 dark:border-dark-700">
            <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
              Transaction History
            </h2>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-gray-50 dark:bg-dark-900">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Date
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Type
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Description
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Amount
                  </th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                    Balance
                  </th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 dark:divide-dark-700">
                {transactions.length > 0 ? (
                  transactions.map((transaction) => (
                    <tr
                      key={transaction.id}
                      className="hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors"
                    >
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                        {formatDate(transaction.created_at)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className="px-2 py-1 text-xs font-semibold rounded-full bg-gray-100 dark:bg-dark-700 text-gray-800 dark:text-gray-300">
                          {transaction.type.toUpperCase()}
                        </span>
                      </td>
                      <td className="px-6 py-4 text-sm text-gray-600 dark:text-gray-400">
                        {transaction.description}
                      </td>
                      <td
                        className={`px-6 py-4 whitespace-nowrap text-right text-sm font-semibold ${getTransactionTypeColor(
                          transaction.type
                        )}`}
                      >
                        {getTransactionSign(transaction.type)}
                        {formatCurrency(Math.abs(transaction.amount))}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm text-gray-900 dark:text-white">
                        {formatCurrency(transaction.balance_after)}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td
                      colSpan={5}
                      className="px-6 py-12 text-center text-gray-500 dark:text-gray-400"
                    >
                      No transactions yet
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};

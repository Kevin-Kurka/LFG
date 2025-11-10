import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { walletService } from '../services/wallet.service';
import { orderService } from '../services/order.service';
import { sportsbookService } from '../services/sportsbook.service';
import { WalletBalance, Order, Bet } from '../types';
import { formatCurrency, formatTimeAgo } from '../utils/format';
import { ErrorMessage } from '../components/ErrorMessage';
import SkeletonStats from '../components/SkeletonStats';
import SkeletonCard from '../components/SkeletonCard';

export const Dashboard: React.FC = () => {
  const { user } = useAuth();
  const [balance, setBalance] = useState<WalletBalance | null>(null);
  const [recentOrders, setRecentOrders] = useState<Order[]>([]);
  const [recentBets, setRecentBets] = useState<Bet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      setLoading(true);
      const [balanceData, ordersData, betsData] = await Promise.all([
        walletService.getBalance(),
        orderService.getOrders().catch(() => []),
        sportsbookService.getBets().catch(() => []),
      ]);

      setBalance(balanceData);
      setRecentOrders(ordersData.slice(0, 5));
      setRecentBets(betsData.slice(0, 5));
    } catch (err: any) {
      setError(err.message || 'Failed to load dashboard data');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Welcome back, {user?.username}!
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Here's your trading overview and recent activity
          </p>
        </div>

        {error && <ErrorMessage message={error} onRetry={loadDashboardData} />}

        {loading ? (
          <>
            <SkeletonStats count={3} className="mb-8" />
            <div className="grid md:grid-cols-2 gap-6 mb-8">
              <SkeletonCard />
              <SkeletonCard />
            </div>
          </>
        ) : (
          <>
            <div className="grid md:grid-cols-3 gap-6 mb-8">
          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">
                Wallet Balance
              </h3>
              <svg className="w-8 h-8 text-primary-600" fill="currentColor" viewBox="0 0 20 20">
                <path d="M4 4a2 2 0 00-2 2v1h16V6a2 2 0 00-2-2H4z" />
                <path
                  fillRule="evenodd"
                  d="M18 9H2v5a2 2 0 002 2h12a2 2 0 002-2V9zM4 13a1 1 0 011-1h1a1 1 0 110 2H5a1 1 0 01-1-1zm5-1a1 1 0 100 2h1a1 1 0 100-2H9z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
            <p className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
              {balance ? formatCurrency(balance.balance) : '$0.00'}
            </p>
            <Link
              to="/wallet"
              className="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 font-medium"
            >
              View wallet →
            </Link>
          </div>

          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">
                Active Orders
              </h3>
              <svg className="w-8 h-8 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z" />
                <path
                  fillRule="evenodd"
                  d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 4a1 1 0 000 2h.01a1 1 0 100-2H7zm3 0a1 1 0 000 2h3a1 1 0 100-2h-3zm-3 4a1 1 0 100 2h.01a1 1 0 100-2H7zm3 0a1 1 0 100 2h3a1 1 0 100-2h-3z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
            <p className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
              {recentOrders.filter((o) => o.status === 'open').length}
            </p>
            <Link
              to="/markets"
              className="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 font-medium"
            >
              View markets →
            </Link>
          </div>

          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">
                Pending Bets
              </h3>
              <svg className="w-8 h-8 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fillRule="evenodd"
                  d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
                  clipRule="evenodd"
                />
              </svg>
            </div>
            <p className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
              {recentBets.filter((b) => b.status === 'pending').length}
            </p>
            <Link
              to="/bets"
              className="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 font-medium"
            >
              View bets →
            </Link>
          </div>
        </div>

        <div className="grid md:grid-cols-2 gap-6 mb-8">
          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg overflow-hidden">
            <div className="p-6 border-b border-gray-200 dark:border-dark-700">
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
                Recent Orders
              </h2>
            </div>
            <div className="p-6">
              {recentOrders.length > 0 ? (
                <div className="space-y-4">
                  {recentOrders.map((order) => (
                    <div
                      key={order.id}
                      className="flex items-center justify-between p-4 bg-gray-50 dark:bg-dark-900 rounded-lg"
                    >
                      <div>
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                          {order.order_type.toUpperCase()} {order.quantity} @ {formatCurrency(order.price)}
                        </p>
                        <p className="text-xs text-gray-500 dark:text-gray-400">
                          {formatTimeAgo(order.created_at)}
                        </p>
                      </div>
                      <span
                        className={`px-2 py-1 text-xs font-semibold rounded-full ${
                          order.status === 'filled'
                            ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400'
                            : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400'
                        }`}
                      >
                        {order.status}
                      </span>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-center text-gray-500 dark:text-gray-400 py-8">
                  No recent orders
                </p>
              )}
            </div>
          </div>

          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg overflow-hidden">
            <div className="p-6 border-b border-gray-200 dark:border-dark-700">
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Recent Bets</h2>
            </div>
            <div className="p-6">
              {recentBets.length > 0 ? (
                <div className="space-y-4">
                  {recentBets.map((bet) => (
                    <div
                      key={bet.id}
                      className="flex items-center justify-between p-4 bg-gray-50 dark:bg-dark-900 rounded-lg"
                    >
                      <div>
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                          {bet.outcome} - {formatCurrency(bet.stake)}
                        </p>
                        <p className="text-xs text-gray-500 dark:text-gray-400">
                          {bet.provider_name} • {formatTimeAgo(bet.placed_at)}
                        </p>
                      </div>
                      <span
                        className={`px-2 py-1 text-xs font-semibold rounded-full ${
                          bet.status === 'won'
                            ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400'
                            : bet.status === 'lost'
                            ? 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400'
                            : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400'
                        }`}
                      >
                        {bet.status}
                      </span>
                    </div>
                  ))}
                </div>
              ) : (
                <p className="text-center text-gray-500 dark:text-gray-400 py-8">No recent bets</p>
              )}
            </div>
          </div>
        </div>

        <div className="grid md:grid-cols-3 gap-6">
          <Link
            to="/arbitrage"
            className="bg-gradient-to-br from-green-500 to-green-600 rounded-lg shadow-lg p-6 text-white hover:shadow-xl transition-shadow"
          >
            <div className="text-4xl mb-3">💰</div>
            <h3 className="text-xl font-semibold mb-2">Find Arbitrage</h3>
            <p className="text-green-100">Discover guaranteed profit opportunities</p>
          </Link>

          <Link
            to="/hedges"
            className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg shadow-lg p-6 text-white hover:shadow-xl transition-shadow"
          >
            <div className="text-4xl mb-3">🛡️</div>
            <h3 className="text-xl font-semibold mb-2">Hedge Bets</h3>
            <p className="text-blue-100">Lock in profits on existing bets</p>
          </Link>

          <Link
            to="/link-account"
            className="bg-gradient-to-br from-purple-500 to-purple-600 rounded-lg shadow-lg p-6 text-white hover:shadow-xl transition-shadow"
          >
            <div className="text-4xl mb-3">🔗</div>
            <h3 className="text-xl font-semibold mb-2">Link Accounts</h3>
            <p className="text-purple-100">Connect your sportsbook accounts</p>
          </Link>
        </div>
          </>
        )}
      </div>
    </div>
  );
};

import React, { useState, useEffect } from 'react';
import { sportsbookService } from '../services/sportsbook.service';
import { SportsbookAccount } from '../types';
import { SportsbookAccountCard } from '../components/SportsbookAccountCard';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { ErrorMessage } from '../components/ErrorMessage';
import { handleApiError } from '../services/api';

const SPORTSBOOK_PROVIDERS = [
  { id: '1', name: 'DraftKings', logo: '👑' },
  { id: '2', name: 'FanDuel', logo: '🎯' },
  { id: '3', name: 'BetMGM', logo: '🦁' },
  { id: '4', name: 'Caesars', logo: '🏛️' },
  { id: '5', name: 'PointsBet', logo: '📍' },
];

export const LinkAccount: React.FC = () => {
  const [accounts, setAccounts] = useState<SportsbookAccount[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [showLinkForm, setShowLinkForm] = useState(false);
  const [linkError, setLinkError] = useState('');
  const [linkLoading, setLinkLoading] = useState(false);

  const [selectedProvider, setSelectedProvider] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  useEffect(() => {
    loadAccounts();
  }, []);

  const loadAccounts = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await sportsbookService.getAccounts();
      setAccounts(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load accounts');
    } finally {
      setLoading(false);
    }
  };

  const handleLinkAccount = async (e: React.FormEvent) => {
    e.preventDefault();
    setLinkError('');
    setLinkLoading(true);

    try {
      await sportsbookService.linkAccount(selectedProvider, username, password);
      setShowLinkForm(false);
      setSelectedProvider('');
      setUsername('');
      setPassword('');
      loadAccounts();
    } catch (err) {
      setLinkError(handleApiError(err));
    } finally {
      setLinkLoading(false);
    }
  };

  const handleDeleteAccount = async (id: string) => {
    if (!confirm('Are you sure you want to remove this account?')) return;

    try {
      await sportsbookService.deleteAccount(id);
      loadAccounts();
    } catch (err: any) {
      setError(err.message || 'Failed to delete account');
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Link Sportsbook Accounts
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Connect your sportsbook accounts to track bets and find opportunities
          </p>
        </div>

        <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4 mb-6">
          <div className="flex items-start">
            <svg
              className="w-5 h-5 text-blue-600 dark:text-blue-400 mt-0.5 mr-3"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                clipRule="evenodd"
              />
            </svg>
            <div>
              <p className="text-sm font-medium text-blue-800 dark:text-blue-400">
                Your credentials are encrypted and secure
              </p>
              <p className="text-xs text-blue-600 dark:text-blue-500 mt-1">
                We use AES-256-GCM encryption to protect your sportsbook login information
              </p>
            </div>
          </div>
        </div>

        {error && <ErrorMessage message={error} onRetry={loadAccounts} />}

        {loading ? (
          <LoadingSpinner size="lg" text="Loading accounts..." />
        ) : (
          <>
            <div className="mb-6">
              <button
                onClick={() => setShowLinkForm(!showLinkForm)}
                className="px-6 py-3 bg-primary-600 hover:bg-primary-700 text-white font-semibold rounded-lg transition-colors"
              >
                {showLinkForm ? 'Cancel' : '+ Link New Account'}
              </button>
            </div>

            {showLinkForm && (
              <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-6 mb-6">
                <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
                  Link New Account
                </h2>

                {linkError && <ErrorMessage message={linkError} />}

                <form onSubmit={handleLinkAccount} className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Sportsbook Provider
                    </label>
                    <select
                      required
                      value={selectedProvider}
                      onChange={(e) => setSelectedProvider(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    >
                      <option value="">Select a provider...</option>
                      {SPORTSBOOK_PROVIDERS.map((provider) => (
                        <option key={provider.id} value={provider.id}>
                          {provider.logo} {provider.name}
                        </option>
                      ))}
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Username / Email
                    </label>
                    <input
                      type="text"
                      required
                      value={username}
                      onChange={(e) => setUsername(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                      placeholder="your-username"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Password
                    </label>
                    <input
                      type="password"
                      required
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                      placeholder="••••••••"
                    />
                  </div>

                  <button
                    type="submit"
                    disabled={linkLoading}
                    className="w-full px-4 py-3 bg-primary-600 hover:bg-primary-700 disabled:bg-primary-400 text-white font-semibold rounded-lg transition-colors"
                  >
                    {linkLoading ? 'Linking Account...' : 'Link Account'}
                  </button>
                </form>
              </div>
            )}

            <div>
              <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-4">
                Linked Accounts ({accounts.length})
              </h2>

              {accounts.length > 0 ? (
                <div className="grid md:grid-cols-2 gap-6">
                  {accounts.map((account) => (
                    <SportsbookAccountCard
                      key={account.id}
                      account={account}
                      onDelete={handleDeleteAccount}
                    />
                  ))}
                </div>
              ) : (
                <div className="bg-white dark:bg-dark-800 rounded-lg shadow p-12 text-center">
                  <svg
                    className="mx-auto h-12 w-12 text-gray-400"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"
                    />
                  </svg>
                  <h3 className="mt-4 text-lg font-medium text-gray-900 dark:text-white">
                    No accounts linked
                  </h3>
                  <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                    Link your sportsbook accounts to unlock arbitrage and hedge features
                  </p>
                </div>
              )}
            </div>
          </>
        )}
      </div>
    </div>
  );
};

import React, { useState, useEffect } from 'react';
import { marketService } from '../services/market.service';
import { Market } from '../types';
import { MarketCard } from '../components/MarketCard';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { ErrorMessage } from '../components/ErrorMessage';

export const Markets: React.FC = () => {
  const [markets, setMarkets] = useState<Market[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    loadMarkets();
  }, [statusFilter]);

  const loadMarkets = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await marketService.getMarkets(statusFilter || undefined);
      setMarkets(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load markets');
    } finally {
      setLoading(false);
    }
  };

  const filteredMarkets = markets.filter((market) =>
    market.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
    market.description.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Prediction Markets
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Trade on real-world events with transparent order books
          </p>
        </div>

        <div className="bg-white dark:bg-dark-800 rounded-lg shadow p-6 mb-6">
          <div className="grid md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Search Markets
              </label>
              <input
                type="text"
                placeholder="Search by title or description..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Filter by Status
              </label>
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              >
                <option value="">All Markets</option>
                <option value="open">Open</option>
                <option value="closed">Closed</option>
                <option value="resolved">Resolved</option>
              </select>
            </div>
          </div>
        </div>

        {error && <ErrorMessage message={error} onRetry={loadMarkets} />}

        {loading ? (
          <LoadingSpinner size="lg" text="Loading markets..." />
        ) : (
          <>
            <div className="mb-4 text-sm text-gray-600 dark:text-gray-400">
              Showing {filteredMarkets.length} market{filteredMarkets.length !== 1 ? 's' : ''}
            </div>

            {filteredMarkets.length > 0 ? (
              <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredMarkets.map((market) => (
                  <MarketCard key={market.id} market={market} />
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
                    d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <h3 className="mt-4 text-lg font-medium text-gray-900 dark:text-white">
                  No markets found
                </h3>
                <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                  {searchTerm
                    ? 'Try adjusting your search or filters'
                    : 'Check back later for new markets'}
                </p>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

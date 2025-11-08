import React, { useState, useEffect } from 'react';
import { sportsbookService } from '../services/sportsbook.service';
import { ArbitrageOpportunity } from '../types';
import { ArbitrageCard } from '../components/ArbitrageCard';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { ErrorMessage } from '../components/ErrorMessage';

export const Arbitrage: React.FC = () => {
  const [opportunities, setOpportunities] = useState<ArbitrageOpportunity[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [minProfit, setMinProfit] = useState(1);

  useEffect(() => {
    loadOpportunities();
  }, [minProfit]);

  const loadOpportunities = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await sportsbookService.getArbitrageOpportunities(minProfit);
      setOpportunities(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load arbitrage opportunities');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Arbitrage Opportunities
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Guaranteed profit opportunities across sportsbooks
          </p>
        </div>

        <div className="bg-white dark:bg-dark-800 rounded-lg shadow p-6 mb-6">
          <div className="max-w-md">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Minimum Profit (%)
            </label>
            <input
              type="number"
              step="0.1"
              min="0"
              value={minProfit}
              onChange={(e) => setMinProfit(parseFloat(e.target.value))}
              className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
              Filter opportunities with at least this profit percentage
            </p>
          </div>
        </div>

        {error && <ErrorMessage message={error} onRetry={loadOpportunities} />}

        {loading ? (
          <LoadingSpinner size="lg" text="Finding arbitrage opportunities..." />
        ) : (
          <>
            <div className="mb-6">
              <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg p-4">
                <div className="flex items-center">
                  <svg
                    className="w-5 h-5 text-green-600 dark:text-green-400 mr-3"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                      clipRule="evenodd"
                    />
                  </svg>
                  <div>
                    <p className="text-sm font-medium text-green-800 dark:text-green-400">
                      Found {opportunities.length} arbitrage opportunit
                      {opportunities.length !== 1 ? 'ies' : 'y'}
                    </p>
                    <p className="text-xs text-green-600 dark:text-green-500 mt-1">
                      These opportunities guarantee profit regardless of outcome
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {opportunities.length > 0 ? (
              <div className="grid lg:grid-cols-2 gap-6">
                {opportunities.map((opportunity) => (
                  <ArbitrageCard key={opportunity.id} opportunity={opportunity} />
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
                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <h3 className="mt-4 text-lg font-medium text-gray-900 dark:text-white">
                  No arbitrage opportunities
                </h3>
                <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                  Try lowering the minimum profit threshold or check back later
                </p>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

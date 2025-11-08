import React, { useState, useEffect } from 'react';
import { sportsbookService } from '../services/sportsbook.service';
import { HedgeOpportunity } from '../types';
import { HedgeCalculator } from '../components/HedgeCalculator';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { ErrorMessage } from '../components/ErrorMessage';

export const Hedges: React.FC = () => {
  const [hedges, setHedges] = useState<HedgeOpportunity[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadHedges();
  }, []);

  const loadHedges = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await sportsbookService.getHedgeOpportunities();
      setHedges(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load hedge opportunities');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
            Hedge Opportunities
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Lock in guaranteed profit on your existing bets
          </p>
        </div>

        {error && <ErrorMessage message={error} onRetry={loadHedges} />}

        {loading ? (
          <LoadingSpinner size="lg" text="Finding hedge opportunities..." />
        ) : (
          <>
            {hedges.length > 0 ? (
              <>
                <div className="mb-6 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
                  <div className="flex items-center">
                    <svg
                      className="w-5 h-5 text-blue-600 dark:text-blue-400 mr-3"
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
                        Found {hedges.length} hedge opportunit{hedges.length !== 1 ? 'ies' : 'y'}
                      </p>
                      <p className="text-xs text-blue-600 dark:text-blue-500 mt-1">
                        Place these hedge bets to guarantee profit regardless of outcome
                      </p>
                    </div>
                  </div>
                </div>

                <div className="grid lg:grid-cols-2 gap-6">
                  {hedges.map((hedge) => (
                    <HedgeCalculator key={hedge.id} hedge={hedge} />
                  ))}
                </div>
              </>
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
                    d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                  />
                </svg>
                <h3 className="mt-4 text-lg font-medium text-gray-900 dark:text-white">
                  No hedge opportunities
                </h3>
                <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                  Track some bets to see hedge opportunities
                </p>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

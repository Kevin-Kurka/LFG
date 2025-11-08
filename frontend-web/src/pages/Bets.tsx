import React, { useState, useEffect } from 'react';
import { sportsbookService } from '../services/sportsbook.service';
import { Bet } from '../types';
import { BetTracker } from '../components/BetTracker';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { ErrorMessage } from '../components/ErrorMessage';

export const Bets: React.FC = () => {
  const [bets, setBets] = useState<Bet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('');

  useEffect(() => {
    loadBets();
  }, [statusFilter]);

  const loadBets = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await sportsbookService.getBets(statusFilter || undefined);
      setBets(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load bets');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">My Bets</h1>
          <p className="text-gray-600 dark:text-gray-400">
            Track and manage all your bets across sportsbooks
          </p>
        </div>

        <div className="bg-white dark:bg-dark-800 rounded-lg shadow p-6 mb-6">
          <div className="max-w-md">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Filter by Status
            </label>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            >
              <option value="">All Bets</option>
              <option value="pending">Pending</option>
              <option value="won">Won</option>
              <option value="lost">Lost</option>
              <option value="void">Void</option>
            </select>
          </div>
        </div>

        {error && <ErrorMessage message={error} onRetry={loadBets} />}

        {loading ? (
          <LoadingSpinner size="lg" text="Loading bets..." />
        ) : (
          <BetTracker bets={bets} />
        )}
      </div>
    </div>
  );
};

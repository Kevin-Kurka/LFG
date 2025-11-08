import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { sportsbookService } from '../services/sportsbook.service';
import { SportsEvent } from '../types';
import { OddsComparison } from '../components/OddsComparison';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { ErrorMessage } from '../components/ErrorMessage';
import { formatDate } from '../utils/format';

export const EventDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [event, setEvent] = useState<SportsEvent | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    if (id) {
      loadEvent();
    }
  }, [id]);

  const loadEvent = async () => {
    try {
      setLoading(true);
      const data = await sportsbookService.getEvent(id!);
      setEvent(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load event');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <LoadingSpinner size="lg" text="Loading event..." />;
  }

  if (error || !event) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-dark-900 p-8">
        <ErrorMessage message={error || 'Event not found'} />
      </div>
    );
  }

  const marketTypes = Array.from(new Set(event.odds?.map((o) => o.market_type) || []));

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-6 mb-6">
          <div className="flex items-start justify-between mb-4">
            <div>
              <span
                className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full mb-3 ${
                  event.status === 'live'
                    ? 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400'
                    : event.status === 'upcoming'
                    ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400'
                    : 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'
                }`}
              >
                {event.status.toUpperCase()}
              </span>
              <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
                {event.home_team} vs {event.away_team}
              </h1>
              <p className="text-gray-600 dark:text-gray-400">{event.sport_name}</p>
            </div>
          </div>

          <div className="flex items-center text-sm text-gray-600 dark:text-gray-400">
            <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
              <path
                fillRule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
                clipRule="evenodd"
              />
            </svg>
            {formatDate(event.start_time)}
          </div>
        </div>

        <div className="space-y-6">
          {marketTypes.map((marketType) => {
            const oddsForMarket = event.odds?.filter((o) => o.market_type === marketType) || [];
            return (
              <OddsComparison key={marketType} odds={oddsForMarket} marketType={marketType} />
            );
          })}

          {(!event.odds || event.odds.length === 0) && (
            <div className="bg-white dark:bg-dark-800 rounded-lg shadow p-12 text-center">
              <p className="text-gray-500 dark:text-gray-400">
                No odds available for this event yet
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

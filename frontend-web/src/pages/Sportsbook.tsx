import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { sportsbookService } from '../services/sportsbook.service';
import { SportsEvent } from '../types';
import { useOddsFormat } from '../context/OddsFormatContext';
import { formatDate, formatOdds } from '../utils/format';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { ErrorMessage } from '../components/ErrorMessage';

export const Sportsbook: React.FC = () => {
  const [events, setEvents] = useState<SportsEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('upcoming');
  const { oddsFormat, setOddsFormat } = useOddsFormat();

  useEffect(() => {
    loadEvents();
  }, [statusFilter]);

  const loadEvents = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await sportsbookService.getEvents(
        undefined,
        statusFilter || undefined
      );
      setEvents(data);
    } catch (err: any) {
      setError(err.message || 'Failed to load events');
    } finally {
      setLoading(false);
    }
  };

  const getBestOdds = (event: SportsEvent, outcome: string) => {
    const oddsForOutcome = event.odds?.filter((o) => o.outcome === outcome) || [];
    if (oddsForOutcome.length === 0) return null;
    return oddsForOutcome.reduce((best, current) =>
      current.odds_decimal > best.odds_decimal ? current : best
    );
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">Sportsbook</h1>
          <p className="text-gray-600 dark:text-gray-400">
            Compare odds across all major sportsbooks
          </p>
        </div>

        <div className="bg-white dark:bg-dark-800 rounded-lg shadow p-6 mb-6">
          <div className="grid md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Status Filter
              </label>
              <select
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              >
                <option value="">All Events</option>
                <option value="upcoming">Upcoming</option>
                <option value="live">Live</option>
                <option value="finished">Finished</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                Odds Format
              </label>
              <select
                value={oddsFormat}
                onChange={(e) => setOddsFormat(e.target.value as any)}
                className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              >
                <option value="american">American</option>
                <option value="decimal">Decimal</option>
                <option value="fractional">Fractional</option>
              </select>
            </div>
          </div>
        </div>

        {error && <ErrorMessage message={error} onRetry={loadEvents} />}

        {loading ? (
          <LoadingSpinner size="lg" text="Loading events..." />
        ) : (
          <>
            <div className="mb-4 text-sm text-gray-600 dark:text-gray-400">
              Showing {events.length} event{events.length !== 1 ? 's' : ''}
            </div>

            {events.length > 0 ? (
              <div className="space-y-4">
                {events.map((event) => {
                  const homeOdds = getBestOdds(event, event.home_team);
                  const awayOdds = getBestOdds(event, event.away_team);

                  return (
                    <Link
                      key={event.id}
                      to={`/sportsbook/${event.id}`}
                      className="block bg-white dark:bg-dark-800 rounded-lg shadow hover:shadow-lg transition-shadow border border-gray-200 dark:border-dark-700 overflow-hidden"
                    >
                      <div className="p-6">
                        <div className="flex items-start justify-between mb-4">
                          <div className="flex-1">
                            <div className="flex items-center mb-2">
                              <span
                                className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full mr-3 ${
                                  event.status === 'live'
                                    ? 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400'
                                    : event.status === 'upcoming'
                                    ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400'
                                    : 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'
                                }`}
                              >
                                {event.status.toUpperCase()}
                              </span>
                              <span className="text-sm text-gray-500 dark:text-gray-400">
                                {event.sport_name}
                              </span>
                            </div>
                            <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                              {event.home_team} vs {event.away_team}
                            </h3>
                            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
                              {formatDate(event.start_time)}
                            </p>
                          </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                          <div className="p-4 bg-gray-50 dark:bg-dark-900 rounded-lg">
                            <p className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">
                              {event.home_team}
                            </p>
                            {homeOdds ? (
                              <div>
                                <p className="text-2xl font-bold text-gray-900 dark:text-white">
                                  {formatOdds(homeOdds.odds_decimal, oddsFormat)}
                                </p>
                                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                                  Best: {homeOdds.provider_name}
                                </p>
                              </div>
                            ) : (
                              <p className="text-sm text-gray-400">No odds available</p>
                            )}
                          </div>

                          <div className="p-4 bg-gray-50 dark:bg-dark-900 rounded-lg">
                            <p className="text-sm font-medium text-gray-600 dark:text-gray-400 mb-2">
                              {event.away_team}
                            </p>
                            {awayOdds ? (
                              <div>
                                <p className="text-2xl font-bold text-gray-900 dark:text-white">
                                  {formatOdds(awayOdds.odds_decimal, oddsFormat)}
                                </p>
                                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                                  Best: {awayOdds.provider_name}
                                </p>
                              </div>
                            ) : (
                              <p className="text-sm text-gray-400">No odds available</p>
                            )}
                          </div>
                        </div>

                        {event.odds && event.odds.length > 0 && (
                          <div className="mt-4 pt-4 border-t border-gray-200 dark:border-dark-700">
                            <p className="text-sm text-gray-600 dark:text-gray-400">
                              {Array.from(new Set(event.odds.map((o) => o.provider_name))).length}{' '}
                              sportsbook{event.odds.length > 1 ? 's' : ''} available • Click to
                              compare all odds
                            </p>
                          </div>
                        )}
                      </div>
                    </Link>
                  );
                })}
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
                  No events found
                </h3>
                <p className="mt-2 text-sm text-gray-500 dark:text-gray-400">
                  Check back later for upcoming events
                </p>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

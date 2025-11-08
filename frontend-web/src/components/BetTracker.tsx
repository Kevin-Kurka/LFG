import React from 'react';
import { Bet } from '../types';
import { formatCurrency, formatDate } from '../utils/format';
import { formatOdds } from '../utils/odds';
import { useOddsFormat } from '../context/OddsFormatContext';

interface BetTrackerProps {
  bets: Bet[];
}

export const BetTracker: React.FC<BetTrackerProps> = ({ bets }) => {
  const { oddsFormat } = useOddsFormat();

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'won':
        return 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400';
      case 'lost':
        return 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/20 dark:text-yellow-400';
      case 'void':
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400';
      default:
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400';
    }
  };

  const calculatePnL = (bet: Bet) => {
    if (bet.status === 'won' && bet.payout) {
      return bet.payout - bet.stake;
    } else if (bet.status === 'lost') {
      return -bet.stake;
    }
    return 0;
  };

  const totalStaked = bets.reduce((sum, bet) => sum + bet.stake, 0);
  const totalPnL = bets.reduce((sum, bet) => sum + calculatePnL(bet), 0);

  return (
    <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg overflow-hidden">
      <div className="p-4 border-b border-gray-200 dark:border-dark-700">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Bet Tracker</h3>
          <div className="flex items-center space-x-4 text-sm">
            <div>
              <span className="text-gray-500 dark:text-gray-400">Total Staked: </span>
              <span className="font-semibold text-gray-900 dark:text-white">
                {formatCurrency(totalStaked)}
              </span>
            </div>
            <div>
              <span className="text-gray-500 dark:text-gray-400">P&L: </span>
              <span
                className={`font-semibold ${
                  totalPnL >= 0
                    ? 'text-green-600 dark:text-green-400'
                    : 'text-red-600 dark:text-red-400'
                }`}
              >
                {totalPnL >= 0 ? '+' : ''}
                {formatCurrency(totalPnL)}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50 dark:bg-dark-900">
            <tr>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Event
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Provider
              </th>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Outcome
              </th>
              <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Odds
              </th>
              <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Stake
              </th>
              <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Potential
              </th>
              <th className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Status
              </th>
              <th className="px-4 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                P&L
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200 dark:divide-dark-700">
            {bets.length > 0 ? (
              bets.map((bet) => {
                const pnl = calculatePnL(bet);
                return (
                  <tr
                    key={bet.id}
                    className="hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors"
                  >
                    <td className="px-4 py-3 whitespace-nowrap">
                      <div className="text-sm font-medium text-gray-900 dark:text-white">
                        {bet.event_description}
                      </div>
                      <div className="text-xs text-gray-500 dark:text-gray-400">
                        {formatDate(bet.placed_at)}
                      </div>
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                      {bet.provider_name}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-sm text-gray-900 dark:text-white">
                      {bet.outcome}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-right text-sm font-medium text-gray-900 dark:text-white">
                      {formatOdds(bet.odds_decimal, oddsFormat)}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-right text-sm text-gray-900 dark:text-white">
                      {formatCurrency(bet.stake)}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-right text-sm text-gray-900 dark:text-white">
                      {formatCurrency(bet.potential_payout)}
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-center">
                      <span
                        className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(
                          bet.status
                        )}`}
                      >
                        {bet.status.toUpperCase()}
                      </span>
                    </td>
                    <td className="px-4 py-3 whitespace-nowrap text-right text-sm font-semibold">
                      {bet.status !== 'pending' && bet.status !== 'void' && (
                        <span
                          className={
                            pnl >= 0
                              ? 'text-green-600 dark:text-green-400'
                              : 'text-red-600 dark:text-red-400'
                          }
                        >
                          {pnl >= 0 ? '+' : ''}
                          {formatCurrency(pnl)}
                        </span>
                      )}
                    </td>
                  </tr>
                );
              })
            ) : (
              <tr>
                <td colSpan={8} className="px-4 py-8 text-center text-gray-500 dark:text-gray-400">
                  No bets tracked yet
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

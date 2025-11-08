import React from 'react';
import { EventOdds } from '../types';
import { useOddsFormat } from '../context/OddsFormatContext';
import { formatOdds } from '../utils/odds';

interface OddsComparisonProps {
  odds: EventOdds[];
  marketType: string;
}

export const OddsComparison: React.FC<OddsComparisonProps> = ({ odds, marketType }) => {
  const { oddsFormat } = useOddsFormat();

  const groupedByOutcome = odds.reduce((acc, odd) => {
    if (!acc[odd.outcome]) {
      acc[odd.outcome] = [];
    }
    acc[odd.outcome].push(odd);
    return acc;
  }, {} as Record<string, EventOdds[]>);

  const outcomes = Object.keys(groupedByOutcome);
  const providers = Array.from(new Set(odds.map((o) => o.provider_name)));

  const getBestOdds = (outcomeOdds: EventOdds[]) => {
    return outcomeOdds.reduce((best, current) =>
      current.odds_decimal > best.odds_decimal ? current : best
    );
  };

  return (
    <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg overflow-hidden">
      <div className="p-4 border-b border-gray-200 dark:border-dark-700">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
          Odds Comparison - {marketType}
        </h3>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50 dark:bg-dark-900">
            <tr>
              <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                Sportsbook
              </th>
              {outcomes.map((outcome) => (
                <th
                  key={outcome}
                  className="px-4 py-3 text-center text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                >
                  {outcome}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200 dark:divide-dark-700">
            {providers.map((provider) => (
              <tr key={provider} className="hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors">
                <td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">
                  {provider}
                </td>
                {outcomes.map((outcome) => {
                  const providerOdds = groupedByOutcome[outcome].find(
                    (o) => o.provider_name === provider
                  );
                  const bestOdds = getBestOdds(groupedByOutcome[outcome]);
                  const isBest = providerOdds?.id === bestOdds.id;

                  return (
                    <td
                      key={outcome}
                      className={`px-4 py-3 whitespace-nowrap text-center text-sm font-semibold ${
                        isBest
                          ? 'text-green-600 dark:text-green-400 bg-green-50 dark:bg-green-900/20'
                          : 'text-gray-900 dark:text-white'
                      }`}
                    >
                      {providerOdds ? (
                        <span
                          className={`inline-block px-3 py-1 rounded ${
                            isBest ? 'ring-2 ring-green-500' : ''
                          }`}
                        >
                          {formatOdds(providerOdds.odds_decimal, oddsFormat)}
                        </span>
                      ) : (
                        <span className="text-gray-400 dark:text-gray-600">-</span>
                      )}
                    </td>
                  );
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="p-4 bg-gray-50 dark:bg-dark-900 border-t border-gray-200 dark:border-dark-700">
        <p className="text-xs text-gray-500 dark:text-gray-400">
          <span className="inline-block w-4 h-4 bg-green-50 dark:bg-green-900/20 border-2 border-green-500 rounded mr-2"></span>
          Best odds highlighted in green
        </p>
      </div>
    </div>
  );
};

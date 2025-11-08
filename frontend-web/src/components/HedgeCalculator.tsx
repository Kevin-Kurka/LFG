import React from 'react';
import { HedgeOpportunity } from '../types';
import { formatCurrency, formatPercent, formatDate } from '../utils/format';
import { formatOdds } from '../utils/odds';
import { useOddsFormat } from '../context/OddsFormatContext';

interface HedgeCalculatorProps {
  hedge: HedgeOpportunity;
}

export const HedgeCalculator: React.FC<HedgeCalculatorProps> = ({ hedge }) => {
  const { oddsFormat } = useOddsFormat();

  return (
    <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg border border-blue-500 overflow-hidden">
      <div className="p-4 bg-blue-50 dark:bg-blue-900/20 border-b border-blue-200 dark:border-blue-800">
        <div className="flex items-center justify-between mb-2">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            {hedge.original_event}
          </h3>
          <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-semibold bg-blue-600 text-white">
            {formatPercent(hedge.profit_percentage)} Profit
          </span>
        </div>
      </div>

      <div className="p-4 space-y-4">
        <div className="grid grid-cols-2 gap-4">
          <div className="p-3 bg-gray-50 dark:bg-dark-900 rounded-lg">
            <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              Original Bet
            </h4>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-500 dark:text-gray-400">Outcome:</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {hedge.original_outcome}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500 dark:text-gray-400">Stake:</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {formatCurrency(hedge.original_stake)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500 dark:text-gray-400">Odds:</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {formatOdds(hedge.original_odds, oddsFormat)}
                </span>
              </div>
              <div className="flex justify-between pt-2 border-t border-gray-200 dark:border-dark-700">
                <span className="text-gray-500 dark:text-gray-400">Potential Payout:</span>
                <span className="font-semibold text-green-600 dark:text-green-400">
                  {formatCurrency(hedge.potential_payout)}
                </span>
              </div>
            </div>
          </div>

          <div className="p-3 bg-gray-50 dark:bg-dark-900 rounded-lg">
            <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
              Hedge Bet
            </h4>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-500 dark:text-gray-400">Provider:</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {hedge.hedge_provider_name}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500 dark:text-gray-400">Outcome:</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {hedge.hedge_outcome}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500 dark:text-gray-400">Required Stake:</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {formatCurrency(hedge.hedge_stake)}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-500 dark:text-gray-400">Odds:</span>
                <span className="font-medium text-gray-900 dark:text-white">
                  {formatOdds(hedge.hedge_odds, oddsFormat)}
                </span>
              </div>
            </div>
          </div>
        </div>

        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">
                Guaranteed Profit (Either Outcome)
              </p>
              <p className="text-3xl font-bold text-green-600 dark:text-green-400">
                {formatCurrency(hedge.guaranteed_profit)}
              </p>
            </div>
            <div className="text-right">
              <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">ROI</p>
              <p className="text-2xl font-semibold text-gray-900 dark:text-white">
                {formatPercent(hedge.profit_percentage)}
              </p>
            </div>
          </div>
        </div>

        <button className="w-full px-4 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition-colors">
          Place Hedge Bet on {hedge.hedge_provider_name}
        </button>
      </div>
    </div>
  );
};

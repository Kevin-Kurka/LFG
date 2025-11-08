import React, { useState } from 'react';
import { ArbitrageOpportunity } from '../types';
import { formatPercent, formatCurrency, formatDate } from '../utils/format';
import { formatOdds } from '../utils/odds';
import { useOddsFormat } from '../context/OddsFormatContext';

interface ArbitrageCardProps {
  opportunity: ArbitrageOpportunity;
}

export const ArbitrageCard: React.FC<ArbitrageCardProps> = ({ opportunity }) => {
  const { oddsFormat } = useOddsFormat();
  const [totalStake, setTotalStake] = useState(100);

  const calculateStakes = () => {
    return opportunity.legs.map((leg) => ({
      ...leg,
      stake: totalStake * leg.stake_percentage,
      potentialReturn: (totalStake * leg.stake_percentage) * leg.odds_decimal,
    }));
  };

  const stakes = calculateStakes();
  const guaranteedProfit = stakes[0].potentialReturn - totalStake;

  return (
    <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg border-2 border-green-500 overflow-hidden">
      <div className="p-4 bg-green-50 dark:bg-green-900/20 border-b border-green-200 dark:border-green-800">
        <div className="flex items-center justify-between mb-2">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            {opportunity.event_description}
          </h3>
          <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-semibold bg-green-600 text-white">
            {formatPercent(opportunity.profit_percentage)} ROI
          </span>
        </div>
        <div className="flex items-center text-sm text-gray-600 dark:text-gray-400">
          <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clipRule="evenodd" />
          </svg>
          {formatDate(opportunity.start_time)}
        </div>
      </div>

      <div className="p-4">
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            Total Stake
          </label>
          <input
            type="number"
            value={totalStake}
            onChange={(e) => setTotalStake(Number(e.target.value))}
            className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            min="1"
          />
        </div>

        <div className="space-y-3 mb-4">
          {stakes.map((leg, index) => (
            <div
              key={index}
              className="p-3 bg-gray-50 dark:bg-dark-900 rounded-lg border border-gray-200 dark:border-dark-700"
            >
              <div className="flex items-center justify-between mb-2">
                <div>
                  <p className="text-sm font-medium text-gray-900 dark:text-white">
                    {leg.provider_name}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">{leg.outcome}</p>
                </div>
                <span className="text-sm font-semibold text-primary-600 dark:text-primary-400">
                  {formatOdds(leg.odds_decimal, oddsFormat)}
                </span>
              </div>
              <div className="grid grid-cols-2 gap-2 text-sm">
                <div>
                  <span className="text-gray-500 dark:text-gray-400">Stake: </span>
                  <span className="font-medium text-gray-900 dark:text-white">
                    {formatCurrency(leg.stake)}
                  </span>
                </div>
                <div>
                  <span className="text-gray-500 dark:text-gray-400">Return: </span>
                  <span className="font-medium text-gray-900 dark:text-white">
                    {formatCurrency(leg.potentialReturn)}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>

        <div className="p-4 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-gray-600 dark:text-gray-400">Guaranteed Profit</p>
              <p className="text-2xl font-bold text-green-600 dark:text-green-400">
                {formatCurrency(guaranteedProfit)}
              </p>
            </div>
            <div className="text-right">
              <p className="text-sm text-gray-600 dark:text-gray-400">Total Return</p>
              <p className="text-xl font-semibold text-gray-900 dark:text-white">
                {formatCurrency(stakes[0].potentialReturn)}
              </p>
            </div>
          </div>
        </div>

        <button className="w-full mt-4 px-4 py-3 bg-primary-600 hover:bg-primary-700 text-white font-semibold rounded-lg transition-colors">
          Place Arbitrage Bets
        </button>
      </div>
    </div>
  );
};

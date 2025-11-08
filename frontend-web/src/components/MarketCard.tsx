import React from 'react';
import { Link } from 'react-router-dom';
import { Market } from '../types';
import { formatDate, formatPercent, formatCurrency } from '../utils/format';

interface MarketCardProps {
  market: Market;
}

export const MarketCard: React.FC<MarketCardProps> = ({ market }) => {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'open':
        return 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400';
      case 'closed':
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400';
      case 'resolved':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/20 dark:text-blue-400';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400';
    }
  };

  return (
    <Link
      to={`/markets/${market.id}`}
      className="block bg-white dark:bg-dark-800 rounded-lg shadow hover:shadow-lg transition-shadow border border-gray-200 dark:border-dark-700 overflow-hidden"
    >
      <div className="p-6">
        <div className="flex items-start justify-between mb-3">
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2 line-clamp-2">
              {market.title}
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400 line-clamp-2">
              {market.description}
            </p>
          </div>
          <span
            className={`ml-3 inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(
              market.status
            )}`}
          >
            {market.status.toUpperCase()}
          </span>
        </div>

        <div className="flex items-center justify-between text-sm text-gray-500 dark:text-gray-400 mb-4">
          <div className="flex items-center">
            <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
              <path
                fillRule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
                clipRule="evenodd"
              />
            </svg>
            {formatDate(market.end_time)}
          </div>
          <span className="px-2 py-1 bg-gray-100 dark:bg-dark-700 rounded text-xs font-medium">
            {market.category}
          </span>
        </div>

        {market.contracts && market.contracts.length > 0 && (
          <div className="grid grid-cols-2 gap-3">
            {market.contracts.map((contract) => (
              <div
                key={contract.id}
                className="p-3 bg-gray-50 dark:bg-dark-900 rounded-lg border border-gray-200 dark:border-dark-700"
              >
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                    {contract.outcome}
                  </span>
                  <span className="text-lg font-bold text-primary-600 dark:text-primary-400">
                    {formatPercent(contract.current_price * 100, 0)}
                  </span>
                </div>
                <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  Vol: {formatCurrency(contract.total_volume)}
                </div>
              </div>
            ))}
          </div>
        )}

        {market.status === 'resolved' && market.outcome && (
          <div className="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
            <p className="text-sm">
              <span className="text-gray-600 dark:text-gray-400">Resolved: </span>
              <span className="font-semibold text-blue-600 dark:text-blue-400">
                {market.outcome}
              </span>
            </p>
          </div>
        )}
      </div>
    </Link>
  );
};

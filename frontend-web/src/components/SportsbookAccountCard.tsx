import React from 'react';
import { SportsbookAccount } from '../types';
import { formatCurrency, formatTimeAgo } from '../utils/format';

interface SportsbookAccountCardProps {
  account: SportsbookAccount;
  onDelete: (id: string) => void;
}

export const SportsbookAccountCard: React.FC<SportsbookAccountCardProps> = ({
  account,
  onDelete,
}) => {
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400';
      case 'inactive':
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400';
      case 'error':
        return 'bg-red-100 text-red-800 dark:bg-red-900/20 dark:text-red-400';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400';
    }
  };

  return (
    <div className="bg-white dark:bg-dark-800 rounded-lg shadow border border-gray-200 dark:border-dark-700 p-6">
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            {account.provider_name}
          </h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">{account.username}</p>
        </div>
        <span
          className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(
            account.status
          )}`}
        >
          {account.status.toUpperCase()}
        </span>
      </div>

      <div className="space-y-2 mb-4">
        <div className="flex justify-between text-sm">
          <span className="text-gray-600 dark:text-gray-400">Balance:</span>
          <span className="font-semibold text-gray-900 dark:text-white">
            {formatCurrency(account.balance)}
          </span>
        </div>
        {account.last_synced && (
          <div className="flex justify-between text-sm">
            <span className="text-gray-600 dark:text-gray-400">Last Synced:</span>
            <span className="text-gray-900 dark:text-white">
              {formatTimeAgo(account.last_synced)}
            </span>
          </div>
        )}
      </div>

      <div className="flex space-x-2">
        <button className="flex-1 px-4 py-2 text-sm font-medium text-primary-600 dark:text-primary-400 bg-primary-50 dark:bg-primary-900/20 hover:bg-primary-100 dark:hover:bg-primary-900/30 rounded-lg transition-colors">
          Sync Now
        </button>
        <button
          onClick={() => onDelete(account.id)}
          className="px-4 py-2 text-sm font-medium text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-900/20 hover:bg-red-100 dark:hover:bg-red-900/30 rounded-lg transition-colors"
        >
          Remove
        </button>
      </div>
    </div>
  );
};

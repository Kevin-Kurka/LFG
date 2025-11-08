import React, { useState, useEffect } from 'react';
import { useWebSocket } from '../hooks/useWebSocket';
import { NotificationMessage } from '../types';

export const LiveUpdates: React.FC = () => {
  const [notifications, setNotifications] = useState<NotificationMessage[]>([]);
  const [isExpanded, setIsExpanded] = useState(false);

  useWebSocket('trade', (data) => {
    const notification: NotificationMessage = {
      type: 'trade',
      data,
      timestamp: new Date().toISOString(),
    };
    setNotifications((prev) => [notification, ...prev].slice(0, 10));
  });

  useWebSocket('arbitrage', (data) => {
    const notification: NotificationMessage = {
      type: 'arbitrage',
      data,
      timestamp: new Date().toISOString(),
    };
    setNotifications((prev) => [notification, ...prev].slice(0, 10));
  });

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'trade':
        return '💱';
      case 'arbitrage':
        return '📊';
      case 'hedge':
        return '🛡️';
      default:
        return '🔔';
    }
  };

  const formatNotificationMessage = (notification: NotificationMessage) => {
    switch (notification.type) {
      case 'trade':
        return `Trade executed: ${notification.data.quantity} @ $${notification.data.price}`;
      case 'arbitrage':
        return `New arbitrage opportunity: ${notification.data.profit_percentage}% ROI`;
      default:
        return 'New notification';
    }
  };

  if (notifications.length === 0) return null;

  return (
    <div className="fixed bottom-4 right-4 z-50 w-96">
      <div className="bg-white dark:bg-dark-800 rounded-lg shadow-2xl border border-gray-200 dark:border-dark-700 overflow-hidden">
        <div
          className="p-4 bg-primary-600 text-white cursor-pointer flex items-center justify-between"
          onClick={() => setIsExpanded(!isExpanded)}
        >
          <div className="flex items-center">
            <div className="w-2 h-2 bg-green-400 rounded-full mr-2 animate-pulse"></div>
            <span className="font-semibold">Live Updates</span>
          </div>
          <div className="flex items-center">
            <span className="text-sm mr-2">{notifications.length}</span>
            <svg
              className={`w-5 h-5 transform transition-transform ${
                isExpanded ? 'rotate-180' : ''
              }`}
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                clipRule="evenodd"
              />
            </svg>
          </div>
        </div>

        {isExpanded && (
          <div className="max-h-96 overflow-y-auto">
            {notifications.map((notification, index) => (
              <div
                key={index}
                className="p-4 border-b border-gray-200 dark:border-dark-700 hover:bg-gray-50 dark:hover:bg-dark-700 transition-colors"
              >
                <div className="flex items-start">
                  <span className="text-2xl mr-3">{getNotificationIcon(notification.type)}</span>
                  <div className="flex-1">
                    <p className="text-sm text-gray-900 dark:text-white">
                      {formatNotificationMessage(notification)}
                    </p>
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                      {new Date(notification.timestamp).toLocaleTimeString()}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

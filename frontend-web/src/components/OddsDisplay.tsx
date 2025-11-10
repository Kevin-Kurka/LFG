import React, { useState, useEffect } from 'react';
import { formatOdds } from '../utils/format';

interface OddsDisplayProps {
  odds: number;
  format: 'american' | 'decimal' | 'fractional';
  previousOdds?: number;
  providerName?: string;
  size?: 'sm' | 'md' | 'lg';
  showProvider?: boolean;
  className?: string;
}

const OddsDisplay: React.FC<OddsDisplayProps> = ({
  odds,
  format,
  previousOdds,
  providerName,
  size = 'md',
  showProvider = true,
  className = '',
}) => {
  const [showChange, setShowChange] = useState(false);
  const [change, setChange] = useState<'up' | 'down' | null>(null);

  useEffect(() => {
    if (previousOdds && previousOdds !== odds) {
      const newChange = odds > previousOdds ? 'up' : 'down';
      setChange(newChange);
      setShowChange(true);

      // Hide the indicator after 3 seconds
      const timer = setTimeout(() => {
        setShowChange(false);
      }, 3000);

      return () => clearTimeout(timer);
    }
  }, [odds, previousOdds]);

  const sizeClasses = {
    sm: 'text-base',
    md: 'text-2xl',
    lg: 'text-3xl',
  };

  const changeColors = {
    up: 'text-green-600 dark:text-green-400',
    down: 'text-red-600 dark:text-red-400',
  };

  return (
    <div className={className}>
      <div className="flex items-center gap-2">
        <p className={`${sizeClasses[size]} font-bold text-gray-900 dark:text-white`}>
          {formatOdds(odds, format)}
        </p>

        {showChange && change && (
          <div
            className={`flex items-center gap-1 animate-bounce-in ${changeColors[change]}`}
            role="status"
            aria-label={`Odds ${change === 'up' ? 'increased' : 'decreased'}`}
          >
            {change === 'up' ? (
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fillRule="evenodd"
                  d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"
                  clipRule="evenodd"
                />
              </svg>
            ) : (
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fillRule="evenodd"
                  d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"
                  clipRule="evenodd"
                />
              </svg>
            )}
            <span className="text-sm font-semibold">
              {change === 'up' ? 'UP' : 'DOWN'}
            </span>
          </div>
        )}
      </div>

      {showProvider && providerName && (
        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
          Best: {providerName}
        </p>
      )}
    </div>
  );
};

export default OddsDisplay;

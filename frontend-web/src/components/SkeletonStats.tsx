import React from 'react';

interface SkeletonStatsProps {
  count?: number;
  className?: string;
}

const SkeletonStats: React.FC<SkeletonStatsProps> = ({
  count = 3,
  className = '',
}) => {
  return (
    <div className={`grid grid-cols-1 md:grid-cols-${count} gap-6 ${className}`}>
      {Array.from({ length: count }).map((_, index) => (
        <div
          key={index}
          className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6"
        >
          <div className="animate-pulse">
            {/* Icon placeholder */}
            <div className="flex items-center justify-between mb-4">
              <div className="h-10 w-10 bg-gray-300 dark:bg-gray-600 rounded-full"></div>
              <div className="h-6 w-6 bg-gray-300 dark:bg-gray-600 rounded"></div>
            </div>

            {/* Label */}
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2 mb-3"></div>

            {/* Value */}
            <div className="h-8 bg-gray-300 dark:bg-gray-600 rounded w-3/4 mb-2"></div>

            {/* Change indicator */}
            <div className="flex items-center gap-2">
              <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-16"></div>
              <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-20"></div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default SkeletonStats;

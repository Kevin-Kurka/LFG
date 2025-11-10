import React from 'react';

interface SkeletonCardProps {
  className?: string;
}

const SkeletonCard: React.FC<SkeletonCardProps> = ({ className = '' }) => {
  return (
    <div className={`bg-white dark:bg-gray-800 rounded-lg shadow-md p-6 ${className}`}>
      <div className="animate-pulse">
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-1/3"></div>
          <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-16"></div>
        </div>

        {/* Title */}
        <div className="h-6 bg-gray-300 dark:bg-gray-600 rounded w-3/4 mb-3"></div>

        {/* Content lines */}
        <div className="space-y-2 mb-4">
          <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-full"></div>
          <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-5/6"></div>
          <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-4/6"></div>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between pt-4 border-t border-gray-200 dark:border-gray-700">
          <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-24"></div>
          <div className="h-8 bg-gray-300 dark:bg-gray-600 rounded w-20"></div>
        </div>
      </div>
    </div>
  );
};

export default SkeletonCard;

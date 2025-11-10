import React from 'react';

interface SuccessCheckmarkProps {
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

const SuccessCheckmark: React.FC<SuccessCheckmarkProps> = ({
  size = 'md',
  className = '',
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6',
  };

  return (
    <svg
      className={`${sizeClasses[size]} text-green-600 dark:text-green-400 animate-bounce-in ${className}`}
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      role="img"
      aria-label="Success"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={3}
        d="M5 13l4 4L19 7"
        className="animate-checkmark"
        style={{
          strokeDasharray: 100,
          strokeDashoffset: 100,
        }}
      />
    </svg>
  );
};

export default SuccessCheckmark;

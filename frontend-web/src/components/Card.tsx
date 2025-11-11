import React, { ReactNode } from 'react';

interface CardProps {
  children: ReactNode;
  title?: string;
  subtitle?: string;
  headerAction?: ReactNode;
  footer?: ReactNode;
  variant?: 'default' | 'bordered' | 'elevated' | 'outlined';
  padding?: 'none' | 'sm' | 'md' | 'lg';
  hover?: boolean;
  className?: string;
  onClick?: () => void;
}

const Card: React.FC<CardProps> = ({
  children,
  title,
  subtitle,
  headerAction,
  footer,
  variant = 'default',
  padding = 'md',
  hover = false,
  className = '',
  onClick,
}) => {
  const variants = {
    default: 'bg-white dark:bg-dark-800 rounded-lg shadow',
    bordered: 'bg-white dark:bg-dark-800 rounded-lg border border-gray-200 dark:border-dark-700',
    elevated: 'bg-white dark:bg-dark-800 rounded-lg shadow-lg',
    outlined: 'bg-transparent rounded-lg border-2 border-gray-300 dark:border-dark-600',
  };

  const paddings = {
    none: '',
    sm: 'p-4',
    md: 'p-6',
    lg: 'p-8',
  };

  const hoverClass = hover
    ? 'transition-all duration-200 hover:shadow-lg hover:-translate-y-0.5 cursor-pointer'
    : '';

  const clickableClass = onClick ? 'cursor-pointer' : '';

  return (
    <div
      className={`${variants[variant]} ${hoverClass} ${clickableClass} ${className}`}
      onClick={onClick}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
      onKeyPress={(e) => {
        if (onClick && (e.key === 'Enter' || e.key === ' ')) {
          e.preventDefault();
          onClick();
        }
      }}
    >
      {(title || subtitle || headerAction) && (
        <div className={`flex items-start justify-between ${padding !== 'none' ? 'pb-4' : ''} border-b border-gray-200 dark:border-dark-700`}>
          <div className={paddings[padding]}>
            {title && (
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                {title}
              </h3>
            )}
            {subtitle && (
              <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
                {subtitle}
              </p>
            )}
          </div>
          {headerAction && <div className={paddings[padding]}>{headerAction}</div>}
        </div>
      )}

      <div className={title || subtitle || headerAction ? paddings[padding] : paddings[padding]}>
        {children}
      </div>

      {footer && (
        <div className={`${paddings[padding]} pt-4 border-t border-gray-200 dark:border-dark-700`}>
          {footer}
        </div>
      )}
    </div>
  );
};

export default Card;

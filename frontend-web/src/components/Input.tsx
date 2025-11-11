import React, { InputHTMLAttributes, ReactNode } from 'react';
import SuccessCheckmark from './SuccessCheckmark';

interface InputProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> {
  label?: string;
  error?: string;
  success?: boolean;
  helperText?: string;
  leftIcon?: ReactNode;
  rightIcon?: ReactNode;
  fullWidth?: boolean;
  inputSize?: 'sm' | 'md' | 'lg';
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  (
    {
      label,
      error,
      success,
      helperText,
      leftIcon,
      rightIcon,
      fullWidth = false,
      inputSize = 'md',
      className = '',
      id,
      ...props
    },
    ref
  ) => {
    const inputId = id || `input-${Math.random().toString(36).substr(2, 9)}`;

    const sizeClasses = {
      sm: 'px-3 py-2 text-sm',
      md: 'px-4 py-3 text-base',
      lg: 'px-5 py-4 text-lg',
    };

    const getBorderColor = () => {
      if (error) return 'border-red-500 focus:ring-red-500';
      if (success) return 'border-green-500 focus:ring-green-500';
      return 'border-gray-300 dark:border-dark-600 focus:ring-primary-500';
    };

    const showSuccessCheckmark = success && !error && !rightIcon;

    return (
      <div className={fullWidth ? 'w-full' : ''}>
        {label && (
          <label
            htmlFor={inputId}
            className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
          >
            {label}
          </label>
        )}

        <div className="relative">
          {leftIcon && (
            <div className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 dark:text-gray-400">
              {leftIcon}
            </div>
          )}

          <input
            ref={ref}
            id={inputId}
            className={`
              ${fullWidth ? 'w-full' : ''}
              ${sizeClasses[inputSize]}
              ${leftIcon ? 'pl-10' : ''}
              ${rightIcon || showSuccessCheckmark ? 'pr-12' : ''}
              border rounded-lg
              bg-white dark:bg-dark-700
              text-gray-900 dark:text-white
              placeholder-gray-400 dark:placeholder-gray-500
              transition-colors
              ${getBorderColor()}
              focus:ring-2 focus:border-transparent
              disabled:bg-gray-100 dark:disabled:bg-dark-800 disabled:cursor-not-allowed
              ${className}
            `.trim().replace(/\s+/g, ' ')}
            aria-invalid={!!error}
            aria-describedby={error ? `${inputId}-error` : helperText ? `${inputId}-helper` : undefined}
            {...props}
          />

          {showSuccessCheckmark && (
            <div className="absolute right-3 top-1/2 -translate-y-1/2">
              <SuccessCheckmark />
            </div>
          )}

          {rightIcon && (
            <div className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 dark:text-gray-400">
              {rightIcon}
            </div>
          )}
        </div>

        {error && (
          <p id={`${inputId}-error`} className="mt-1 text-sm text-red-600 dark:text-red-400" role="alert">
            {error}
          </p>
        )}

        {!error && helperText && (
          <p id={`${inputId}-helper`} className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {helperText}
          </p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';

export default Input;

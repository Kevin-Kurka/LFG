import React from 'react';

interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'danger' | 'success' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  disabled?: boolean;
  fullWidth?: boolean;
  type?: 'button' | 'submit' | 'reset';
  onClick?: () => void;
  children: React.ReactNode;
  className?: string;
}

const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  loading = false,
  disabled = false,
  fullWidth = false,
  type = 'button',
  onClick,
  children,
  className = '',
}) => {
  const baseClasses = `
    relative overflow-hidden
    inline-flex items-center justify-center gap-2
    font-semibold rounded-lg
    focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2
    active:scale-95 active:brightness-90
    transition-all duration-200 ease-in-out
    disabled:cursor-not-allowed disabled:opacity-50 disabled:active:scale-100
  `;

  const variants = {
    primary: `
      bg-primary-600 hover:bg-primary-700 text-white
      focus-visible:ring-primary-500
      shadow-sm hover:shadow-md
    `,
    secondary: `
      bg-gray-200 hover:bg-gray-300 text-gray-900
      dark:bg-gray-700 dark:hover:bg-gray-600 dark:text-white
      focus-visible:ring-gray-500
      shadow-sm hover:shadow-md
    `,
    danger: `
      bg-red-600 hover:bg-red-700 text-white
      focus-visible:ring-red-500
      shadow-sm hover:shadow-md
    `,
    success: `
      bg-green-600 hover:bg-green-700 text-white
      focus-visible:ring-green-500
      shadow-sm hover:shadow-md
    `,
    ghost: `
      bg-transparent hover:bg-gray-100 text-gray-700
      dark:hover:bg-gray-800 dark:text-gray-300
      focus-visible:ring-gray-500
    `,
  };

  const sizes = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2.5 text-base',
    lg: 'px-6 py-3.5 text-lg',
  };

  const widthClass = fullWidth ? 'w-full' : '';

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled || loading}
      className={`
        ${baseClasses}
        ${variants[variant]}
        ${sizes[size]}
        ${widthClass}
        ${className}
      `.trim().replace(/\s+/g, ' ')}
    >
      {loading && (
        <svg
          className="animate-spin h-5 w-5"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          />
        </svg>
      )}
      <span className={loading ? 'opacity-0' : ''}>{children}</span>
    </button>
  );
};

export default Button;

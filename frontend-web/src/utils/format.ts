import { format, formatDistanceToNow } from 'date-fns';

export const formatCurrency = (amount: number): string => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
};

export const formatNumber = (num: number, decimals: number = 2): string => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  }).format(num);
};

export const formatPercent = (value: number, decimals: number = 2): string => {
  return `${value.toFixed(decimals)}%`;
};

export const formatDate = (date: string | Date): string => {
  return format(new Date(date), 'MMM d, yyyy h:mm a');
};

export const formatDateShort = (date: string | Date): string => {
  return format(new Date(date), 'MMM d, yyyy');
};

export const formatTimeAgo = (date: string | Date): string => {
  return formatDistanceToNow(new Date(date), { addSuffix: true });
};

export const truncateAddress = (address: string, chars: number = 6): string => {
  return `${address.slice(0, chars)}...${address.slice(-chars)}`;
};

import { OddsFormat } from '../types';

export const convertOdds = {
  americanToDecimal(american: number): number {
    if (american > 0) {
      return (american / 100) + 1;
    } else {
      return (100 / Math.abs(american)) + 1;
    }
  },

  decimalToAmerican(decimal: number): number {
    if (decimal >= 2) {
      return Math.round((decimal - 1) * 100);
    } else {
      return Math.round(-100 / (decimal - 1));
    }
  },

  decimalToFractional(decimal: number): string {
    const numerator = decimal - 1;
    const gcd = (a: number, b: number): number => (b ? gcd(b, a % b) : a);
    const denominator = 100;
    const num = Math.round(numerator * denominator);
    const divisor = gcd(num, denominator);
    return `${num / divisor}/${denominator / divisor}`;
  },

  fractionalToDecimal(fractional: string): number {
    const [num, denom] = fractional.split('/').map(Number);
    return (num / denom) + 1;
  },
};

export const formatOdds = (
  oddsDecimal: number,
  format: OddsFormat = 'american'
): string => {
  switch (format) {
    case 'american':
      const american = convertOdds.decimalToAmerican(oddsDecimal);
      return american > 0 ? `+${american}` : `${american}`;
    case 'decimal':
      return oddsDecimal.toFixed(2);
    case 'fractional':
      return convertOdds.decimalToFractional(oddsDecimal);
    default:
      return oddsDecimal.toFixed(2);
  }
};

export const calculateImpliedProbability = (oddsDecimal: number): number => {
  return (1 / oddsDecimal) * 100;
};

export const calculatePayout = (stake: number, oddsDecimal: number): number => {
  return stake * oddsDecimal;
};

export const calculateProfit = (stake: number, oddsDecimal: number): number => {
  return calculatePayout(stake, oddsDecimal) - stake;
};

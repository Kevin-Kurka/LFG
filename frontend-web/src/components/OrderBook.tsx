import React from 'react';
import { OrderBook as OrderBookType } from '../types';
import { formatNumber, formatCurrency } from '../utils/format';

interface OrderBookProps {
  orderBook: OrderBookType;
}

export const OrderBook: React.FC<OrderBookProps> = ({ orderBook }) => {
  const maxQuantity = Math.max(
    ...orderBook.bids.map((b) => b.quantity),
    ...orderBook.asks.map((a) => a.quantity)
  );

  const calculateBarWidth = (quantity: number) => {
    return (quantity / maxQuantity) * 100;
  };

  return (
    <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg overflow-hidden">
      <div className="p-4 border-b border-gray-200 dark:border-dark-700">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Order Book</h3>
      </div>

      <div className="grid grid-cols-2 divide-x divide-gray-200 dark:divide-dark-700">
        <div className="p-4">
          <div className="mb-3">
            <h4 className="text-sm font-semibold text-green-600 dark:text-green-400 mb-2">
              BIDS (Buy Orders)
            </h4>
            <div className="grid grid-cols-3 gap-2 text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">
              <div>Price</div>
              <div className="text-right">Quantity</div>
              <div className="text-right">Total</div>
            </div>
          </div>

          <div className="space-y-1">
            {orderBook.bids.length > 0 ? (
              orderBook.bids.slice(0, 15).map((bid, index) => (
                <div key={index} className="relative">
                  <div
                    className="absolute inset-0 bg-green-500/10 dark:bg-green-400/10"
                    style={{ width: `${calculateBarWidth(bid.quantity)}%` }}
                  />
                  <div className="relative grid grid-cols-3 gap-2 text-sm py-1 px-2">
                    <div className="text-green-600 dark:text-green-400 font-medium">
                      {formatCurrency(bid.price)}
                    </div>
                    <div className="text-right text-gray-900 dark:text-white">
                      {formatNumber(bid.quantity, 0)}
                    </div>
                    <div className="text-right text-gray-600 dark:text-gray-400">
                      {formatCurrency(bid.total)}
                    </div>
                  </div>
                </div>
              ))
            ) : (
              <p className="text-sm text-gray-500 dark:text-gray-400 text-center py-4">
                No buy orders
              </p>
            )}
          </div>
        </div>

        <div className="p-4">
          <div className="mb-3">
            <h4 className="text-sm font-semibold text-red-600 dark:text-red-400 mb-2">
              ASKS (Sell Orders)
            </h4>
            <div className="grid grid-cols-3 gap-2 text-xs font-medium text-gray-500 dark:text-gray-400 mb-2">
              <div>Price</div>
              <div className="text-right">Quantity</div>
              <div className="text-right">Total</div>
            </div>
          </div>

          <div className="space-y-1">
            {orderBook.asks.length > 0 ? (
              orderBook.asks.slice(0, 15).map((ask, index) => (
                <div key={index} className="relative">
                  <div
                    className="absolute inset-0 bg-red-500/10 dark:bg-red-400/10"
                    style={{ width: `${calculateBarWidth(ask.quantity)}%` }}
                  />
                  <div className="relative grid grid-cols-3 gap-2 text-sm py-1 px-2">
                    <div className="text-red-600 dark:text-red-400 font-medium">
                      {formatCurrency(ask.price)}
                    </div>
                    <div className="text-right text-gray-900 dark:text-white">
                      {formatNumber(ask.quantity, 0)}
                    </div>
                    <div className="text-right text-gray-600 dark:text-gray-400">
                      {formatCurrency(ask.total)}
                    </div>
                  </div>
                </div>
              ))
            ) : (
              <p className="text-sm text-gray-500 dark:text-gray-400 text-center py-4">
                No sell orders
              </p>
            )}
          </div>
        </div>
      </div>

      <div className="p-4 bg-gray-50 dark:bg-dark-900 border-t border-gray-200 dark:border-dark-700">
        <div className="flex justify-between text-sm">
          <div>
            <span className="text-gray-500 dark:text-gray-400">Spread: </span>
            <span className="font-medium text-gray-900 dark:text-white">
              {orderBook.asks.length > 0 && orderBook.bids.length > 0
                ? formatCurrency(orderBook.asks[0].price - orderBook.bids[0].price)
                : 'N/A'}
            </span>
          </div>
          <div>
            <span className="text-gray-500 dark:text-gray-400">Mid Price: </span>
            <span className="font-medium text-gray-900 dark:text-white">
              {orderBook.asks.length > 0 && orderBook.bids.length > 0
                ? formatCurrency((orderBook.asks[0].price + orderBook.bids[0].price) / 2)
                : 'N/A'}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
};

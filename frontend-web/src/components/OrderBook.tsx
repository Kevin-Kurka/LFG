import React, { useState, useEffect } from 'react';
import { OrderBook as OrderBookType } from '../types';
import { formatNumber, formatCurrency } from '../utils/format';
import Tooltip from './Tooltip';

interface OrderBookProps {
  orderBook: OrderBookType;
  onPriceClick?: (price: number, side: 'buy' | 'sell') => void;
  onQuantityClick?: (quantity: number) => void;
}

export const OrderBook: React.FC<OrderBookProps> = ({
  orderBook,
  onPriceClick,
  onQuantityClick
}) => {
  const [flashingBids, setFlashingBids] = useState<Set<number>>(new Set());
  const [flashingAsks, setFlashingAsks] = useState<Set<number>>(new Set());
  const [prevBids, setPrevBids] = useState(orderBook.bids);
  const [prevAsks, setPrevAsks] = useState(orderBook.asks);

  const maxQuantity = Math.max(
    ...orderBook.bids.map((b) => b.quantity),
    ...orderBook.asks.map((a) => a.quantity)
  );

  const calculateBarWidth = (quantity: number) => {
    return (quantity / maxQuantity) * 100;
  };

  // Detect new orders and trigger flash animation
  useEffect(() => {
    const newFlashingBids = new Set<number>();
    const newFlashingAsks = new Set<number>();

    // Check for new bid orders
    orderBook.bids.forEach((bid, index) => {
      const prevBid = prevBids[index];
      if (!prevBid || bid.price !== prevBid.price || bid.quantity !== prevBid.quantity) {
        newFlashingBids.add(index);
      }
    });

    // Check for new ask orders
    orderBook.asks.forEach((ask, index) => {
      const prevAsk = prevAsks[index];
      if (!prevAsk || ask.price !== prevAsk.price || ask.quantity !== prevAsk.quantity) {
        newFlashingAsks.add(index);
      }
    });

    if (newFlashingBids.size > 0) {
      setFlashingBids(newFlashingBids);
      setTimeout(() => setFlashingBids(new Set()), 500);
    }

    if (newFlashingAsks.size > 0) {
      setFlashingAsks(newFlashingAsks);
      setTimeout(() => setFlashingAsks(new Set()), 500);
    }

    setPrevBids(orderBook.bids);
    setPrevAsks(orderBook.asks);
  }, [orderBook]);

  const handleRowClick = (price: number, quantity: number, side: 'buy' | 'sell') => {
    if (onPriceClick) {
      onPriceClick(price, side);
    }
    if (onQuantityClick) {
      onQuantityClick(quantity);
    }
  };

  const renderBidRow = (bid: any, index: number) => {
    const isFlashing = flashingBids.has(index);
    const isClickable = !!onPriceClick || !!onQuantityClick;

    return (
      <Tooltip
        key={index}
        content={
          isClickable ? (
            <div className="text-xs">
              <div>Click to fill order form</div>
              <div className="text-gray-400 mt-1">
                Price: {formatCurrency(bid.price)} × Qty: {formatNumber(bid.quantity, 0)}
              </div>
            </div>
          ) : (
            `Total: ${formatCurrency(bid.total)}`
          )
        }
        position="right"
      >
        <div
          className={`relative transition-all duration-200 ${
            isClickable ? 'cursor-pointer hover:bg-green-500/20 dark:hover:bg-green-400/20' : ''
          } ${isFlashing ? 'animate-bounce-in' : ''}`}
          onClick={() => isClickable && handleRowClick(bid.price, bid.quantity, 'sell')}
          role={isClickable ? 'button' : undefined}
          tabIndex={isClickable ? 0 : undefined}
          onKeyPress={(e) => {
            if (isClickable && (e.key === 'Enter' || e.key === ' ')) {
              handleRowClick(bid.price, bid.quantity, 'sell');
            }
          }}
        >
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
      </Tooltip>
    );
  };

  const renderAskRow = (ask: any, index: number) => {
    const isFlashing = flashingAsks.has(index);
    const isClickable = !!onPriceClick || !!onQuantityClick;

    return (
      <Tooltip
        key={index}
        content={
          isClickable ? (
            <div className="text-xs">
              <div>Click to fill order form</div>
              <div className="text-gray-400 mt-1">
                Price: {formatCurrency(ask.price)} × Qty: {formatNumber(ask.quantity, 0)}
              </div>
            </div>
          ) : (
            `Total: ${formatCurrency(ask.total)}`
          )
        }
        position="left"
      >
        <div
          className={`relative transition-all duration-200 ${
            isClickable ? 'cursor-pointer hover:bg-red-500/20 dark:hover:bg-red-400/20' : ''
          } ${isFlashing ? 'animate-bounce-in' : ''}`}
          onClick={() => isClickable && handleRowClick(ask.price, ask.quantity, 'buy')}
          role={isClickable ? 'button' : undefined}
          tabIndex={isClickable ? 0 : undefined}
          onKeyPress={(e) => {
            if (isClickable && (e.key === 'Enter' || e.key === ' ')) {
              handleRowClick(ask.price, ask.quantity, 'buy');
            }
          }}
        >
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
      </Tooltip>
    );
  };

  return (
    <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg overflow-hidden">
      <div className="p-4 border-b border-gray-200 dark:border-dark-700">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Order Book</h3>
          {(onPriceClick || onQuantityClick) && (
            <span className="text-xs text-gray-500 dark:text-gray-400">
              Click orders to auto-fill
            </span>
          )}
        </div>
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
              orderBook.bids.slice(0, 15).map((bid, index) => renderBidRow(bid, index))
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
              orderBook.asks.slice(0, 15).map((ask, index) => renderAskRow(ask, index))
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

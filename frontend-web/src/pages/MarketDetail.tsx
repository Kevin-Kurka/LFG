import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { marketService } from '../services/market.service';
import { orderService } from '../services/order.service';
import { Market, OrderBook as OrderBookType } from '../types';
import { OrderBook } from '../components/OrderBook';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { ErrorMessage } from '../components/ErrorMessage';
import { formatDate, formatPercent, formatCurrency } from '../utils/format';
import { handleApiError } from '../services/api';
import { useWebSocket } from '../hooks/useWebSocket';

export const MarketDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [market, setMarket] = useState<Market | null>(null);
  const [orderBook, setOrderBook] = useState<OrderBookType | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedContract, setSelectedContract] = useState<string>('');
  const [orderType, setOrderType] = useState<'buy' | 'sell'>('buy');
  const [price, setPrice] = useState('');
  const [quantity, setQuantity] = useState('');
  const [orderLoading, setOrderLoading] = useState(false);
  const [orderError, setOrderError] = useState('');
  const [orderSuccess, setOrderSuccess] = useState('');

  useEffect(() => {
    if (id) {
      loadMarketData();
    }
  }, [id]);

  useEffect(() => {
    if (selectedContract && market) {
      loadOrderBook();
    }
  }, [selectedContract]);

  useWebSocket('market_update', (data) => {
    if (data.market_id === id) {
      loadMarketData();
    }
  });

  useWebSocket('trade', (data) => {
    if (selectedContract && data.contract_id === selectedContract) {
      loadOrderBook();
    }
  });

  const loadMarketData = async () => {
    try {
      setLoading(true);
      const data = await marketService.getMarket(id!);
      setMarket(data);
      if (data.contracts && data.contracts.length > 0 && !selectedContract) {
        setSelectedContract(data.contracts[0].id);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load market');
    } finally {
      setLoading(false);
    }
  };

  const loadOrderBook = async () => {
    try {
      const data = await marketService.getOrderBook(id!, selectedContract);
      setOrderBook(data);
    } catch (err: any) {
      console.error('Failed to load order book:', err);
    }
  };

  const handlePlaceOrder = async (e: React.FormEvent) => {
    e.preventDefault();
    setOrderError('');
    setOrderSuccess('');
    setOrderLoading(true);

    try {
      await orderService.placeOrder(
        selectedContract,
        orderType,
        parseFloat(price),
        parseInt(quantity)
      );
      setOrderSuccess('Order placed successfully!');
      setPrice('');
      setQuantity('');
      loadOrderBook();
    } catch (err) {
      setOrderError(handleApiError(err));
    } finally {
      setOrderLoading(false);
    }
  };

  if (loading) {
    return <LoadingSpinner size="lg" text="Loading market..." />;
  }

  if (error || !market) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-dark-900 p-8">
        <ErrorMessage message={error || 'Market not found'} />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-6 mb-6">
          <div className="flex items-start justify-between mb-4">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">
                {market.title}
              </h1>
              <p className="text-gray-600 dark:text-gray-400 mb-4">{market.description}</p>
            </div>
            <span
              className={`px-3 py-1 text-sm font-semibold rounded-full ${
                market.status === 'open'
                  ? 'bg-green-100 text-green-800 dark:bg-green-900/20 dark:text-green-400'
                  : 'bg-gray-100 text-gray-800 dark:bg-gray-900/20 dark:text-gray-400'
              }`}
            >
              {market.status.toUpperCase()}
            </span>
          </div>

          <div className="grid grid-cols-3 gap-4 text-sm">
            <div>
              <span className="text-gray-500 dark:text-gray-400">Category:</span>
              <span className="ml-2 font-medium text-gray-900 dark:text-white">
                {market.category}
              </span>
            </div>
            <div>
              <span className="text-gray-500 dark:text-gray-400">End Time:</span>
              <span className="ml-2 font-medium text-gray-900 dark:text-white">
                {formatDate(market.end_time)}
              </span>
            </div>
            {market.resolution_time && (
              <div>
                <span className="text-gray-500 dark:text-gray-400">Resolved:</span>
                <span className="ml-2 font-medium text-gray-900 dark:text-white">
                  {formatDate(market.resolution_time)}
                </span>
              </div>
            )}
          </div>
        </div>

        <div className="grid lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 space-y-6">
            {orderBook && <OrderBook orderBook={orderBook} />}

            {market.contracts && market.contracts.length > 0 && (
              <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-6">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                  Market Contracts
                </h3>
                <div className="grid grid-cols-2 gap-4">
                  {market.contracts.map((contract) => (
                    <button
                      key={contract.id}
                      onClick={() => setSelectedContract(contract.id)}
                      className={`p-4 rounded-lg border-2 transition-colors ${
                        selectedContract === contract.id
                          ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                          : 'border-gray-200 dark:border-dark-700 hover:border-gray-300 dark:hover:border-dark-600'
                      }`}
                    >
                      <div className="flex items-center justify-between mb-2">
                        <span className="font-semibold text-gray-900 dark:text-white">
                          {contract.outcome}
                        </span>
                        <span className="text-2xl font-bold text-primary-600 dark:text-primary-400">
                          {formatPercent(contract.current_price * 100, 0)}
                        </span>
                      </div>
                      <div className="text-sm text-gray-500 dark:text-gray-400">
                        Volume: {formatCurrency(contract.total_volume)}
                      </div>
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>

          <div>
            <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-6 sticky top-4">
              <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">
                Place Order
              </h3>

              {orderSuccess && (
                <div className="mb-4 p-3 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-lg text-sm text-green-800 dark:text-green-400">
                  {orderSuccess}
                </div>
              )}

              {orderError && <ErrorMessage message={orderError} />}

              <form onSubmit={handlePlaceOrder} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Order Type
                  </label>
                  <div className="grid grid-cols-2 gap-2">
                    <button
                      type="button"
                      onClick={() => setOrderType('buy')}
                      className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                        orderType === 'buy'
                          ? 'bg-green-600 text-white'
                          : 'bg-gray-100 dark:bg-dark-700 text-gray-700 dark:text-gray-300'
                      }`}
                    >
                      Buy
                    </button>
                    <button
                      type="button"
                      onClick={() => setOrderType('sell')}
                      className={`px-4 py-2 rounded-lg font-medium transition-colors ${
                        orderType === 'sell'
                          ? 'bg-red-600 text-white'
                          : 'bg-gray-100 dark:bg-dark-700 text-gray-700 dark:text-gray-300'
                      }`}
                    >
                      Sell
                    </button>
                  </div>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Price per Share
                  </label>
                  <input
                    type="number"
                    step="0.01"
                    min="0"
                    max="1"
                    required
                    value={price}
                    onChange={(e) => setPrice(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    placeholder="0.50"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Quantity
                  </label>
                  <input
                    type="number"
                    min="1"
                    required
                    value={quantity}
                    onChange={(e) => setQuantity(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 dark:border-dark-600 rounded-lg bg-white dark:bg-dark-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    placeholder="100"
                  />
                </div>

                {price && quantity && (
                  <div className="p-3 bg-gray-50 dark:bg-dark-900 rounded-lg">
                    <div className="flex justify-between text-sm mb-1">
                      <span className="text-gray-600 dark:text-gray-400">Total Cost:</span>
                      <span className="font-semibold text-gray-900 dark:text-white">
                        {formatCurrency(parseFloat(price) * parseInt(quantity))}
                      </span>
                    </div>
                  </div>
                )}

                <button
                  type="submit"
                  disabled={orderLoading || !selectedContract}
                  className="w-full px-4 py-3 bg-primary-600 hover:bg-primary-700 disabled:bg-primary-400 text-white font-semibold rounded-lg transition-colors"
                >
                  {orderLoading ? 'Placing Order...' : 'Place Order'}
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

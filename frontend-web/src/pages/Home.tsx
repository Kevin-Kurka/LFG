import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export const Home: React.FC = () => {
  const { user } = useAuth();

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 via-white to-blue-50 dark:from-dark-900 dark:via-dark-800 dark:to-dark-900">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
        <div className="text-center mb-16">
          <h1 className="text-6xl font-bold text-gray-900 dark:text-white mb-6">
            LFG Platform
          </h1>
          <p className="text-2xl text-gray-600 dark:text-gray-300 mb-4">
            Prediction Marketplace & Sportsbook Intelligence
          </p>
          <p className="text-lg text-gray-500 dark:text-gray-400 max-w-3xl mx-auto mb-8">
            Trade on prediction markets, compare odds across sportsbooks, discover arbitrage
            opportunities, and manage your bets all in one place
          </p>

          {!user && (
            <div className="flex justify-center space-x-4">
              <Link
                to="/register"
                className="px-8 py-4 bg-primary-600 hover:bg-primary-700 text-white text-lg font-semibold rounded-lg transition-colors shadow-lg"
              >
                Get Started
              </Link>
              <Link
                to="/login"
                className="px-8 py-4 bg-white dark:bg-dark-800 hover:bg-gray-50 dark:hover:bg-dark-700 text-gray-900 dark:text-white text-lg font-semibold rounded-lg transition-colors shadow-lg border border-gray-300 dark:border-dark-600"
              >
                Sign In
              </Link>
            </div>
          )}

          {user && (
            <Link
              to="/dashboard"
              className="inline-block px-8 py-4 bg-primary-600 hover:bg-primary-700 text-white text-lg font-semibold rounded-lg transition-colors shadow-lg"
            >
              Go to Dashboard
            </Link>
          )}
        </div>

        <div className="grid md:grid-cols-3 gap-8 mb-16">
          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-8 text-center">
            <div className="text-5xl mb-4">📊</div>
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-3">
              Prediction Markets
            </h3>
            <p className="text-gray-600 dark:text-gray-400">
              Trade on real-world events with decentralized prediction markets. Buy and sell
              outcome shares with transparent order books.
            </p>
          </div>

          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-8 text-center">
            <div className="text-5xl mb-4">⚡</div>
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-3">
              Odds Comparison
            </h3>
            <p className="text-gray-600 dark:text-gray-400">
              Compare odds across all major sportsbooks in real-time. Always find the best lines
              and maximize your value.
            </p>
          </div>

          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-8 text-center">
            <div className="text-5xl mb-4">💰</div>
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-3">
              Arbitrage Detection
            </h3>
            <p className="text-gray-600 dark:text-gray-400">
              Automatically discover arbitrage opportunities across sportsbooks. Lock in guaranteed
              profits with calculated stake recommendations.
            </p>
          </div>

          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-8 text-center">
            <div className="text-5xl mb-4">🛡️</div>
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-3">
              Hedge Opportunities
            </h3>
            <p className="text-gray-600 dark:text-gray-400">
              Find hedging opportunities for your existing bets. Calculate optimal hedge stakes to
              guarantee profit regardless of outcome.
            </p>
          </div>

          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-8 text-center">
            <div className="text-5xl mb-4">📈</div>
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-3">
              Bet Tracking
            </h3>
            <p className="text-gray-600 dark:text-gray-400">
              Track all your bets across multiple sportsbooks in one dashboard. Monitor P&L, view
              analytics, and optimize your strategy.
            </p>
          </div>

          <div className="bg-white dark:bg-dark-800 rounded-lg shadow-lg p-8 text-center">
            <div className="text-5xl mb-4">🔔</div>
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-3">
              Real-Time Updates
            </h3>
            <p className="text-gray-600 dark:text-gray-400">
              Stay informed with WebSocket-powered live updates. Get instant notifications for
              trades, odds changes, and new opportunities.
            </p>
          </div>
        </div>

        <div className="bg-white dark:bg-dark-800 rounded-lg shadow-xl p-12 text-center">
          <h2 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">
            Ready to get started?
          </h2>
          <p className="text-lg text-gray-600 dark:text-gray-400 mb-8">
            Join thousands of traders using LFG to make smarter bets
          </p>
          {!user && (
            <Link
              to="/register"
              className="inline-block px-8 py-4 bg-primary-600 hover:bg-primary-700 text-white text-lg font-semibold rounded-lg transition-colors shadow-lg"
            >
              Create Free Account
            </Link>
          )}
        </div>
      </div>
    </div>
  );
};

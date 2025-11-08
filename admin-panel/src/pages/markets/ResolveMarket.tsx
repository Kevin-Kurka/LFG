import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { marketService } from '../../services/marketService';
import { Market, Outcome } from '../../types';
import Card from '../../components/common/Card';
import Button from '../../components/common/Button';

const ResolveMarket: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [market, setMarket] = useState<Market | null>(null);
  const [outcomes, setOutcomes] = useState<Outcome[]>([]);
  const [selectedOutcome, setSelectedOutcome] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (id) {
      loadMarket();
    }
  }, [id]);

  const loadMarket = async () => {
    if (!id) return;

    setLoading(true);
    try {
      const [marketData, outcomesData] = await Promise.all([
        marketService.getMarketById(id),
        marketService.getMarketOutcomes(id),
      ]);
      setMarket(marketData);
      setOutcomes(outcomesData);
    } catch (error) {
      console.error('Error loading market:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleResolve = async () => {
    if (!id || !selectedOutcome) return;

    if (!window.confirm('Are you sure you want to resolve this market? This action cannot be undone.')) {
      return;
    }

    setSubmitting(true);
    try {
      await marketService.resolveMarket(id, selectedOutcome);
      alert('Market resolved successfully');
      navigate(`/markets/${id}`);
    } catch (error) {
      console.error('Error resolving market:', error);
      alert('Failed to resolve market');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (!market) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 dark:text-gray-400">Market not found</p>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Resolve Market</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          Select the winning outcome for: {market.title}
        </p>
      </div>

      <Card>
        <div className="space-y-6">
          <div className="p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
            <p className="text-sm text-yellow-800 dark:text-yellow-200">
              <strong>Warning:</strong> Resolving a market is irreversible. Make sure you select the correct winning outcome.
            </p>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-4">
              Select Winning Outcome
            </label>
            <div className="space-y-3">
              {outcomes.map((outcome) => (
                <div
                  key={outcome.id}
                  onClick={() => setSelectedOutcome(outcome.id)}
                  className={`p-4 rounded-lg border-2 cursor-pointer transition-colors ${
                    selectedOutcome === outcome.id
                      ? 'border-primary-500 bg-primary-50 dark:bg-primary-900/20'
                      : 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'
                  }`}
                >
                  <div className="flex items-center">
                    <input
                      type="radio"
                      checked={selectedOutcome === outcome.id}
                      onChange={() => setSelectedOutcome(outcome.id)}
                      className="h-4 w-4 text-primary-600 focus:ring-primary-500"
                    />
                    <div className="ml-3 flex-1">
                      <div className="font-medium text-gray-900 dark:text-white">
                        {outcome.name}
                      </div>
                      <div className="text-sm text-gray-500 dark:text-gray-400">
                        Current probability: {(outcome.probability * 100).toFixed(1)}% | Volume: ${outcome.total_volume.toLocaleString()}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="flex space-x-4">
            <Button
              onClick={handleResolve}
              disabled={!selectedOutcome || submitting}
              variant="success"
            >
              {submitting ? 'Resolving...' : 'Resolve Market'}
            </Button>
            <Button
              variant="secondary"
              onClick={() => navigate(`/markets/${id}`)}
            >
              Cancel
            </Button>
          </div>
        </div>
      </Card>
    </div>
  );
};

export default ResolveMarket;

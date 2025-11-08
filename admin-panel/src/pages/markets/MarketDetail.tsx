import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { marketService } from '../../services/marketService';
import { Market, Outcome } from '../../types';
import Card from '../../components/common/Card';
import Badge from '../../components/common/Badge';
import Button from '../../components/common/Button';
import { format } from 'date-fns';

const MarketDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [market, setMarket] = useState<Market | null>(null);
  const [outcomes, setOutcomes] = useState<Outcome[]>([]);
  const [loading, setLoading] = useState(true);

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

  const handleCancelMarket = async () => {
    if (!id || !window.confirm('Are you sure you want to cancel this market?')) return;

    try {
      await marketService.cancelMarket(id);
      loadMarket();
    } catch (error) {
      console.error('Error cancelling market:', error);
      alert('Failed to cancel market');
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

  const getStatusBadge = (status: string) => {
    const variants: any = {
      ACTIVE: 'success',
      CLOSED: 'warning',
      RESOLVED: 'info',
      CANCELLED: 'danger',
    };
    return <Badge variant={variants[status] || 'default'}>{status}</Badge>;
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{market.title}</h1>
          <div className="mt-2 flex items-center space-x-4">
            {getStatusBadge(market.status)}
            <span className="text-gray-600 dark:text-gray-400">{market.category}</span>
          </div>
        </div>
        <div className="flex space-x-2">
          {market.status === 'ACTIVE' && (
            <>
              <Button onClick={() => navigate(`/markets/${id}/resolve`)}>
                Resolve Market
              </Button>
              <Button variant="danger" onClick={handleCancelMarket}>
                Cancel Market
              </Button>
            </>
          )}
          <Button variant="secondary" onClick={() => navigate('/markets')}>
            Back to Markets
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <Card title="Market Details">
            <dl className="space-y-4">
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Description</dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">{market.description}</dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Resolution Source</dt>
                <dd className="mt-1 text-sm text-gray-900 dark:text-white">{market.resolution_source}</dd>
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">End Date</dt>
                  <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                    {format(new Date(market.end_date), 'PPpp')}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Created</dt>
                  <dd className="mt-1 text-sm text-gray-900 dark:text-white">
                    {format(new Date(market.created_at), 'PPpp')}
                  </dd>
                </div>
              </div>
            </dl>
          </Card>

          <Card title="Outcomes">
            <div className="space-y-3">
              {outcomes.map((outcome) => (
                <div
                  key={outcome.id}
                  className="flex justify-between items-center p-4 bg-gray-50 dark:bg-gray-700 rounded-lg"
                >
                  <div>
                    <div className="font-medium text-gray-900 dark:text-white">
                      {outcome.name}
                      {outcome.is_winning && (
                        <Badge variant="success" className="ml-2">Winner</Badge>
                      )}
                    </div>
                    <div className="text-sm text-gray-500 dark:text-gray-400">
                      Probability: {(outcome.probability * 100).toFixed(1)}%
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="font-medium text-gray-900 dark:text-white">
                      ${outcome.total_volume.toLocaleString()}
                    </div>
                    <div className="text-sm text-gray-500 dark:text-gray-400">Volume</div>
                  </div>
                </div>
              ))}
            </div>
          </Card>
        </div>

        <div className="space-y-6">
          <Card title="Statistics">
            <dl className="space-y-4">
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Total Volume</dt>
                <dd className="mt-1 text-2xl font-bold text-gray-900 dark:text-white">
                  ${(market.total_volume || 0).toLocaleString()}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500 dark:text-gray-400">Number of Outcomes</dt>
                <dd className="mt-1 text-2xl font-bold text-gray-900 dark:text-white">
                  {outcomes.length}
                </dd>
              </div>
            </dl>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default MarketDetail;

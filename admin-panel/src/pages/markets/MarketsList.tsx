import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { marketService } from '../../services/marketService';
import { Market, MarketFilters } from '../../types';
import Table from '../../components/common/Table';
import Badge from '../../components/common/Badge';
import Button from '../../components/common/Button';
import Select from '../../components/common/Select';
import Input from '../../components/common/Input';
import Pagination from '../../components/common/Pagination';
import { format } from 'date-fns';

const MarketsList: React.FC = () => {
  const navigate = useNavigate();
  const [markets, setMarkets] = useState<Market[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [filters, setFilters] = useState<MarketFilters>({});

  useEffect(() => {
    loadMarkets();
  }, [currentPage, filters]);

  const loadMarkets = async () => {
    setLoading(true);
    try {
      const response = await marketService.getMarkets(filters, currentPage, 20);
      setMarkets(response.data);
      setTotalPages(response.total_pages);
    } catch (error) {
      console.error('Error loading markets:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const variants: any = {
      ACTIVE: 'success',
      CLOSED: 'warning',
      RESOLVED: 'info',
      CANCELLED: 'danger',
    };
    return <Badge variant={variants[status] || 'default'}>{status}</Badge>;
  };

  const columns = [
    {
      key: 'title',
      label: 'Market Title',
      render: (market: Market) => (
        <div className="max-w-md">
          <div className="font-medium text-gray-900 dark:text-white">{market.title}</div>
          <div className="text-sm text-gray-500 dark:text-gray-400">{market.category}</div>
        </div>
      ),
    },
    {
      key: 'status',
      label: 'Status',
      render: (market: Market) => getStatusBadge(market.status),
    },
    {
      key: 'total_volume',
      label: 'Volume',
      render: (market: Market) => (
        <span className="font-medium">${(market.total_volume || 0).toLocaleString()}</span>
      ),
    },
    {
      key: 'end_date',
      label: 'End Date',
      render: (market: Market) => format(new Date(market.end_date), 'MMM dd, yyyy'),
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (market: Market) => (
        <div className="flex space-x-2">
          <Button
            size="sm"
            variant="secondary"
            onClick={() => navigate(`/markets/${market.id}`)}
          >
            View
          </Button>
          {market.status === 'ACTIVE' && (
            <Button
              size="sm"
              variant="primary"
              onClick={() => navigate(`/markets/${market.id}/resolve`)}
            >
              Resolve
            </Button>
          )}
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Markets</h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            Manage prediction markets
          </p>
        </div>
        <Button onClick={() => navigate('/markets/create')}>
          Create Market
        </Button>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <Input
            placeholder="Search markets..."
            value={filters.search || ''}
            onChange={(e) => setFilters({ ...filters, search: e.target.value })}
          />
          <Select
            value={filters.status || ''}
            onChange={(e) => setFilters({ ...filters, status: e.target.value })}
            options={[
              { value: '', label: 'All Statuses' },
              { value: 'ACTIVE', label: 'Active' },
              { value: 'CLOSED', label: 'Closed' },
              { value: 'RESOLVED', label: 'Resolved' },
              { value: 'CANCELLED', label: 'Cancelled' },
            ]}
          />
          <Select
            value={filters.category || ''}
            onChange={(e) => setFilters({ ...filters, category: e.target.value })}
            options={[
              { value: '', label: 'All Categories' },
              { value: 'SPORTS', label: 'Sports' },
              { value: 'POLITICS', label: 'Politics' },
              { value: 'CRYPTO', label: 'Crypto' },
              { value: 'ENTERTAINMENT', label: 'Entertainment' },
            ]}
          />
        </div>

        <Table
          columns={columns}
          data={markets}
          keyExtractor={(market) => market.id}
          loading={loading}
          emptyMessage="No markets found"
        />

        <Pagination
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={setCurrentPage}
        />
      </div>
    </div>
  );
};

export default MarketsList;

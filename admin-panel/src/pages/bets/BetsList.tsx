import React, { useEffect, useState } from 'react';
import { sportsbookService } from '../../services/sportsbookService';
import { Bet, BetFilters } from '../../types';
import Table from '../../components/common/Table';
import Badge from '../../components/common/Badge';
import Select from '../../components/common/Select';
import Pagination from '../../components/common/Pagination';
import { format } from 'date-fns';

const BetsList: React.FC = () => {
  const [bets, setBets] = useState<Bet[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [filters, setFilters] = useState<BetFilters>({});

  useEffect(() => {
    loadBets();
  }, [currentPage, filters]);

  const loadBets = async () => {
    setLoading(true);
    try {
      const response = await sportsbookService.getBets(filters, currentPage, 20);
      setBets(response.data);
      setTotalPages(response.total_pages);
    } catch (error) {
      console.error('Error loading bets:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const variants: any = {
      PENDING: 'warning',
      WON: 'success',
      LOST: 'danger',
      CANCELLED: 'default',
      VOID: 'default',
    };
    return <Badge variant={variants[status] || 'default'}>{status}</Badge>;
  };

  const columns = [
    {
      key: 'id',
      label: 'Bet ID',
      render: (bet: Bet) => (
        <span className="text-sm font-mono text-gray-700 dark:text-gray-300">
          {bet.id.slice(0, 8)}...
        </span>
      ),
    },
    {
      key: 'user',
      label: 'User',
      render: (bet: Bet) => (
        <span className="text-sm text-gray-700 dark:text-gray-300">
          {bet.user?.username || bet.user_id.slice(0, 8)}
        </span>
      ),
    },
    {
      key: 'event',
      label: 'Event',
      render: (bet: Bet) => (
        <div className="max-w-xs">
          <div className="font-medium text-gray-900 dark:text-white truncate">
            {bet.event ? `${bet.event.home_team} vs ${bet.event.away_team}` : 'N/A'}
          </div>
          <div className="text-sm text-gray-500 dark:text-gray-400">{bet.market_type}</div>
        </div>
      ),
    },
    {
      key: 'selection',
      label: 'Selection',
      render: (bet: Bet) => (
        <span className="font-medium text-gray-900 dark:text-white">{bet.selection}</span>
      ),
    },
    {
      key: 'stake',
      label: 'Stake',
      render: (bet: Bet) => (
        <span className="font-medium">${bet.stake.toFixed(2)}</span>
      ),
    },
    {
      key: 'odds',
      label: 'Odds',
      render: (bet: Bet) => (
        <span className="font-medium text-primary-600 dark:text-primary-400">
          {bet.odds.toFixed(2)}
        </span>
      ),
    },
    {
      key: 'potential_payout',
      label: 'Potential Payout',
      render: (bet: Bet) => (
        <span className="font-medium text-green-600 dark:text-green-400">
          ${bet.potential_payout.toFixed(2)}
        </span>
      ),
    },
    {
      key: 'status',
      label: 'Status',
      render: (bet: Bet) => getStatusBadge(bet.status),
    },
    {
      key: 'placed_at',
      label: 'Placed',
      render: (bet: Bet) => format(new Date(bet.placed_at), 'MMM dd, HH:mm'),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Bets Management</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          View and manage all sportsbook bets
        </p>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <Select
            value={filters.status || ''}
            onChange={(e) => setFilters({ ...filters, status: e.target.value })}
            options={[
              { value: '', label: 'All Statuses' },
              { value: 'PENDING', label: 'Pending' },
              { value: 'WON', label: 'Won' },
              { value: 'LOST', label: 'Lost' },
              { value: 'CANCELLED', label: 'Cancelled' },
              { value: 'VOID', label: 'Void' },
            ]}
          />
        </div>

        <Table
          columns={columns}
          data={bets}
          keyExtractor={(bet) => bet.id}
          loading={loading}
          emptyMessage="No bets found"
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

export default BetsList;

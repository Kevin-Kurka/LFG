import React, { useEffect, useState } from 'react';
import { arbitrageService } from '../../services/arbitrageService';
import { ArbitrageOpportunity, HedgeOpportunity } from '../../types';
import Card from '../../components/common/Card';
import Table from '../../components/common/Table';
import Badge from '../../components/common/Badge';
import Button from '../../components/common/Button';
import Select from '../../components/common/Select';
import Pagination from '../../components/common/Pagination';
import { format } from 'date-fns';

const ArbitrageList: React.FC = () => {
  const [arbitrages, setArbitrages] = useState<ArbitrageOpportunity[]>([]);
  const [hedges, setHedges] = useState<HedgeOpportunity[]>([]);
  const [loading, setLoading] = useState(true);
  const [arbPage, setArbPage] = useState(1);
  const [arbTotalPages, setArbTotalPages] = useState(1);
  const [hedgePage, setHedgePage] = useState(1);
  const [hedgeTotalPages, setHedgeTotalPages] = useState(1);
  const [statusFilter, setStatusFilter] = useState('ACTIVE');

  useEffect(() => {
    loadData();
  }, [arbPage, hedgePage, statusFilter]);

  const loadData = async () => {
    setLoading(true);
    try {
      const [arbResponse, hedgeResponse] = await Promise.all([
        arbitrageService.getArbitrageOpportunities(statusFilter, arbPage, 10),
        arbitrageService.getHedgeOpportunities(statusFilter, hedgePage, 10),
      ]);
      setArbitrages(arbResponse.data);
      setArbTotalPages(arbResponse.total_pages);
      setHedges(hedgeResponse.data);
      setHedgeTotalPages(hedgeResponse.total_pages);
    } catch (error) {
      console.error('Error loading arbitrage data:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const variants: any = {
      ACTIVE: 'success',
      EXPIRED: 'danger',
      EXECUTED: 'info',
    };
    return <Badge variant={variants[status] || 'default'}>{status}</Badge>;
  };

  const handleExecuteArbitrage = async (id: string) => {
    if (!window.confirm('Are you sure you want to execute this arbitrage opportunity?')) return;

    try {
      await arbitrageService.executeArbitrage(id);
      alert('Arbitrage executed successfully');
      loadData();
    } catch (error) {
      console.error('Error executing arbitrage:', error);
      alert('Failed to execute arbitrage');
    }
  };

  const handleExecuteHedge = async (id: string) => {
    if (!window.confirm('Are you sure you want to execute this hedge opportunity?')) return;

    try {
      await arbitrageService.executeHedge(id);
      alert('Hedge executed successfully');
      loadData();
    } catch (error) {
      console.error('Error executing hedge:', error);
      alert('Failed to execute hedge');
    }
  };

  const arbColumns = [
    {
      key: 'sport',
      label: 'Sport / Event',
      render: (arb: ArbitrageOpportunity) => (
        <div>
          <div className="font-medium text-gray-900 dark:text-white">{arb.sport}</div>
          <div className="text-sm text-gray-500 dark:text-gray-400">
            {arb.selections.length} selections
          </div>
        </div>
      ),
    },
    {
      key: 'profit',
      label: 'Profit',
      render: (arb: ArbitrageOpportunity) => (
        <div>
          <div className="font-bold text-green-600 dark:text-green-400">
            +{arb.profit_percentage.toFixed(2)}%
          </div>
          <div className="text-sm text-gray-600 dark:text-gray-400">
            ${arb.guaranteed_profit.toFixed(2)}
          </div>
        </div>
      ),
    },
    {
      key: 'stake',
      label: 'Total Stake',
      render: (arb: ArbitrageOpportunity) => (
        <span className="font-medium">${arb.total_stake.toFixed(2)}</span>
      ),
    },
    {
      key: 'status',
      label: 'Status',
      render: (arb: ArbitrageOpportunity) => getStatusBadge(arb.status),
    },
    {
      key: 'detected_at',
      label: 'Detected',
      render: (arb: ArbitrageOpportunity) => format(new Date(arb.detected_at), 'MMM dd, HH:mm'),
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (arb: ArbitrageOpportunity) => (
        arb.status === 'ACTIVE' && (
          <Button
            size="sm"
            variant="success"
            onClick={() => handleExecuteArbitrage(arb.id)}
          >
            Execute
          </Button>
        )
      ),
    },
  ];

  const hedgeColumns = [
    {
      key: 'event',
      label: 'Event',
      render: (hedge: HedgeOpportunity) => (
        <div>
          <div className="font-medium text-gray-900 dark:text-white">Event {hedge.event_id.slice(0, 8)}...</div>
          <div className="text-sm text-gray-500 dark:text-gray-400">{hedge.hedge_selection}</div>
        </div>
      ),
    },
    {
      key: 'profit',
      label: 'Guaranteed Profit',
      render: (hedge: HedgeOpportunity) => (
        <span className="font-bold text-green-600 dark:text-green-400">
          ${hedge.guaranteed_profit.toFixed(2)}
        </span>
      ),
    },
    {
      key: 'stake',
      label: 'Hedge Stake',
      render: (hedge: HedgeOpportunity) => (
        <span className="font-medium">${hedge.hedge_stake.toFixed(2)}</span>
      ),
    },
    {
      key: 'odds',
      label: 'Hedge Odds',
      render: (hedge: HedgeOpportunity) => (
        <span className="font-medium text-primary-600 dark:text-primary-400">
          {hedge.hedge_odds.toFixed(2)}
        </span>
      ),
    },
    {
      key: 'status',
      label: 'Status',
      render: (hedge: HedgeOpportunity) => getStatusBadge(hedge.status),
    },
    {
      key: 'detected_at',
      label: 'Detected',
      render: (hedge: HedgeOpportunity) => format(new Date(hedge.detected_at), 'MMM dd, HH:mm'),
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (hedge: HedgeOpportunity) => (
        hedge.status === 'ACTIVE' && (
          <Button
            size="sm"
            variant="success"
            onClick={() => handleExecuteHedge(hedge.id)}
          >
            Execute
          </Button>
        )
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Arbitrage & Hedge Opportunities</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          Monitor and execute profit opportunities
        </p>
      </div>

      <div className="flex justify-between items-center">
        <Select
          value={statusFilter}
          onChange={(e) => setStatusFilter(e.target.value)}
          options={[
            { value: 'ACTIVE', label: 'Active Only' },
            { value: '', label: 'All Statuses' },
            { value: 'EXECUTED', label: 'Executed' },
            { value: 'EXPIRED', label: 'Expired' },
          ]}
          className="w-48"
        />
      </div>

      <Card title="Arbitrage Opportunities">
        <Table
          columns={arbColumns}
          data={arbitrages}
          keyExtractor={(arb) => arb.id}
          loading={loading}
          emptyMessage="No arbitrage opportunities found"
        />
        <Pagination
          currentPage={arbPage}
          totalPages={arbTotalPages}
          onPageChange={setArbPage}
        />
      </Card>

      <Card title="Hedge Opportunities">
        <Table
          columns={hedgeColumns}
          data={hedges}
          keyExtractor={(hedge) => hedge.id}
          loading={loading}
          emptyMessage="No hedge opportunities found"
        />
        <Pagination
          currentPage={hedgePage}
          totalPages={hedgeTotalPages}
          onPageChange={setHedgePage}
        />
      </Card>
    </div>
  );
};

export default ArbitrageList;

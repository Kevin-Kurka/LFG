import React, { useEffect, useState } from 'react';
import { orderService } from '../../services/orderService';
import { Order, OrderFilters } from '../../types';
import Table from '../../components/common/Table';
import Badge from '../../components/common/Badge';
import Select from '../../components/common/Select';
import Pagination from '../../components/common/Pagination';
import { format } from 'date-fns';

const OrdersList: React.FC = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [filters, setFilters] = useState<OrderFilters>({});

  useEffect(() => {
    loadOrders();
  }, [currentPage, filters]);

  const loadOrders = async () => {
    setLoading(true);
    try {
      const response = await orderService.getOrders(filters, currentPage, 20);
      setOrders(response.data);
      setTotalPages(response.total_pages);
    } catch (error) {
      console.error('Error loading orders:', error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadge = (status: string) => {
    const variants: any = {
      OPEN: 'info',
      FILLED: 'success',
      CANCELLED: 'danger',
      PARTIALLY_FILLED: 'warning',
    };
    return <Badge variant={variants[status] || 'default'}>{status}</Badge>;
  };

  const getOrderTypeBadge = (type: string) => {
    return (
      <Badge variant={type === 'BUY' ? 'success' : 'danger'}>
        {type}
      </Badge>
    );
  };

  const columns = [
    {
      key: 'id',
      label: 'Order ID',
      render: (order: Order) => (
        <span className="text-sm font-mono text-gray-700 dark:text-gray-300">
          {order.id.slice(0, 8)}...
        </span>
      ),
    },
    {
      key: 'type',
      label: 'Type',
      render: (order: Order) => getOrderTypeBadge(order.order_type),
    },
    {
      key: 'market',
      label: 'Market',
      render: (order: Order) => (
        <div className="max-w-xs">
          <div className="font-medium text-gray-900 dark:text-white truncate">
            {order.market?.title || 'N/A'}
          </div>
        </div>
      ),
    },
    {
      key: 'quantity',
      label: 'Quantity',
      render: (order: Order) => (
        <div>
          <div className="font-medium text-gray-900 dark:text-white">{order.quantity}</div>
          {order.filled_quantity && order.filled_quantity > 0 && (
            <div className="text-sm text-gray-500 dark:text-gray-400">
              Filled: {order.filled_quantity}
            </div>
          )}
        </div>
      ),
    },
    {
      key: 'price',
      label: 'Price',
      render: (order: Order) => (
        <span className="font-medium">${order.price.toFixed(2)}</span>
      ),
    },
    {
      key: 'status',
      label: 'Status',
      render: (order: Order) => getStatusBadge(order.status),
    },
    {
      key: 'created_at',
      label: 'Created',
      render: (order: Order) => format(new Date(order.created_at), 'MMM dd, HH:mm'),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Orders & Trades</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          View and manage all orders on the platform
        </p>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <Select
            value={filters.status || ''}
            onChange={(e) => setFilters({ ...filters, status: e.target.value })}
            options={[
              { value: '', label: 'All Statuses' },
              { value: 'OPEN', label: 'Open' },
              { value: 'FILLED', label: 'Filled' },
              { value: 'CANCELLED', label: 'Cancelled' },
              { value: 'PARTIALLY_FILLED', label: 'Partially Filled' },
            ]}
          />
          <Select
            value={filters.order_type || ''}
            onChange={(e) => setFilters({ ...filters, order_type: e.target.value })}
            options={[
              { value: '', label: 'All Types' },
              { value: 'BUY', label: 'Buy' },
              { value: 'SELL', label: 'Sell' },
            ]}
          />
        </div>

        <Table
          columns={columns}
          data={orders}
          keyExtractor={(order) => order.id}
          loading={loading}
          emptyMessage="No orders found"
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

export default OrdersList;

import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { sportsbookService } from '../../services/sportsbookService';
import { SportsbookProvider, SportsEvent } from '../../types';
import Card from '../../components/common/Card';
import Badge from '../../components/common/Badge';
import Button from '../../components/common/Button';
import Table from '../../components/common/Table';
import { format } from 'date-fns';

const SportsbookDashboard: React.FC = () => {
  const navigate = useNavigate();
  const [providers, setProviders] = useState<SportsbookProvider[]>([]);
  const [events, setEvents] = useState<SportsEvent[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [providersData, eventsData] = await Promise.all([
        sportsbookService.getProviders(),
        sportsbookService.getEvents(undefined, 'LIVE', 1),
      ]);
      setProviders(providersData);
      setEvents(eventsData.data);
    } catch (error) {
      console.error('Error loading sportsbook data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleToggleProvider = async (id: string, enabled: boolean) => {
    try {
      await sportsbookService.toggleProvider(id, !enabled);
      loadData();
    } catch (error) {
      console.error('Error toggling provider:', error);
      alert('Failed to update provider status');
    }
  };

  const handleSyncProvider = async (id: string) => {
    try {
      await sportsbookService.syncProvider(id);
      alert('Provider sync initiated');
    } catch (error) {
      console.error('Error syncing provider:', error);
      alert('Failed to sync provider');
    }
  };

  const providerColumns = [
    {
      key: 'name',
      label: 'Provider',
      render: (provider: SportsbookProvider) => (
        <span className="font-medium text-gray-900 dark:text-white">{provider.name}</span>
      ),
    },
    {
      key: 'enabled',
      label: 'Status',
      render: (provider: SportsbookProvider) => (
        <Badge variant={provider.enabled ? 'success' : 'danger'}>
          {provider.enabled ? 'Active' : 'Disabled'}
        </Badge>
      ),
    },
    {
      key: 'last_sync',
      label: 'Last Sync',
      render: (provider: SportsbookProvider) =>
        provider.last_sync ? format(new Date(provider.last_sync), 'PPp') : 'Never',
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (provider: SportsbookProvider) => (
        <div className="flex space-x-2">
          <Button
            size="sm"
            variant={provider.enabled ? 'danger' : 'success'}
            onClick={() => handleToggleProvider(provider.id, provider.enabled)}
          >
            {provider.enabled ? 'Disable' : 'Enable'}
          </Button>
          <Button
            size="sm"
            variant="secondary"
            onClick={() => handleSyncProvider(provider.id)}
          >
            Sync
          </Button>
        </div>
      ),
    },
  ];

  const eventColumns = [
    {
      key: 'teams',
      label: 'Event',
      render: (event: SportsEvent) => (
        <div>
          <div className="font-medium text-gray-900 dark:text-white">
            {event.home_team} vs {event.away_team}
          </div>
          <div className="text-sm text-gray-500 dark:text-gray-400">
            {format(new Date(event.start_time), 'PPp')}
          </div>
        </div>
      ),
    },
    {
      key: 'status',
      label: 'Status',
      render: (event: SportsEvent) => {
        const variants: any = {
          SCHEDULED: 'default',
          LIVE: 'success',
          FINISHED: 'info',
          CANCELLED: 'danger',
        };
        return <Badge variant={variants[event.status]}>{event.status}</Badge>;
      },
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (event: SportsEvent) => (
        <Button
          size="sm"
          variant="secondary"
          onClick={() => navigate(`/sportsbook/events/${event.id}`)}
        >
          View
        </Button>
      ),
    },
  ];

  if (loading) {
    return (
      <div className="flex justify-center items-center h-full">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Sportsbook Management</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          Manage sportsbook providers and events
        </p>
      </div>

      <Card title="Sportsbook Providers">
        <Table
          columns={providerColumns}
          data={providers}
          keyExtractor={(provider) => provider.id}
          emptyMessage="No providers configured"
        />
      </Card>

      <Card
        title="Live Events"
        action={
          <Button size="sm" onClick={() => navigate('/sportsbook/events')}>
            View All Events
          </Button>
        }
      >
        <Table
          columns={eventColumns}
          data={events}
          keyExtractor={(event) => event.id}
          emptyMessage="No live events"
        />
      </Card>
    </div>
  );
};

export default SportsbookDashboard;

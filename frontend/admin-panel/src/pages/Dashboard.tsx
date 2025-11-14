import { useEffect, useState } from 'react';
import { api } from '../services/api';
import { DashboardStats } from '../types';

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadStats();
  }, []);

  async function loadStats() {
    try {
      setLoading(true);
      const data = await api.getDashboardStats();
      setStats(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load stats');
    } finally {
      setLoading(false);
    }
  }

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="error">
        <strong>Error:</strong> {error}
      </div>
    );
  }

  return (
    <div>
      <h2 style={{ marginBottom: '20px' }}>Dashboard</h2>

      <div className="stats-grid">
        <div className="stat-card">
          <h3>Total Markets</h3>
          <div className="value">{stats?.total_markets || 0}</div>
        </div>

        <div className="stat-card">
          <h3>Active Markets</h3>
          <div className="value">{stats?.active_markets || 0}</div>
        </div>

        <div className="stat-card">
          <h3>Total Users</h3>
          <div className="value">{stats?.total_users || 0}</div>
        </div>

        <div className="stat-card">
          <h3>Active Orders</h3>
          <div className="value">{stats?.active_orders || 0}</div>
        </div>
      </div>

      <div className="card">
        <h3 style={{ marginBottom: '16px' }}>Overview</h3>
        <p>
          Welcome to the LFG Admin Panel. This dashboard provides an overview of
          your prediction market platform's key metrics.
        </p>
        <p style={{ marginTop: '12px' }}>
          Use the navigation above to manage markets, view user activity, and
          monitor platform performance.
        </p>
      </div>
    </div>
  );
}

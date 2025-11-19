import React, { useEffect, useState } from 'react';
import { apiClient } from '../api/client';
import './Dashboard.css';

interface HealthStatus {
  status: string;
}

interface Market {
  id: string;
  question: string;
  category: string;
  end_date: string;
  status: string;
  total_volume?: number;
}

export function Dashboard() {
  const [health, setHealth] = useState<HealthStatus | null>(null);
  const [markets, setMarkets] = useState<Market[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    setLoading(true);
    setError(null);

    // Check API health
    const healthResponse = await apiClient.checkHealth();
    if (healthResponse.data) {
      setHealth(healthResponse.data);
    } else {
      setError(healthResponse.error || 'Failed to connect to API');
      setLoading(false);
      return;
    }

    // Load markets
    const marketsResponse = await apiClient.getMarkets(1, 10);
    if (marketsResponse.data) {
      setMarkets(marketsResponse.data.markets || []);
    } else {
      setError(marketsResponse.error || 'Failed to load markets');
    }

    setLoading(false);
  };

  if (loading) {
    return (
      <div className="dashboard loading">
        <div className="spinner"></div>
        <p>Loading dashboard...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="dashboard error">
        <div className="error-card">
          <h2>⚠️ Connection Error</h2>
          <p>{error}</p>
          <button onClick={loadDashboardData}>Retry</button>
        </div>
      </div>
    );
  }

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <h1>LFG Platform - Admin Dashboard</h1>
        <div className="health-status">
          <span className={`status-indicator ${health?.status === 'healthy' ? 'healthy' : 'unhealthy'}`}></span>
          <span>API: {health?.status || 'Unknown'}</span>
        </div>
      </header>

      <div className="dashboard-stats">
        <div className="stat-card">
          <h3>Active Markets</h3>
          <div className="stat-value">{markets.length}</div>
        </div>
        <div className="stat-card">
          <h3>System Status</h3>
          <div className="stat-value">✅ Online</div>
        </div>
        <div className="stat-card">
          <h3>API Gateway</h3>
          <div className="stat-value">Port 8000</div>
        </div>
        <div className="stat-card">
          <h3>Database</h3>
          <div className="stat-value">PostgreSQL</div>
        </div>
      </div>

      <div className="markets-section">
        <h2>Recent Markets</h2>
        {markets.length > 0 ? (
          <table className="markets-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Question</th>
                <th>Category</th>
                <th>Status</th>
                <th>End Date</th>
              </tr>
            </thead>
            <tbody>
              {markets.map((market) => (
                <tr key={market.id}>
                  <td><code>{market.id.substring(0, 8)}</code></td>
                  <td>{market.question}</td>
                  <td><span className="category-badge">{market.category}</span></td>
                  <td><span className={`status-badge ${market.status}`}>{market.status}</span></td>
                  <td>{new Date(market.end_date).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <div className="empty-state">
            <p>No markets found. The database might need to be seeded with test data.</p>
            <button onClick={() => window.location.reload()}>Refresh</button>
          </div>
        )}
      </div>

      <div className="services-section">
        <h2>Microservices Status</h2>
        <div className="services-grid">
          <ServiceCard name="API Gateway" port="8000" status="running" />
          <ServiceCard name="User Service" port="9080" status="running" />
          <ServiceCard name="Wallet Service" port="9081" status="running" />
          <ServiceCard name="Order Service" port="9082" status="running" />
          <ServiceCard name="Market Service" port="9083" status="running" />
          <ServiceCard name="Credit Exchange" port="9084" status="running" />
          <ServiceCard name="Notification Service" port="9085" status="running" />
          <ServiceCard name="Matching Engine" port="50051" status="running" type="gRPC" />
        </div>
      </div>
    </div>
  );
}

interface ServiceCardProps {
  name: string;
  port: string;
  status: 'running' | 'stopped' | 'error';
  type?: 'HTTP' | 'gRPC';
}

function ServiceCard({ name, port, status, type = 'HTTP' }: ServiceCardProps) {
  return (
    <div className={`service-card ${status}`}>
      <div className="service-header">
        <h4>{name}</h4>
        <span className={`service-status ${status}`}>●</span>
      </div>
      <div className="service-details">
        <span className="service-type">{type}</span>
        <span className="service-port">:{port}</span>
      </div>
    </div>
  );
}

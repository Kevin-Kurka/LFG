import { useEffect, useState } from 'react';
import { api } from '../services/api';
import { Market } from '../types';

export default function Markets() {
  const [markets, setMarkets] = useState<Market[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<string>('');
  const [search, setSearch] = useState<string>('');

  useEffect(() => {
    loadMarkets();
  }, [filter]);

  async function loadMarkets() {
    try {
      setLoading(true);
      const params: any = {};
      if (filter) params.status = filter;
      if (search) params.search = search;

      const data = await api.getMarkets(params);
      setMarkets(data.markets);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load markets');
    } finally {
      setLoading(false);
    }
  }

  function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    loadMarkets();
  }

  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  }

  function getStatusBadgeClass(status: string) {
    switch (status) {
      case 'OPEN':
        return 'badge-success';
      case 'CLOSED':
        return 'badge-warning';
      case 'RESOLVED':
        return 'badge-danger';
      default:
        return '';
    }
  }

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner" />
      </div>
    );
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h2>Markets</h2>
      </div>

      <div className="card" style={{ marginBottom: '20px' }}>
        <div style={{ display: 'flex', gap: '16px', alignItems: 'flex-end' }}>
          <form onSubmit={handleSearch} style={{ flex: 1 }}>
            <label style={{ display: 'block', marginBottom: '8px', fontSize: '14px', fontWeight: 500 }}>
              Search Markets
            </label>
            <div style={{ display: 'flex', gap: '8px' }}>
              <input
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Search by question..."
                style={{
                  flex: 1,
                  padding: '10px 12px',
                  border: '1px solid #ddd',
                  borderRadius: '4px',
                  fontSize: '14px',
                }}
              />
              <button type="submit" className="btn btn-primary">
                Search
              </button>
            </div>
          </form>

          <div>
            <label style={{ display: 'block', marginBottom: '8px', fontSize: '14px', fontWeight: 500 }}>
              Filter by Status
            </label>
            <select
              value={filter}
              onChange={(e) => setFilter(e.target.value)}
              style={{
                padding: '10px 12px',
                border: '1px solid #ddd',
                borderRadius: '4px',
                fontSize: '14px',
                minWidth: '150px',
              }}
            >
              <option value="">All</option>
              <option value="OPEN">Open</option>
              <option value="CLOSED">Closed</option>
              <option value="RESOLVED">Resolved</option>
            </select>
          </div>
        </div>
      </div>

      {error && (
        <div className="error">
          <strong>Error:</strong> {error}
          <button onClick={loadMarkets} className="btn btn-secondary" style={{ marginLeft: '16px' }}>
            Retry
          </button>
        </div>
      )}

      <div className="card">
        {markets.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '40px', color: '#666' }}>
            No markets found
          </div>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>Question</th>
                <th>Category</th>
                <th>Status</th>
                <th>Resolution Date</th>
                <th>Created</th>
                <th>Outcome</th>
              </tr>
            </thead>
            <tbody>
              {markets.map((market) => (
                <tr key={market.id}>
                  <td>
                    <strong>{market.question}</strong>
                    {market.description && (
                      <div style={{ fontSize: '13px', color: '#666', marginTop: '4px' }}>
                        {market.description.slice(0, 100)}
                        {market.description.length > 100 && '...'}
                      </div>
                    )}
                  </td>
                  <td>
                    <span className="badge" style={{ background: '#f0f0f0', color: '#666' }}>
                      {market.category}
                    </span>
                  </td>
                  <td>
                    <span className={`badge ${getStatusBadgeClass(market.status)}`}>
                      {market.status}
                    </span>
                  </td>
                  <td>{formatDate(market.resolution_date)}</td>
                  <td>{formatDate(market.created_at)}</td>
                  <td>
                    {market.outcome ? (
                      <strong style={{ color: '#155724' }}>{market.outcome}</strong>
                    ) : (
                      <span style={{ color: '#999' }}>Pending</span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <div style={{ marginTop: '16px', color: '#666', fontSize: '14px' }}>
        Showing {markets.length} market{markets.length !== 1 ? 's' : ''}
      </div>
    </div>
  );
}

import React, { useState, useEffect } from 'react';
import cryptoService, {
  CryptoDeposit,
  CryptoWithdrawal,
  ExchangeRate,
  CryptoMonitoringJob,
} from '../../services/cryptoService';

const CryptoManagement: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'overview' | 'deposits' | 'withdrawals' | 'monitoring'>('overview');
  const [deposits, setDeposits] = useState<CryptoDeposit[]>([]);
  const [withdrawals, setWithdrawals] = useState<CryptoWithdrawal[]>([]);
  const [pendingDeposits, setPendingDeposits] = useState<CryptoDeposit[]>([]);
  const [pendingWithdrawals, setPendingWithdrawals] = useState<CryptoWithdrawal[]>([]);
  const [rates, setRates] = useState<ExchangeRate[]>([]);
  const [monitoringJobs, setMonitoringJobs] = useState<CryptoMonitoringJob[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 30000); // Refresh every 30 seconds
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [depositsData, withdrawalsData, pendingDeps, pendingWiths, ratesData, monitoringData] = await Promise.all([
        cryptoService.getAllDeposits(50),
        cryptoService.getAllWithdrawals(50),
        cryptoService.getPendingDeposits(),
        cryptoService.getPendingWithdrawals(),
        cryptoService.getExchangeRates(),
        cryptoService.getMonitoringJobs(),
      ]);

      setDeposits(depositsData);
      setWithdrawals(withdrawalsData);
      setPendingDeposits(pendingDeps);
      setPendingWithdrawals(pendingWiths);
      setRates(ratesData);
      setMonitoringJobs(monitoringData);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load crypto data');
    } finally {
      setLoading(false);
    }
  };

  const handleApproveWithdrawal = async (id: string) => {
    if (!confirm('Are you sure you want to approve this withdrawal?')) return;

    try {
      await cryptoService.approveWithdrawal(id);
      setSuccess('Withdrawal approved successfully');
      await loadData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to approve withdrawal');
    }
  };

  const handleRejectWithdrawal = async (id: string) => {
    const reason = prompt('Enter rejection reason:');
    if (!reason) return;

    try {
      await cryptoService.rejectWithdrawal(id, reason);
      setSuccess('Withdrawal rejected successfully');
      await loadData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to reject withdrawal');
    }
  };

  const handleRestartMonitoring = async (currency: string) => {
    try {
      await cryptoService.restartMonitoring(currency);
      setSuccess(`Monitoring restarted for ${currency}`);
      await loadData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to restart monitoring');
    }
  };

  const calculateStats = () => {
    const totalDepositsUSD = deposits.reduce((sum, d) => sum + d.amount_credits, 0);
    const totalWithdrawalsUSD = withdrawals.reduce((sum, w) => sum + w.amount_credits, 0);
    const deposits24h = deposits.filter(
      (d) => new Date(d.created_at).getTime() > Date.now() - 24 * 60 * 60 * 1000
    ).length;
    const withdrawals24h = withdrawals.filter(
      (w) => new Date(w.created_at).getTime() > Date.now() - 24 * 60 * 60 * 1000
    ).length;

    return {
      totalDepositsUSD,
      totalWithdrawalsUSD,
      deposits24h,
      withdrawals24h,
      pendingDepositsCount: pendingDeposits.length,
      pendingWithdrawalsCount: pendingWithdrawals.length,
    };
  };

  const stats = calculateStats();

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-3xl font-bold">Cryptocurrency Management</h1>
        <p className="text-gray-600 mt-2">Monitor and manage crypto deposits, withdrawals, and exchange rates</p>
      </div>

      {/* Alert Messages */}
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
          <button onClick={() => setError('')} className="float-right font-bold">
            ×
          </button>
        </div>
      )}
      {success && (
        <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
          {success}
          <button onClick={() => setSuccess('')} className="float-right font-bold">
            ×
          </button>
        </div>
      )}

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-4 mb-6">
        <div className="bg-white rounded-lg shadow p-4">
          <div className="text-sm text-gray-600">Total Deposits</div>
          <div className="text-2xl font-bold">{deposits.length}</div>
          <div className="text-xs text-gray-500">${stats.totalDepositsUSD.toFixed(2)}</div>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="text-sm text-gray-600">Total Withdrawals</div>
          <div className="text-2xl font-bold">{withdrawals.length}</div>
          <div className="text-xs text-gray-500">${stats.totalWithdrawalsUSD.toFixed(2)}</div>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="text-sm text-gray-600">Pending Deposits</div>
          <div className="text-2xl font-bold text-yellow-600">{stats.pendingDepositsCount}</div>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="text-sm text-gray-600">Pending Withdrawals</div>
          <div className="text-2xl font-bold text-yellow-600">{stats.pendingWithdrawalsCount}</div>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="text-sm text-gray-600">Deposits (24h)</div>
          <div className="text-2xl font-bold">{stats.deposits24h}</div>
        </div>
        <div className="bg-white rounded-lg shadow p-4">
          <div className="text-sm text-gray-600">Withdrawals (24h)</div>
          <div className="text-2xl font-bold">{stats.withdrawals24h}</div>
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200 mb-6">
        <nav className="-mb-px flex space-x-8">
          {(['overview', 'deposits', 'withdrawals', 'monitoring'] as const).map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`py-4 px-1 border-b-2 font-medium text-sm capitalize ${
                activeTab === tab
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
            >
              {tab}
            </button>
          ))}
        </nav>
      </div>

      {/* Overview Tab */}
      {activeTab === 'overview' && (
        <div className="space-y-6">
          {/* Exchange Rates */}
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-xl font-semibold mb-4">Exchange Rates</h2>
            <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
              {rates.map((rate) => (
                <div key={rate.currency} className="border rounded p-4">
                  <div className="text-sm text-gray-600">{rate.currency}</div>
                  <div className="text-2xl font-bold">${rate.usd_rate.toFixed(2)}</div>
                  <div className="text-xs text-gray-500">{rate.provider}</div>
                  <div className="text-xs text-gray-400">
                    {new Date(rate.last_updated).toLocaleTimeString()}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Pending Items */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Pending Deposits */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold mb-4">
                Pending Deposits ({pendingDeposits.length})
              </h2>
              {pendingDeposits.length === 0 ? (
                <p className="text-gray-500">No pending deposits</p>
              ) : (
                <div className="space-y-2">
                  {pendingDeposits.slice(0, 5).map((deposit) => (
                    <div key={deposit.id} className="border rounded p-3">
                      <div className="flex justify-between items-start">
                        <div>
                          <div className="font-semibold">
                            {deposit.amount_crypto.toFixed(8)} {deposit.currency}
                          </div>
                          <div className="text-sm text-gray-600">
                            {deposit.confirmations}/{deposit.required_confirmations} confirmations
                          </div>
                          <div className="text-xs text-gray-500 font-mono">
                            {deposit.tx_hash.substring(0, 20)}...
                          </div>
                        </div>
                        <span className="px-2 py-1 bg-yellow-100 text-yellow-800 rounded text-xs">
                          {deposit.status}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Pending Withdrawals */}
            <div className="bg-white rounded-lg shadow p-6">
              <h2 className="text-xl font-semibold mb-4">
                Pending Withdrawals ({pendingWithdrawals.length})
              </h2>
              {pendingWithdrawals.length === 0 ? (
                <p className="text-gray-500">No pending withdrawals</p>
              ) : (
                <div className="space-y-2">
                  {pendingWithdrawals.slice(0, 5).map((withdrawal) => (
                    <div key={withdrawal.id} className="border rounded p-3">
                      <div className="flex justify-between items-start mb-2">
                        <div>
                          <div className="font-semibold">
                            {withdrawal.amount_crypto.toFixed(8)} {withdrawal.currency}
                          </div>
                          <div className="text-sm text-gray-600">
                            {withdrawal.amount_credits.toFixed(2)} Credits
                          </div>
                          <div className="text-xs text-gray-500 font-mono">
                            To: {withdrawal.to_address.substring(0, 20)}...
                          </div>
                        </div>
                        <span className="px-2 py-1 bg-yellow-100 text-yellow-800 rounded text-xs">
                          {withdrawal.status}
                        </span>
                      </div>
                      <div className="flex space-x-2">
                        <button
                          onClick={() => handleApproveWithdrawal(withdrawal.id)}
                          className="text-xs bg-green-600 text-white px-3 py-1 rounded hover:bg-green-700"
                        >
                          Approve
                        </button>
                        <button
                          onClick={() => handleRejectWithdrawal(withdrawal.id)}
                          className="text-xs bg-red-600 text-white px-3 py-1 rounded hover:bg-red-700"
                        >
                          Reject
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Deposits Tab */}
      {activeTab === 'deposits' && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">All Deposits</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left">Date</th>
                  <th className="px-4 py-2 text-left">User ID</th>
                  <th className="px-4 py-2 text-left">Currency</th>
                  <th className="px-4 py-2 text-right">Amount</th>
                  <th className="px-4 py-2 text-right">Credits</th>
                  <th className="px-4 py-2 text-left">Status</th>
                  <th className="px-4 py-2 text-left">Confirmations</th>
                  <th className="px-4 py-2 text-left">TX Hash</th>
                </tr>
              </thead>
              <tbody>
                {deposits.map((deposit) => (
                  <tr key={deposit.id} className="border-t hover:bg-gray-50">
                    <td className="px-4 py-2">{new Date(deposit.created_at).toLocaleString()}</td>
                    <td className="px-4 py-2 font-mono text-xs">{deposit.user_id.substring(0, 8)}...</td>
                    <td className="px-4 py-2">{deposit.currency}</td>
                    <td className="px-4 py-2 text-right">{deposit.amount_crypto.toFixed(8)}</td>
                    <td className="px-4 py-2 text-right">{deposit.amount_credits.toFixed(2)}</td>
                    <td className="px-4 py-2">
                      <span
                        className={`px-2 py-1 rounded text-xs ${
                          deposit.status === 'CREDITED'
                            ? 'bg-green-100 text-green-800'
                            : 'bg-yellow-100 text-yellow-800'
                        }`}
                      >
                        {deposit.status}
                      </span>
                    </td>
                    <td className="px-4 py-2">
                      {deposit.confirmations}/{deposit.required_confirmations}
                    </td>
                    <td className="px-4 py-2 font-mono text-xs">{deposit.tx_hash.substring(0, 16)}...</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Withdrawals Tab */}
      {activeTab === 'withdrawals' && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">All Withdrawals</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left">Date</th>
                  <th className="px-4 py-2 text-left">User ID</th>
                  <th className="px-4 py-2 text-left">Currency</th>
                  <th className="px-4 py-2 text-right">Amount</th>
                  <th className="px-4 py-2 text-right">Credits</th>
                  <th className="px-4 py-2 text-left">Status</th>
                  <th className="px-4 py-2 text-left">To Address</th>
                  <th className="px-4 py-2 text-left">TX Hash</th>
                </tr>
              </thead>
              <tbody>
                {withdrawals.map((withdrawal) => (
                  <tr key={withdrawal.id} className="border-t hover:bg-gray-50">
                    <td className="px-4 py-2">{new Date(withdrawal.created_at).toLocaleString()}</td>
                    <td className="px-4 py-2 font-mono text-xs">{withdrawal.user_id.substring(0, 8)}...</td>
                    <td className="px-4 py-2">{withdrawal.currency}</td>
                    <td className="px-4 py-2 text-right">{withdrawal.amount_crypto.toFixed(8)}</td>
                    <td className="px-4 py-2 text-right">{withdrawal.amount_credits.toFixed(2)}</td>
                    <td className="px-4 py-2">
                      <span
                        className={`px-2 py-1 rounded text-xs ${
                          withdrawal.status === 'CONFIRMED' || withdrawal.status === 'SENT'
                            ? 'bg-green-100 text-green-800'
                            : withdrawal.status === 'FAILED'
                            ? 'bg-red-100 text-red-800'
                            : 'bg-yellow-100 text-yellow-800'
                        }`}
                      >
                        {withdrawal.status}
                      </span>
                    </td>
                    <td className="px-4 py-2 font-mono text-xs">{withdrawal.to_address.substring(0, 16)}...</td>
                    <td className="px-4 py-2 font-mono text-xs">
                      {withdrawal.tx_hash ? withdrawal.tx_hash.substring(0, 16) + '...' : 'N/A'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Monitoring Tab */}
      {activeTab === 'monitoring' && (
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold mb-4">Blockchain Monitoring</h2>
          <div className="space-y-4">
            {monitoringJobs.map((job) => (
              <div key={job.id} className="border rounded p-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h3 className="font-semibold text-lg">{job.currency}</h3>
                    <div className="text-sm text-gray-600 mt-1">
                      Last scanned block: {job.last_scanned_block}
                    </div>
                    <div className="text-sm text-gray-600">
                      Last scan: {new Date(job.last_scan_at).toLocaleString()}
                    </div>
                    {job.last_error && (
                      <div className="text-sm text-red-600 mt-1">
                        Error: {job.last_error} (Count: {job.error_count})
                      </div>
                    )}
                  </div>
                  <div className="flex items-center space-x-3">
                    <span
                      className={`px-3 py-1 rounded text-sm ${
                        job.is_running ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                      }`}
                    >
                      {job.is_running ? 'Running' : 'Stopped'}
                    </span>
                    <button
                      onClick={() => handleRestartMonitoring(job.currency)}
                      className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 text-sm"
                    >
                      Restart
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default CryptoManagement;

import React, { useState, useEffect } from 'react';
import cryptoService, { CryptoWallet, CryptoDeposit, CryptoWithdrawal, ExchangeRate } from '../services/crypto.service';

const SUPPORTED_CURRENCIES = ['BTC', 'ETH', 'USDC', 'USDT', 'LTC'];

const CURRENCY_NAMES: Record<string, string> = {
  BTC: 'Bitcoin',
  ETH: 'Ethereum',
  USDC: 'USD Coin',
  USDT: 'Tether',
  LTC: 'Litecoin',
};

const Crypto: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'wallets' | 'deposit' | 'withdraw' | 'history'>('wallets');
  const [wallets, setWallets] = useState<CryptoWallet[]>([]);
  const [deposits, setDeposits] = useState<CryptoDeposit[]>([]);
  const [withdrawals, setWithdrawals] = useState<CryptoWithdrawal[]>([]);
  const [rates, setRates] = useState<Record<string, ExchangeRate>>({});
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // Deposit form
  const [selectedDepositCurrency, setSelectedDepositCurrency] = useState('BTC');

  // Withdrawal form
  const [withdrawCurrency, setWithdrawCurrency] = useState('BTC');
  const [withdrawAmount, setWithdrawAmount] = useState('');
  const [withdrawAddress, setWithdrawAddress] = useState('');
  const [withdrawPreview, setWithdrawPreview] = useState<any>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      const [walletsData, depositsData, withdrawalsData, ratesData] = await Promise.all([
        cryptoService.getWallets(),
        cryptoService.getDeposits(),
        cryptoService.getWithdrawals(),
        cryptoService.getExchangeRates(),
      ]);
      setWallets(walletsData);
      setDeposits(depositsData);
      setWithdrawals(withdrawalsData);
      setRates(ratesData);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load crypto data');
    } finally {
      setLoading(false);
    }
  };

  const createWallet = async (currency: string) => {
    try {
      setLoading(true);
      setError('');
      const result = await cryptoService.createWallet(currency);
      setSuccess(`${currency} wallet created successfully!`);
      if (result.mnemonic) {
        alert(`IMPORTANT: Save your recovery phrase:\n\n${result.mnemonic}\n\nThis will only be shown once!`);
      }
      await loadData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create wallet');
    } finally {
      setLoading(false);
    }
  };

  const handleWithdraw = async () => {
    if (!withdrawAmount || !withdrawAddress) {
      setError('Please fill in all fields');
      return;
    }

    try {
      setLoading(true);
      setError('');
      await cryptoService.requestWithdrawal(withdrawCurrency, parseFloat(withdrawAmount), withdrawAddress);
      setSuccess('Withdrawal request submitted successfully!');
      setWithdrawAmount('');
      setWithdrawAddress('');
      setWithdrawPreview(null);
      await loadData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to request withdrawal');
    } finally {
      setLoading(false);
    }
  };

  const previewWithdrawal = async () => {
    if (!withdrawAmount || parseFloat(withdrawAmount) <= 0) {
      setError('Please enter a valid amount');
      return;
    }

    try {
      const conversion = await cryptoService.convertAmount(
        withdrawCurrency,
        parseFloat(withdrawAmount),
        'credits_to_crypto'
      );
      setWithdrawPreview(conversion);
      setError('');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to calculate conversion');
    }
  };

  const simulateDeposit = async () => {
    if (!selectedDepositCurrency) return;

    const amount = prompt(`Enter ${selectedDepositCurrency} amount to simulate:`);
    if (!amount) return;

    try {
      setLoading(true);
      setError('');
      await cryptoService.simulateDeposit(selectedDepositCurrency, parseFloat(amount));
      setSuccess(`Test deposit of ${amount} ${selectedDepositCurrency} created!`);
      await loadData();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to simulate deposit');
    } finally {
      setLoading(false);
    }
  };

  const getWalletForCurrency = (currency: string): CryptoWallet | undefined => {
    return wallets.find((w) => w.currency === currency);
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    setSuccess('Copied to clipboard!');
    setTimeout(() => setSuccess(''), 2000);
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-6">Cryptocurrency Wallet</h1>

      {/* Alert Messages */}
      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          {error}
          <button onClick={() => setError('')} className="float-right font-bold">×</button>
        </div>
      )}
      {success && (
        <div className="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded mb-4">
          {success}
          <button onClick={() => setSuccess('')} className="float-right font-bold">×</button>
        </div>
      )}

      {/* Exchange Rates */}
      <div className="bg-white rounded-lg shadow-md p-6 mb-6">
        <h2 className="text-xl font-semibold mb-4">Current Exchange Rates</h2>
        <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
          {SUPPORTED_CURRENCIES.map((currency) => (
            <div key={currency} className="text-center p-4 bg-gray-50 rounded">
              <div className="text-sm text-gray-600">{currency}</div>
              <div className="text-lg font-bold">
                ${rates[currency]?.usd_rate?.toFixed(currency === 'BTC' ? 0 : 2) || '...'}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200 mb-6">
        <nav className="-mb-px flex space-x-8">
          {(['wallets', 'deposit', 'withdraw', 'history'] as const).map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`py-4 px-1 border-b-2 font-medium text-sm capitalize ${
                activeTab === tab
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`}
            >
              {tab}
            </button>
          ))}
        </nav>
      </div>

      {/* Wallets Tab */}
      {activeTab === 'wallets' && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-semibold mb-4">Your Crypto Wallets</h2>
          <div className="space-y-4">
            {SUPPORTED_CURRENCIES.map((currency) => {
              const wallet = getWalletForCurrency(currency);
              return (
                <div key={currency} className="border rounded-lg p-4">
                  <div className="flex justify-between items-center">
                    <div>
                      <h3 className="text-lg font-semibold">
                        {CURRENCY_NAMES[currency]} ({currency})
                      </h3>
                      {wallet ? (
                        <>
                          <p className="text-sm text-gray-600 mt-1">
                            Address:{' '}
                            <code className="bg-gray-100 px-2 py-1 rounded text-xs">
                              {wallet.address}
                            </code>
                            <button
                              onClick={() => copyToClipboard(wallet.address)}
                              className="ml-2 text-blue-600 text-xs hover:underline"
                            >
                              Copy
                            </button>
                          </p>
                          <p className="text-sm text-gray-600 mt-1">
                            Balance: {wallet.balance_crypto.toFixed(8)} {currency}
                          </p>
                        </>
                      ) : (
                        <p className="text-sm text-gray-500 mt-1">No wallet created</p>
                      )}
                    </div>
                    <div>
                      {!wallet ? (
                        <button
                          onClick={() => createWallet(currency)}
                          disabled={loading}
                          className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:bg-gray-400"
                        >
                          Create Wallet
                        </button>
                      ) : (
                        <span className="text-green-600 font-semibold">✓ Active</span>
                      )}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Deposit Tab */}
      {activeTab === 'deposit' && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-semibold mb-4">Deposit Crypto</h2>
          <div className="mb-6">
            <label className="block text-sm font-medium mb-2">Select Currency</label>
            <select
              value={selectedDepositCurrency}
              onChange={(e) => setSelectedDepositCurrency(e.target.value)}
              className="w-full border rounded px-3 py-2"
            >
              {SUPPORTED_CURRENCIES.map((currency) => (
                <option key={currency} value={currency}>
                  {CURRENCY_NAMES[currency]} ({currency})
                </option>
              ))}
            </select>
          </div>

          {(() => {
            const wallet = getWalletForCurrency(selectedDepositCurrency);
            if (!wallet) {
              return (
                <div className="text-center py-8">
                  <p className="text-gray-600 mb-4">You need to create a wallet first</p>
                  <button
                    onClick={() => createWallet(selectedDepositCurrency)}
                    disabled={loading}
                    className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700"
                  >
                    Create {selectedDepositCurrency} Wallet
                  </button>
                </div>
              );
            }

            return (
              <div>
                <div className="bg-gray-50 p-6 rounded-lg mb-4">
                  <h3 className="font-semibold mb-2">Deposit Address</h3>
                  <div className="flex items-center space-x-2">
                    <code className="bg-white px-4 py-2 rounded flex-1 font-mono text-sm">
                      {wallet.address}
                    </code>
                    <button
                      onClick={() => copyToClipboard(wallet.address)}
                      className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700"
                    >
                      Copy
                    </button>
                  </div>
                  <p className="text-sm text-gray-600 mt-2">
                    Send {selectedDepositCurrency} to this address. Funds will be credited after confirmations.
                  </p>
                  <p className="text-sm text-yellow-600 mt-1">
                    ⚠️ Only send {selectedDepositCurrency} to this address. Other assets will be lost.
                  </p>
                </div>

                {/* Simulate Deposit for Testing */}
                <div className="border-t pt-4">
                  <p className="text-sm text-gray-600 mb-2">For testing purposes:</p>
                  <button
                    onClick={simulateDeposit}
                    className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700"
                  >
                    Simulate Deposit (Dev Only)
                  </button>
                </div>

                {/* Pending Deposits */}
                {deposits.filter(d => d.currency === selectedDepositCurrency && d.status !== 'CREDITED').length > 0 && (
                  <div className="mt-6">
                    <h3 className="font-semibold mb-2">Pending Deposits</h3>
                    {deposits
                      .filter(d => d.currency === selectedDepositCurrency && d.status !== 'CREDITED')
                      .map((deposit) => (
                        <div key={deposit.id} className="bg-yellow-50 p-4 rounded mb-2">
                          <p className="text-sm">
                            {deposit.amount_crypto} {deposit.currency} -{' '}
                            {deposit.confirmations}/{deposit.required_confirmations} confirmations
                          </p>
                          <p className="text-xs text-gray-600">TX: {deposit.tx_hash}</p>
                        </div>
                      ))}
                  </div>
                )}
              </div>
            );
          })()}
        </div>
      )}

      {/* Withdraw Tab */}
      {activeTab === 'withdraw' && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-semibold mb-4">Withdraw Crypto</h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2">Currency</label>
              <select
                value={withdrawCurrency}
                onChange={(e) => {
                  setWithdrawCurrency(e.target.value);
                  setWithdrawPreview(null);
                }}
                className="w-full border rounded px-3 py-2"
              >
                {SUPPORTED_CURRENCIES.map((currency) => (
                  <option key={currency} value={currency}>
                    {CURRENCY_NAMES[currency]} ({currency})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">Amount (Credits)</label>
              <input
                type="number"
                value={withdrawAmount}
                onChange={(e) => {
                  setWithdrawAmount(e.target.value);
                  setWithdrawPreview(null);
                }}
                placeholder="0.00"
                className="w-full border rounded px-3 py-2"
                step="0.01"
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-2">Destination Address</label>
              <input
                type="text"
                value={withdrawAddress}
                onChange={(e) => setWithdrawAddress(e.target.value)}
                placeholder={`Enter ${withdrawCurrency} address`}
                className="w-full border rounded px-3 py-2 font-mono text-sm"
              />
            </div>

            <button
              onClick={previewWithdrawal}
              disabled={!withdrawAmount || loading}
              className="w-full bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:bg-gray-400"
            >
              Calculate Withdrawal
            </button>

            {withdrawPreview && (
              <div className="bg-gray-50 p-4 rounded">
                <h3 className="font-semibold mb-2">Withdrawal Preview</h3>
                <div className="space-y-1 text-sm">
                  <div className="flex justify-between">
                    <span>Amount:</span>
                    <span>{withdrawAmount} Credits</span>
                  </div>
                  <div className="flex justify-between">
                    <span>You will receive:</span>
                    <span className="font-semibold">
                      {withdrawPreview.result.toFixed(8)} {withdrawCurrency}
                    </span>
                  </div>
                  <div className="flex justify-between text-gray-600">
                    <span>Exchange rate:</span>
                    <span>1 {withdrawCurrency} = ${withdrawPreview.exchange_rate.toFixed(2)}</span>
                  </div>
                  <div className="flex justify-between text-gray-600">
                    <span>Fees (1% + network):</span>
                    <span>~{(parseFloat(withdrawAmount) * 0.01).toFixed(2)} Credits</span>
                  </div>
                </div>

                <button
                  onClick={handleWithdraw}
                  disabled={loading || !withdrawAddress}
                  className="w-full mt-4 bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 disabled:bg-gray-400"
                >
                  Confirm Withdrawal
                </button>
              </div>
            )}
          </div>
        </div>
      )}

      {/* History Tab */}
      {activeTab === 'history' && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-semibold mb-4">Transaction History</h2>

          {/* Deposits */}
          <div className="mb-6">
            <h3 className="text-lg font-semibold mb-3">Deposits</h3>
            {deposits.length === 0 ? (
              <p className="text-gray-500">No deposits yet</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-4 py-2 text-left">Date</th>
                      <th className="px-4 py-2 text-left">Currency</th>
                      <th className="px-4 py-2 text-right">Amount</th>
                      <th className="px-4 py-2 text-right">Credits</th>
                      <th className="px-4 py-2 text-left">Status</th>
                      <th className="px-4 py-2 text-left">TX Hash</th>
                    </tr>
                  </thead>
                  <tbody>
                    {deposits.map((deposit) => (
                      <tr key={deposit.id} className="border-t">
                        <td className="px-4 py-2">{new Date(deposit.created_at).toLocaleDateString()}</td>
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
                        <td className="px-4 py-2 font-mono text-xs">{deposit.tx_hash.substring(0, 16)}...</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>

          {/* Withdrawals */}
          <div>
            <h3 className="text-lg font-semibold mb-3">Withdrawals</h3>
            {withdrawals.length === 0 ? (
              <p className="text-gray-500">No withdrawals yet</p>
            ) : (
              <div className="overflow-x-auto">
                <table className="w-full text-sm">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-4 py-2 text-left">Date</th>
                      <th className="px-4 py-2 text-left">Currency</th>
                      <th className="px-4 py-2 text-right">Amount</th>
                      <th className="px-4 py-2 text-right">Credits</th>
                      <th className="px-4 py-2 text-left">Status</th>
                      <th className="px-4 py-2 text-left">To Address</th>
                    </tr>
                  </thead>
                  <tbody>
                    {withdrawals.map((withdrawal) => (
                      <tr key={withdrawal.id} className="border-t">
                        <td className="px-4 py-2">{new Date(withdrawal.created_at).toLocaleDateString()}</td>
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
                        <td className="px-4 py-2 font-mono text-xs">
                          {withdrawal.to_address.substring(0, 16)}...
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default Crypto;

import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { ThemeProvider } from './contexts/ThemeContext';
import ProtectedRoute from './components/common/ProtectedRoute';
import Layout from './components/layout/Layout';

// Pages
import Login from './pages/auth/Login';
import Dashboard from './pages/dashboard/Dashboard';
import MarketsList from './pages/markets/MarketsList';
import CreateMarket from './pages/markets/CreateMarket';
import MarketDetail from './pages/markets/MarketDetail';
import ResolveMarket from './pages/markets/ResolveMarket';
import SportsbookDashboard from './pages/sportsbook/SportsbookDashboard';
import UsersList from './pages/users/UsersList';
import UserDetail from './pages/users/UserDetail';
import OrdersList from './pages/orders/OrdersList';
import BetsList from './pages/bets/BetsList';
import ArbitrageList from './pages/arbitrage/ArbitrageList';

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <Router>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route
              path="/*"
              element={
                <ProtectedRoute>
                  <Layout>
                    <Routes>
                      <Route path="/" element={<Dashboard />} />
                      <Route path="/markets" element={<MarketsList />} />
                      <Route path="/markets/create" element={<CreateMarket />} />
                      <Route path="/markets/:id" element={<MarketDetail />} />
                      <Route path="/markets/:id/resolve" element={<ResolveMarket />} />
                      <Route path="/sportsbook" element={<SportsbookDashboard />} />
                      <Route path="/users" element={<UsersList />} />
                      <Route path="/users/:id" element={<UserDetail />} />
                      <Route path="/orders" element={<OrdersList />} />
                      <Route path="/bets" element={<BetsList />} />
                      <Route path="/arbitrage" element={<ArbitrageList />} />
                      <Route path="*" element={<Navigate to="/" replace />} />
                    </Routes>
                  </Layout>
                </ProtectedRoute>
              }
            />
          </Routes>
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;

import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AuthProvider } from './context/AuthContext';
import { ThemeProvider } from './context/ThemeContext';
import { OddsFormatProvider } from './context/OddsFormatContext';
import { Navbar } from './components/Navbar';
import { ProtectedRoute } from './components/ProtectedRoute';
import { LiveUpdates } from './components/LiveUpdates';

import { Home } from './pages/Home';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { Dashboard } from './pages/Dashboard';
import { Markets } from './pages/Markets';
import { MarketDetail } from './pages/MarketDetail';
import { Sportsbook } from './pages/Sportsbook';
import { EventDetail } from './pages/EventDetail';
import { Arbitrage } from './pages/Arbitrage';
import { Hedges } from './pages/Hedges';
import { Bets } from './pages/Bets';
import { LinkAccount } from './pages/LinkAccount';
import { Profile } from './pages/Profile';
import { Wallet } from './pages/Wallet';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: 1,
      staleTime: 30000,
    },
  },
});

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <ThemeProvider>
          <OddsFormatProvider>
            <AuthProvider>
              <div className="min-h-screen bg-gray-50 dark:bg-dark-900">
                <Navbar />
                <Routes>
                  <Route path="/" element={<Home />} />
                  <Route path="/login" element={<Login />} />
                  <Route path="/register" element={<Register />} />

                  <Route
                    path="/dashboard"
                    element={
                      <ProtectedRoute>
                        <Dashboard />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/markets"
                    element={
                      <ProtectedRoute>
                        <Markets />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/markets/:id"
                    element={
                      <ProtectedRoute>
                        <MarketDetail />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/sportsbook"
                    element={
                      <ProtectedRoute>
                        <Sportsbook />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/sportsbook/:id"
                    element={
                      <ProtectedRoute>
                        <EventDetail />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/arbitrage"
                    element={
                      <ProtectedRoute>
                        <Arbitrage />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/hedges"
                    element={
                      <ProtectedRoute>
                        <Hedges />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/bets"
                    element={
                      <ProtectedRoute>
                        <Bets />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/link-account"
                    element={
                      <ProtectedRoute>
                        <LinkAccount />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/profile"
                    element={
                      <ProtectedRoute>
                        <Profile />
                      </ProtectedRoute>
                    }
                  />

                  <Route
                    path="/wallet"
                    element={
                      <ProtectedRoute>
                        <Wallet />
                      </ProtectedRoute>
                    }
                  />

                  <Route path="*" element={<Navigate to="/" replace />} />
                </Routes>
                <LiveUpdates />
              </div>
            </AuthProvider>
          </OddsFormatProvider>
        </ThemeProvider>
      </BrowserRouter>
    </QueryClientProvider>
  );
}

export default App;

import { BrowserRouter, Routes, Route, Link, useLocation } from 'react-router-dom'
import Dashboard from './pages/Dashboard'
import Markets from './pages/Markets'
import './index.css'

function Header() {
  const location = useLocation();

  return (
    <div className="header">
      <h1>LFG Admin Panel</h1>
      <nav className="nav">
        <Link to="/" className={location.pathname === '/' ? 'active' : ''}>
          Dashboard
        </Link>
        <Link to="/markets" className={location.pathname === '/markets' ? 'active' : ''}>
          Markets
        </Link>
      </nav>
    </div>
  );
}

function App() {
  return (
    <BrowserRouter>
      <Header />
      <div className="container">
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/markets" element={<Markets />} />
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App

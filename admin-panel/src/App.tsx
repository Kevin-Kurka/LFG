import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import './App.css';

function App() {
  return (
    <Router>
      <div className="App">
        <nav className="sidebar">
          <ul>
            <li><Link to="/">Dashboard</Link></li>
            <li><Link to="/users">Users</Link></li>
            <li><Link to="/markets">Markets</Link></li>
          </ul>
        </nav>
        <main className="content">
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/users" element={<Users />} />
            <Route path="/markets" element={<Markets />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

function Dashboard() {
  return <h2>Dashboard</h2>;
}

function Users() {
  return <h2>Users</h2>;
}

function Markets() {
  return <h2>Markets</h2>;
}

export default App;

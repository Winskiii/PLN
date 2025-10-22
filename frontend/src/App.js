import { useState, useEffect } from 'react';
import CsvTable from './CsvTable';
import logo from './logo.svg';
import './App.css';

function App() {
  const [health, setHealth] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Check backend health
    fetch('/api/health')
      .then(res => res.json())
      .then(data => {
        setHealth(data);
        setLoading(false);
      })
      .catch(err => {
        console.error('Backend health check failed:', err);
        setLoading(false);
      });
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <h1>PLN Secure Web Application</h1>
        <p>
          Built with <strong>React + Vite</strong> and <strong>Go Backend</strong>
        </p>
        
        <div className="status-card">
          <h3>Backend Status</h3>
          {loading ? (
            <p>Checking backend...</p>
          ) : health ? (
            <div className="status-success">
              <p>✅ Status: {health.status}</p>
              <p>🕐 Time: {new Date(health.timestamp).toLocaleString()}</p>
            </div>
          ) : (
            <div className="status-error">
              <p>❌ Backend not connected</p>
              <p>Make sure Go backend is running on port 8080</p>
            </div>
          )}
        </div>

        <CsvTable />

        <div className="features">
          <h3>🔒 Security Features (SSDLC)</h3>
          <ul>
            <li>Secure Headers (CSP, HSTS, X-Frame-Options)</li>
            <li>JWT Authentication</li>
            <li>Rate Limiting</li>
            <li>Input Validation</li>
            <li>CORS Protection</li>
            <li>SQL Injection Prevention</li>
          </ul>
        </div>
      </header>
    </div>
  );
}

export default App;

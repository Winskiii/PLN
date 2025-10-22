import { useEffect, useState } from 'react';

export default function CsvTable() {
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch('/api/csv')
      .then(res => res.json())
      .then(data => {
        setRows(data);
        setLoading(false);
      })
      .catch(err => {
        setError('Gagal mengambil data CSV');
        setLoading(false);
      });
  }, []);

  if (loading) return <div>Loading CSV data...</div>;
  if (error) return <div style={{color:'red'}}>{error}</div>;
  if (!rows.length) return <div>Tidak ada data CSV</div>;

  const headers = Object.keys(rows[0]);

  return (
    <div className="csv-table-container">
      <h2>Data CSV</h2>
      <table className="csv-table">
        <thead>
          <tr>
            {headers.map(h => <th key={h}>{h}</th>)}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, i) => (
            <tr key={i}>
              {headers.map(h => <td key={h}>{row[h]}</td>)}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

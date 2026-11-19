import { useCallback, useEffect, useState } from "react";
import { AddSecret, ListRecords, RemoveKey } from "./bindings/backstage.lisitede.com/hub/vaultservice.js";

function formatTime(ts) {
  const d = new Date(ts);
  if (Number.isNaN(d.getTime())) return String(ts);
  const p = (n, w = 2) => String(n).padStart(w, "0");
  return (
    `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ` +
    `${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}.${p(d.getMilliseconds(), 3)}`
  );
}

function App() {
  const [records, setRecords] = useState([]);
  const [newSecret, setNewSecret] = useState(null);
  const [error, setError] = useState(null);

  const refresh = useCallback(() => {
    ListRecords()
      .then((rs) => {
        setRecords(rs || []);
        setError(null);
      })
      .catch((err) => setError(String(err)));
  }, []);

  useEffect(() => {
    refresh();
    const timer = setInterval(refresh, 2000);
    return () => clearInterval(timer);
  }, [refresh]);

  const handleAdd = useCallback(() => {
    AddSecret()
      .then((secret) => {
        setNewSecret(secret);
      })
      .catch((err) => setError(String(err)));
  }, []);

  const handleDelete = useCallback((secret) => {
    RemoveKey(secret)
      .then(() => {
        if (secret === newSecret) setNewSecret(null);
        refresh();
      })
      .catch((err) => setError(String(err)));
  }, [newSecret, refresh]);

  return (
    <div className="app">
      <div className="vault-header">
        <h1>Vault</h1>
        {newSecret && (
          <span className="notice">
            新 secret：<code>{newSecret}</code>
          </span>
        )}
        <button className="add" onClick={handleAdd}>
          +
        </button>
      </div>
      <hr />
      {error ? (
        <p className="error">记录读取失败：{error}</p>
      ) : records.length === 0 ? (
        <p className="empty">暂无审计记录</p>
      ) : (
        <table>
          <thead>
            <tr>
              <th>时间</th>
              <th>Role</th>
              <th>Secret</th>
              <th>n</th>
              <th>k</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {records.map((r, i) => (
              <tr key={i}>
                <td>{formatTime(r.timestamp)}</td>
                <td>{r.role}</td>
                <td>{r.secret}</td>
                <td>{r.n}</td>
                <td>{r.k}</td>
                <td className="actions">
                  <button className="delete" onClick={() => handleDelete(r.secret)}>
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default App;

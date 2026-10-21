import { useCallback, useEffect, useState } from "react";
import { GetStatus } from "./bindings/backstage.lisitede.com/hub/hubservice.js";

function App() {
  const [status, setStatus] = useState(null);
  const [error, setError] = useState(null);

  const refresh = useCallback(() => {
    GetStatus()
      .then((s) => {
        setStatus(s);
        setError(null);
      })
      .catch((err) => setError(String(err)));
  }, []);

  useEffect(() => {
    refresh();
    const timer = setInterval(refresh, 5000);
    return () => clearInterval(timer);
  }, [refresh]);

  return (
    <div className="app">
      <h1>Backstage Hub</h1>
      {error ? (
        <p className="error">状态读取失败：{error}</p>
      ) : status ? (
        <dl className="status">
          <dt>应用版本</dt>
          <dd>{status.version}</dd>
          <dt>HTTP 服务地址</dt>
          <dd>
            <a href={status.baseUrl} target="_blank" rel="noreferrer">
              {status.baseUrl}
            </a>
          </dd>
          <dt>运行状态</dt>
          <dd className={status.running ? "ok" : "down"}>
            {status.running ? "运行中" : "已停止"}
          </dd>
        </dl>
      ) : (
        <p>加载中…</p>
      )}
      <button onClick={refresh}>刷新</button>
    </div>
  );
}

export default App;

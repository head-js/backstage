const tasks = [
  {
    id: "rust.log.timestamp",
    name: "Timestamp Log",
    interval: "5s",
    status: "running",
    lastRunAt: "just now",
    lastDurationMs: 8,
  },
  {
    id: "rust.http.doget",
    name: "HTTP DoGet",
    interval: "30s",
    status: "idle",
    lastRunAt: "24s ago",
    lastDurationMs: 412,
  },
  {
    id: "desktop.cleanup.cache",
    name: "Cache Cleanup",
    interval: "manual",
    status: "skipped",
    lastRunAt: "not run",
    lastDurationMs: 0,
  },
];

const statusClass: Record<string, string> = {
  running: "badge badge-success",
  idle: "badge badge-info",
  failed: "badge badge-error",
  skipped: "badge badge-warning",
};

export function Scheduler() {
  return (
    <div className="space-y-4">
      <div>
        <h2 className="text-base font-semibold">Scheduler</h2>
        <p className="text-sm text-base-content/60">Mock task state for the current UI validation stage.</p>
      </div>

      <div className="rounded-md border border-base-300 bg-base-100">
        <div className="overflow-x-auto">
          <table className="table table-sm">
            <thead>
              <tr>
                <th>Task</th>
                <th>Interval</th>
                <th>Status</th>
                <th>Last run</th>
                <th>Duration</th>
              </tr>
            </thead>
            <tbody>
              {tasks.map((task) => (
                <tr key={task.id}>
                  <td>
                    <div className="font-medium">{task.name}</div>
                    <div className="font-mono text-xs text-base-content/55">{task.id}</div>
                  </td>
                  <td>{task.interval}</td>
                  <td>
                    <span className={statusClass[task.status]}>{task.status}</span>
                  </td>
                  <td>{task.lastRunAt}</td>
                  <td>{task.lastDurationMs ? `${task.lastDurationMs}ms` : "-"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

const stats = [
  { label: "Runtime", value: "Ready", detail: "Tauri shell online" },
  { label: "Commands", value: "1", detail: "ai hasshin" },
  { label: "Tasks", value: "2", detail: "mock status view" },
  { label: "Theme", value: "Store", detail: "LazyStore backed" },
];

export function Dashboard() {
  return (
    <div className="space-y-4">
      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        {stats.map((stat) => (
          <article key={stat.label} className="card rounded-md border border-base-300 bg-base-100">
            <div className="card-body p-4">
              <div className="text-xs font-medium uppercase text-base-content/55">{stat.label}</div>
              <div className="text-2xl font-semibold">{stat.value}</div>
              <div className="text-sm text-base-content/60">{stat.detail}</div>
            </div>
          </article>
        ))}
      </section>

      <section className="rounded-md border border-base-300 bg-base-100">
        <div className="border-b border-base-300 px-4 py-3">
          <h2 className="text-sm font-semibold">Console Snapshot</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="table table-sm">
            <tbody>
              <tr>
                <td className="font-medium">Window mode</td>
                <td>Opaque dashboard shell, 800x600 minimum</td>
                <td>
                  <span className="badge badge-success badge-sm">ready</span>
                </td>
              </tr>
              <tr>
                <td className="font-medium">System bridge</td>
                <td>get_system_info command</td>
                <td>
                  <span className="badge badge-info badge-sm">live</span>
                </td>
              </tr>
              <tr>
                <td className="font-medium">Command bridge</td>
                <td>Fixed ai hasshin runner</td>
                <td>
                  <span className="badge badge-info badge-sm">live</span>
                </td>
              </tr>
              <tr>
                <td className="font-medium">Scheduler</td>
                <td>Rust runtime running, UI state mocked</td>
                <td>
                  <span className="badge badge-warning badge-sm">mock</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>
  );
}

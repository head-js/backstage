import { useEffect, useState } from "react";
import { listBlames } from "../lib/tauri";

interface BlameRow {
  id?: string;
  blameId?: string;
  name?: string;
  appId?: string;
  planId?: string;
  context?: string;
  updatedAt?: string;
}

export function Blames() {
  const [rows, setRows] = useState<BlameRow[]>([]);
  const [selectedRow, setSelectedRow] = useState<BlameRow | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  function openRow(row: BlameRow) {
    setSelectedRow(row);
  }

  function closeDrawer() {
    setSelectedRow(null);
  }

  function formatDrawerSubtitle(row: BlameRow | null) {
    if (!row) {
      return "";
    }
    return [row.appId, row.planId, row.blameId].filter(Boolean).join(" / ");
  }

  useEffect(() => {
    let cancelled = false;

    async function load() {
      setLoading(true);
      setError(null);
      try {
        const data = await listBlames();
        console.log("tauri_edge /blames", data);
        if (!cancelled) {
          setRows(Array.isArray(data) ? data : []);
        }
      } catch (nextError) {
        console.error("Failed to load blames:", nextError);
        if (!cancelled) {
          setRows([]);
          setError(String(nextError));
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    }

    load();
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <div className="drawer drawer-end">
      <input
        id="blame-detail-drawer"
        type="checkbox"
        className="drawer-toggle"
        checked={selectedRow !== null}
        onChange={(event) => {
          if (!event.currentTarget.checked) {
            closeDrawer();
          }
        }}
      />
      <div className="drawer-content space-y-4">
        {error && (
          <div className="rounded-md border border-error/30 bg-error/10 px-4 py-3 text-sm text-error">
            {error}
          </div>
        )}

        <section className="rounded-md border border-base-300 bg-base-100">
          <div className="overflow-x-auto">
            <table className="table table-sm">
              <thead>
                <tr>
                  <th>App Id</th>
                  <th>Plan Id</th>
                  <th>Name</th>
                  <th>Updated At</th>
                </tr>
              </thead>
              <tbody>
                {loading && (
                  <tr>
                    <td colSpan={4}>
                      <span className="loading loading-spinner loading-sm" />
                    </td>
                  </tr>
                )}
                {!loading && rows.map((row, index) => (
                  <tr
                    key={row.id ?? index}
                    className="cursor-pointer hover:bg-base-200"
                    role="button"
                    tabIndex={0}
                    onClick={() => openRow(row)}
                    onKeyDown={(event) => {
                      if (event.key === "Enter" || event.key === " ") {
                        event.preventDefault();
                        openRow(row);
                      }
                    }}
                  >
                    <td>{row.appId ?? ""}</td>
                    <td>{row.planId ?? ""}</td>
                    <td className="font-medium">{row.name ?? ""}</td>
                    <td>{row.updatedAt ?? ""}</td>
                  </tr>
                ))}
                {!loading && rows.length === 0 && (
                  <tr>
                    <td colSpan={4} className="text-base-content/55">
                      No blames found.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </section>
      </div>
      <div className="drawer-side z-30 overflow-x-hidden">
        <label htmlFor="blame-detail-drawer" aria-label="close detail" className="drawer-overlay" />
        <aside className="app-scrollbar flex h-full w-full flex-col overflow-y-auto border-l border-base-300 bg-base-100 sm:w-[70vw] sm:max-w-[70vw]">
          <div className="flex min-h-14 items-center justify-between gap-3 border-b border-base-300 px-4 py-2">
            <h2 className="min-w-0 flex-1">
              <span className="block truncate text-sm font-semibold leading-5">{selectedRow?.name ?? "Blame Title"}</span>
              <span className="block truncate font-mono text-xs leading-4 text-base-content/55">
                {formatDrawerSubtitle(selectedRow)}
              </span>
            </h2>
            <button type="button" className="btn btn-ghost btn-sm shrink-0" onClick={closeDrawer}>
              Close
            </button>
          </div>
          <div className="p-4">
            <pre className="app-scrollbar max-h-[calc(100vh-9rem)] whitespace-pre-wrap break-words rounded bg-base-200 p-3 text-xs leading-5">{selectedRow?.context ?? ""}</pre>
          </div>
        </aside>
      </div>
    </div>
  );
}

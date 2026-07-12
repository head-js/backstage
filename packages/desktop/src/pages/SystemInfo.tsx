import { useCallback, useEffect, useState } from "react";
import { getSystemInfo } from "../lib/tauri";
import type { SystemInfo as SystemInfoData } from "../types/system";

export function SystemInfo() {
  const [info, setInfo] = useState<SystemInfoData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      setInfo(await getSystemInfo());
    } catch (err) {
      setError(String(err));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between gap-3">
        <div>
          <h2 className="text-base font-semibold">System Info</h2>
          <p className="text-sm text-base-content/60">Native OS metadata from Tauri.</p>
        </div>
        <button className="btn btn-primary btn-sm" onClick={load} disabled={loading}>
          {loading ? "Loading" : "Retry"}
        </button>
      </div>

      {error && <div className="alert alert-error text-sm">{error}</div>}

      <div className="rounded-md border border-base-300 bg-base-100">
        <div className="overflow-x-auto">
          <table className="table">
            <tbody>
              {["platform", "version", "arch", "family", "hostname"].map((key) => (
                <tr key={key}>
                  <td className="w-40 text-sm font-semibold capitalize">{key}</td>
                  <td className="font-mono text-sm">
                    {loading && !info ? <span className="loading loading-dots loading-sm" /> : info?.[key as keyof SystemInfoData]}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

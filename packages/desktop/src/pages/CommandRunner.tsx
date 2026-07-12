import { useState } from "react";
import { runAiHasshin } from "../lib/tauri";
import type { CommandResult } from "../types/command";

export function CommandRunner() {
  const [result, setResult] = useState<CommandResult | null>(null);
  const [loading, setLoading] = useState(false);

  const run = async () => {
    setLoading(true);
    setResult(null);
    try {
      setResult(await runAiHasshin());
    } catch (error) {
      setResult({
        stdout: "",
        stderr: "",
        exit_code: null,
        success: false,
        error: String(error),
        duration_ms: 0,
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-base font-semibold">Command Runner</h2>
          <p className="text-sm text-base-content/60">Runs the fixed host command: ai hasshin.</p>
        </div>
        <button className="btn btn-primary btn-sm" onClick={run} disabled={loading}>
          {loading ? "Running" : "Run"}
        </button>
      </div>

      <div className="rounded-md border border-base-300 bg-base-100 p-4">
        <div className="mb-3 flex flex-wrap items-center gap-2">
          <span className="badge badge-outline">ai hasshin</span>
          {result && (
            <span className={result.success ? "badge badge-success" : "badge badge-error"}>
              {result.success ? "success" : "failed"}
            </span>
          )}
          {result && <span className="badge badge-ghost">{result.duration_ms}ms</span>}
          {loading && <span className="loading loading-spinner loading-sm" />}
        </div>

        {result ? (
          <div className="grid gap-3">
            <OutputPanel title="stdout" value={result.stdout} />
            <OutputPanel title="stderr" value={result.stderr} />
            <OutputPanel title="error" value={result.error ?? ""} />
            <div className="text-sm text-base-content/65">exit_code: {result.exit_code ?? "N/A"}</div>
          </div>
        ) : (
          <div className="text-sm text-base-content/60">No command output yet.</div>
        )}
      </div>
    </div>
  );
}

function OutputPanel({ title, value }: { title: string; value: string }) {
  return (
    <section>
      <div className="mb-1 text-xs font-semibold uppercase text-base-content/55">{title}</div>
      <pre className="app-scrollbar max-h-48 overflow-auto rounded bg-base-200 p-3 text-xs leading-5">
        {value || "(empty)"}
      </pre>
    </section>
  );
}

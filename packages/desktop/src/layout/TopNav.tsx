import { useLocation } from "react-router-dom";
import { findMenuItem } from "../router/menu";

export function TopNav() {
  const location = useLocation();
  const current = findMenuItem(location.pathname);

  return (
    <header className="navbar min-h-14 border-b border-base-300 bg-base-100 px-4">
      <div className="flex-1">
        <label
          htmlFor="app-drawer"
          className="btn btn-ghost btn-square btn-sm mr-2 lg:hidden"
          aria-label="Open menu"
          title="Open menu"
        >
          <span className="flex w-4 flex-col gap-1">
            <span className="h-0.5 rounded bg-current" />
            <span className="h-0.5 rounded bg-current" />
            <span className="h-0.5 rounded bg-current" />
          </span>
        </label>
        <div>
          <h1 className="text-base font-semibold leading-tight">{current?.label ?? "Console"}</h1>
          <p className="text-xs text-base-content/55">{location.pathname}</p>
        </div>
      </div>
    </header>
  );
}

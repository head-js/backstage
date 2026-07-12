import { Outlet } from "react-router-dom";
import { Sidebar } from "./Sidebar";
import { TopNav } from "./TopNav";

export function AppShell() {
  return (
    <div className="drawer h-screen lg:drawer-open">
      <input id="app-drawer" type="checkbox" className="drawer-toggle" />
      <div className="drawer-content flex min-h-0 flex-col bg-base-200">
        <TopNav />
        <main className="app-scrollbar min-h-0 flex-1 overflow-y-auto p-4">
          <Outlet />
        </main>
      </div>
      <div className="drawer-side z-40">
        <label htmlFor="app-drawer" aria-label="close sidebar" className="drawer-overlay" />
        <Sidebar />
      </div>
    </div>
  );
}

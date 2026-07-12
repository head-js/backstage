import { menuGroups } from "../router/menu";

interface SidebarProps {
  activePath?: string | null;
  onNavigate?: (path: string) => void;
}

export function Sidebar({ activePath, onNavigate }: SidebarProps) {
  return (
    <aside className="flex h-screen w-52 shrink-0 flex-col border-r border-base-300 bg-transparent">
      <div className="flex h-14 items-center border-b border-base-300 bg-yellow-100 px-4">
        <div>
          <div className="text-sm font-semibold leading-tight">backstage desktop</div>
          <div className="text-xs text-base-content/60">Desktop console</div>
        </div>
      </div>
      <nav className="app-scrollbar flex-1 overflow-y-auto px-3 py-4">
        {menuGroups.map((group) => (
          <section key={group.label} className="mb-5">
            <div className="mb-2 px-3 text-[11px] font-semibold uppercase tracking-wide text-base-content/45">
              {group.label}
            </div>
            <ul className="menu gap-1 p-0">
              {group.items.map((item) => (
                <li key={item.path}>
                  <button
                    type="button"
                    onClick={() => onNavigate?.(item.path)}
                    className={activePath === item.path ? "active text-sm font-medium" : "text-sm font-medium"}
                  >
                    {item.label}
                  </button>
                </li>
              ))}
            </ul>
          </section>
        ))}
      </nav>
    </aside>
  );
}

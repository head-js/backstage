import { useEffect, useState } from "react";
import { emitTo } from "@tauri-apps/api/event";
import { getCurrentWindow, Window } from "@tauri-apps/api/window";
import { HashRouter } from "react-router-dom";
import {
  APPEARANCE_CHANGED_EVENT,
  loadAppearance,
  type AppearanceSettings,
} from "./lib/store";
import { Sidebar } from "./layout/Sidebar";

const currentWindow = getCurrentWindow();

function delay(ms: number) {
  return new Promise((resolve) => window.setTimeout(resolve, ms));
}

async function showWorkspace(workspace: Window | null, path: string) {
  if (!workspace) {
    return;
  }

  try {
    await workspace.setFocusable(false);
    await workspace.show();
    await emitTo("workspace", "desktop:navigate", path);
    await delay(50);
  } finally {
    await workspace.setFocusable(true);
  }

  await currentWindow.setFocus();
}

export default function AppMenu() {
  const [activeMenuPath, setActiveMenuPath] = useState<string | null>(null);

  useEffect(() => {
    loadAppearance()
      .then((appearance) => {
        document.documentElement.setAttribute("data-theme", appearance.theme);
      })
      .catch((error) => {
        console.error("Failed to load appearance settings:", error);
      });
  }, []);

  useEffect(() => {
    const handleAppearanceChanged = (event: Event) => {
      const nextTheme = (event as CustomEvent<AppearanceSettings>).detail.theme;
      document.documentElement.setAttribute("data-theme", nextTheme);
    };

    window.addEventListener(APPEARANCE_CHANGED_EVENT, handleAppearanceChanged);
    return () => window.removeEventListener(APPEARANCE_CHANGED_EVENT, handleAppearanceChanged);
  }, []);

  return (
    <HashRouter>
      <Sidebar
        activePath={activeMenuPath}
        onNavigate={async (path) => {
          const workspace = await Window.getByLabel("workspace");
          if (activeMenuPath === path) {
            setActiveMenuPath(null);
            await workspace?.hide();
            return;
          }

          setActiveMenuPath(path);
          await showWorkspace(workspace, path);
        }}
      />
    </HashRouter>
  );
}

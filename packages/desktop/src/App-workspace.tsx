import { useEffect } from "react";
import { listen } from "@tauri-apps/api/event";
import { HashRouter } from "react-router-dom";
import {
  APPEARANCE_CHANGED_EVENT,
  loadAppearance,
  type AppearanceSettings,
} from "./lib/store";
import { AppRoutes } from "./router";

export default function AppWorkspace() {
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

  useEffect(() => {
    let unlisten: (() => void) | undefined;

    void listen<string>("desktop:navigate", (event) => {
      window.location.hash = event.payload;
    }).then((dispose) => {
      unlisten = dispose;
    });

    return () => unlisten?.();
  }, []);

  return (
    <HashRouter>
      <div className="h-screen min-w-0 bg-base-200">
        <main className="app-scrollbar h-full overflow-y-auto p-4">
          <AppRoutes />
        </main>
      </div>
    </HashRouter>
  );
}

import { useEffect } from "react";
import { HashRouter } from "react-router-dom";
import {
  APPEARANCE_CHANGED_EVENT,
  loadAppearance,
  type AppearanceSettings,
} from "./lib/store";
import { AppRoutes } from "./router";

export default function App() {
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
      <AppRoutes />
    </HashRouter>
  );
}

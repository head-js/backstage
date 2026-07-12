import { LazyStore } from "@tauri-apps/plugin-store";

export type ThemePreference = "light" | "dark";

export interface AppearanceSettings {
  theme: ThemePreference;
}

export const APPEARANCE_CHANGED_EVENT = "desktop:appearance-changed";

const settings = new LazyStore("settings.json");

function normalizeTheme(theme: unknown): ThemePreference {
  return theme === "dark" || theme === "light" ? theme : "light";
}

export async function loadAppearance(): Promise<AppearanceSettings> {
  const value = await settings.get<{ theme?: unknown }>("appearance");
  return {
    theme: normalizeTheme(value?.theme),
  };
}

export async function saveAppearance(appearance: AppearanceSettings) {
  await settings.set("appearance", appearance);
  await settings.save();
  window.dispatchEvent(new CustomEvent<AppearanceSettings>(APPEARANCE_CHANGED_EVENT, { detail: appearance }));
}

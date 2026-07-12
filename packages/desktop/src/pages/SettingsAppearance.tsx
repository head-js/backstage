import { useEffect, useState } from "react";
import {
  loadAppearance,
  saveAppearance,
  type ThemePreference,
} from "../lib/store";

const themes: { label: string; value: ThemePreference; note: string }[] = [
  { label: "Light", value: "light", note: "daisyUI default light theme." },
  { label: "Dark", value: "dark", note: "daisyUI default dark theme." },
];

export function SettingsAppearance() {
  const [theme, setTheme] = useState<ThemePreference>("light");
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadAppearance()
      .then((appearance) => {
        setTheme(appearance.theme);
      })
      .catch((error) => console.error("Failed to load appearance settings:", error));
  }, []);

  const updateTheme = async (nextTheme: ThemePreference) => {
    if (nextTheme === theme || saving) {
      return;
    }

    setSaving(true);
    setError(null);
    try {
      await saveAppearance({ theme: nextTheme });
      setTheme(nextTheme);
    } catch (err) {
      setError(String(err));
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-4">
      <div>
        <h2 className="text-base font-semibold">Appearance Settings</h2>
        <p className="text-sm text-base-content/60">Theme preference is persisted with LazyStore.</p>
      </div>

      <div className="rounded-md border border-base-300 bg-base-100 p-4">
        <div className="mb-3 flex items-center justify-between">
          <div className="text-sm font-semibold">Theme</div>
          {saving && <span className="loading loading-spinner loading-sm" />}
        </div>
        {error && <div className="alert alert-error mb-3 py-2 text-sm">{error}</div>}
        <div className="grid gap-2 md:grid-cols-2">
          {themes.map((item) => (
            <label
              key={item.value}
              className="flex cursor-pointer items-start gap-3 rounded border border-base-300 p-3 hover:bg-base-200"
            >
              <input
                type="radio"
                className="radio radio-primary radio-sm mt-1"
                checked={theme === item.value}
                disabled={saving}
                onChange={() => updateTheme(item.value)}
              />
              <span>
                <span className="block text-sm font-medium">{item.label}</span>
                <span className="block text-xs text-base-content/60">{item.note}</span>
              </span>
            </label>
          ))}
        </div>
      </div>
    </div>
  );
}

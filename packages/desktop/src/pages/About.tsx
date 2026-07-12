export function About() {
  return (
    <div className="space-y-4">
      <div>
        <h2 className="text-base font-semibold">About</h2>
        <p className="text-sm text-base-content/60">desktop dashboard demo.</p>
      </div>

      <div className="rounded-md border border-base-300 bg-base-100">
        <table className="table">
          <tbody>
            <tr>
              <td className="w-40 font-medium">Application</td>
              <td>desktop</td>
            </tr>
            <tr>
              <td className="font-medium">Version</td>
              <td>0.1.0</td>
            </tr>
            <tr>
              <td className="font-medium">Stack</td>
              <td>Tauri 2, React 18, TypeScript, Tailwind CSS, daisyUI</td>
            </tr>
            <tr>
              <td className="font-medium">Current scope</td>
              <td>System info, fixed command runner, mock scheduler, theme settings</td>
            </tr>
            <tr>
              <td className="font-medium">Window capability</td>
              <td>Tray menu supports show, hide, quit, and always-on-top toggle</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
}

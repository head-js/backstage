import { Navigate, Route, Routes } from "react-router-dom";
import { About } from "../pages/About";
import { Agenda } from "../pages/Agenda";
import { CommandRunner } from "../pages/CommandRunner";
import { Dashboard } from "../pages/Dashboard";
import { NotFound } from "../pages/NotFound";
import { Scheduler } from "../pages/Scheduler";
import { SettingsAppearance } from "../pages/SettingsAppearance";
import { SystemInfo } from "../pages/SystemInfo";

export function AppRoutes() {
  return (
    <Routes>
      <Route index element={<Navigate to="/dashboard" replace />} />
      <Route path="/dashboard" element={<Dashboard />} />
      <Route path="/agenda" element={<Agenda />} />
      <Route path="/system/info" element={<SystemInfo />} />
      <Route path="/system/command" element={<CommandRunner />} />
      <Route path="/system/scheduler" element={<Scheduler />} />
      <Route path="/settings/appearance" element={<SettingsAppearance />} />
      <Route path="/about" element={<About />} />
      <Route path="*" element={<NotFound />} />
    </Routes>
  );
}

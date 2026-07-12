import { invoke } from "@tauri-apps/api/core";
import type { CommandResult } from "../types/command";
import type { SystemInfo } from "../types/system";

export function getSystemInfo() {
  return invoke<SystemInfo>("get_system_info");
}

export function runAiHasshin() {
  return invoke<CommandResult>("run_command");
}

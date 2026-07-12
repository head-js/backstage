export interface CommandResult {
  stdout: string;
  stderr: string;
  exit_code: number | null;
  success: boolean;
  error: string | null;
  duration_ms: number;
}

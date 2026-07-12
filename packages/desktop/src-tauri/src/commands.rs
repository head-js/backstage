use serde::Serialize;
use std::time::Duration;

#[derive(Serialize)]
pub struct SystemInfo {
    pub platform: String,
    pub version: String,
    pub arch: String,
    pub family: String,
    pub hostname: String,
}

#[tauri::command]
pub async fn get_system_info() -> SystemInfo {
    let platform = tauri_plugin_os::platform().to_string();
    let version = tauri_plugin_os::version().to_string();
    let arch = tauri_plugin_os::arch().to_string();
    let family = tauri_plugin_os::family().to_string();
    let hostname = tauri_plugin_os::hostname();
    SystemInfo {
        platform,
        version,
        arch,
        family,
        hostname,
    }
}

#[derive(Serialize)]
pub struct CommandResult {
    pub stdout: String,
    pub stderr: String,
    pub exit_code: Option<i32>,
    pub success: bool,
    pub error: Option<String>,
    pub duration_ms: u64,
}

#[tauri::command]
pub async fn run_command() -> CommandResult {
    let command = "ai hasshin";
    let start = std::time::Instant::now();

    let result = tokio::time::timeout(Duration::from_secs(30), async {
        #[cfg(target_os = "windows")]
        {
            tokio::process::Command::new("pwsh")
                .args(["-NoProfile", "-NonInteractive", "-Command", command])
                .output()
                .await
        }
        #[cfg(target_os = "macos")]
        {
            tokio::process::Command::new("/bin/bash")
                .args(["-lc", command])
                .output()
                .await
        }
        #[cfg(not(any(target_os = "windows", target_os = "macos")))]
        {
            Err(std::io::Error::new(
                std::io::ErrorKind::Unsupported,
                "unsupported platform",
            ))
        }
    })
    .await;

    let duration_ms = start.elapsed().as_millis() as u64;

    match result {
        Ok(Ok(output)) => CommandResult {
            stdout: String::from_utf8_lossy(&output.stdout).into_owned(),
            stderr: String::from_utf8_lossy(&output.stderr).into_owned(),
            exit_code: output.status.code(),
            success: output.status.success(),
            error: None,
            duration_ms,
        },
        Ok(Err(e)) => {
            let error_msg = if e.kind() == std::io::ErrorKind::NotFound {
                #[cfg(target_os = "windows")]
                {
                    format!("PowerShell 7 not found: {}", e)
                }
                #[cfg(not(target_os = "windows"))]
                {
                    format!("Bash not found: {}", e)
                }
            } else {
                format!("{}", e)
            };
            CommandResult {
                stdout: String::new(),
                stderr: String::new(),
                exit_code: None,
                success: false,
                error: Some(error_msg),
                duration_ms,
            }
        }
        Err(_) => CommandResult {
            stdout: String::new(),
            stderr: String::new(),
            exit_code: None,
            success: false,
            error: Some("Command timed out after 30 seconds".to_string()),
            duration_ms,
        },
    }
}
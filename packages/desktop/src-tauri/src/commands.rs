use crate::cli::parse_cli_output;
use serde::{Deserialize, Serialize};
use std::time::Duration;

#[derive(Debug, PartialEq, Eq)]
enum EdgeAction {
    ListBlames,
}

#[derive(Deserialize)]
struct GiteaBlame {
    id: String,
    name: String,
    gitea: GiteaBlameGiteaExtra,
}

#[derive(Default, Deserialize)]
struct GiteaBlameGiteaExtra {
    owner: Option<String>,
    repo: Option<String>,
    body: Option<String>,
    #[serde(rename(deserialize = "updatedAt"))]
    updated_at: Option<String>,
}

#[derive(Serialize)]
#[serde(rename_all(serialize = "camelCase"))]
struct EdgeBlame {
    id: String,
    blame_id: String,
    name: String,
    app_id: String,
    plan_id: String,
    context: String,
    updated_at: String,
}

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

#[tauri::command]
pub async fn tauri_edge(
    method: String,
    pathname: String,
    query: serde_json::Value,
) -> Result<serde_json::Value, String> {
    let _query = query;
    let action = route_edge_action(&method, &pathname)?;
    run_edge_action(action).await
}

fn route_edge_action(method: &str, pathname: &str) -> Result<EdgeAction, String> {
    let normalized_method = method.to_uppercase();

    match (normalized_method.as_str(), pathname) {
        ("GET", "/blames") => Ok(EdgeAction::ListBlames),
        _ => Err(format!(
            "unsupported edge route: {} {}",
            normalized_method, pathname
        )),
    }
}

async fn run_edge_action(action: EdgeAction) -> Result<serde_json::Value, String> {
    match action {
        EdgeAction::ListBlames => run_list_blames().await,
    }
}

async fn run_list_blames() -> Result<serde_json::Value, String> {
    let result = tokio::time::timeout(Duration::from_secs(30), async {
        tokio::process::Command::new("backstage-gitea")
            .args(["plan", "GET", "/blames"])
            .output()
            .await
    })
    .await;

    match result {
        Ok(Ok(output)) => {
            let stdout = String::from_utf8_lossy(&output.stdout).into_owned();
            let stderr = String::from_utf8_lossy(&output.stderr).into_owned();
            transform_list_blames_output(&stdout, &stderr, output.status.code())
        }
        Ok(Err(e)) => {
            if e.kind() == std::io::ErrorKind::NotFound {
                Err(format!("backstage-gitea not found: {}", e))
            } else {
                Err(format!("failed to run backstage-gitea: {}", e))
            }
        }
        Err(_) => Err("backstage-gitea plan GET /blames timed out after 30 seconds".to_string()),
    }
}

fn transform_list_blames_output(
    stdout: &str,
    stderr: &str,
    exit_code: Option<i32>,
) -> Result<serde_json::Value, String> {
    let blames: Vec<GiteaBlame> = parse_cli_output(stdout, stderr, exit_code)?;
    let edge_blames: Vec<EdgeBlame> = blames.into_iter().map(map_gitea_blame).collect();

    serde_json::to_value(edge_blames).map_err(|e| format!("failed to serialize edge blames: {}", e))
}

fn map_gitea_blame(blame: GiteaBlame) -> EdgeBlame {
    let gitea = blame.gitea;
    let blame_id = blame.id;
    let app_id = gitea.owner.unwrap_or_default();
    let plan_id = gitea.repo.unwrap_or_default();

    EdgeBlame {
        id: format!("{}-{}-{}", app_id, plan_id, blame_id),
        blame_id,
        name: blame.name,
        app_id,
        plan_id,
        context: gitea.body.unwrap_or_default(),
        updated_at: format_updated_at(gitea.updated_at.as_deref().unwrap_or_default()),
    }
}

fn format_updated_at(updated_at: &str) -> String {
    match (updated_at.get(5..10), updated_at.get(11..19)) {
        (Some(date), Some(time)) => format!("{} {}", date, time),
        _ => updated_at.to_string(),
    }
}

#[cfg(test)]
#[path = "commands_tests.rs"]
mod tests;

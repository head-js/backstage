use serde::de::DeserializeOwned;

pub(crate) fn parse_cli_output<T: DeserializeOwned>(
    stdout: &str,
    stderr: &str,
    exit_code: Option<i32>,
) -> Result<T, String> {
    let trimmed_stderr = stderr.trim();
    if !trimmed_stderr.is_empty() {
        return Err(format!(
            "CLI wrote to stderr with exit_code {:?}: {}",
            exit_code, trimmed_stderr
        ));
    }

    if exit_code != Some(0) {
        return Err(format!(
            "CLI exited with code {:?}: {}",
            exit_code,
            summarize_output(stdout)
        ));
    }

    serde_json::from_str(stdout).map_err(|e| {
        format!(
            "failed to parse CLI stdout as JSON: {}; stdout: {}",
            e,
            summarize_output(stdout)
        )
    })
}

fn summarize_output(value: &str) -> String {
    const MAX_LEN: usize = 500;
    let trimmed = value.trim();
    if trimmed.len() <= MAX_LEN {
        return trimmed.to_string();
    }
    format!("{}...", &trimmed[..MAX_LEN])
}

#[cfg(test)]
#[path = "cli_tests.rs"]
mod tests;

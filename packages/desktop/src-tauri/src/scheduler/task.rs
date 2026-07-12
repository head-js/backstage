use std::time::{Duration, SystemTime, UNIX_EPOCH};

pub fn run_timestamp_log() -> tokio::task::JoinHandle<()> {
    tokio::spawn(async move {
        let ts = SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default()
            .as_secs();
        tracing::info!("[scheduler][rust.log.timestamp] timestamp={ts}");
    })
}

pub fn run_http_doget() -> tokio::task::JoinHandle<()> {
    tokio::spawn(async move {
        let url = "https://httpbin.org";
        tracing::info!("[scheduler][rust.http.doget] GET {url} started");

        let client = reqwest::Client::builder()
            .timeout(Duration::from_secs(10))
            .build();

        let client = match client {
            Ok(c) => c,
            Err(e) => {
                tracing::error!("[scheduler][rust.http.doget] failed: client build error: {e}");
                return;
            }
        };

        let start = SystemTime::now();
        match client.get(url).send().await {
            Ok(resp) => {
                let elapsed = start.elapsed().unwrap_or_default().as_millis();
                let status = resp.status().as_u16();
                let body = resp.text().await.unwrap_or_default();
                let body_bytes = body.len();
                tracing::info!(
                    "[scheduler][rust.http.doget] status={status} elapsed_ms={elapsed} body_bytes={body_bytes}"
                );
            }
            Err(e) => {
                tracing::error!("[scheduler][rust.http.doget] failed: {e}");
            }
        }
    })
}

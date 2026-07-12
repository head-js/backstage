use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::time::Duration;

use crate::scheduler::task;

pub fn start() {
    spawn_loop("rust.log.timestamp", Duration::from_secs(5), task::run_timestamp_log);
    spawn_loop("rust.http.doget", Duration::from_secs(30), task::run_http_doget);
}

fn spawn_loop(task_id: &'static str, interval: Duration, job: fn() -> tokio::task::JoinHandle<()>) {
    let running = Arc::new(AtomicBool::new(false));

    tauri::async_runtime::spawn(async move {
        dispatch(task_id, Arc::clone(&running), job);

        let mut ticker = tokio::time::interval(interval);
        loop {
            ticker.tick().await;
            dispatch(task_id, Arc::clone(&running), job);
        }
    });
}

/// 尝试派发任务，若同一任务仍在运行则跳过。
///
/// 通过原子 running 标记防止并发堆积（如慢 HTTP 请求未完成时不再重复派发）。
fn dispatch(task_id: &'static str, running: Arc<AtomicBool>, job: fn() -> tokio::task::JoinHandle<()>) {
    if running.compare_exchange(false, true, Ordering::Acquire, Ordering::Relaxed).is_err() {
        tracing::warn!("[scheduler][{task_id}] skipped (already running)");
        return;
    }

    tokio::spawn(async move {
        job().await.ok();
        running.store(false, Ordering::Release);
    });
}

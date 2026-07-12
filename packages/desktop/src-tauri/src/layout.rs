#[cfg(target_os = "macos")]
use objc2_app_kit::NSWindow;
use std::{thread, time::Duration};
use tauri::{LogicalSize, PhysicalPosition};

const CONTENT_WIDTH: f64 = 992.0;
const MENU_FIRST_ITEM_TOP_OFFSET: f64 = 96.0;
#[cfg(target_os = "macos")]
const MENU_TITLEBAR_FALLBACK_HEIGHT: f64 = 34.0;

#[cfg(target_os = "macos")]
fn menu_titlebar_top_inset(menu: &tauri::WebviewWindow) -> f64 {
    let Ok(ns_window) = menu.ns_window() else {
        return MENU_TITLEBAR_FALLBACK_HEIGHT;
    };

    unsafe {
        let ns_window = ns_window as *mut NSWindow;
        let frame = (*ns_window).frame();
        let content = (*ns_window).contentLayoutRect();
        let inset = frame.size.height - (content.origin.y + content.size.height);
        if inset > 0.0 {
            inset
        } else {
            MENU_TITLEBAR_FALLBACK_HEIGHT
        }
    }
}

#[cfg(not(target_os = "macos"))]
fn menu_titlebar_top_inset(menu: &tauri::WebviewWindow) -> f64 {
    let Ok(outer_position) = menu.outer_position() else {
        return 0.0;
    };
    let Ok(inner_position) = menu.inner_position() else {
        return 0.0;
    };
    (inner_position.y - outer_position.y).max(0) as f64
}

pub(crate) fn sync_workspace_window(menu: &tauri::WebviewWindow, workspace: &tauri::WebviewWindow) {
    let Ok(outer_position) = menu.outer_position() else {
        return;
    };
    let Ok(outer_size) = menu.outer_size() else {
        return;
    };
    let Ok(scale_factor) = menu.scale_factor() else {
        return;
    };
    let content_width = (CONTENT_WIDTH * scale_factor).round() as i32;
    let bottom = outer_position.y + outer_size.height as i32;
    let top_inset = ((menu_titlebar_top_inset(menu) + MENU_FIRST_ITEM_TOP_OFFSET) * scale_factor)
        .round() as i32;
    let workspace_top = outer_position.y + top_inset;
    let content_height = (bottom - workspace_top).max(1) as f64 / scale_factor;

    let _ = workspace.set_position(PhysicalPosition::new(
        outer_position.x - content_width,
        workspace_top,
    ));
    let _ = workspace.set_size(LogicalSize::new(CONTENT_WIDTH, content_height));
}

// 防御性同步：在窗口事件未覆盖状态变化时，纠正 workspace 的可见性和布局。
pub(crate) fn start_workspace_follow_loop(
    menu: tauri::WebviewWindow,
    workspace: tauri::WebviewWindow,
) {
    thread::spawn(move || loop {
        let menu_hidden = !menu.is_visible().unwrap_or(false);
        let menu_minimized = menu.is_minimized().unwrap_or(false);

        if menu_hidden || menu_minimized {
            let _ = workspace.hide();
        } else if workspace.is_visible().unwrap_or(false) {
            sync_workspace_window(&menu, &workspace);
        }

        thread::sleep(Duration::from_millis(150));
    });
}

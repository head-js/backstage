mod cli;
mod commands;
mod layout;
mod scheduler;

#[cfg(target_os = "macos")]
use objc2_app_kit::NSWindow;

use tauri::{
    menu::{CheckMenuItemBuilder, MenuBuilder, MenuItemBuilder},
    tray::{MouseButton, MouseButtonState, TrayIconBuilder, TrayIconEvent},
    Manager, WindowEvent,
};

#[cfg(target_os = "windows")]
use windows::{
    Win32::Foundation::{COLORREF, HWND},
    Win32::UI::WindowsAndMessaging::{
        GetWindowLongPtrW, SetLayeredWindowAttributes, SetWindowLongPtrW, GWL_EXSTYLE, LWA_ALPHA,
        WS_EX_LAYERED,
    },
};

#[cfg(target_os = "windows")]
fn set_window_opacity(hwnd: HWND, opacity: u8) {
    let ex_style = unsafe { GetWindowLongPtrW(hwnd, GWL_EXSTYLE) };
    unsafe { SetWindowLongPtrW(hwnd, GWL_EXSTYLE, ex_style | WS_EX_LAYERED.0 as isize) };
    unsafe {
        let _ = SetLayeredWindowAttributes(hwnd, COLORREF(0), opacity, LWA_ALPHA);
    };
}

#[cfg(target_os = "macos")]
fn set_window_alpha(window: &tauri::WebviewWindow, alpha: f64) {
    let label = window.label().to_string();
    match window.ns_window() {
        Ok(ns_window) => unsafe {
            let ns_window = ns_window as *mut NSWindow;
            (*ns_window).setAlphaValue(alpha);
            let applied_alpha = (*ns_window).alphaValue();
            tracing::info!("[window] {label} native alpha={applied_alpha}");
        },
        Err(error) => {
            tracing::error!("[window] failed to access {label} NSWindow: {error}");
        }
    }
}

fn hide_workspace(workspace: &tauri::WebviewWindow) {
    let _ = workspace.hide();
}

fn show_menu(menu: &tauri::WebviewWindow, workspace: Option<&tauri::WebviewWindow>) {
    let _ = menu.show();
    let _ = menu.unminimize();
    if let Some(workspace) = workspace {
        layout::sync_workspace_window(menu, workspace);
    }
    let _ = menu.set_focus();
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_os::init())
        .plugin(tauri_plugin_store::Builder::default().build())
        .invoke_handler(tauri::generate_handler![
            commands::get_system_info,
            commands::run_command,
            commands::tauri_edge
        ])
        .setup(|app| {
            tracing_subscriber::fmt()
                .with_env_filter(
                    tracing_subscriber::EnvFilter::try_from_default_env()
                        .unwrap_or_else(|_| tracing_subscriber::EnvFilter::new("info")),
                )
                .init();

            scheduler::runtime::start();
            tracing::info!("[scheduler] started");

            let show_item = MenuItemBuilder::with_id("show", "显示窗口").build(app)?;
            let hide_item = MenuItemBuilder::with_id("hide", "隐藏窗口").build(app)?;
            let top_item = CheckMenuItemBuilder::with_id("toggle_top", "置顶")
                .checked(true)
                .build(app)?;
            let quit_item = MenuItemBuilder::with_id("quit", "退出").build(app)?;

            let menu = MenuBuilder::new(app)
                .items(&[&show_item, &hide_item, &top_item, &quit_item])
                .build()?;

            #[cfg(target_os = "macos")]
            if let Some(menu) = app.get_webview_window("menu") {
                set_window_alpha(&menu, 1.0);
            }

            #[cfg(target_os = "macos")]
            if let Some(workspace) = app.get_webview_window("workspace") {
                set_window_alpha(&workspace, 0.9);
            }

            if let (Some(menu_window), Some(workspace)) = (
                app.get_webview_window("menu"),
                app.get_webview_window("workspace"),
            ) {
                layout::sync_workspace_window(&menu_window, &workspace);
                layout::start_workspace_follow_loop(menu_window.clone(), workspace.clone());
                let menu_for_events = menu_window.clone();
                let workspace_for_events = workspace.clone();
                menu_window.on_window_event(move |event| match event {
                    WindowEvent::Moved(_)
                    | WindowEvent::Resized(_)
                    | WindowEvent::ScaleFactorChanged { .. } => {
                        if menu_for_events.is_minimized().unwrap_or(false) {
                            hide_workspace(&workspace_for_events);
                            return;
                        }
                        layout::sync_workspace_window(&menu_for_events, &workspace_for_events);
                    }
                    WindowEvent::CloseRequested { api, .. } => {
                        api.prevent_close();
                        hide_workspace(&workspace_for_events);
                        let _ = menu_for_events.hide();
                    }
                    WindowEvent::Destroyed => {
                        hide_workspace(&workspace_for_events);
                    }
                    _ => {}
                });
            }

            if let Some(workspace) = app.get_webview_window("workspace") {
                hide_workspace(&workspace);
            }

            TrayIconBuilder::new()
                .icon(app.default_window_icon().unwrap().clone())
                .tooltip("desktop")
                .menu(&menu)
                .on_menu_event(|app, event| match event.id.as_ref() {
                    "show" => {
                        if let Some(menu_window) = app.get_webview_window("menu") {
                            let workspace = app.get_webview_window("workspace");
                            show_menu(&menu_window, workspace.as_ref());
                        }
                    }
                    "hide" => {
                        if let Some(workspace) = app.get_webview_window("workspace") {
                            hide_workspace(&workspace);
                        }
                        if let Some(menu_window) = app.get_webview_window("menu") {
                            let _ = menu_window.hide();
                        }
                    }
                    "toggle_top" => {
                        if let Some(menu_window) = app.get_webview_window("menu") {
                            if let Ok(is_top) = menu_window.is_always_on_top() {
                                let next = !is_top;
                                let _ = menu_window.set_always_on_top(next);
                                if let Some(workspace) = app.get_webview_window("workspace") {
                                    let _ = workspace.set_always_on_top(next);
                                }
                            }
                        }
                    }
                    "quit" => {
                        app.exit(0);
                    }
                    _ => {}
                })
                .on_tray_icon_event(|tray, event| {
                    if let TrayIconEvent::Click {
                        button: MouseButton::Left,
                        button_state: MouseButtonState::Up,
                        ..
                    } = event
                    {
                        let app = tray.app_handle();
                        if let Some(menu_window) = app.get_webview_window("menu") {
                            let workspace = app.get_webview_window("workspace");
                            show_menu(&menu_window, workspace.as_ref());
                        }
                    }
                })
                .build(app)?;

            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

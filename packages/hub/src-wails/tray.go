package main

import (
	"embed"
	"runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed build/darwin/trayTemplateIcon.png build/windows/trayIcon.png build/windows/trayIconDark.png
var trayIconFS embed.FS

// 本机工具链对 []byte 形式的 //go:embed 会误报 "imported and not used"，故走 embed.FS
var (
	trayIconDarwin      = mustReadTrayIcon("build/darwin/trayTemplateIcon.png")
	trayIconWindows     = mustReadTrayIcon("build/windows/trayIcon.png")
	trayIconWindowsDark = mustReadTrayIcon("build/windows/trayIconDark.png")
)

func mustReadTrayIcon(name string) []byte {
	data, err := trayIconFS.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return data
}

func setupTray(app *application.App, window *application.WebviewWindow) {
	var quitting bool

	showWindow := func() {
		window.Show()
		window.Focus()
		setDockIconVisible(true)
	}
	hideWindow := func() {
		window.Hide()
		setDockIconVisible(false)
	}

	tray := app.SystemTray.New()
	if runtime.GOOS == "darwin" {
		tray.SetTemplateIcon(trayIconDarwin)
	} else {
		tray.SetIcon(trayIconWindows)
		tray.SetDarkModeIcon(trayIconWindowsDark)
	}
	tray.SetTooltip("Backstage Hub")

	menu := app.NewMenu()
	menu.Add("显示窗口").OnClick(func(*application.Context) {
		showWindow()
	})
	menu.Add("隐藏窗口").OnClick(func(*application.Context) {
		hideWindow()
	})
	alwaysOnTop := menu.AddCheckbox("置顶", true)
	alwaysOnTop.OnClick(func(*application.Context) {
		window.SetAlwaysOnTop(alwaysOnTop.Checked())
	})
	menu.AddSeparator()
	menu.Add("退出").OnClick(func(*application.Context) {
		quitting = true
		app.Quit()
	})
	tray.SetMenu(menu)

	tray.OnClick(showWindow)
	tray.OnRightClick(tray.OpenMenu)

	window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if quitting {
			return
		}
		event.Cancel()
		window.Hide()
		setDockIconVisible(false)
	})

	if runtime.GOOS == "darwin" {
		app.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(*application.ApplicationEvent) {
			showWindow()
		})
	}
}

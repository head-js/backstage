//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit
#include <AppKit/AppKit.h>

static void dockSetVisible(int visible) {
	NSApplicationActivationPolicy policy = visible
		? NSApplicationActivationPolicyRegular
		: NSApplicationActivationPolicyAccessory;
	[[NSApplication sharedApplication] setActivationPolicy:policy];
}
*/
import "C"

import "github.com/wailsapp/wails/v3/pkg/application"

// setDockIconVisible 必须在主线程执行 setActivationPolicy，经 InvokeSync 投递
func setDockIconVisible(visible bool) {
	application.InvokeSync(func() {
		if visible {
			C.dockSetVisible(1)
		} else {
			C.dockSetVisible(0)
		}
	})
}

// +build windows

package utils

import "github.com/lxn/win"

func getMonitorScale() int {
	hwnd := win.GetDesktopWindow()
	return int(win.GetDpiForWindow(hwnd)) / 96
}

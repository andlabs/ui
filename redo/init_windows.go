// 11 july 2014

package ui

import (
	"fmt"
)

var (
	hInstance uintptr
	nCmdShow int

	hDefaultIcon uintptr
	hArrowCursor uintptr
)

func getWinMainParams() (err error) {
	hInstance, err = f_GetModuleHandleW(nil)
	if hInstance == 0 {
		return fmt.Errorf("error getting hInstance: %v", err)
	}

	var info s_STARTUPINFOW

	f_GetStartupInfoW(&info)
	if info.dwFlags & c_STARTF_USESHOWWINDOW != 0 {
		nCmdShow = int(info.wShowWindow)
	} else {
		nCmdShow = c_SW_SHOWDEFAULT
	}

	return nil
}

// TODO move to common_windows.go
var hNULL uintptr = 0

func loadIconsCursors() (err error) {
	hDefaultIcon, err = f_LoadIconW(hNULL, c_IDI_APPLICATION)
	if hDefaultIcon == hNULL {
		return fmt.Errorf("error loading default icon: %v", err)
	}
	hArrowCursor, err = f_LoadCursorW(hNULL, c_IDC_ARROW)
	if hArrowCursor == hNULL {
		return fmt.Errorf("error loading arrow (default) cursor: %v", err)
	}
	return nil
}

func initWindows() error {
	if err := getWinMainParams(); err != nil {
		return fmt.Errorf("error getting WinMain() parameters: %v", err)
	}
	if err := loadIconsCursors(); err != nil {
		return fmt.Errorf("error loading standard/default icons and cursors: %v", err)
	}
	return nil
}

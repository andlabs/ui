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

func initWindows() error {
	if err := getWinMainParams(); err != nil {
		return fmt.Errorf("error getting WinMain() parameters: %v", err)
	}
	return nil
}

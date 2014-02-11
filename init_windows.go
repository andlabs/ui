// 8 february 2014
package main

import (
	"fmt"
//	"syscall"
	"unsafe"
)

var (
	hInstance		HANDLE
	nCmdShow	int
	// TODO font
)

// TODO is this trick documented in MSDN?
func getWinMainhInstance() (err error) {
	r1, _, err := kernel32.NewProc("GetModuleHandleW").Call(uintptr(NULL))
	if r1 == 0 {		// failure
		return err
	}
	hInstance = HANDLE(r1)
	return nil
}

// TODO this is what MinGW-w64's crt (svn revision xxx) does; is it best? is any of this documented anywhere on MSDN?
// TODO I highly doubt Windows API functions ever not fail, so figure out what to do should an error actually occur
func getWinMainnCmdShow() (nCmdShow int, err error) {
	var info struct {
		cb				uint32
		lpReserved		*uint16
		lpDesktop			*uint16
		lpTitle			*uint16
		dwX				uint32
		dwY				uint32
		dwXSize			uint32
		dwYSzie			uint32
		dwXCountChars	uint32
		dwYCountChars	uint32
		dwFillAttribute		uint32
		dwFlags			uint32
		wShowWindow		uint16
		cbReserved2		uint16
		lpReserved2		*byte
		hStdInput			HANDLE
		hStdOutput		HANDLE
		hStdError			HANDLE
	}
	const _STARTF_USESHOWWINDOW = 0x00000001

	// does not fail according to MSDN
	kernel32.NewProc("GetStartupInfoW").Call(uintptr(unsafe.Pointer(&info)))
	if info.dwFlags & _STARTF_USESHOWWINDOW != 0 {
		nCmdShow = int(info.wShowWindow)
		return nil
	}
	nCmdShow = _SW_SHOWDEFAULT
	return nil
}

func doWindowsInit() (err error) {
	err = getWinMainhInstance()
	if err != nil {
		return fmt.Errorf("error getting WinMain hInstance: %v", err)
	}
	err = getWinMainnCmdShow()
	if err != nil {
		return fmt.Errorf("error getting WinMain nCmdShow: %v", err)
	}
	err = registerStdWndClass()
	if err != nil {
		reteurn fmt.Errorf("error registering standard window class: %v", err)
	}
	// TODO others
	return nil		// all ready to go
}

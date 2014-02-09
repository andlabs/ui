// 8 february 2014
package main

import (
//	"syscall"
	"unsafe"
)

// this provides the hInstance and nCmdShow that are normally passed to WinMain()

const (
	STARTF_USESHOWWINDOW = 0x00000001
)

var (
	getModuleHandle = kernel32.NewProc("GetModuleHandleW")
	getStartupInfo = kernel32.NewProc("GetStartupInfoW")
)

// TODO is this trick documented in MSDN?
func getWinMainhInstance() (hInstance HANDLE, err error) {
	r1, _, err := getModuleHandle.Call(uintptr(NULL))
	if r1 == 0 {
		return NULL, err
	}
	return HANDLE(r1), nil
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

	// does not fail according to MSDN
	getStartupInfo.Call(uintptr(unsafe.Pointer(&info)))
	if info.dwFlags & STARTF_USESHOWWINDOW != 0 {
		return int(info.wShowWindow), nil
	}
	return SW_SHOWDEFAULT, nil
}

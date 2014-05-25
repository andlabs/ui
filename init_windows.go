// 8 february 2014

package ui

import (
	"fmt"
//	"syscall"
	"unsafe"
)

var (
	hInstance		_HANDLE
	nCmdShow	int
)

// TODO is this documented?
func getWinMainhInstance() (err error) {
	r1, _, err := kernel32.NewProc("GetModuleHandleW").Call(uintptr(_NULL))
	if r1 == 0 {		// failure
		return err
	}
	hInstance = _HANDLE(r1)
	return nil
}

// this is what MinGW-w64 does (for instance, http://sourceforge.net/p/mingw-w64/code/6604/tree/trunk/mingw-w64-crt/crt/crtexe.c#l320); Burgundy in irc.freenode.net/#winapi said that the Visual C++ runtime does this too
func getWinMainnCmdShow() {
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
		hStdInput			_HANDLE
		hStdOutput		_HANDLE
		hStdError			_HANDLE
	}

	// does not fail according to MSDN
	kernel32.NewProc("GetStartupInfoW").Call(uintptr(unsafe.Pointer(&info)))
	if info.dwFlags & _STARTF_USESHOWWINDOW != 0 {
		nCmdShow = int(info.wShowWindow)
	} else {
		nCmdShow = _SW_SHOWDEFAULT
	}
}

func doWindowsInit() (err error) {
	err = getWinMainhInstance()
	if err != nil {
		return fmt.Errorf("error getting WinMain hInstance: %v", err)
	}
	getWinMainnCmdShow()
	err = initWndClassInfo()
	if err != nil {
		return fmt.Errorf("error initializing standard window class auxiliary info: %v", err)
	}
	err = getStandardWindowFonts()
	if err != nil {
		return fmt.Errorf("error getting standard window fonts: %v", err)
	}
	err = initCommonControls()
	if err != nil {
		return fmt.Errorf("error initializing Common Controls (comctl32.dll): %v", err)
	}
	// others go here
	return nil		// all ready to go
}

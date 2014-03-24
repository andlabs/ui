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
	gdiplusToken	uintptr
)

// TODO is this trick documented in MSDN?
func getWinMainhInstance() (err error) {
	r1, _, err := kernel32.NewProc("GetModuleHandleW").Call(uintptr(_NULL))
	if r1 == 0 {		// failure
		return err
	}
	hInstance = _HANDLE(r1)
	return nil
}

// TODO this is what MinGW-w64's crt (svn revision TODO) does; is it best? is any of this documented anywhere on MSDN?
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
	const _STARTF_USESHOWWINDOW = 0x00000001

	// does not fail according to MSDN
	kernel32.NewProc("GetStartupInfoW").Call(uintptr(unsafe.Pointer(&info)))
	if info.dwFlags & _STARTF_USESHOWWINDOW != 0 {
		nCmdShow = int(info.wShowWindow)
	} else {
		nCmdShow = _SW_SHOWDEFAULT
	}
}

func initGDIPlus() (err error) {
	var gdiplusInit struct {
		GdiplusVersion				uint32
		DebugEventCallback			uintptr
		SuppressBackgroundThread	int32		// originally BOOL
		SuppressExternalCodecs		int32		// originally BOOL
	}

	gdiplusInit.GdiplusVersion = 1		// required
	// TODO suppress external codecs?
	r1, _, err := gdiplus.NewProc("GdiplusStartup").Call(
		uintptr(unsafe.Pointer(&gdiplusToken)),
		uintptr(unsafe.Pointer(&gdiplusInit)),
		uintptr(0))		// we use the GDI+ thread so no need for output info
	if r1 != 0 {			// failure
		return fmt.Errorf("error initializing GDI+ (GDI+ error code %d; windows last error %v)", r1, err)
	}
	return nil
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
	err = initGDIPlus()
	if err != nil {
		return fmt.Errorf("error initializing GDI+ (gdiplus.dll): %v", err)
	}
	// TODO others
	return nil		// all ready to go
}

func doWindowsQuitStuff() {
	gdiplus.NewProc("GdiplusShutdown").Call(gdiplusToken)		// returns void according to MSDN
}

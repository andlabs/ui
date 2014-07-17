// 12 july 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
)

// #cgo LDFLAGS: -luser32 -lkernel32 -lgdi32
// #include "winapi_windows.h"
import "C"

var msgwin C.HWND

func uiinit() error {
	if err := initWindows(); err != nil {
		return fmt.Errorf("error initializing package ui on Windows: %v", err)
	}
	if err := initCommonControls(); err != nil {
		return fmt.Errorf("error initializing comctl32.dll version 6: %v", err)
	}
	if err := makemsgwin(); err != nil {
		return fmt.Errorf("error creating message-only window: %v", err)
	}
	if err := makeWindowWindowClass(); err != nil {
		return fmt.Errorf("error creating Window window class: %v", err)
	}
	return nil
}

func uimsgloop() {
	C.uimsgloop()
}

func uistop() {
	f_PostQuitMessage(0)
}

func issue(req *Request) {
	C.issue(unsafe.Pointer(req))
}

func makemsgwin() error {
	var errmsg *C.char

	err := C.makemsgwin(&errmsg)
	if err != 0 || errmsg != nil {
		return fmt.Errorf("%s: %v", C.GoString(errmsg), syscall.Errno(err))
	}
	return nil
}

//export xperform
func xperform(xreq unsafe.Pointer) {
	req := (*Request)(xreq)
	perform(req)
}

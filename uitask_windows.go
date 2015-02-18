// 12 july 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
)

// #cgo CFLAGS: --std=c99
// #cgo LDFLAGS: -luser32 -lkernel32 -lgdi32 -luxtheme -lmsimg32 -lcomdlg32 -lole32 -loleaut32 -loleacc -luuid
// #include "winapi_windows.h"
import "C"

var msgwin C.HWND

func uiinit() error {
	var errmsg *C.char

	errcode := C.initWindows(&errmsg)
	if errcode != 0 || errmsg != nil {
		return fmt.Errorf("error initializing package ui on Windows: %s: %v", C.GoString(errmsg), syscall.Errno(errcode))
	}
	errmsg = nil
	errcode = C.initCommonControls(&errmsg)
	if errcode != 0 || errmsg != nil {
		return fmt.Errorf("error initializing comctl32.dll: %s: %v", C.GoString(errmsg), syscall.Errno(errcode))
	}
	if err := makemsgwin(); err != nil {
		return fmt.Errorf("error creating message-only window: %v", err)
	}
	if err := makeWindowWindowClass(); err != nil {
		return fmt.Errorf("error creating Window window class: %v", err)
	}
	if err := makeAreaWindowClass(); err != nil {
		return fmt.Errorf("error creating Area window class: %v", err)
	}
	// this depends on the common controls having been initialized already
	C.doInitTable()
	return nil
}

func uimsgloop() {
	C.uimsgloop()
}

func uistop() {
	C.PostQuitMessage(0)
}

func issue(f *func()) {
	C.issue(unsafe.Pointer(f))
}

func makemsgwin() error {
	var errmsg *C.char

	err := C.makemsgwin(&errmsg)
	if err != 0 || errmsg != nil {
		return fmt.Errorf("%s: %v", C.GoString(errmsg), syscall.Errno(err))
	}
	return nil
}

//export doissue
func doissue(fp unsafe.Pointer) {
	perform(fp)
}

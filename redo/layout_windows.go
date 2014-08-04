// 4 august 2014

package ui

// TODO clean this up relative to window_windows.go

import (
	"fmt"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type layout struct {
	hwnd		C.HWND

	closing		*event

	*sizer
}

func makeContainerWindowClass() error {
	var errmsg *C.char

	err := C.makeContainerWindowClass(&errmsg)
	if err != 0 || errmsg != nil {
		return fmt.Errorf("%s: %v", C.GoString(errmsg), syscall.Errno(err))
	}
	return nil
}

func newLayout(title string, width int, height int, child C.BOOL, control Control) *layout {
	l := &layout{
		sizer:		new(sizer),
	}
	hwnd := C.newContainer(unsafe.Pointer(l))
	if hwnd != l.hwnd {
		panic(fmt.Errorf("inconsistency: hwnd returned by CreateWindowEx() (%p) and hwnd stored in container (%p) differ", hwnd, l.hwnd))
	}
	l.child = control
	l.child.setParent(&controlParent{l.hwnd})
	return l
}

func (l *layout) setParent(p *controlParent) {
	C.controlSetParent(l.hwnd, p.hwnd)
}

//export storeContainerHWND
func storeContainerHWND(data unsafe.Pointer, hwnd C.HWND) {
	l := (*layout)(data)
	l.hwnd = hwnd
}

//export containerResize
func containerResize(data unsafe.Pointer, r *C.RECT) {
	l := (*layout)(data)
	// the origin of the window's content area is always (0, 0), but let's use the values from the RECT just to be safe
	l.resize(int(r.left), int(r.top), int(r.right - r.left), int(r.bottom - r.top))
}

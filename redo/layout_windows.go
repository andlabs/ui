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

func newLayout(title string, width int, height int, child C.BOOL, control Control) *layout {
	l := &layout{
		// hwnd set in WM_CREATE handler
		closing:		newEvent(),
		sizer:		new(sizer),
	}
	hwnd := C.newWindow(toUTF16(title), C.int(width), C.int(height), child, unsafe.Pointer(l))
	if hwnd != l.hwnd {
		panic(fmt.Errorf("inconsistency: hwnd returned by CreateWindowEx() (%p) and hwnd stored in window/layout (%p) differ", hwnd, l.hwnd))
	}
	l.child = control
	l.child.setParent(&controlParent{l.hwnd})
	return l
}

func (l *layout) setParent(p *controlParent) {
	C.controlSetParent(l.hwnd, p.hwnd)
}

//export storeWindowHWND
func storeWindowHWND(data unsafe.Pointer, hwnd C.HWND) {
	l := (*layout)(data)
	l.hwnd = hwnd
}

//export windowResize
func windowResize(data unsafe.Pointer, r *C.RECT) {
	l := (*layout)(data)
	// the origin of the window's content area is always (0, 0), but let's use the values from the RECT just to be safe
	l.resize(int(r.left), int(r.top), int(r.right - r.left), int(r.bottom - r.top))
}

//export windowClosing
func windowClosing(data unsafe.Pointer) {
	l := (*layout)(data)
	close := l.closing.fire()
	if close {
		C.windowClose(l.hwnd)
	}
}

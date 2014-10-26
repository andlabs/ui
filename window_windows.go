// 12 july 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type window struct {
	hwnd        C.HWND
	shownbefore bool

	closing *event

	child			Control
	margined		bool
}

func makeWindowWindowClass() error {
	var errmsg *C.char

	err := C.makeWindowWindowClass(&errmsg)
	if err != 0 || errmsg != nil {
		return fmt.Errorf("%s: %v", C.GoString(errmsg), syscall.Errno(err))
	}
	return nil
}

func newWindow(title string, width int, height int, control Control) *window {
	w := &window{
		closing:   newEvent(),
		child:	control,
	}
	w.hwnd = C.newWindow(toUTF16(title), C.int(width), C.int(height), unsafe.Pointer(w))
	hresult := C.EnableThemeDialogTexture(w.hwnd, C.ETDT_ENABLE|C.ETDT_USETABTEXTURE)
	if hresult != C.S_OK {
		panic(fmt.Errorf("error setting tab background texture on Window; HRESULT: 0x%X", hresult))
	}
	w.child.setParent(&controlParent{w.hwnd})
	return w
}

func (w *window) Title() string {
	return getWindowText(w.hwnd)
}

func (w *window) SetTitle(title string) {
	C.setWindowText(w.hwnd, toUTF16(title))
}

func (w *window) Show() {
	if !w.shownbefore {
		C.ShowWindow(w.hwnd, C.nCmdShow)
		C.updateWindow(w.hwnd)
		w.shownbefore = true
	} else {
		C.ShowWindow(w.hwnd, C.SW_SHOW)
	}
}

func (w *window) Hide() {
	C.ShowWindow(w.hwnd, C.SW_HIDE)
}

func (w *window) Close() {
	C.windowClose(w.hwnd)
}

func (w *window) OnClosing(e func() bool) {
	w.closing.setbool(e)
}

func (w *window) Margined() bool {
	return w.margined
}

func (w *window) SetMargined(margined bool) {
	w.margined = margined
}

//export windowResize
func windowResize(data unsafe.Pointer, r *C.RECT) {
	w := (*window)(data)
	d := beginResize(w.hwnd)
	if w.margined {
		marginRectDLU(r, marginDialogUnits, marginDialogUnits, marginDialogUnits, marginDialogUnits, d)
	}
	w.child.resize(int(r.left), int (r.top), int(r.right - r.left), int(r.bottom - r.top), d)
}

//export windowClosing
func windowClosing(data unsafe.Pointer) {
	w := (*window)(data)
	close := w.closing.fire()
	if close {
		C.windowClose(w.hwnd)
	}
}

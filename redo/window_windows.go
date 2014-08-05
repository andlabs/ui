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
	hwnd		C.HWND
	shownbefore	bool

	closing		*event

	*container
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
		// hwnd set in WM_CREATE handler
		closing:		newEvent(),
		container:		newContainer(control),
	}
	hwnd := C.newWindow(toUTF16(title), C.int(width), C.int(height), unsafe.Pointer(w))
	if hwnd != w.hwnd {
		panic(fmt.Errorf("inconsistency: hwnd returned by CreateWindowEx() (%p) and hwnd stored in Window (%p) differ", hwnd, w.hwnd))
	}
	hresult := C.EnableThemeDialogTexture(w.hwnd, C.ETDT_ENABLE | C.ETDT_USETABTEXTURE)
	if hresult != C.S_OK {
		panic(fmt.Errorf("error setting tab background texture on Window; HRESULT: 0x%X", hresult))
	}
	w.container.setParent(&controlParent{w.hwnd})
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

//export storeWindowHWND
func storeWindowHWND(data unsafe.Pointer, hwnd C.HWND) {
	w := (*window)(data)
	w.hwnd = hwnd
}

//export windowResize
func windowResize(data unsafe.Pointer, r *C.RECT) {
	w := (*window)(data)
	// the origin of the window's content area is always (0, 0), but let's use the values from the RECT just to be safe
	w.container.move(r)
}

//export windowClosing
func windowClosing(data unsafe.Pointer) {
	w := (*window)(data)
	close := w.closing.fire()
	if close {
		C.windowClose(w.hwnd)
	}
}

// 12 july 2014

package ui

import (
	"fmt"
	"syscall"
)

// #include "winapi_windows.h"
import "C"

type window struct {
	*layout
	shownbefore	bool
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
		layout:	newLayout(title, width, height, C.FALSE, control),
	}
	// TODO keep?
	hresult := C.EnableThemeDialogTexture(w.hwnd, C.ETDT_ENABLE | C.ETDT_USETABTEXTURE)
	if hresult != C.S_OK {
		panic(fmt.Errorf("error setting tab background texture on Window; HRESULT: 0x%X", hresult))
	}
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

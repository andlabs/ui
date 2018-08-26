// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Window is a Control that represents a top-level window.
// A Window contains one child Control that occupies the
// entirety of the window. Though a Window is a Control,
// a Window cannot be the child of another Control.
type Window struct {
	ControlBase
	w	*C.uiWindow
	child		Control
	onClosing		func(w *Window) bool
}

// NewWindow creates a new Window.
func NewWindow(title string, width int, height int, hasMenubar bool) *Window {
	w := new(Window)

	ctitle := C.CString(title)
	w.w = C.uiNewWindow(ctitle, C.int(width), C.int(height), frombool(hasMenubar))
	freestr(ctitle)

	C.pkguiWindowOnClosing(w.w)

	w.ControlBase = NewControlBase(w, uintptr(unsafe.Pointer(w.w)))
	return w
}

// Destroy destroys the Window. If the Window has a child,
// Destroy calls Destroy on that as well.
func (w *Window) Destroy() {
	w.Hide()		// first hide the window, in case anything in the below if statement forces an immediate redraw
	if w.child != nil {
		c := w.child
		w.SetChild(nil)
		c.Destroy()
	}
	w.ControlBase.Destroy()
}

// Title returns the Window's title.
func (w *Window) Title() string {
	ctitle := C.uiWindowTitle(w.w)
	title := C.GoString(ctitle)
	C.uiFreeText(ctitle)
	return title
}

// SetTitle sets the Window's title to title.
func (w *Window) SetTitle(title string) {
	ctitle := C.CString(title)
	C.uiWindowSetTitle(w.w, ctitle)
	freestr(ctitle)
}

// TODO ContentSize
// TODO SetContentSize
// TODO Fullscreen
// TODO SetFullscreen
// TODO OnContentSizeChanged

// OnClosing registers f to be run when the user clicks the Window's
// close button. Only one function can be registered at a time.
// If f returns true, the window is destroyed with the Destroy method.
// If f returns false, or if OnClosing is never called, the window is not
// destroyed and is kept visible.
func (w *Window) OnClosing(f func(*Window) bool) {
	w.onClosing = f
}

//export pkguiDoWindowOnClosing
func pkguiDoWindowOnClosing(ww *C.uiWindow, data unsafe.Pointer) C.int {
	w := ControlFromLibui(uintptr(unsafe.Pointer(ww))).(*Window)
	if w.onClosing == nil {
		return 0
	}
	if w.onClosing(w) {
		w.Destroy()
	}
	return 0
}

// Borderless returns whether the Window is borderless.
func (w *Window) Borderless() bool {
	return tobool(C.uiWindowBorderless(w.w))
}

// SetBorderless sets the Window to be borderless or not.
func (w *Window) SetBorderless(borderless bool) {
	C.uiWindowSetBorderless(w.w, frombool(borderless))
}

// SetChild sets the Window's child to child. If child is nil, the Window
// will not have a child.
func (w *Window) SetChild(child Control) {
	w.child = child
	c := (*C.uiControl)(nil)
	if w.child != nil {
		c = touiControl(w.child.LibuiControl())
	}
	C.uiWindowSetChild(w.w, c)
}

// Margined returns whether the Window has margins around its child.
func (w *Window) Margined() bool {
	return tobool(C.uiWindowMargined(w.w))
}

// SetMargined controls whether the Window has margins around its
// child. The size of the margins are determined by the OS and its
// best practices.
func (w *Window) SetMargined(margined bool) {
	C.uiWindowSetMargined(w.w, frombool(margined))
}

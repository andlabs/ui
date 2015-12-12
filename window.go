// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern int doOnClosing(uiWindow *, void *);
// static inline void realuiWindowOnClosing(uiWindow *w)
// {
// 	uiWindowOnClosing(w, doOnClosing, NULL);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var windows = make(map[*C.uiWindow]*Window)

// Window is a Control that represents a top-level window.
// A Window contains one child Control that occupies the
// entirety of the window. Though a Window is a Control,
// a Window cannot be the child of another Control.
type Window struct {
	c	*C.uiControl
	w	*C.uiWindow

	child		Control

	onClosing		func(w *Window) bool
}

// NewWindow creates a new Window.
func NewWindow(title string, width int, height int, hasMenubar bool) *Window {
	w := new(Window)

	ctitle := C.CString(title)
	// TODO wait why did I make these ints and not intmax_ts?
	w.w = C.uiNewWindow(ctitle, C.int(width), C.int(height), frombool(hasMenubar))
	w.c = (*C.uiControl)(unsafe.Pointer(w.w))
	freestr(ctitle)

	C.realuiWindowOnClosing(w.w)
	windows[w.w] = w

	return w
}

// Destroy destroys the Window. If the Window has a child,
// Destroy calls Destroy on that as well.
func (w *Window) Destroy() {
	// first hide ourselves
	w.Hide()
	// get rid of the child
	if w.child != nil {
		c := w.child
		w.SetChild(nil)
		c.Destroy()
	}
	// unregister events
	delete(windows, w.w)
	// and finally destroy ourselves
	C.uiControlDestroy(w.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (w *Window) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(w.c))
}

// Handle returns the OS-level handle associated with this Window.
// On Windows this is an HWND of a libui-internal class.
// On GTK+ this is a pointer to a GtkWindow.
// On OS X this is a pointer to a NSWindow.
func (w *Window) Handle() uintptr {
	return uintptr(C.uiControlHandle(w.c))
}

// Show shows the Window. It uses the OS conception of "presenting"
// the Window, whatever that may be on a given OS.
func (w *Window) Show() {
	C.uiControlShow(w.c)
}

// Hide hides the Window.
func (w *Window) Hide() {
	C.uiControlHide(w.c)
}

// Enable enables the Window.
func (w *Window) Enable() {
	C.uiControlEnable(w.c)
}

// Disable disables the Window.
func (w *Window) Disable() {
	C.uiControlDisable(w.c)
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

// OnClosing registers f to be run when the user clicks the Window's
// close button. Only one function can be registered at a time.
// If f returns true, the window is destroyed with the Destroy method.
// If f returns false, or if OnClosing is never called, the window is not
// destroyed and is kept visible.
func (w *Window) OnClosing(f func(*Window) bool) {
	w.onClosing = f
}

//export doOnClosing
func doOnClosing(ww *C.uiWindow, data unsafe.Pointer) C.int {
	w := windows[ww]
	if w.onClosing == nil {
		return 0
	}
	if w.onClosing(w) {
		w.Destroy()
	}
	return 0
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

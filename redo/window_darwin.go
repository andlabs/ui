// 8 july 2014

package ui

import (
	"unsafe"
"fmt"
)

// #include "objc_darwin.h"
import "C"

type window struct {
	id		C.id

	child		Control

	closing	*event

	spaced	bool
}

func newWindow(title string, width int, height int) *window {
	id := C.newWindow(C.intptr_t(width), C.intptr_t(height))
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.windowSetTitle(id, ctitle)
	w := &window{
		id:		id,
		closing:	newEvent(),
	}
	C.windowSetDelegate(id, unsafe.Pointer(w))
	return w
}

func (w *window) SetControl(control Control) {
	if w.child != nil {		// unparent existing control
		w.child.unparent()
	}
	control.unparent()
	control.parent(w)
	w.child = control
}

func (w *window) Title() string {
	return C.GoString(C.windowTitle(w.id))
}

func (w *window) SetTitle(title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.windowSetTitle(w.id, ctitle)
}

func (w *window) Show() {
	C.windowShow(w.id)
}

func (w *window) Hide() {
	C.windowHide(w.id)
}

func (w *window) Close() {
	C.windowClose(w.id)
}

func (w *window) OnClosing(e func() bool) {
	w.closing.setbool(e)
}

//export windowClosing
func windowClosing(xw unsafe.Pointer) C.BOOL {
	w := (*window)(unsafe.Pointer(xw))
	close := w.closing.fire()
	if close {
		return C.YES
	}
	return C.NO
}

//export windowResized
func windowResized(xw unsafe.Pointer, width C.uintptr_t, height C.uintptr_t) {
	w := (*window)(unsafe.Pointer(xw))
	w.doresize(int(width), int(height))
	fmt.Printf("new size %d x %d\n", width, height)
}

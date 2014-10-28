// 8 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type window struct {
	id C.id

	closing *event

	child			Control
	container		*container
}

func newWindow(title string, width int, height int, control Control) *window {
	id := C.newWindow(C.intptr_t(width), C.intptr_t(height))
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.windowSetTitle(id, ctitle)
	w := &window{
		id:        id,
		closing:   newEvent(),
		child:		control,
	}
	C.windowSetDelegate(w.id, unsafe.Pointer(w))
	w.container = newContainer(w.child.resize)
	w.child.setParent(w.container.parent())
	C.windowSetContentView(w.id, w.container.id)
	// trigger an initial resize
	return w
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
	// TODO we need a dummy resize here because things might not be in the right place
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

func (w *window) Margined() bool {
	return w.container.margined
}

func (w *window) SetMargined(margined bool) {
	w.container.margined = margined
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

// no need for windowResized; the child container takes care of that

// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type widgetbase struct {
	id		C.id
	notnew	bool		// to prevent unparenting a new control
	floating	bool
}

func newWidget(id C.id) *widgetbase {
	return &widgetbase{
		id:	id,
	}
}

// these few methods are embedded by all the various Controls since they all will do the same thing

func (w *widgetbase) unparent() {
	if w.notnew {
		// redrawing the old window handled by C.unparent()
		C.unparent(w.id)
		w.floating = true
	}
}

func (w *widgetbase) parent(win *window) {
	// redrawing the new window handled by C.parent()
	C.parent(w.id, win.id, toBOOL(w.floating))
	w.floating = false
	w.notnew = true
}

type button struct {
	*widgetbase
	clicked		*event
}

func finishNewButton(id C.id, text string) *button {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	b := &button{
		widgetbase:	newWidget(id),
		clicked:		newEvent(),
	}
	C.buttonSetText(b.id, ctext)
	C.buttonSetDelegate(b.id, unsafe.Pointer(b))
	return b
}

func newButton(text string) *button {
	return finishNewButton(C.newButton(), text)
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

//export buttonClicked
func buttonClicked(xb unsafe.Pointer) {
	b := (*button)(unsafe.Pointer(xb))
	b.clicked.fire()
	println("button clicked")
}

func (b *button) Text() string {
	return C.GoString(C.buttonText(b.id))
}

func (b *button) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.buttonSetText(b.id, ctext)
}

type checkbox struct {
	*button
}

func newCheckbox(text string) *checkbox {
	return &checkbox{
		button:	finishNewButton(C.newCheckbox(), text),
	}
}

// we don't need to define our own event here; we can just reuse Button's
// (it's all target-action anyway)

type (c *checkbox) Checked() bool {
	return fromBOOL(C.checkboxChecked(c.id))
}

type (c *checkbox) SetChecked(checked bool) {
	C.checkboxSetChecked(c.id, toBOOL(checked))
}

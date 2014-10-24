// 24 october 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type textbox struct {
	*scroller
}

func newTextbox() Textbox {
	id := C.newTextbox()
	t := &textbox{
		scroller:		newScroller(id, true),		// border on Textbox (TODO confirm type)
	}
	// TODO preferred size
	return t
}

func (t *textbox) Text() string {
	return C.GoString(C.textboxText(t.id))
}

func (t *textbox) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.textboxSetText(t.id, ctext)
}

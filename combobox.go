// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Combobox is a Control that represents a drop-down list of strings
// that the user can choose one of at any time. For a Combobox that
// users can type values into, see EditableCombobox.
type Combobox struct {
	ControlBase
	c	*C.uiCombobox
	onSelected		func(*Combobox)
}

// NewCombobox creates a new Combobox.
func NewCombobox() *Combobox {
	c := new(Combobox)

	c.c = C.uiNewCombobox()

	C.pkguiComboboxOnSelected(c.c)

	c.ControlBase = NewControlBase(c, uintptr(unsafe.Pointer(c.c)))
	return c
}

// Append adds the named item to the end of the Combobox.
func (c *Combobox) Append(text string) {
	ctext := C.CString(text)
	C.uiComboboxAppend(c.c, ctext)
	freestr(ctext)
}

// Selected returns the index of the currently selected item in the
// Combobox, or -1 if nothing is selected.
func (c *Combobox) Selected() int {
	return int(C.uiComboboxSelected(c.c))
}

// SetSelected sets the currently selected item in the Combobox
// to index. If index is -1 no item will be selected.
func (c *Combobox) SetSelected(index int) {
	C.uiComboboxSetSelected(c.c, C.int(index))
}

// OnSelected registers f to be run when the user selects an item in
// the Combobox. Only one function can be registered at a time.
func (c *Combobox) OnSelected(f func(*Combobox)) {
	c.onSelected = f
}

//export pkguiDoComboboxOnSelected
func pkguiDoComboboxOnSelected(cc *C.uiCombobox, data unsafe.Pointer) {
	c := ControlFromLibui(uintptr(unsafe.Pointer(cc))).(*Combobox)
	if c.onSelected != nil {
		c.onSelected(c)
	}
}

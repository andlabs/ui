// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Checkbox is a Control that represents a box with a text label at its
// side. When the user clicks the checkbox, a check mark will appear
// in the box; clicking it again removes the check.
type Checkbox struct {
	ControlBase
	c	*C.uiCheckbox
	onToggled		func(*Checkbox)
}

// NewCheckbox creates a new Checkbox with the given text as its
// label.
func NewCheckbox(text string) *Checkbox {
	c := new(Checkbox)

	ctext := C.CString(text)
	c.c = C.uiNewCheckbox(ctext)
	freestr(ctext)

	C.pkguiCheckboxOnToggled(c.c)

	c.ControlBase = NewControlBase(c, uintptr(unsafe.Pointer(c.c)))
	return c
}

// Text returns the Checkbox's text.
func (c *Checkbox) Text() string {
	ctext := C.uiCheckboxText(c.c)
	text := C.GoString(ctext)
	C.uiFreeText(ctext)
	return text
}

// SetText sets the Checkbox's text to text.
func (c *Checkbox) SetText(text string) {
	ctext := C.CString(text)
	C.uiCheckboxSetText(c.c, ctext)
	freestr(ctext)
}

// OnToggled registers f to be run when the user clicks the Checkbox.
// Only one function can be registered at a time.
func (c *Checkbox) OnToggled(f func(*Checkbox)) {
	c.onToggled = f
}

//export pkguiDoCheckboxOnToggled
func pkguiDoCheckboxOnToggled(cc *C.uiCheckbox, data unsafe.Pointer) {
	c := ControlFromLibui(uintptr(unsafe.Pointer(cc))).(*Checkbox)
	if c.onToggled != nil {
		c.onToggled(c)
	}
}

// Checked returns whether the Checkbox is checked.
func (c *Checkbox) Checked() bool {
	return tobool(C.uiCheckboxChecked(c.c))
}

// SetChecked sets whether the Checkbox is checked.
func (c *Checkbox) SetChecked(checked bool) {
	C.uiCheckboxSetChecked(c.c, frombool(checked))
}

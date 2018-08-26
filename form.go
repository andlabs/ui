// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Form is a Control that holds a group of Controls vertically
// with labels next to each. By default, each control has its
// preferred height; if a control is marked "stretchy", it will take
// whatever space is left over. If multiple controls are marked
// stretchy, they will be given equal shares of the leftover space.
// There can also be space between each control ("padding").
type Form struct {
	ControlBase
	f	*C.uiForm
	children	[]Control
}

// NewForm creates a new horizontal Form.
func NewForm() *Form {
	f := new(Form)

	f.f = C.uiNewForm()

	f.ControlBase = NewControlBase(f, uintptr(unsafe.Pointer(f.f)))
	return f
}

// Destroy destroys the Form. If the Form has children,
// Destroy calls Destroy on those Controls as well.
func (f *Form) Destroy() {
	for len(f.children) != 0 {
		c := f.children[0]
		f.Delete(0)
		c.Destroy()
	}
	f.ControlBase.Destroy()
}

// Append adds the given control to the end of the Form.
func (f *Form) Append(label string, child Control, stretchy bool) {
	clabel := C.CString(label)
	defer freestr(clabel)
	c := touiControl(child.LibuiControl())
	C.uiFormAppend(f.f, clabel, c, frombool(stretchy))
	f.children = append(f.children, child)
}

// Delete deletes the nth control of the Form.
func (f *Form) Delete(n int) {
	f.children = append(f.children[:n], f.children[n + 1:]...)
	C.uiFormDelete(f.f, C.int(n))
}

// Padded returns whether there is space between each control
// of the Form.
func (f *Form) Padded() bool {
	return tobool(C.uiFormPadded(f.f))
}

// SetPadded controls whether there is space between each control
// of the Form. The size of the padding is determined by the OS and
// its best practices.
func (f *Form) SetPadded(padded bool) {
	C.uiFormSetPadded(f.f, frombool(padded))
}

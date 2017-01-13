// 14 june 2016

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

type Form struct {
	c	*C.uiControl
	f	*C.uiForm

	children	[]Control
}

// NewForm creates a new Form.
func NewForm() *Form {
	f := new(Form)

	f.f = C.uiNewForm()
	f.c = (*C.uiControl)(unsafe.Pointer(f.f))

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
	C.uiControlDestroy(f.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Form. This is only used by package ui itself and should
// not be called by programs.
func (f *Form) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(f.c))
}

// Handle returns the OS-level handle associated with this Form.
func (f *Form) Handle() uintptr {
	return uintptr(C.uiControlHandle(f.c))
}

// Show shows the Form.
func (f *Form) Show() {
	C.uiControlShow(f.c)
}

// Hide hides the Form.
func (f *Form) Hide() {
	C.uiControlHide(f.c)
}

// Enable enables the Form.
func (f *Form) Enable() {
	C.uiControlEnable(f.c)
}

// Disable disables the Form.
func (f *Form) Disable() {
	C.uiControlDisable(f.c)
}

// Append adds the given control to the end of the Form.
func (f *Form) Append(label string, child Control, stretchy bool) {
	clabel := C.CString(label)

	c := (*C.uiControl)(nil)
	if child != nil {
		c = touiControl(child.LibuiControl())
	}

	C.uiFormAppend(f.f, clabel, c, frombool(stretchy))
	freestr(clabel)

	f.children = append(f.children, child)
}

// Delete deletes the nth control of the Form.
func (f *Form) Delete(n int) {
	f.children = append(f.children[:n], f.children[n + 1:]...)
	//C.uiFormDelete(f.f, C.int(n))
}

// TODO: InsertAt

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

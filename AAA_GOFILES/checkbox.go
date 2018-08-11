// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doCheckboxOnToggled(uiCheckbox *, void *);
// static inline void realuiCheckboxOnToggled(uiCheckbox *c)
// {
// 	uiCheckboxOnToggled(c, doCheckboxOnToggled, NULL);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var checkboxes = make(map[*C.uiCheckbox]*Checkbox)

// Checkbox is a Control that represents a box with a text label at its
// side. When the user clicks the checkbox, a check mark will appear
// in the box; clicking it again removes the check.
type Checkbox struct {
	co	*C.uiControl
	c	*C.uiCheckbox

	onToggled		func(*Checkbox)
}

// NewCheckbox creates a new Checkbox with the given text as its label.
func NewCheckbox(text string) *Checkbox {
	c := new(Checkbox)

	ctext := C.CString(text)
	c.c = C.uiNewCheckbox(ctext)
	c.co = (*C.uiControl)(unsafe.Pointer(c.c))
	freestr(ctext)

	C.realuiCheckboxOnToggled(c.c)
	checkboxes[c.c] = c

	return c
}

// Destroy destroys the Checkbox.
func (c *Checkbox) Destroy() {
	delete(checkboxes, c.c)
	C.uiControlDestroy(c.co)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (c *Checkbox) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(c.co))
}

// Handle returns the OS-level handle associated with this Checkbox.
// On Windows this is an HWND of a standard Windows API BUTTON
// class (as provided by Common Controls version 6).
// On GTK+ this is a pointer to a GtkCheckButton.
// On OS X this is a pointer to a NSButton.
func (c *Checkbox) Handle() uintptr {
	return uintptr(C.uiControlHandle(c.co))
}

// Show shows the Checkbox.
func (c *Checkbox) Show() {
	C.uiControlShow(c.co)
}

// Hide hides the Checkbox.
func (c *Checkbox) Hide() {
	C.uiControlHide(c.co)
}

// Enable enables the Checkbox.
func (c *Checkbox) Enable() {
	C.uiControlEnable(c.co)
}

// Disable disables the Checkbox.
func (c *Checkbox) Disable() {
	C.uiControlDisable(c.co)
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

//export doCheckboxOnToggled
func doCheckboxOnToggled(cc *C.uiCheckbox, data unsafe.Pointer) {
	c := checkboxes[cc]
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

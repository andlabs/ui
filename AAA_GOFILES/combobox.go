// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doComboboxOnSelected(uiCombobox *, void *);
// static inline void realuiComboboxOnSelected(uiCombobox *c)
// {
// 	uiComboboxOnSelected(c, doComboboxOnSelected, NULL);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var comboboxes = make(map[*C.uiCombobox]*Combobox)

// Combobox is a Control that represents a drop-down list of strings
// that the user can choose one of at any time. An editable
// Combobox also has an entry field that the user can type an alternate
// choice into.
type Combobox struct {
	co	*C.uiControl
	c	*C.uiCombobox

	onSelected		func(*Combobox)
}

// NewCombobox creates a new Combobox.
// This Combobox is not editable.
func NewCombobox() *Combobox {
	c := new(Combobox)

	c.c = C.uiNewCombobox()
	c.co = (*C.uiControl)(unsafe.Pointer(c.c))

	C.realuiComboboxOnSelected(c.c)
	comboboxes[c.c] = c

	return c
}
/*TODO
// NewEditableCombobox creates a new editable Combobox.
func NewEditableCombobox() *Combobox {
	c := new(Combobox)

	c.c = C.uiNewEditableCombobox()
	c.co = (*C.uiControl)(unsafe.Pointer(c.c))

	C.realuiComboboxOnSelected(c.c)
	comboboxes[c.c] = c

	return c
}
*/
// Destroy destroys the Combobox.
func (c *Combobox) Destroy() {
	delete(comboboxes, c.c)
	C.uiControlDestroy(c.co)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (c *Combobox) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(c.co))
}

// Handle returns the OS-level handle associated with this Combobox.
// On Windows this is an HWND of a standard Windows API COMBOBOX
// class (as provided by Common Controls version 6).
// On GTK+ this is a pointer to a GtkComboBoxText.
// On OS X this is a pointer to a NSComboBox for editable Comboboxes
// and to a NSPopUpButton for noneditable Comboboxes.
func (c *Combobox) Handle() uintptr {
	return uintptr(C.uiControlHandle(c.co))
}

// Show shows the Combobox.
func (c *Combobox) Show() {
	C.uiControlShow(c.co)
}

// Hide hides the Combobox.
func (c *Combobox) Hide() {
	C.uiControlHide(c.co)
}

// Enable enables the Combobox.
func (c *Combobox) Enable() {
	C.uiControlEnable(c.co)
}

// Disable disables the Combobox.
func (c *Combobox) Disable() {
	C.uiControlDisable(c.co)
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

// SetChecked sets the currently select item in the Combobox
// to index. If index is -1 no item will be selected.
func (c *Combobox) SetSelected(index int) {
	C.uiComboboxSetSelected(c.c, C.intmax_t(index))
}

// OnSelected registers f to be run when the user selects an item in
// the Combobox. Only one function can be registered at a time.
func (c *Combobox) OnSelected(f func(*Combobox)) {
	c.onSelected = f
}

//export doComboboxOnSelected
func doComboboxOnSelected(cc *C.uiCombobox, data unsafe.Pointer) {
	c := comboboxes[cc]
	if c.onSelected != nil {
		c.onSelected(c)
	}
}

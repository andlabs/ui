// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doEditableComboboxOnChanged(uiCombobox *, void *);
// static inline void realuiEditableComboboxOnChanged(uiCombobox *c)
// {
// 	uiEditableComboboxOnChanged(c, doEditableComboboxOnChanged, NULL);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var editableComboboxes = make(map[*C.uiEditableCombobox]*Combobox)

// EditableCombobox is a Control that represents a drop-down list
// of strings that the user can choose one of at any time. It also has
// an entry field that the user can type an alternate choice into.
type EditableCombobox struct {
	co	*C.uiControl
	c	*C.uiEditableCombobox

	onChanged		func(*EditableCombobox)
}

// NewEditableCombobox creates a new EditableCombobox.
func NewEditableCombobox() *EditableCombobox {
	c := new(EditableCombobox)

	c.c = C.uiNewEditableCombobox()
	c.co = (*C.uiControl)(unsafe.Pointer(c.c))

	C.realuiEditableComboboxOnChanged(c.c)
	editableComboboxes[c.c] = c

	return c
}

// Destroy destroys the EditableCombobox.
func (c *Combobox) Destroy() {
	delete(editableComboboxes, c.c)
	C.uiControlDestroy(c.co)
}

// LibuiControl returns the libui uiControl pointer that backs
// the EditableCombobox. This is only used by package ui itself and
// should not be called by programs.
func (c *Combobox) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(c.co))
}

// Handle returns the OS-level handle associated with this EditableCombobox.
// On Windows this is an HWND of a standard Windows API COMBOBOX
// class (as provided by Common Controls version 6).
// On GTK+ this is a pointer to a GtkComboBoxText.
// On OS X this is a pointer to a NSComboBox.
func (c *Combobox) Handle() uintptr {
	return uintptr(C.uiControlHandle(c.co))
}

// Show shows the EditableCombobox.
func (c *Combobox) Show() {
	C.uiControlShow(c.co)
}

// Hide hides the EditableCombobox.
func (c *Combobox) Hide() {
	C.uiControlHide(c.co)
}

// Enable enables the EditableCombobox.
func (c *Combobox) Enable() {
	C.uiControlEnable(c.co)
}

// Disable disables the EditableCombobox.
func (c *Combobox) Disable() {
	C.uiControlDisable(c.co)
}

// Append adds the named item to the end of the EditableCombobox.
func (c *Combobox) Append(text string) {
	ctext := C.CString(text)
	C.uiComboboxAppend(c.c, ctext)
	freestr(ctext)
}

// Text returns the text in the entry of the EditableCombobox, which
// could be one of the choices in the list if the user has selected one.
func (c *Combobox) Text() string {
	ctext := C.uiEditableComboboxText(c.c)
	text := C.GoString(ctext)
	C.uiFreeText(ctext)
	return text
}

// SetText sets the text in the entry of the EditableCombobox.
func (c *Combobox) SetText(index int) {
	ctext := C.CString(text)
	C.uiEditableComboboxSetText(c.c, ctext)
	freestr(ctext)
}

// OnChanged registers f to be run when the user selects an item in
// the Combobox. Only one function can be registered at a time.
func (c *Combobox) OnChanged(f func(*Combobox)) {
	c.onChanged = f
}

//export doComboboxOnChanged
func doComboboxOnChanged(cc *C.uiCombobox, data unsafe.Pointer) {
	c := editableComboboxes[cc]
	if c.onChanged != nil {
		c.onChanged(c)
	}
}

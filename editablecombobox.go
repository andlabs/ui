// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// EditableCombobox is a Control that represents a drop-down list
// of strings that the user can choose one of at any time. It also has
// an entry field that the user can type an alternate choice into.
type EditableCombobox struct {
	ControlBase
	c	*C.uiEditableCombobox
	onChanged		func(*EditableCombobox)
}

// NewEditableCombobox creates a new EditableCombobox.
func NewEditableCombobox() *EditableCombobox {
	c := new(EditableCombobox)

	c.c = C.uiNewEditableCombobox()

	C.pkguiEditableComboboxOnChanged(c.c)

	c.ControlBase = NewControlBase(c, uintptr(unsafe.Pointer(c.c)))
	return c
}

// Append adds the named item to the end of the EditableCombobox.
func (e *EditableCombobox) Append(text string) {
	ctext := C.CString(text)
	C.uiEditableComboboxAppend(e.c, ctext)
	freestr(ctext)
}

// Text returns the text in the entry of the EditableCombobox, which
// could be one of the choices in the list if the user has selected one.
func (e *EditableCombobox) Text() string {
	ctext := C.uiEditableComboboxText(e.c)
	text := C.GoString(ctext)
	C.uiFreeText(ctext)
	return text
}

// SetText sets the text in the entry of the EditableCombobox.
func (e *EditableCombobox) SetText(text string) {
	ctext := C.CString(text)
	C.uiEditableComboboxSetText(e.c, ctext)
	freestr(ctext)
}

// OnChanged registers f to be run when the user either selects an
// item or changes the text in the EditableCombobox. Only one
// function can be registered at a time.
func (e *EditableCombobox) OnChanged(f func(*EditableCombobox)) {
	e.onChanged = f
}

//export pkguiDoEditableComboboxOnChanged
func pkguiDoEditableComboboxOnChanged(cc *C.uiEditableCombobox, data unsafe.Pointer) {
	e := ControlFromLibui(uintptr(unsafe.Pointer(cc))).(*EditableCombobox)
	if e.onChanged != nil {
		e.onChanged(e)
	}
}

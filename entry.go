// 12 december 2015

// TODO typing in entry in OS X crashes libui
// I've had similar issues with checkboxes on libui
// something's wrong with NSMapTable

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Entry is a Control that represents a space that the user can
// type a single line of text into.
type Entry struct {
	ControlBase
	e	*C.uiEntry
	onChanged		func(*Entry)
}

func finishNewEntry(ee *C.uiEntry) *Entry {
	e := new(Entry)

	e.e = ee

	C.pkguiEntryOnChanged(e.e)

	e.ControlBase = NewControlBase(e, uintptr(unsafe.Pointer(e.e)))
	return e
}

// NewEntry creates a new Entry.
func NewEntry() *Entry {
	return finishNewEntry(C.uiNewEntry())
}

// NewPasswordEntry creates a new Entry whose contents are
// visibly obfuscated, suitable for passwords.
func NewPasswordEntry() *Entry {
	return finishNewEntry(C.uiNewPasswordEntry())
}

// NewSearchEntry creates a new Entry suitable for searching with.
// Changed events may, depending on the system, be delayed
// with a search Entry, to produce a smoother user experience.
func NewSearchEntry() *Entry {
	return finishNewEntry(C.uiNewSearchEntry())
}

// Text returns the Entry's text.
func (e *Entry) Text() string {
	ctext := C.uiEntryText(e.e)
	text := C.GoString(ctext)
	C.uiFreeText(ctext)
	return text
}

// SetText sets the Entry's text to text.
func (e *Entry) SetText(text string) {
	ctext := C.CString(text)
	C.uiEntrySetText(e.e, ctext)
	freestr(ctext)
}

// OnChanged registers f to be run when the user makes a change to
// the Entry. Only one function can be registered at a time.
func (e *Entry) OnChanged(f func(*Entry)) {
	e.onChanged = f
}

//export pkguiDoEntryOnChanged
func pkguiDoEntryOnChanged(ee *C.uiEntry, data unsafe.Pointer) {
	e := ControlFromLibui(uintptr(unsafe.Pointer(ee))).(*Entry)
	if e.onChanged != nil {
		e.onChanged(e)
	}
}

// ReadOnly returns whether the Entry can be changed.
func (e *Entry) ReadOnly() bool {
	return tobool(C.uiEntryReadOnly(e.e))
}

// SetReadOnly sets whether the Entry can be changed.
func (e *Entry) SetReadOnly(ro bool) {
	C.uiEntrySetReadOnly(e.e, frombool(ro))
}

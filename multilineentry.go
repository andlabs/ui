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

// MultilineEntry is a Control that represents a space that the user
// can type multiple lines of text into.
type MultilineEntry struct {
	ControlBase
	e	*C.uiMultilineEntry
	onChanged		func(*MultilineEntry)
}

func finishNewMultilineEntry(ee *C.uiMultilineEntry) *MultilineEntry {
	m := new(MultilineEntry)

	m.e = ee

	C.pkguiMultilineEntryOnChanged(m.e)

	m.ControlBase = NewControlBase(m, uintptr(unsafe.Pointer(m.e)))
	return m
}

// NewMultilineEntry creates a new MultilineEntry.
// The MultilineEntry soft-word-wraps and has no horizontal
// scrollbar.
func NewMultilineEntry() *MultilineEntry {
	return finishNewMultilineEntry(C.uiNewMultilineEntry())
}

// NewNonWrappingMultilineEntry creates a new MultilineEntry.
// The MultilineEntry does not word-wrap and thus has horizontal
// scrollbar.
func NewNonWrappingMultilineEntry() *MultilineEntry {
	return finishNewMultilineEntry(C.uiNewNonWrappingMultilineEntry())
}

// Text returns the MultilineEntry's text.
func (m *MultilineEntry) Text() string {
	ctext := C.uiMultilineEntryText(m.e)
	text := C.GoString(ctext)
	C.uiFreeText(ctext)
	return text
}

// SetText sets the MultilineEntry's text to text.
func (m *MultilineEntry) SetText(text string) {
	ctext := C.CString(text)
	C.uiMultilineEntrySetText(m.e, ctext)
	freestr(ctext)
}

// Append adds text to the end of the MultilineEntry's text.
// TODO selection and scroll behavior
func (m *MultilineEntry) Append(text string) {
	ctext := C.CString(text)
	C.uiMultilineEntryAppend(m.e, ctext)
	freestr(ctext)
}

// OnChanged registers f to be run when the user makes a change to
// the MultilineEntry. Only one function can be registered at a time.
func (m *MultilineEntry) OnChanged(f func(*MultilineEntry)) {
	m.onChanged = f
}

//export pkguiDoMultilineEntryOnChanged
func pkguiDoMultilineEntryOnChanged(ee *C.uiMultilineEntry, data unsafe.Pointer) {
	m := ControlFromLibui(uintptr(unsafe.Pointer(ee))).(*MultilineEntry)
	if m.onChanged != nil {
		m.onChanged(m)
	}
}

// ReadOnly returns whether the MultilineEntry can be changed.
func (m *MultilineEntry) ReadOnly() bool {
	return tobool(C.uiMultilineEntryReadOnly(m.e))
}

// SetReadOnly sets whether the MultilineEntry can be changed.
func (m *MultilineEntry) SetReadOnly(ro bool) {
	C.uiMultilineEntrySetReadOnly(m.e, frombool(ro))
}

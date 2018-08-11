// 12 december 2015

// TODO typing in entry in OS X crashes libui
// I've had similar issues with checkboxes on libui
// something's wrong with NSMapTable

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doEntryOnChanged(uiEntry *, void *);
// static inline void realuiEntryOnChanged(uiEntry *b)
// {
// 	uiEntryOnChanged(b, doEntryOnChanged, NULL);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var entries = make(map[*C.uiEntry]*Entry)

// Entry is a Control that represents a space that the user can
// type a single line of text into.
type Entry struct {
	c	*C.uiControl
	e	*C.uiEntry

	onChanged		func(*Entry)
}

// NewEntry creates a new Entry.
func NewEntry() *Entry {
	e := new(Entry)

	e.e = C.uiNewEntry()
	e.c = (*C.uiControl)(unsafe.Pointer(e.e))

	C.realuiEntryOnChanged(e.e)
	entries[e.e] = e

	return e
}

// Destroy destroys the Entry.
func (e *Entry) Destroy() {
	delete(entries, e.e)
	C.uiControlDestroy(e.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (e *Entry) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(e.c))
}

// Handle returns the OS-level handle associated with this Entry.
// On Windows this is an HWND of a standard Windows API EDIT
// class (as provided by Common Controls version 6).
// On GTK+ this is a pointer to a GtkEntry.
// On OS X this is a pointer to a NSTextField.
func (e *Entry) Handle() uintptr {
	return uintptr(C.uiControlHandle(e.c))
}

// Show shows the Entry.
func (e *Entry) Show() {
	C.uiControlShow(e.c)
}

// Hide hides the Entry.
func (e *Entry) Hide() {
	C.uiControlHide(e.c)
}

// Enable enables the Entry.
func (e *Entry) Enable() {
	C.uiControlEnable(e.c)
}

// Disable disables the Entry.
func (e *Entry) Disable() {
	C.uiControlDisable(e.c)
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

//export doEntryOnChanged
func doEntryOnChanged(ee *C.uiEntry, data unsafe.Pointer) {
	e := entries[ee]
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

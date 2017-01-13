// multiline entry

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doMultilineEntryOnChanged(uiMultilineEntry*, void *);
// static inline void realuiMultilineEntryOnChanged(uiMultilineEntry *e)
// {
// 	 uiMultilineEntryOnChanged(e, doMultilineEntryOnChanged, NULL);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var mEntries = make(map[*C.uiMultilineEntry]*MultilineEntry)

// MultilineEntry is a Control that represents a space that the user can
// type multiple lines of text into.
type MultilineEntry struct {
	c *C.uiControl
	e *C.uiMultilineEntry

	onChanged func(*MultilineEntry)
}

// NewMultilineEntry creates a new MultilineEntry.
func NewMultilineEntry() *MultilineEntry {
	e := new(MultilineEntry)

	e.e = C.uiNewMultilineEntry()
	e.c = (*C.uiControl)(unsafe.Pointer(e.e))

	C.realuiMultilineEntryOnChanged(e.e)
	mEntries[e.e] = e

	return e
}

// NewMultilineNonWrappingEntry creates a new MultilineEntry.
func NewMultilineNonWrappingEntry() *MultilineEntry {
	e := new(MultilineEntry)

	e.e = C.uiNewNonWrappingMultilineEntry()
	e.c = (*C.uiControl)(unsafe.Pointer(e.e))

	C.realuiMultilineEntryOnChanged(e.e)
	mEntries[e.e] = e

	return e
}

// Destroy destroys the MultilineEntry.
func (e *MultilineEntry) Destroy() {
	delete(mEntries, e.e)
	C.uiControlDestroy(e.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (e *MultilineEntry) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(e.c))
}

// Handle returns the OS-level handle associated with this MultilineEntry.
func (e *MultilineEntry) Handle() uintptr {
	return uintptr(C.uiControlHandle(e.c))
}

// Show shows the MultilineEntry.
func (e *MultilineEntry) Show() {
	C.uiControlShow(e.c)
}

// Hide hides the MultilineEntry.
func (e *MultilineEntry) Hide() {
	C.uiControlHide(e.c)
}

// Enable enables the MultilineEntry.
func (e *MultilineEntry) Enable() {
	C.uiControlEnable(e.c)
}

// Disable disables the MultilineEntry.
func (e *MultilineEntry) Disable() {
	C.uiControlDisable(e.c)
}

// Text returns the MultilineEntry's text.
func (e *MultilineEntry) Text() string {
	ctext := C.uiMultilineEntryText(e.e)
	text := C.GoString(ctext)
	C.uiFreeText(ctext)
	return text
}

// SetText sets the MultilineEntry's text to text.
func (e *MultilineEntry) SetText(text string) {
	ctext := C.CString(text)
	C.uiMultilineEntrySetText(e.e, ctext)
	freestr(ctext)
}

// Append text to the MultilineEntry's text.
func (e *MultilineEntry) Append(text string) {
	ctext := C.CString(text)
	C.uiMultilineEntryAppend(e.e, ctext)
	freestr(ctext)
}

// OnChanged registers f to be run when the user makes a change to
// the MultilineEntry. Only one function can be registered at a time.
func (e *MultilineEntry) OnChanged(f func(*MultilineEntry)) {
	e.onChanged = f
}

//export doMultilineEntryOnChanged
func doMultilineEntryOnChanged(ee *C.uiMultilineEntry, data unsafe.Pointer) {
	e := mEntries[ee]
	if e.onChanged != nil {
		e.onChanged(e)
	}
}

// ReadOnly returns whether the MultilineEntry can be changed.
func (e *MultilineEntry) ReadOnly() bool {
	return tobool(C.uiMultilineEntryReadOnly(e.e))
}

// SetReadOnly sets whether the MultilineEntry can be changed.
func (e *MultilineEntry) SetReadOnly(ro bool) {
	C.uiMultilineEntrySetReadOnly(e.e, frombool(ro))
}

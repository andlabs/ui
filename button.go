// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doButtonOnClicked(uiButton *, void *);
// static inline void realuiButtonOnClicked(uiButton *b)
// {
// 	uiButtonOnClicked(b, doButtonOnClicked, NULL);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var buttons = make(map[*C.uiButton]*Button)

// Button is a Control that represents a button that the user can
// click to perform an action. A Button has a text label that should
// describe what the button does.
type Button struct {
	c	*C.uiControl
	b	*C.uiButton

	onClicked		func(*Button)
}

// NewButton creates a new Button with the given text as its label.
func NewButton(text string) *Button {
	b := new(Button)

	ctext := C.CString(text)
	b.b = C.uiNewButton(ctext)
	b.c = (*C.uiControl)(unsafe.Pointer(b.b))
	freestr(ctext)

	C.realuiButtonOnClicked(b.b)
	buttons[b.b] = b

	return b
}

// Destroy destroys the Button.
func (b *Button) Destroy() {
	delete(buttons, b.b)
	C.uiControlDestroy(b.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Button. This is only used by package ui itself and should
// not be called by programs.
func (b *Button) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(b.c))
}

// Handle returns the OS-level handle associated with this Button.
// On Windows this is an HWND of a standard Windows API BUTTON
// class (as provided by Common Controls version 6).
// On GTK+ this is a pointer to a GtkButton.
// On OS X this is a pointer to a NSButton.
func (b *Button) Handle() uintptr {
	return uintptr(C.uiControlHandle(b.c))
}

// Show shows the Button.
func (b *Button) Show() {
	C.uiControlShow(b.c)
}

// Hide hides the Button.
func (b *Button) Hide() {
	C.uiControlHide(b.c)
}

// Enable enables the Button.
func (b *Button) Enable() {
	C.uiControlEnable(b.c)
}

// Disable disables the Button.
func (b *Button) Disable() {
	C.uiControlDisable(b.c)
}

// Text returns the Button's text.
func (b *Button) Text() string {
	ctext := C.uiButtonText(b.b)
	text := C.GoString(ctext)
	C.uiFreeText(ctext)
	return text
}

// SetText sets the Button's text to text.
func (b *Button) SetText(text string) {
	ctext := C.CString(text)
	C.uiButtonSetText(b.b, ctext)
	freestr(ctext)
}

// OnClicked registers f to be run when the user clicks the Button.
// Only one function can be registered at a time.
func (b *Button) OnClicked(f func(*Button)) {
	b.onClicked = f
}

//export doButtonOnClicked
func doButtonOnClicked(bb *C.uiButton, data unsafe.Pointer) {
	b := buttons[bb]
	if b.onClicked != nil {
		b.onClicked(b)
	}
}

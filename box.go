// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

// Box is a Control that holds a group of Controls horizontally
// or vertically. If horizontally, then all controls have the same
// height. If vertically, then all controls have the same width.
// By default, each control has its preferred width (horizontal)
// or height (vertical); if a control is marked "stretchy", it will
// take whatever space is left over. If multiple controls are marked
// stretchy, they will be given equal shares of the leftover space.
// There can also be space between each control ("padding").
type Box struct {
	c	*C.uiControl
	b	*C.uiBox

	children	[]Control
}

// NewHorizontalBox creates a new horizontal Box.
func NewHorizontalBox() *Box {
	b := new(Box)

	b.b = C.uiNewHorizontalBox()
	b.c = (*C.uiControl)(unsafe.Pointer(b.b))

	return b
}

// NewVerticalBox creates a new vertical Box.
func NewVerticalBox() *Box {
	b := new(Box)

	b.b = C.uiNewVerticalBox()
	b.c = (*C.uiControl)(unsafe.Pointer(b.b))

	return b
}

// Destroy destroys the Box. If the Box has children,
// Destroy calls Destroy on those Controls as well.
func (b *Box) Destroy() {
	for len(b.children) != 0 {
		c := b.children[0]
		b.Delete(0)
		c.Destroy()
	}
	C.uiControlDestroy(b.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Box. This is only used by package ui itself and should
// not be called by programs.
func (b *Box) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(b.c))
}

// Handle returns the OS-level handle associated with this Box.
// On Windows this is an HWND of a libui-internal class.
// On GTK+ this is a pointer to a GtkBox.
// On OS X this is a pointer to a NSView.
func (b *Box) Handle() uintptr {
	return uintptr(C.uiControlHandle(b.c))
}

// Show shows the Box.
func (b *Box) Show() {
	C.uiControlShow(b.c)
}

// Hide hides the Box.
func (b *Box) Hide() {
	C.uiControlHide(b.c)
}

// Enable enables the Box.
func (b *Box) Enable() {
	C.uiControlEnable(b.c)
}

// Disable disables the Box.
func (b *Box) Disable() {
	C.uiControlDisable(b.c)
}

// Append adds the given control to the end of the Box.
func (b *Box) Append(child Control, stretchy bool) {
	c := (*C.uiControl)(nil)
	if child != nil {
		c = touiControl(child.LibuiControl())
	}
	C.uiBoxAppend(b.b, c, frombool(stretchy))
	b.children = append(b.children, child)
}

// Delete deletes the nth control of the Box.
func (b *Box) Delete(n int) {
	b.children = append(b.children[:n], b.children[n + 1:]...)
	// TODO why is this uintmax_t instead of intmax_t
	C.uiBoxDelete(b.b, C.uintmax_t(n))
}

// Padded returns whether there is space between each control
// of the Box.
func (b *Box) Padded() bool {
	return tobool(C.uiBoxPadded(b.b))
}

// SetPadded controls whether there is space between each control
// of the Box. The size of the padding is determined by the OS and
// its best practices.
func (b *Box) SetPadded(padded bool) {
	C.uiBoxSetPadded(b.b, frombool(padded))
}

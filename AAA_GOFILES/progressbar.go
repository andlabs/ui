// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

// ProgressBar is a Control that represents a horizontal bar that
// is filled in progressively over time as a process completes.
type ProgressBar struct {
	c	*C.uiControl
	p	*C.uiProgressBar
}

// NewProgressBar creates a new ProgressBar.
func NewProgressBar() *ProgressBar {
	p := new(ProgressBar)

	p.p = C.uiNewProgressBar()
	p.c = (*C.uiControl)(unsafe.Pointer(p.p))

	return p
}

// Destroy destroys the ProgressBar.
func (p *ProgressBar) Destroy() {
	C.uiControlDestroy(p.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (p *ProgressBar) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(p.c))
}

// Handle returns the OS-level handle associated with this ProgressBar.
// On Windows this is an HWND of a standard Windows API
// PROGRESS_CLASS class (as provided by Common Controls
// version 6).
// On GTK+ this is a pointer to a GtkProgressBar.
// On OS X this is a pointer to a NSProgressIndicator.
func (p *ProgressBar) Handle() uintptr {
	return uintptr(C.uiControlHandle(p.c))
}

// Show shows the ProgressBar.
func (p *ProgressBar) Show() {
	C.uiControlShow(p.c)
}

// Hide hides the ProgressBar.
func (p *ProgressBar) Hide() {
	C.uiControlHide(p.c)
}

// Enable enables the ProgressBar.
func (p *ProgressBar) Enable() {
	C.uiControlEnable(p.c)
}

// Disable disables the ProgressBar.
func (p *ProgressBar) Disable() {
	C.uiControlDisable(p.c)
}

// TODO Value

// SetValue sets the ProgressBar's currently displayed percentage
// to value. value must be between 0 and 100 inclusive.
func (p *ProgressBar) SetValue(value int) {
	C.uiProgressBarSetValue(p.p, C.int(value))
}

// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// ProgressBar is a Control that represents a horizontal bar that
// is filled in progressively over time as a process completes.
type ProgressBar struct {
	ControlBase
	p	*C.uiProgressBar
}

// NewProgressBar creates a new ProgressBar.
func NewProgressBar() *ProgressBar {
	p := new(ProgressBar)

	p.p = C.uiNewProgressBar()

	p.ControlBase = NewControlBase(p, uintptr(unsafe.Pointer(p.p)))
	return p
}

// Value returns the value currently shown in the ProgressBar.
func (p *ProgressBar) Value() int {
	return int(C.uiProgressBarValue(p.p))
}

// SetValue sets the ProgressBar's currently displayed percentage
// to value. value must be between 0 and 100 inclusive, or -1 for
// an indeterminate progressbar.
func (p *ProgressBar) SetValue(value int) {
	C.uiProgressBarSetValue(p.p, C.int(value))
}

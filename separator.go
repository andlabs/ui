// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Separator is a Control that represents a horizontal line that
// visually separates controls.
type Separator struct {
	ControlBase
	s	*C.uiSeparator
}

// NewHorizontalSeparator creates a new horizontal Separator.
func NewHorizontalSeparator() *Separator {
	s := new(Separator)

	s.s = C.uiNewHorizontalSeparator()

	s.ControlBase = NewControlBase(s, uintptr(unsafe.Pointer(s.s)))
	return s
}

// NewVerticalSeparator creates a new vertical Separator.
func NewVerticalSeparator() *Separator {
	s := new(Separator)

	s.s = C.uiNewVerticalSeparator()

	s.ControlBase = NewControlBase(s, uintptr(unsafe.Pointer(s.s)))
	return s
}

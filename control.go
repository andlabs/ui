// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

// Control represents a GUI control. It provdes methods
// common to all Controls.
// 
// To create a new Control, implement the control on
// the libui side, then provide access to that control on
// the Go side via an implementation of Control as
// described.
type Control interface {
	// Destroy destroys the Control.
	// 
	// Implementations should do any necessary cleanup,
	// then call LibuiControlDestroy.
	Destroy()

	// LibuiControl returns the libui uiControl pointer that backs
	// the Control. This is only used by package ui itself and should
	// not be called by programs.
	LibuiControl() uintptr

	// Handle returns the OS-level handle that backs the
	// Control. On OSs that use reference counting for
	// controls, Handle does not increment the reference
	// count; you are sharing package ui's reference.
	// 
	// Implementations should call LibuiControlHandle and
	// document exactly what kind of handle is returned.
	Handle() uintptr

	// Show shows the Control.
	// 
	// Implementations should call LibuiControlShow.
	Show()

	// Hide shows the Control. Hidden controls do not participate
	// in layout (that is, Box, Grid, etc. does not reserve space for
	// hidden controls).
	// 
	// Implementations should call LibuiControlHide.
	Hide()

	// Enable enables the Control.
	// 
	// Implementations should call LibuiControlEnable.
	Enable()

	// Disable disables the Control.
	// 
	// Implementations should call LibuiControlDisable.
	Disable()
}

func touiControl(c uintptr) *C.uiControl {
	return (*C.uiControl)(unsafe.Pointer(c))
}

// LibuiControlDestroy allows implementations of Control
// to call the libui function uiControlDestroy.
func LibuiControlDestroy(c uintptr) {
	C.uiControlDestroy(touiControl(c))
}

// LibuiControlHandle allows implementations of Control
// to call the libui function uiControlHandle.
func LibuiControlHandle(c uintptr) uintptr {
	return uintptr(C.uiControlHandle(touiControl(c)))
}

// LibuiControlShow allows implementations of Control
// to call the libui function uiControlShow.
func LibuiControlShow(c uintptr) {
	C.uiControlShow(touiControl(c))
}

// LibuiControlHide allows implementations of Control
// to call the libui function uiControlHide.
func LibuiControlHide(c uintptr) {
	C.uiControlHide(touiControl(c))
}

// LibuiControlEnable allows implementations of Control
// to call the libui function uiControlEnable.
func LibuiControlEnable(c uintptr) {
	C.uiControlEnable(touiControl(c))
}

// LibuiControlDisable allows implementations of Control
// to call the libui function uiControlDisable.
func LibuiControlDisable(c uintptr) {
	C.uiControlDisable(touiControl(c))
}

// LibuiFreeText allows implementations of Control
// to call the libui function uiFreeText.
func LibuiFreeText(c uintptr) {
	C.uiFreeText((*C.char)(unsafe.Pointer(c)))
}

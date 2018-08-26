// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// no need to lock this; only the GUI thread can access it
var controls = make(map[*C.uiControl]Control)

// Control represents a GUI control. It provdes methods
// common to all Controls.
// 
// The preferred way to create new Controls is to use
// ControlBase; see ControlBase below.
type Control interface {
	// LibuiControl returns the uiControl pointer for the Control.
	// This is intended for use when adding a control to a
	// container.
	LibuiControl() uintptr

	// Destroy destroys the Control.
	Destroy()

	// Handle returns the OS-level handle that backs the
	// Control. On OSs that use reference counting for
	// controls, Handle does not increment the reference
	// count; you are sharing package ui's reference.
	Handle() uintptr

	// Visible returns whether the Control is visible.
	Visible() bool

	// Show shows the Control.
	Show()

	// Hide shows the Control. Hidden controls do not participate
	// in layout (that is, Box, Grid, etc. does not reserve space for
	// hidden controls).
	Hide()

	// Enabled returns whether the Control is enabled.
	Enabled() bool

	// Enable enables the Control.
	Enable()

	// Disable disables the Control.
	Disable()
}

// ControlBase is an implementation of Control that provides
// all the methods that Control requires. To use it, embed a
// ControlBase (not a *ControlBase) into your structure, then
// assign the result of NewControlBase to that field:
// 
// 	type MyControl struct {
// 		ui.ControlBase
// 		c *C.MyControl
// 	}
// 	
// 	func NewMyControl() *MyControl {
// 		m := &NewMyControl{
// 			c: C.newMyControl(),
// 		}
// 		m.ControlBase = ui.NewControlBase(m, uintptr(unsafe.Pointer(c)))
// 		return m
// 	}
type ControlBase struct {
	iface		Control
	c		*C.uiControl
}

// NewControlBase creates a new ControlBase. See the
// documentation of ControlBase for an example.
// NewControl should only be called once per instance of Control.
func NewControlBase(iface Control, c uintptr) ControlBase {
	b := ControlBase{
		iface:	iface,
		c:		(*C.uiControl)(unsafe.Pointer(c)),
	}
	controls[b.c] = b.iface
	return b
}

func (c *ControlBase) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(c.c))
}

func (c *ControlBase) Destroy() {
	delete(controls, c.c)
	C.uiControlDestroy(c.c)
}

func (c *ControlBase) Handle() uintptr {
	return uintptr(C.uiControlHandle(c.c))
}

func (c *ControlBase) Visible() bool {
	return tobool(C.uiControlVisible(c.c))
}

func (c *ControlBase) Show() {
	C.uiControlShow(c.c)
}

func (c *ControlBase) Hide() {
	C.uiControlHide(c.c)
}

func (c *ControlBase) Enabled() bool {
	return tobool(C.uiControlEnabled(c.c))
}

func (c *ControlBase) Enable() {
	C.uiControlEnable(c.c)
}

func (c *ControlBase) Disable() {
	C.uiControlDisable(c.c)
}

// ControlFromLibui returns the Control associated with a libui
// uiControl. This is intended for implementing event handlers
// on the Go side, to prevent sharing Go pointers with C.
// This function only works on Controls that use ControlBase.
func ControlFromLibui(c uintptr) Control {
	// comma-ok form to avoid creating nil entries
	cc, _ := controls[(*C.uiControl)(unsafe.Pointer(c))]
	return cc
}

func touiControl(c uintptr) *C.uiControl {
	return (*C.uiControl)(unsafe.Pointer(c))
}

// LibuiFreeText allows implementations of Control
// to call the libui function uiFreeText.
func LibuiFreeText(c uintptr) {
	C.uiFreeText((*C.char)(unsafe.Pointer(c)))
}

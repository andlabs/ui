// 16 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// Area is a Control that represents a blank canvas that a program
// can draw on as it wishes. Areas also receive keyboard and mouse
// events, and programs can react to those as they see fit. Drawing
// and event handling are handled through an instance of a type
// that implements AreaHandler that every Area has; see AreaHandler
// for details.
// 
// There are two types of areas. Non-scrolling areas are rectangular
// and have no scrollbars. Programs can draw on and get mouse
// events from any point in the Area, and the size of the Area is
// decided by package ui itself, according to the layout of controls
// in the Window the Area is located in and the size of said Window.
// There is no way to query the Area's size or be notified when its
// size changes; instead, you are given the area size as part of the
// draw and mouse event handlers, for use solely within those
// handlers.
// 
// Scrolling areas have horziontal and vertical scrollbars. The amount
// that can be scrolled is determined by the area's size, which is
// decided by the programmer (both when creating the Area and by
// a call to SetSize). Only a portion of the Area is visible at any time;
// drawing and mouse events are automatically adjusted to match
// what portion is visible, so you do not have to worry about scrolling
// in your event handlers. AreaHandler has more information.
// 
// The internal coordinate system of an Area is points, which are
// floating-point and device-independent. For more details, see
// AreaHandler. The size of a scrolling Area must be an exact integer
// number of points (that is, you cannot have an Area that is 32.5
// points tall) and thus the parameters to NewScrollingArea and
// SetSize are ints. All other instances of points in parameters and
// structures (including sizes of drawn objects) are float64s.
type Area struct {
	ControlBase
	a	*C.uiArea
	ah	*C.uiAreaHandler
	scrolling	bool
}

// NewArea creates a new non-scrolling Area.
func NewArea(handler AreaHandler) *Area {
	a := new(Area)
	a.scrolling = false
	a.ah = registerAreaHandler(handler)

	a.a = C.uiNewArea(a.ah)

	a.ControlBase = NewControlBase(a, uintptr(unsafe.Pointer(a.a)))
	return a
}

// NewScrollingArea creates a new scrolling Area of the given size,
// in points.
func NewScrollingArea(handler AreaHandler, width int, height int) *Area {
	a := new(Area)
	a.scrolling = true
	a.ah = registerAreaHandler(handler)

	a.a = C.uiNewScrollingArea(a.ah, C.int(width), C.int(height))

	a.ControlBase = NewControlBase(a, uintptr(unsafe.Pointer(a.a)))
	return a
}

// Destroy destroys the Area.
func (a *Area) Destroy() {
	unregisterAreaHandler(a.ah)
	a.ControlBase.Destroy()
}

// SetSize sets the size of a scrolling Area to the given size, in points.
// SetSize panics if called on a non-scrolling Area.
func (a *Area) SetSize(width int, height int) {
	if !a.scrolling {
		panic("attempt to call SetSize on non-scrolling Area")
	}
	C.uiAreaSetSize(a.a, C.int(width), C.int(height))
}

// QueueRedrawAll queues the entire Area for redraw.
// The Area is not redrawn before this function returns; it is
// redrawn when next possible.
func (a *Area) QueueRedrawAll() {
	C.uiAreaQueueRedrawAll(a.a)
}

// ScrollTo scrolls the Area to show the given rectangle; what this
// means is implementation-defined, but you can safely assume
// that as much of the given rectangle as possible will be visible
// after this call. (TODO verify this on OS X) ScrollTo panics if called
// on a non-scrolling Area.
func (a *Area) ScrollTo(x float64, y float64, width float64, height float64) {
	if !a.scrolling {
		panic("attempt to call ScrollTo on non-scrolling Area")
	}
	C.uiAreaScrollTo(a.a, C.double(x), C.double(y), C.double(width), C.double(height))
}

// TODO BeginUserWindowMove
// TODO BeginUserWindowResize

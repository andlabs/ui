// 16 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

// no need to lock this; only the GUI thread can access it
var areas = make(map[*C.uiArea]*Area)

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
	c	*C.uiControl
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
	a.c = (*C.uiControl)(unsafe.Pointer(a.a))

	areas[a.a] = a

	return a
}

// NewScrollingArea creates a new scrolling Area of the given size,
// in points.
func NewScrollingArea(handler AreaHandler, width int, height int) *Area {
	a := new(Area)
	a.scrolling = true
	a.ah = registerAreaHandler(handler)

	a.a = C.uiNewScrollingArea(a.ah, C.intmax_t(width), C.intmax_t(height))
	a.c = (*C.uiControl)(unsafe.Pointer(a.a))

	areas[a.a] = a

	return a
}

// Destroy destroys the Area.
func (a *Area) Destroy() {
	delete(areas, a.a)
	C.uiControlDestroy(a.c)
	unregisterAreaHandler(a.ah)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Area. This is only used by package ui itself and should
// not be called by programs.
func (a *Area) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(a.c))
}

// Handle returns the OS-level handle associated with this Area.
// On Windows this is an HWND of a libui-internal class.
// On GTK+ this is a pointer to a GtkScrolledWindow with a
// GtkViewport as its child. The child of the viewport is the
// GtkDrawingArea that provides the Area itself.
// On OS X this is a pointer to a NSScrollView whose document view
// is the NSView that provides the Area itself.
func (a *Area) Handle() uintptr {
	return uintptr(C.uiControlHandle(a.c))
}

// Show shows the Area.
func (a *Area) Show() {
	C.uiControlShow(a.c)
}

// Hide hides the Area.
func (a *Area) Hide() {
	C.uiControlHide(a.c)
}

// Enable enables the Area.
func (a *Area) Enable() {
	C.uiControlEnable(a.c)
}

// Disable disables the Area.
func (a *Area) Disable() {
	C.uiControlDisable(a.c)
}

// SetSize sets the size of a scrolling Area to the given size, in points.
// SetSize panics if called on a non-scrolling Area.
func (a *Area) SetSize(width int, height int) {
	if !a.scrolling {
		panic("attempt to call SetSize on non-scrolling Area")
	}
	C.uiAreaSetSize(a.a, C.intmax_t(width), C.intmax_t(height))
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

// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

/*
On Windows, container controls are just regular controls that notify their parent when the user wants to do things; changing the contents of a switching container (such as a tab control) must be done manually.

We'll create a dummy window using the container window class for each tab page. This makes showing and hiding tabs a matter of showing and hiding one control.
*/

type tab struct {
	*controlSingleHWND
	tabs		[]*container
	children	[]Control
}

func newTab() Tab {
	hwnd := C.newControl(C.xWC_TABCONTROL,
		C.TCS_TOOLTIPS|C.WS_TABSTOP,
		0) // don't set WS_EX_CONTROLPARENT here; see uitask_windows.c
	t := &tab{
		controlSingleHWND:		newControlSingleHWND(hwnd),
	}
	t.fpreferredSize = t.preferredSize
	t.fresize = t.resize
	// count tabs as 1 tab stop; the actual number of tab stops varies
	C.controlSetControlFont(t.hwnd)
	C.setTabSubclass(t.hwnd, unsafe.Pointer(t))
	return t
}

// TODO margined
func (t *tab) Append(name string, control Control) {
	c := newContainer()
	control.setParent(&controlParent{c.hwnd})
	c.setParent(&controlParent{t.hwnd})
	t.tabs = append(t.tabs, c)
	t.children = append(t.children, control)
	// initially hide tab 1..n controls; if we don't, they'll appear over other tabs, resulting in weird behavior
	if len(t.tabs) != 1 {
		t.tabs[len(t.tabs)-1].hide()
	}
	C.tabAppend(t.hwnd, toUTF16(name))
}

//export tabChanging
func tabChanging(data unsafe.Pointer, current C.LRESULT) {
	t := (*tab)(data)
	t.tabs[int(current)].hide()
}

//export tabChanged
func tabChanged(data unsafe.Pointer, new C.LRESULT) {
	t := (*tab)(data)
	t.tabs[int(new)].show()
}

//export tabTabHasChildren
func tabTabHasChildren(data unsafe.Pointer, which C.LRESULT) C.BOOL {
	t := (*tab)(data)
	if len(t.tabs) == 0 { // currently no tabs
		return C.FALSE
	}
	if t.children[int(which)].nTabStops() > 0 {
		return C.TRUE
	}
	return C.FALSE
}

func (t *tab) preferredSize(d *sizing) (width, height int) {
	for _, c := range t.children {
		w, h := c.preferredSize(d)
		if width < w {
			width = w
		}
		if height < h {
			height = h
		}
	}
	return width, height + int(C.tabGetTabHeight(t.hwnd))
}

// a tab control contains other controls; size appropriately
func (t *tab) resize(x int, y int, width int, height int, d *sizing) {
	// first, chain up to the container base to keep the Z-order correct
	// TODO use a variable for this
	t.controlSingleHWND.resize(x, y, width, height, d)

	// now resize the children
	var r C.RECT

	// figure out what the rect for each child is...
	// the tab contents are children of the tab itself, so ignore x and y, which are relative to the window!
	r.left = C.LONG(0)
	r.top = C.LONG(0)
	r.right = C.LONG(width)
	r.bottom = C.LONG(height)
	C.tabGetContentRect(t.hwnd, &r)
	// and resize tabs
	// don't resize just the current tab; resize all tabs!
	for i, _ := range t.tabs {
		// because each widget is actually a child of the Window, the origin is the one we calculated above
		t.tabs[i].resize(int(r.left), int(r.top), int(r.right - r.left), int(r.bottom - r.top), d)
		// TODO get the actual client rect
		t.children[i].resize(int(0), int(0), int(r.right - r.left), int(r.bottom - r.top), d)
	}
}

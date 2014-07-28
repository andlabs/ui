// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

/*
On Windows, container controls are just regular controls; their children have to be children of the parent window, and changing the contents of a switching container (such as a tab control) must be done manually.

TODO
- make sure all tabs cannot be deselected (that is, make sure the current tab can never have index -1)
- make sure tabs initially show the right control
- for some reason the text entry tabs show the checkbox tab until the checkbox tab is clicked, THEN they show their proper contents
*/

type tab struct {
	*widgetbase
	tabs			[]*container
}

func newTab() Tab {
	w := newWidget(C.xWC_TABCONTROL,
		C.TCS_TOOLTIPS | C.WS_TABSTOP,
		0)
	t := &tab{
		widgetbase:	w,
	}
	C.controlSetControlFont(w.hwnd)
	C.setTabSubclass(w.hwnd, unsafe.Pointer(t))
	return t
}

func (t *tab) setParent(win C.HWND) {
	t.widgetbase.setParent(win)
	for _, c := range t.tabs {
		c.child.setParent(win)
	}
}

func (t *tab) Append(name string, control Control) {
	c := new(container)
	t.tabs = append(t.tabs, c)
	c.child = control
	if t.parent != nil {
		c.child.setParent(t.parent)
	}
	C.tabAppend(t.hwnd, toUTF16(name))
}

//export tabChanging
func tabChanging(data unsafe.Pointer, current C.LRESULT) {
	t := (*tab)(data)
	t.tabs[int(current)].child.containerHide()
}

//export tabChanged
func tabChanged(data unsafe.Pointer, new C.LRESULT) {
	t := (*tab)(data)
	t.tabs[int(new)].child.containerShow()
}

// a tab control contains other controls; size appropriately
func (t *tab) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	var r C.RECT

	// figure out what the rect for each child is...
	r.left = C.LONG(x)				// load structure with the window's rect
	r.top = C.LONG(y)
	r.right = C.LONG(x + width)
	r.bottom = C.LONG(y + height)
	C.tabGetContentRect(t.hwnd, &r)
	// and allocate
	// don't allocate to just the current tab; allocate to all tabs!
	for _, c := range t.tabs {
		// because each widget is actually a child of the Window, the origin is the one we calculated above
		c.resize(int(r.left), int(r.top), int(r.right - r.left), int(r.bottom - r.top))
	}
	// and now allocate the tab control itself
	return t.widgetbase.allocate(x, y, width, height, d)
}

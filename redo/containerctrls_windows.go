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
- see if we can safely make the controls children of the tab control itself or if that would just screw our subclassing
*/

type tab struct {
	*controlbase
	tabs				[]*sizer
	supersetParent		func(p *controlParent)
	superallocate		func(x int, y int, width int, height int, d *sizing) []*allocation
}

func newTab() Tab {
	c := newControl(C.xWC_TABCONTROL,
		C.TCS_TOOLTIPS | C.WS_TABSTOP,
		0)
	t := &tab{
		controlbase:	c,
	}
	t.supersetParent = t.fsetParent
	t.fsetParent = t.tabsetParent
	t.superallocate = t.fallocate
	t.fallocate = t.taballocate
	C.controlSetControlFont(t.hwnd)
	C.setTabSubclass(t.hwnd, unsafe.Pointer(t))
	return t
}

func (t *tab) tabsetParent(p *controlParent) {
	t.supersetParent(p)
	for _, c := range t.tabs {
		c.child.setParent(p)
	}
}

func (t *tab) Append(name string, control Control) {
	s := new(sizer)
	t.tabs = append(t.tabs, s)
	s.child = control
	if t.parent != nil {
		s.child.setParent(&controlParent{t.parent})
	}
	// initially hide tab 1..n controls; if we don't, they'll appear over other tabs, resulting in weird behavior
	if len(t.tabs) != 1 {
		s.child.containerHide()
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
// TODO change this to commitResize()
func (t *tab) taballocate(x int, y int, width int, height int, d *sizing) []*allocation {
	var r C.RECT

	// figure out what the rect for each child is...
	r.left = C.LONG(x)				// load structure with the window's rect
	r.top = C.LONG(y)
	r.right = C.LONG(x + width)
	r.bottom = C.LONG(y + height)
	C.tabGetContentRect(t.hwnd, &r)
	// and allocate
	// don't allocate to just the current tab; allocate to all tabs!
	for _, s := range t.tabs {
		// because each widget is actually a child of the Window, the origin is the one we calculated above
		s.resize(int(r.left), int(r.top), int(r.right - r.left), int(r.bottom - r.top))
	}
	// and now allocate the tab control itself
	return t.superallocate(x, y, width, height, d)
}

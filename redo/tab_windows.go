// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

/*
On Windows, container controls are just regular controls that notify their parent when the user wants to do things; changing the contents of a switching container (such as a tab control) must be done manually.

We'll create a dummy window using the pre-existing Window window class for each tab page. This makes showing and hiding tabs a matter of showing and hiding one control.
*/

type tab struct {
	_hwnd	C.HWND
	tabs		[]*container
}

func newTab() Tab {
	hwnd := C.newControl(C.xWC_TABCONTROL,
		C.TCS_TOOLTIPS | C.WS_TABSTOP,
		0)		// don't set WS_EX_CONTROLPARENT here; see uitask_windows.c
	t := &tab{
		_hwnd:	hwnd,
	}
	C.controlSetControlFont(t._hwnd)
	C.setTabSubclass(t._hwnd, unsafe.Pointer(t))
	return t
}

func (t *tab) Append(name string, control Control) {
	c := newContainer(control)
	c.setParent(t._hwnd)
	t.tabs = append(t.tabs, c)
	// initially hide tab 1..n controls; if we don't, they'll appear over other tabs, resulting in weird behavior
	if len(t.tabs) != 1 {
		t.tabs[len(t.tabs) - 1].hide()
	}
	C.tabAppend(t._hwnd, toUTF16(name))
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
	if len(t.tabs) == 0 {		// currently no tabs
		return C.FALSE
	}
	if t.tabs[int(which)].nchildren > 0 {
		return C.TRUE
	}
	return C.FALSE
}

func (t *tab) hwnd() C.HWND {
	return t._hwnd
}

func (t *tab) setParent(p *controlParent) {
	basesetParent(t, p)
}

func (t *tab) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

func (t *tab) preferredSize(d *sizing) (width, height int) {
	for _, s := range t.tabs {
		w, h := s.child.preferredSize(d)
		if width < w {
			width = w
		}
		if height < h {
			height = h
		}
	}
	return width, height + int(C.tabGetTabHeight(t._hwnd))
}

// a tab control contains other controls; size appropriately
func (t *tab) commitResize(c *allocation, d *sizing) {
	var r C.RECT

	// figure out what the rect for each child is...
	// the tab contents are children of the tab itself, so ignore c.x and c.y, which are relative to the window!
	r.left = C.LONG(0)
	r.top = C.LONG(0)
	r.right = C.LONG(c.width)
	r.bottom = C.LONG(c.height)
	C.tabGetContentRect(t._hwnd, &r)
	// and resize tabs
	// don't resize just the current tab; resize all tabs!
	for _, c := range t.tabs {
		// because each widget is actually a child of the Window, the origin is the one we calculated above
		c.move(&r)
	}
	// and now resize the tab control itself
	basecommitResize(t, c, d)
}

func (t *tab) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}

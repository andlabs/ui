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
	_hwnd	C.HWND
	tabs		[]*sizer
	parent	C.HWND
}

func newTab() Tab {
	hwnd := C.newControl(C.xWC_TABCONTROL,
		C.TCS_TOOLTIPS | C.WS_TABSTOP,
		0)
	t := &tab{
		_hwnd:	hwnd,
	}
	C.controlSetControlFont(t._hwnd)
	C.setTabSubclass(t._hwnd, unsafe.Pointer(t))
	return t
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
	C.tabAppend(t._hwnd, toUTF16(name))
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

func (t *tab) hwnd() C.HWND {
	return t._hwnd
}

func (t *tab) setParent(p *controlParent) {
	basesetParent(t, p)
	for _, c := range t.tabs {
		c.child.setParent(p)
	}
	t.parent = p.hwnd
}

// TODO actually write this
func (t *tab) containerShow() {
	basecontainerShow(t)
}

// TODO actually write this
func (t *tab) containerHide() {
	basecontainerHide(t)
}

func (t *tab) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

func (t *tab) preferredSize(d *sizing) (width, height int) {
	// TODO only consider the size of the current tab?
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
	r.left = C.LONG(c.x)				// load structure with the window's rect
	r.top = C.LONG(c.y)
	r.right = C.LONG(c.x + c.width)
	r.bottom = C.LONG(c.y + c.height)
	C.tabGetContentRect(t._hwnd, &r)
	// and resize tabs
	// don't resize just the current tab; resize all tabs!
	for _, s := range t.tabs {
		// because each widget is actually a child of the Window, the origin is the one we calculated above
		s.resize(int(r.left), int(r.top), int(r.right - r.left), int(r.bottom - r.top))
	}
	// and now resize the tab control itself
	basecommitResize(t, c, d)
}

func (t *tab) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}

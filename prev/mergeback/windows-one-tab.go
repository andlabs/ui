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

TODO
- make sure all tabs cannot be deselected (that is, make sure the current tab can never have index -1)
*/

type tab struct {
	_hwnd		C.HWND
	tabs			[]*container
	switchrect	C.RECT		// size that new tab should take when switching to it
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
	c := newContainer(control)
	c.setParent(&controlParent{t._hwnd})
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
	// resize the new tab...
	t.tabs[int(new)].move(&t.switchrect)
	// ...then show
	t.tabs[int(new)].show()
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
	// the tab contents are children of the tab itself, so ignore c.x and c.y, which are relative to the window!
	r.left = C.LONG(0)
	r.top = C.LONG(0)
	r.right = C.LONG(c.width)
	r.bottom = C.LONG(c.height)
	C.tabGetContentRect(t._hwnd, &r)
	// and resize tabs
	// resize only the current tab; we trigger a resize on a tab change to make sure things look correct
	if len(t.tabs) > 0 {
		t.tabs[C.SendMessageW(t._hwnd, C.TCM_GETCURSEL, 0, 0)].move(&r)
	}
	// save the tab size so we can
	t.switchrect = r
	// and now resize the tab control itself
	basecommitResize(t, c, d)
}

func (t *tab) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}

// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

/*
On Windows, container controls are just regular controls that notify their parent when the user wants to do things; changing the contents of a switching container (such as a tab control) must be done manually.
*/

// TODO FIGURE OUT HOW OR WHY THIS IS PARTIALLY MARGINED WTF
// it's probably something I did a while ago and forgot but still wow

type tab struct {
	*controlSingleHWND
	children		[]Control
	chainresize	func(x int, y int, width int, height int, d *sizing)
}

func newTab() Tab {
	hwnd := C.newControl(C.xWC_TABCONTROL,
		C.TCS_TOOLTIPS|C.WS_TABSTOP,
		0) // don't set WS_EX_CONTROLPARENT here; see uitask_windows.c
	t := &tab{
		controlSingleHWND:		newControlSingleHWND(hwnd),
	}
	t.fpreferredSize = t.xpreferredSize
	t.chainresize = t.fresize
	t.fresize = t.xresize
	// count tabs as 1 tab stop; the actual number of tab stops varies
	C.controlSetControlFont(t.hwnd)
	C.setTabSubclass(t.hwnd, unsafe.Pointer(t))
	return t
}

// TODO margined
func (t *tab) Append(name string, control Control) {
	control.setParent(&controlParent{t.hwnd})
	t.children = append(t.children, control)
	// initially hide tab 1..n controls; if we don't, they'll appear over other tabs, resulting in weird behavior
	if len(t.children) != 1 {
		t.children[len(t.children)-1].containerHide()
	}
	C.tabAppend(t.hwnd, toUTF16(name))
}

//export tabChanging
func tabChanging(data unsafe.Pointer, current C.LRESULT) {
	t := (*tab)(data)
	t.children[int(current)].containerHide()
}

//export tabChanged
func tabChanged(data unsafe.Pointer, new C.LRESULT) {
	t := (*tab)(data)
	t.children[int(new)].containerShow()
}

//export tabTabHasChildren
func tabTabHasChildren(data unsafe.Pointer, which C.LRESULT) C.BOOL {
	t := (*tab)(data)
	if len(t.children) == 0 { // currently no tabs
		return C.FALSE
	}
	if t.children[int(which)].nTabStops() > 0 {
		return C.TRUE
	}
	return C.FALSE
}

func (t *tab) xpreferredSize(d *sizing) (width, height int) {
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

// no need to resize the other controls; we do that in tabResized() which is called by the tab subclass handler
func (t *tab) xresize(x int, y int, width int, height int, d *sizing) {
	// just chain up to the container base to keep the Z-order correct
	t.chainresize(x, y, width, height, d)
}

//export tabResized
func tabResized(data unsafe.Pointer, r C.RECT) {
	t := (*tab)(data)
	if len(t.children) == 0 {		// nothing to do
		return
	}
	d := beginResize(t.hwnd)
	// only need to resize the current tab; we resize new tabs when the tab changes in tabChanged() above
	// because each widget is actually a child of the Window, the origin is the one we calculated above
	for i := 0; i < len(t.children); i++ {
		t.children[i].resize(int(r.left), int(r.top), int(r.right - r.left), int(r.bottom - r.top), d)
	}
}

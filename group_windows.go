// 15 august 2014

package ui

// #include "winapi_windows.h"
import "C"

type group struct {
	*controlSingleHWNDWithText
	child			Control
	margined		bool
	chainresize	func(x int, y int, width int, height int, d *sizing)
}

func newGroup(text string, control Control) Group {
	hwnd := C.newControl(buttonclass,
		C.BS_GROUPBOX,
		C.WS_EX_CONTROLPARENT)
	g := &group{
		controlSingleHWNDWithText:		newControlSingleHWNDWithText(hwnd),
		child:		control,
	}
	g.fpreferredSize = g.xpreferredSize
	g.chainresize = g.fresize
	g.fresize = g.xresize
	g.fnTabStops = control.nTabStops		// groupbox itself is not tabbable but the contents might be
	g.SetText(text)
	C.controlSetControlFont(g.hwnd)
	control.setParent(&controlParent{g.hwnd})
	return g
}

func (g *group) Text() string {
	return g.text()
}

func (g *group) SetText(text string) {
	g.setText(text)
}

func (g *group) Margined() bool {
	return g.margined
}

func (g *group) SetMargined(margined bool) {
	g.margined = margined
}

const (
	groupXMargin       = 6
	groupYMarginTop    = 11 // note this value /includes the groupbox label/
	groupYMarginBottom = 7
)

func (g *group) xpreferredSize(d *sizing) (width, height int) {
	var r C.RECT

	width, height = g.child.preferredSize(d)
	if width < int(g.textlen) { // if the text is longer, try not to truncate
		width = int(g.textlen)
	}
	r.left = 0
	r.top = 0
	r.right = C.LONG(width)
	r.bottom = C.LONG(height)
	// use negative numbers to increase the size of the rectangle
	if g.margined {
		marginRectDLU(&r, -groupYMarginTop, -groupYMarginBottom, -groupXMargin, -groupXMargin, d)
	} else {
		// unforutnately, as mentioned above, the size of a groupbox includes the label and border
		// 1 character cell (4DLU x, 8DLU y) on each side (but only 3DLU on the bottom) should be enough to make up for that; TODO is not, we can change it
		// TODO make these named constants
		marginRectDLU(&r, -8, -3, -4, -4, d)
	}
	return int(r.right - r.left), int(r.bottom - r.top)
}

func (g *group) xresize(x int, y int, width int, height int, d *sizing) {
	// first, chain up to the container base to keep the Z-order correct
	g.chainresize(x, y, width, height, d)

	// now resize the child container
	var r C.RECT

	// pretend that the client area of the group box only includes the actual empty space
	// container will handle the necessary adjustments properly
	r.left = 0
	r.top = 0
	r.right = C.LONG(width)
	r.bottom = C.LONG(height)
	if g.margined {
		// see above
		marginRectDLU(&r, groupYMarginTop, groupYMarginBottom, groupXMargin, groupXMargin, d)
	} else {
		marginRectDLU(&r, 8, 3, 4, 4, d)
	}
	g.child.resize(int(r.left), int(r.top), int(r.right - r.left), int(r.bottom - r.top), d)
}

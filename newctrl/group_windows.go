// 15 august 2014

package ui

// #include "winapi_windows.h"
import "C"

type group struct {
	*controlSingleHWNDWithText
	child			Control
	margined		bool
}

func newGroup(text string, control Control) Group {
	hwnd := C.newControl(buttonclass,
		C.BS_GROUPBOX,
		C.WS_EX_CONTROLPARENT)
	g := &group{
		controlSingleHWNDWithText:		newControlSingleHWNDWithText(hwnd),
		child:		control,
	}
	g.fpreferredSize = g.preferredSize
	g.fresize = g.resize
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

const (
	groupXMargin       = 6
	groupYMarginTop    = 11 // note this value /includes the groupbox label/
	groupYMarginBottom = 7
)

func (g *group) preferredSize(d *sizing) (width, height int) {
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
		// unforutnately, as mentioned above, the size of a groupbox includes the label and border
		// 1DLU on each side should be enough to make up for that; if not, we can change it
		// TODO make these named constants
		marginRectDLU(&r, -1, -1, -1, -1, d)
	} else {
		marginRectDLU(&r, -groupYMarginTop, -groupYMarginBottom, -groupXMargin, -groupXMargin, d)
	}
	return int(r.right - r.left), int(r.bottom - r.top)
}

func (g *group) resize(x int, y int, width int, height int, d *sizing) {
	// first, chain up to the container base to keep the Z-order correct
	// TODO use a variable for this
	g.controlSingleHWNDWithText.resize(x, y, width, height, d)

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
		marginRectDLU(&r, 1, 1, 1, 1, d)
	} else {
		marginRectDLU(&r, groupYMarginTop, groupYMarginBottom, groupXMargin, groupXMargin, d)
	}
	g.child.resize(int(r.left), int(r.top), int(r.right - r.left), int(r.bottom - r.top), d)
}

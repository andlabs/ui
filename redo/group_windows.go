// 15 august 2014

package ui

// #include "winapi_windows.h"
import "C"

type group struct {
	_hwnd	C.HWND
	_textlen	C.LONG

	*container
}

func newGroup(text string, control Control) Group {
	hwnd := C.newControl(buttonclass,
		C.BS_GROUPBOX,
		0)
	g := &group{
		_hwnd:		hwnd,
		container:		newContainer(control),
	}
	g.SetText(text)
	C.controlSetControlFont(g._hwnd)
	g.container.setParent(g._hwnd)
	g.container.isGroup = true
	return g
}

func (g *group) Text() string {
	return baseText(g)
}

func (g *group) SetText(text string) {
	baseSetText(g, text)
}

func (g *group) hwnd() C.HWND {
	return g._hwnd
}

func (g *group) textlen() C.LONG {
	return g._textlen
}

func (g *group) settextlen(len C.LONG) {
	g._textlen = len
}

func (g *group) setParent(p *controlParent) {
	basesetParent(g, p)
}

func (g *group) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(g, x, y, width, height, d)
}

func (g *group) preferredSize(d *sizing) (width, height int) {
	width, height = g.child.preferredSize(d)
	if width < int(g._textlen) {		// if the text is longer, try not to truncate
		width = int(g._textlen)
	}
	// the two margin constants come from container_windows.go
	return width, height + fromdlgunitsY(groupYMarginTop, d) + fromdlgunitsY(groupYMarginBottom, d)
}

func (g *group) commitResize(c *allocation, d *sizing) {
	var r C.RECT

	// pretend that the client area of the group box only includes the actual empty space
	// container will handle the necessary adjustments properly
	r.left = 0
	r.top = 0
	r.right = C.LONG(c.width)
	r.bottom = C.LONG(c.height)
	g.container.move(&r)
	basecommitResize(g, c, d)
}

func (g *group) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(g, d)
}

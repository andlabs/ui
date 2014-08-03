// 15 july 2014

package ui

// #include "winapi_windows.h"
import "C"

type label struct {
	*controlbase
	standalone	bool
}

var labelclass = toUTF16("STATIC")

func finishNewLabel(text string, standalone bool) *label {
	c := newControl(labelclass,
		// SS_NOPREFIX avoids accelerator translation; SS_LEFTNOWORDWRAP clips text past the end
		// controls are vertically aligned to the top by default (thanks Xeek in irc.freenode.net/#winapi)
		C.SS_NOPREFIX | C.SS_LEFTNOWORDWRAP,
		0)
	c.setText(text)
	C.controlSetControlFont(c.hwnd)
	l := &label{
		controlbase:	c,
		standalone:	standalone,
	}
	return l
}

func newLabel(text string) Label {
	return finishNewLabel(text, false)
}

func newStandaloneLabel(text string) Label {
	return finishNewLabel(text, true)
}

func (l *label) Text() string {
	return l.text()
}

func (l *label) SetText(text string) {
	l.setText(text)
}

func (l *label) setParent(p *controlParent) {
	basesetParent(l.controlbase, p)
}

func (l *label) containerShow() {
	basecontainerShow(l.controlbase)
}

func (l *label) containerHide() {
	basecontainerHide(l.controlbase)
}

func (l *label) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(l, x, y, width, height, d)
}

const (
	// via http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	labelHeight = 8
	labelYOffset = 3
	// TODO the label is offset slightly by default...
)

func (l *label) preferredSize(d *sizing) (width, height int) {
	return int(l.textlen), fromdlgunitsY(labelHeight, d)
}

func (l *label) commitResize(c *allocation, d *sizing) {
	if !l.standalone {
		yoff := fromdlgunitsY(labelYOffset, d)
		c.y += yoff
		c.height -= yoff
	}
	basecommitResize(l.controlbase, c, d)
}

func (l *label) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(d)
}

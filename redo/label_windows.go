// 15 july 2014

package ui

// #include "winapi_windows.h"
import "C"

type label struct {
	_hwnd		C.HWND
	_textlen		C.LONG
	standalone	bool
}

var labelclass = toUTF16("STATIC")

func finishNewLabel(text string, standalone bool) *label {
	hwnd := C.newControl(labelclass,
		// SS_NOPREFIX avoids accelerator translation; SS_LEFTNOWORDWRAP clips text past the end
		// controls are vertically aligned to the top by default (thanks Xeek in irc.freenode.net/#winapi)
		C.SS_NOPREFIX | C.SS_LEFTNOWORDWRAP,
		0)
	l := &label{
		_hwnd:		hwnd,
		standalone:	standalone,
	}
	l.SetText(text)
	C.controlSetControlFont(l._hwnd)
	return l
}

func newLabel(text string) Label {
	return finishNewLabel(text, false)
}

func newStandaloneLabel(text string) Label {
	return finishNewLabel(text, true)
}

func (l *label) Text() string {
	return baseText(l)
}

func (l *label) SetText(text string) {
	baseSetText(l, text)
}

func (l *label) hwnd() C.HWND {
	return l._hwnd
}

func (l *label) textlen() C.LONG {
	return l._textlen
}

func (l *label) settextlen(len C.LONG) {
	l._textlen = len
}

func (l *label) setParent(p *controlParent) {
	basesetParent(l, p)
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
	return int(l._textlen), fromdlgunitsY(labelHeight, d)
}

func (l *label) commitResize(c *allocation, d *sizing) {
	if !l.standalone {
		yoff := fromdlgunitsY(labelYOffset, d)
		c.y += yoff
		c.height -= yoff
	}
	c.y -= int(d.internalLeading)
	c.height += int(d.internalLeading)
	basecommitResize(l, c, d)
}

func (l *label) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(l, d)
}

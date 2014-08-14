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
		C.WS_EX_TRANSPARENT)
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
	C.controlSetParent(l.hwnd(), p.c.hwnd)
	// don't increment p.c.nchildren here because Labels aren't tab stops
}

func (l *label) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(l, x, y, width, height, d)
}

const (
	// via http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	labelHeight = 8
	labelYOffset = 3
)

func (l *label) preferredSize(d *sizing) (width, height int) {
	return int(l._textlen), fromdlgunitsY(labelHeight, d)
}

func (l *label) commitResize(c *allocation, d *sizing) {
	if !l.standalone {
		yoff := fromdlgunitsY(labelYOffset, d)
		c.y += yoff
		c.height -= yoff
		// by default, labels are drawn offset by the internal leading (the space reserved for accents on uppercase letters)
		// the above calculation assumes otherwise, so account for the difference
		// there will be enough space left over for the internal leading anyway (at least on the standard fonts)
		// don't do this to standalone labels, otherwise those accents get cut off!
		c.y -= int(d.internalLeading)
		c.height += int(d.internalLeading)
	}
	basecommitResize(l, c, d)
}

func (l *label) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(l, d)
}

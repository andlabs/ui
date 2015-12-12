// 15 july 2014

package ui

// #include "winapi_windows.h"
import "C"

type label struct {
	*controlSingleHWNDWithText
}

var labelclass = toUTF16("STATIC")

func newLabel(text string) Label {
	hwnd := C.newControl(labelclass,
		// SS_NOPREFIX avoids accelerator translation; SS_LEFTNOWORDWRAP clips text past the end
		// controls are vertically aligned to the top by default (thanks Xeek in irc.freenode.net/#winapi)
		C.SS_NOPREFIX|C.SS_LEFTNOWORDWRAP,
		C.WS_EX_TRANSPARENT)
	l := &label{
		controlSingleHWNDWithText:		newControlSingleHWNDWithText(hwnd),
	}
	l.fpreferredSize = l.xpreferredSize
	l.fnTabStops = func() int {
		// labels are not tab stops
		return 0
	}
	l.SetText(text)
	C.controlSetControlFont(l.hwnd)
	return l
}

func (l *label) Text() string {
	return l.text()
}

func (l *label) SetText(text string) {
	l.setText(text)
}

const (
	// via http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	labelHeight  = 8
	labelYOffset = 3
)

func (l *label) xpreferredSize(d *sizing) (width, height int) {
	return int(l.textlen), fromdlgunitsY(labelHeight, d)
}

/*TODO
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
*/

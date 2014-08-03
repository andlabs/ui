// 15 july 2014

package ui

// #include "winapi_windows.h"
import "C"

type textField struct {
	*controlbase
}

var editclass = toUTF16("EDIT")

func startNewTextField(style C.DWORD) *textField {
	c := newControl(editclass,
		style | C.ES_AUTOHSCROLL | C.ES_LEFT | C.ES_NOHIDESEL | C.WS_TABSTOP,
		C.WS_EX_CLIENTEDGE)		// WS_EX_CLIENTEDGE without WS_BORDER will show the canonical visual styles border (thanks to MindChild in irc.efnet.net/#winprog)
	C.controlSetControlFont(c.hwnd)
	t := &textField{
		controlbase:	c,
	}
	return t
}

func newTextField() *textField {
	return startNewTextField(0)
}

func newPasswordField() *textField {
	return startNewTextField(C.ES_PASSWORD)
}

func (t *textField) Text() string {
	return t.text()
}

func (t *textField) SetText(text string) {
	t.setText(text)
}

func (t *textField) setParent(p *controlParent) {
	basesetParent(t.controlbase, p)
}

func (t *textField) containerShow() {
	basecontainerShow(t.controlbase)
}

func (t *textField) containerHide() {
	basecontainerHide(t.controlbase)
}

func (t *textField) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

const (
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	textfieldWidth = 107		// this is actually the shorter progress bar width, but Microsoft only indicates as wide as necessary
	textfieldHeight = 14
)

func (t *textField) preferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(textfieldWidth, d), fromdlgunitsY(textfieldHeight, d)
}

func (t *textField) commitResize(a *allocation, d *sizing) {
	basecommitResize(t.controlbase, a, d)
}

func (t *textField) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(d)
}

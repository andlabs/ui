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
	t.fpreferredSize = t.textfieldpreferredSize
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

const (
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	textfieldWidth = 107		// this is actually the shorter progress bar width, but Microsoft only indicates as wide as necessary
	textfieldHeight = 14
)

func (t *textField) textfieldpreferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(textfieldWidth, d), fromdlgunitsY(textfieldHeight, d)
}

// 23 october 2014

package ui

// #include "winapi_windows.h"
import "C"

type textbox struct {
	*controlSingleHWNDWithText
}

// TODO autohide scrollbars
func newTextbox() Textbox {
	hwnd := C.newControl(editclass,
		// TODO ES_AUTOHSCROLL/ES_AUTOVSCROLL as well?
		// TODO word wrap
		C.ES_LEFT | C.ES_MULTILINE | C.ES_NOHIDESEL | C.ES_WANTRETURN | C.WS_HSCROLL | C.WS_VSCROLL,
		C.WS_EX_CLIENTEDGE)
	t := &textbox{
		controlSingleHWNDWithText:		newControlSingleHWNDWithText(hwnd),
	}
	t.fpreferredSize = t.xpreferredSize
	C.controlSetControlFont(t.hwnd)
	return t
}

func (t *textbox) Text() string {
	return t.text()
}

func (t *textbox) SetText(text string) {
	t.setText(text)
}

// just reuse the preferred textfield width
// TODO allow alternate widths
// TODO current height probably can be better calculated
func (t *textbox) xpreferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(textfieldWidth, d), fromdlgunitsY(textfieldHeight, d) * 3
}

// 15 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type button struct {
	*controlbase
	clicked		*event
}

var buttonclass = toUTF16("BUTTON")

func startNewButton(text string, style C.DWORD) *button {
	c := newControl(buttonclass,
		style | C.WS_TABSTOP,
		0)
	c.setText(text)
	C.controlSetControlFont(c.hwnd)
	b := &button{
		controlbase:	c,
		clicked:		newEvent(),
	}
	return b
}

func newButton(text string) *button {
	b := startNewButton(text, C.BS_PUSHBUTTON)
	C.setButtonSubclass(b.hwnd, unsafe.Pointer(b))
	b.fpreferredSize = b.buttonpreferredSize
	return b
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

func (b *button) Text() string {
	return b.text()
}

func (b *button) SetText(text string) {
	b.setText(text)
}

//export buttonClicked
func buttonClicked(data unsafe.Pointer) {
	b := (*button)(data)
	b.clicked.fire()
	println("button clicked")
}

const (
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	buttonHeight = 14
)

func (b *button) buttonpreferredSize(d *sizing) (width, height int) {
	// common controls 6 thankfully provides a method to grab this...
	var size C.SIZE

	size.cx = 0		// explicitly ask for ideal size
	size.cy = 0
	if C.SendMessageW(b.hwnd, C.BCM_GETIDEALSIZE, 0, C.LPARAM(uintptr(unsafe.Pointer(&size)))) != C.FALSE {
		return int(size.cx), int(size.cy)
	}
	// that failed, fall back
	// don't worry about the error return from GetSystemMetrics(); there's no way to tell (explicitly documented as such)
	xmargins := 2 * int(C.GetSystemMetrics(C.SM_CXEDGE))
	return xmargins + int(b.textlen), fromdlgunitsY(buttonHeight, d)
}

type checkbox struct {
	*button
}

func newCheckbox(text string) *checkbox {
	c := &checkbox{
		// don't use BS_AUTOCHECKBOX here because it creates problems when refocusing (see http://blogs.msdn.com/b/oldnewthing/archive/2014/05/22/10527522.aspx)
		// we'll handle actually toggling the check state ourselves (see controls_windows.c)
		button:	startNewButton(text, C.BS_CHECKBOX),
	}
	c.fpreferredSize = c.checkboxpreferredSize
	C.setCheckboxSubclass(c.hwnd, unsafe.Pointer(c))
	return c
}

func (c *checkbox) Checked() bool {
	if C.checkboxChecked(c.hwnd) == C.FALSE {
		return false
	}
	return true
}

func (c *checkbox) SetChecked(checked bool) {
	if checked {
		C.checkboxSetChecked(c.hwnd, C.TRUE)
		return
	}
	C.checkboxSetChecked(c.hwnd, C.FALSE)
}

//export checkboxToggled
func checkboxToggled(data unsafe.Pointer) {
	c := (*checkbox)(data)
	c.clicked.fire()
	println("checkbox toggled")
}

const (
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	checkboxHeight = 10
	// from http://msdn.microsoft.com/en-us/library/windows/desktop/bb226818%28v=vs.85%29.aspx
	checkboxXFromLeftOfBoxToLeftOfLabel = 12
)

func (c *checkbox) checkboxpreferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(checkboxXFromLeftOfBoxToLeftOfLabel, d) + int(c.textlen),
		fromdlgunitsY(checkboxHeight, d)
}

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

type label struct {
	*controlbase
	standalone			bool
	supercommitResize		func(c *allocation, d *sizing)
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
	l.fpreferredSize = l.labelpreferredSize
	l.supercommitResize = l.fcommitResize
	l.fcommitResize = l.labelcommitResize
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

const (
	// via http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	labelHeight = 8
	labelYOffset = 3
	// TODO the label is offset slightly by default...
)

func (l *label) labelpreferredSize(d *sizing) (width, height int) {
	return int(l.textlen), fromdlgunitsY(labelHeight, d)
}

func (l *label) labelcommitResize(c *allocation, d *sizing) {
	if !l.standalone {
		yoff := fromdlgunitsY(labelYOffset, d)
		c.y += yoff
		c.height -= yoff
	}
	l.supercommitResize(c, d)
}

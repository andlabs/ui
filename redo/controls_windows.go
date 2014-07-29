// 15 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type widgetbase struct {
	hwnd	C.HWND
	parent	C.HWND
}

func newWidget(class C.LPCWSTR, style C.DWORD, extstyle C.DWORD) *widgetbase {
	return &widgetbase{
		hwnd:	C.newWidget(class, style, extstyle),
	}
}

// these few methods are embedded by all the various Controls since they all will do the same thing

func (w *widgetbase) setParent(win C.HWND) {
	C.controlSetParent(w.hwnd, win)
	w.parent = win
}

func (w *widgetbase) containerShow() {
	C.ShowWindow(w.hwnd, C.SW_SHOW)
}

func (w *widgetbase) containerHide() {
	C.ShowWindow(w.hwnd, C.SW_HIDE)
}

// don't embed these as exported; let each Control decide if it should

func (w *widgetbase) text() string {
	return getWindowText(w.hwnd)
}

func (w *widgetbase) settext(text string) {
	C.setWindowText(w.hwnd, toUTF16(text))
}

type button struct {
	*widgetbase
	clicked		*event
}

var buttonclass = toUTF16("BUTTON")

func startNewButton(text string, style C.DWORD) *button {
	w := newWidget(buttonclass,
		style | C.WS_TABSTOP,
		0)
	C.setWindowText(w.hwnd, toUTF16(text))
	C.controlSetControlFont(w.hwnd)
	b := &button{
		widgetbase:	w,
		clicked:		newEvent(),
	}
	return b
}

func newButton(text string) *button {
	b := startNewButton(text, C.BS_PUSHBUTTON)
	C.setButtonSubclass(b.hwnd, unsafe.Pointer(b))
	return b
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

func (b *button) Text() string {
	return b.text()
}

func (b *button) SetText(text string) {
	b.settext(text)
}

//export buttonClicked
func buttonClicked(data unsafe.Pointer) {
	b := (*button)(data)
	b.clicked.fire()
	println("button clicked")
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

type textField struct {
	*widgetbase
}

var editclass = toUTF16("EDIT")

func startNewTextField(style C.DWORD) *textField {
	w := newWidget(editclass,
		style | C.ES_AUTOHSCROLL | C.ES_LEFT | C.ES_NOHIDESEL | C.WS_TABSTOP,
		C.WS_EX_CLIENTEDGE)		// WS_EX_CLIENTEDGE without WS_BORDER will show the canonical visual styles border (thanks to MindChild in irc.efnet.net/#winprog)
	C.controlSetControlFont(w.hwnd)
	return &textField{
		widgetbase:	w,
	}
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
	t.settext(text)
}

type label struct {
	*widgetbase
	standalone	bool
}

var labelclass = toUTF16("STATIC")

func finishNewLabel(text string, standalone bool) *label {
	w := newWidget(labelclass,
		// SS_NOPREFIX avoids accelerator translation; SS_LEFTNOWORDWRAP clips text past the end
		// controls are vertically aligned to the top by default (thanks Xeek in irc.freenode.net/#winapi)
		C.SS_NOPREFIX | C.SS_LEFTNOWORDWRAP,
		0)
	C.setWindowText(w.hwnd, toUTF16(text))
	C.controlSetControlFont(w.hwnd)
	l := &label{
		widgetbase:	w,
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
	l.settext(text)
}

// TODO label commitResize

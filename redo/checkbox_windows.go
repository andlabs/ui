// 15 july 2014

package ui

import (
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

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

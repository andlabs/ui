// 2 march 2014
package ui

import (
	// ...
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
import "C"

// NSAlert styles.
const (
	_NSWarningAlertStyle = 0			// default
	_NSInformationalAlertStyle = 1
	_NSCriticalAlertStyle = 2
)

var (
	_NSAlert = objc_getClass("NSAlert")

	_setMessageText = sel_getUid("setMessageText:")
	_setInformativeText = sel_getUid("setInformativeText:")
	_setAlertStyle = sel_getUid("setAlertStyle:")
	_addButtonWithTitle = sel_getUid("addButtonWithTitle:")
	_runModal = sel_getUid("runModal")
)

func _msgBox(title string, text string, style uintptr, button0 string) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		box := objc_new(_NSAlert)
		// TODO is this appropriate for title?
		C.objc_msgSend_id(box, _setMessageText, toNSString(title))
		C.objc_msgSend_id(box, _setInformativeText, toNSString(text))
		objc_msgSend_uint(box, _setAlertStyle, style)
		C.objc_msgSend_id(box, _addButtonWithTitle, toNSString(button0))
		C.objc_msgSend_noargs(box, _runModal)
		ret <- struct{}{}
	}
	<-ret
}

func msgBox(title string, text string) {
	// TODO _NSInformationalAlertStyle?
	_msgBox(title, text, _NSWarningAlertStyle, "OK")
}

func msgBoxError(title string, text string) {
	_msgBox(title, text, _NSCriticalAlertStyle, "OK")
}

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

func _msgBox(primarytext string, secondarytext string, style uintptr, button0 string) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		box := C.objc_msgSend_noargs(_NSAlert, _new)
		C.objc_msgSend_id(box, _setMessageText, toNSString(primarytext))
		if secondarytext != "" {
			C.objc_msgSend_id(box, _setInformativeText, toNSString(secondarytext))
		}
		C.objc_msgSend_uint(box, _setAlertStyle, C.uintptr_t(style))
		C.objc_msgSend_id(box, _addButtonWithTitle, toNSString(button0))
		C.objc_msgSend_noargs(box, _runModal)
		ret <- struct{}{}
	}
	<-ret
}

func msgBox(primarytext string, secondarytext string) {
	// TODO _NSInformationalAlertStyle?
	_msgBox(primarytext, secondarytext, _NSWarningAlertStyle, "OK")
}

func msgBoxError(primarytext string, secondarytext string) {
	_msgBox(primarytext, secondarytext, _NSCriticalAlertStyle, "OK")
}

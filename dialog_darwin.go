// 2 march 2014

package ui

import (
	// ...
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
// #include "dialog_darwin.h"
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

func _msgBox(primarytext string, secondarytext string, style uintptr) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		primary := toNSString(primarytext)
		secondary := C.id(nil)
		if secondarytext != "" {
			secondary = toNSString(secondarytext)
		}
		switch style {
		case 0:		// normal
			C.msgBox(primary, secondary)
		case 1:		// error
			C.msgBoxError(primary, secondary)
		}
		ret <- struct{}{}
	}
	<-ret
}

func msgBox(primarytext string, secondarytext string) {
	_msgBox(primarytext, secondarytext, 0)
}

func msgBoxError(primarytext string, secondarytext string) {
	_msgBox(primarytext, secondarytext, 1)
}

// 2 march 2014

package ui

import (
	// ...
)

// #include "objc_darwin.h"
import "C"

func _msgBox(parent *Window, primarytext string, secondarytext string, style uintptr) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		var pwin C.id = nil

		if parent != nil {
			pwin = parent.sysData.id
		}
		primary := toNSString(primarytext)
		secondary := C.id(nil)
		if secondarytext != "" {
			secondary = toNSString(secondarytext)
		}
		switch style {
		case 0:		// normal
			C.msgBox(pwin, primary, secondary)
		case 1:		// error
			C.msgBoxError(pwin, primary, secondary)
		}
		ret <- struct{}{}
	}
	<-ret
}

func msgBox(parent *Window, primarytext string, secondarytext string) {
	_msgBox(parent, primarytext, secondarytext, 0)
}

func msgBoxError(parent *Window, primarytext string, secondarytext string) {
	_msgBox(parent, primarytext, secondarytext, 1)
}

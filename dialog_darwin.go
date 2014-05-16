// 2 march 2014

package ui

import (
	// ...
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
import "C"

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

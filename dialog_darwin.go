// 2 march 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

//export dialog_send
func dialog_send(pchan unsafe.Pointer, res C.intptr_t) {
	rchan := (*chan int)(pchan)
	go func() {		// send it in a new goroutine like we do with everything else
		*rchan <- int(res)
	}()
}

func _msgBox(parent *Window, primarytext string, secondarytext string, style uintptr) chan int {
	ret := make(chan int)
	uitask <- func() {
		var pwin C.id = nil

		if parent != dialogWindow {
			pwin = parent.sysData.id
		}
		primary := toNSString(primarytext)
		secondary := C.id(nil)
		if secondarytext != "" {
			secondary = toNSString(secondarytext)
		}
		switch style {
		case 0:		// normal
			C.msgBox(pwin, primary, secondary, unsafe.Pointer(&ret))
		case 1:		// error
			C.msgBoxError(pwin, primary, secondary, unsafe.Pointer(&ret))
		}
	}
	return ret
}

func (w *Window) msgBox(primarytext string, secondarytext string) (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		<-_msgBox(w, primarytext, secondarytext, 0)
		done <- struct{}{}
	}()
	return done
}

func (w *Window) msgBoxError(primarytext string, secondarytext string) (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		<-_msgBox(w, primarytext, secondarytext, 1)
		done <- struct{}{}
	}()
	return done
}

// 28 february 2014
package main

import (
	"fmt"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
// extern void windowShouldClose(id, SEL, id);
// extern void buttonClicked(id, SEL, id);
// extern void gotNotification(id, SEL, id);
import "C"

// TODO move these around later
var (
	_stop = sel_getUid("stop:")
)

//export windowShouldClose
func windowShouldClose(self C.id, sel C.SEL, sender C.id) {
	fmt.Println("-[hello windowShouldClose:]")
	C.objc_msgSend_id(NSApp, _stop, sender)
}

//export buttonClicked
func buttonClicked(self C.id, sel C.SEL, sender C.id) {
	fmt.Println("button clicked; sending notification...")
	notify("button")
}

//export gotNotification
func gotNotification(self C.id, sel C.SEL, object C.id) {
	fmt.Printf("got notification from %s\n", fromNSString(object))
}

func mk(name string, selW C.SEL, selB C.SEL, selN C.SEL) C.id {
	class := newClass(name)
	addDelegateMethod(class, selW, C.windowShouldClose)
	addDelegateMethod(class, selB, C.buttonClicked)
	addDelegateMethod(class, selN, C.gotNotification)
	return objc_getClass(name)
}

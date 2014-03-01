// 27 february 2014
package main

import (
	"fmt"
	"unsafe"
)

// #cgo CFLAGS: -Dqqq
// #cgo LDFLAGS: -lobjc -framework Foundation
// #include <stdlib.h>
// #include "objc_darwin.h"
// extern void windowShouldClose(id, SEL, id);
// extern void buttonClicked(id, SEL, id);
// extern void gotNotification(id, SEL, id);
// /* because cgo doesn't like Nil */
// extern Class NilClass;		/* defined in runtimetest.go due to cgo limitations */
import "C"

var (
	_NSObject_Class = C.object_getClass(_NSObject)
)

func newClass(name string) C.Class {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	c := C.objc_allocateClassPair(_NSObject_Class, cname, 0)
	if c == C.NilClass {
		panic("unable to create Objective-C class " + name)
	}
	return c
}

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

func addOurMethod(class C.Class, sel C.SEL, imp C.IMP) {
//	ty := []C.char{'v', '@', ':', 0}		// according to the example for class_addMethod()
	ty := []C.char{'v', '@', ':', '@', 0}

	// clas methods get stored in the metaclass; the objc_allocateClassPair() docs say this will work
//	metaclass := C.object_getClass(C.id(unsafe.Pointer(class)))
//	ok := C.class_addMethod(metaclass,
	ok := C.class_addMethod(class,
		sel,
		imp,
		&ty[0])
	if ok == C.BOOL(C.NO) {
		panic("unable to add ourMethod")
	}
}

func mk(name string, selW C.SEL, selB C.SEL, selN C.SEL) C.id {
	class := newClass(name)
	addOurMethod(class, selW,
	// using &C.ourMethod causes faults for some reason
		C.IMP(unsafe.Pointer(C.windowShouldClose)))
	C.objc_registerClassPair(class)
	addOurMethod(class, selB,
		C.IMP(unsafe.Pointer(C.buttonClicked)))
	addOurMethod(class, selN,
		C.IMP(unsafe.Pointer(C.gotNotification)))
	return objc_getClass(name)
}

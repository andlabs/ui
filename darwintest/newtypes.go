// 27 february 2014
package main

import (
	"fmt"
	"unsafe"
)

// #cgo CFLAGS: -Dqqq
// #cgo LDFLAGS: -lobjc -framework Foundation
// #include <stdlib.h>
// #include <objc/message.h>
// #include <objc/objc.h>
// #include <objc/runtime.h>
// extern void windowShouldClose(id, SEL, id);
// extern id objc_msgSend_id(id, SEL, id);
// extern void buttonClicked(id, SEL, id);
// extern void gotNotification(id, SEL, id);
// extern id objc_msgSend_id_id_id(id, SEL, id, id, id);
// /* cgo doesn't like nil or Nil */
// extern id objc_msgSend_noargs(id, SEL);
// extern Class NilClass; /* in runtimetest.go because of cgo limitations */
// extern id Nilid;
import "C"

var NSObject = C.object_getClass(objc_getClass("NSObject"))

func newClass(name string) C.Class {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	c := C.objc_allocateClassPair(NSObject, cname, 0)
	if c == C.NilClass {
		panic("unable to create Objective-C class " + name)
	}
	return c
}

//export windowShouldClose
func windowShouldClose(self C.id, sel C.SEL, sender C.id) {
	fmt.Println("-[hello windowShouldClose:]")
	C.objc_msgSend_id(NSApp,
		sel_getUid("stop:"),
		sender)
}

//export buttonClicked
func buttonClicked(self C.id, sel C.SEL, sender C.id) {
	fmt.Println("button clicked; sending notification...")
	notify("button")
}

//export gotNotification
func gotNotification(self C.id, sel C.SEL, note C.id) {
	data := C.objc_msgSend_noargs(note,
		sel_getUid("userInfo"))
	val := C.objc_msgSend_id(data,
		sel_getUid("objectForKey:"),
		notekey)
	source := (*C.char)(unsafe.Pointer(
		C.objc_msgSend_noargs(val,
			sel_getUid("UTF8String"))))
	fmt.Println("got notification from %s",
		C.GoString(source))
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

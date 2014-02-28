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
// extern void ourMethod(id, SEL);
// /* cgo doesn't like Nil */
// extern Class NilClass; /* in runtimetest.go because of cgo limitations */
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

//export ourMethod
func ourMethod(self C.id, sel C.SEL) {
	fmt.Println("hello, world")
}

func addOurMethod(class C.Class, sel C.SEL) {
	ty := []C.char{'v', '@', ':', 0}		// according to the example for class_addMethod()

	// clas methods get stored in the metaclass; the objc_allocateClassPair() docs say this will work
	metaclass := C.object_getClass(C.id(unsafe.Pointer(class)))
	ok := C.class_addMethod(metaclass,
		sel,
		// using &C.ourMethod causes faults for some reason
		C.IMP(unsafe.Pointer(C.ourMethod)),
		&ty[0])
	if ok == C.BOOL(C.NO) {
		panic("unable to add ourMethod")
	}
}

func mk(name string, sel C.SEL) C.id {
	class := newClass(name)
	addOurMethod(class, sel)
	C.objc_registerClassPair(class)
	return objc_getClass(name)
}

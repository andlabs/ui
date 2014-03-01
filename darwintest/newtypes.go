// 27 february 2014
package main

import (
	"unsafe"
)

// #cgo LDFLAGS: -lobjc -framework Foundation
// #include <stdlib.h>
// #include "objc_darwin.h"
// /* because cgo doesn't like Nil */
// Class NilClass = Nil;
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
	C.objc_registerClassPair(c)
	return c
}

// according to errors spit out by cgo, C function pointers are unsafe.Pointer
func addDelegateMethod(class C.Class, sel C.SEL, imp unsafe.Pointer) {
	// maps to void (*)(id, SEL, id)
	ty := []C.char{'v', '@', ':', '@', 0}

	// clas methods get stored in the metaclass; the objc_allocateClassPair() docs say this will work
	// metaclass := C.object_getClass(C.id(unsafe.Pointer(class)))
	// we're adding instance methods, so just class will do
	ok := C.class_addMethod(class,
		sel,
		C.IMP(imp),
		&ty[0])
	if ok == C.BOOL(C.NO) {
		panic("unable to add ourMethod")
	}
}

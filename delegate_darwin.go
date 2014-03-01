// 27 february 2014
package ui

import (
	"fmt"
	"unsafe"
)

/*
This creates a class goAppDelegate that will be used as the delegate for /everything/. Specifically, it:
	- runs uitask requests (uitask:)
	- handles window close events (windowShouldClose:)
	- handles window resize events (windowDidResize: (TODO also windowDidEndLiveResize:?))
	- handles button click events (buttonClick:)
*/

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
// #include "objc_darwin.h"
// extern void appDelegate_uitask(id, SEL, id);		/* from uitask_darwin.go */
// extern BOOL appDelegate_windowShouldClose(id, SEL, id);
// extern void appDelegate_windowDidResize(id, SEL, id);
import "C"

var (
	appDelegate C.id
)

const (
	_goAppDelegate = "goAppDelegate"
)

var (
	_uitask = sel_getUid("uitask:")
	_windowShouldClose = sel_getUid("windowShouldClose:")
	_windowDidResize = sel_getUid("windowDidResize:")
)

func mkAppDelegate() error {
	appdelegateclass, err := makeDelegateClass(_goAppDelegate)
	if err != nil {
		return fmt.Errorf("error creating NSApplication delegate: %v", err)
	}
	err = addDelegateMethod(appdelegateclass, _uitask,
		C.appDelegate_uitask, delegate_void)
	if err != nil {
		return fmt.Errorf("error adding NSApplication delegate uitask: method (to do UI tasks): %v", err)
	}
	err = addDelegateMethod(appdelegateclass, _windowShouldClose,
		C.appDelegate_windowShouldClose, delegate_bool)
	if err != nil {
		return fmt.Errorf("error adding NSApplication delegate windowShouldClose: method (to handle window close button events): %v", err)
	}
	err = addDelegateMethod(appdelegateclass, _windowDidResize,
		C.appDelegate_windowDidResize, delegate_void)
	if err != nil {
		return fmt.Errorf("error adding NSApplication delegate windowDidResize: method (to handle window resize events): %v", err)
	}
	// TODO using objc_new() causes a segfault; find out why
	// TODO make alloc followed by init (I thought NSObject provided its own init?)
	appDelegate = objc_alloc(objc_getClass(_goAppDelegate))
	return nil
}

//export appDelegate_windowShouldClose
func appDelegate_windowShouldClose(self C.id, sel C.SEL, win C.id) C.BOOL {
	sysData := getSysData(win)
	sysData.signal()
	return C.BOOL(C.NO)		// don't close
}

var (
	_object = sel_getUid("object")
)

//export appDelegate_windowDidResize
func appDelegate_windowDidResize(self C.id, sel C.SEL, notification C.id) {
	win := C.objc_msgSend_noargs(notification, _object)
	sysData := getSysData(win)
	r := C.objc_msgSend_stret_rect_noargs(win, _frame)
	if sysData.resize != nil {
		err := sysData.resize(int(r.x), int(r.y), int(r.width), int(r.height))
		if err != nil {
			panic("child resize failed: " + err.Error())
		}
	}
}

// this actually constructs the delegate class

var (
	_NSObject_Class = C.object_getClass(_NSObject)
)

func makeDelegateClass(name string) (C.Class, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	c := C.objc_allocateClassPair(_NSObject_Class, cname, 0)
	if c == C.NilClass {
		return C.NilClass, fmt.Errorf("unable to create Objective-C class %s; reason unknown", name)
	}
	C.objc_registerClassPair(c)
	return c, nil
}

var (
	delegate_void = []C.char{'v', '@', ':', '@', 0}		// void (*)(id, SEL, id)
	delegate_bool = []C.char{'c', '@', ':', '@', 0}		// BOOL (*)(id, SEL, id)
)

// according to errors spit out by cgo, C function pointers are unsafe.Pointer
func addDelegateMethod(class C.Class, sel C.SEL, imp unsafe.Pointer, ty []C.char) error {
	ok := C.class_addMethod(class, sel, C.IMP(imp), &ty[0])
	if ok == C.BOOL(C.NO) {
		// TODO get function name
		return fmt.Errorf("unable to add selector %v/imp %v (reason unknown)", sel, imp)
	}
	return nil
}

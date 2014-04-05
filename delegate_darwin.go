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
	- handles window resize events (windowDidResize:)
	- handles button click events (buttonClicked:)
*/

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
// #include "objc_darwin.h"
// extern void appDelegate_uitask(id, SEL, id);		/* from uitask_darwin.go */
// extern BOOL appDelegate_windowShouldClose(id, SEL, id);
// extern void appDelegate_windowDidResize(id, SEL, id);
// extern void appDelegate_buttonClicked(id, SEL, id);
import "C"

var (
	appDelegate C.id
)

const (
	_goAppDelegate = "goAppDelegate"
)

var (
	_uitask = sel_getUid("uitask:")					// used by uitask_darwin.go
	_buttonClicked = sel_getUid("buttonClicked:")		// used by sysdata_darwin.go
)

var appDelegateSels = []selector{
	selector{"uitask:", uintptr(C.appDelegate_uitask), sel_void_id,
		"performing/dispatching UI events"},
	selector{"windowShouldClose:", uintptr(C.appDelegate_windowShouldClose), sel_bool_id,
		"handling window close button events"},
	selector{"windowDidResize:", uintptr(C.appDelegate_windowDidResize), sel_void_id,
		"handling window resize events"},
	selector{"buttonClicked:", uintptr(C.appDelegate_buttonClicked), sel_bool_id,
		"handling button clicks"},
}

func mkAppDelegate() error {
	err := makeClass(_goAppDelegate, _NSObject, appDelegateSels,
		"application delegate (handles events)")
	if err != nil {
		return err
	}
	appDelegate = C.objc_msgSend_noargs(objc_getClass(_goAppDelegate), _new)
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
	_display = sel_getUid("display")
)

//export appDelegate_windowDidResize
func appDelegate_windowDidResize(self C.id, sel C.SEL, notification C.id) {
	win := C.objc_msgSend_noargs(notification, _object)
	s := getSysData(win)
	wincv := C.objc_msgSend_noargs(win, _contentView)		// we want the content view's size, not the window's; selector defined in sysdata_darwin.go
	r := C.objc_msgSend_stret_rect_noargs(wincv, _frame)
	if s.resize != nil {
		// winheight is used here because (0,0) is the bottom-left corner, not the top-left corner
		s.resizes = s.resizes[0:0]		// set len to 0 without changing cap
		s.resize(0, 0, int(r.width), int(r.height), &s.resizes)
		for _, s := range s.resizes {
			err := s.sysData.setRect(s.x, s.y, s.width, s.height, int(r.height))
			if err != nil {
				panic("child resize failed: " + err.Error())
			}
		}
	}
	C.objc_msgSend_noargs(win, _display)		// redraw everything; TODO only if resize() was called?
}

//export appDelegate_buttonClicked
func appDelegate_buttonClicked(self C.id, sel C.SEL, button C.id) {
	sysData := getSysData(button)
	sysData.signal()
}

// this actually constructs the delegate class

var (
	delegate_void = []C.char{'v', '@', ':', '@', 0}		// void (*)(id, SEL, id)
	delegate_bool = []C.char{'c', '@', ':', '@', 0}		// BOOL (*)(id, SEL, id)
)

// according to errors spit out by cgo, C function pointers are unsafe.Pointer
func addDelegateMethod(class C.Class, sel C.SEL, imp unsafe.Pointer, ty []C.char) error {
	ok := C.class_addMethod(class, sel, C.IMP(imp), &ty[0])
	if ok == C.BOOL(C.NO) {
		// TODO get function name
		return fmt.Errorf("unable to add selector %v/imp %v to class %v (reason unknown)", sel, imp, class)
	}
	return nil
}

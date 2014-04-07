// 27 february 2014

package ui

import (
	// ...
)

/*
This creates a class goAppDelegate that will be used as the delegate for /everything/. Specifically, it:
	- runs uitask requests (uitask:)
	- handles window close events (windowShouldClose:)
	- handles window resize events (windowDidResize:)
	- handles button click events (buttonClicked:)
	- handles the application-global Quit event (such as from the Dock) (applicationShouldTerminate)
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
	selector{"applicationShouldTerminate:", uintptr(C._appDelegate_applicationShouldTerminate), sel_terminatereply_id,
		"handling Quit menu items (such as from the Dock)/the AppQuit channel"},
}

func mkAppDelegate() error {
	id, err := makeClass(_goAppDelegate, _NSObject, appDelegateSels,
		"application delegate (handles events)")
	if err != nil {
		return err
	}
	appDelegate = C.objc_msgSend_noargs(id, _new)
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
	// winheight is used here because (0,0) is the bottom-left corner, not the top-left corner
	s.doResize(0, 0, int(r.width), int(r.height), int(r.height))
	C.objc_msgSend_noargs(win, _display)		// redraw everything; TODO only if resize() was called?
}

//export appDelegate_buttonClicked
func appDelegate_buttonClicked(self C.id, sel C.SEL, button C.id) {
	sysData := getSysData(button)
	sysData.signal()
}

//export appDelegate_applicationShouldTerminate
func appDelegate_applicationShouldTerminate() {
	// asynchronous so as to return control to the event loop
	go func() {
		AppQuit <- struct{}{}
	}()
	// xxx in bleh_darwin.m tells Cocoa not to quit
}

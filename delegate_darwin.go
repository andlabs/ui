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
// /* TODO this goes in objc_darwin.h once I take care of everything else */
// extern id makeAppDelegate(void);
// extern id windowGetContentView(id);
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

func mkAppDelegate() error {
	appDelegate = C.makeAppDelegate()
	return nil
}

//export appDelegate_windowShouldClose
func appDelegate_windowShouldClose(win C.id) {
	sysData := getSysData(win)
	sysData.signal()
}

var (
	_object = sel_getUid("object")
	_display = sel_getUid("display")
)

//export appDelegate_windowDidResize
func appDelegate_windowDidResize(win C.id) {
	s := getSysData(win)
	wincv := C.windowGetContentView(win)		// we want the content view's size, not the window's
	r := C.objc_msgSend_stret_rect_noargs(wincv, _frame)
	// winheight is used here because (0,0) is the bottom-left corner, not the top-left corner
	s.doResize(0, 0, int(r.width), int(r.height), int(r.height))
	C.objc_msgSend_noargs(win, _display)		// redraw everything; TODO only if resize() was called?
}

//export appDelegate_buttonClicked
func appDelegate_buttonClicked(button C.id) {
	sysData := getSysData(button)
	sysData.signal()
}

//export appDelegate_applicationShouldTerminate
func appDelegate_applicationShouldTerminate() {
	// asynchronous so as to return control to the event loop
	go func() {
		AppQuit <- struct{}{}
	}()
}

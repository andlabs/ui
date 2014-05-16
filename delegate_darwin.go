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
import "C"

var (
	appDelegate C.id
)

func makeAppDelegate() {
	appDelegate = C.makeAppDelegate()
}

//export appDelegate_windowShouldClose
func appDelegate_windowShouldClose(win C.id) {
	sysData := getSysData(win)
	sysData.signal()
}

//export appDelegate_windowDidResize
func appDelegate_windowDidResize(win C.id) {
	s := getSysData(win)
	wincv := C.windowGetContentView(win)		// we want the content view's size, not the window's
	r := C.frame(wincv)
	// winheight is used here because (0,0) is the bottom-left corner, not the top-left corner
	s.doResize(0, 0, int(r.width), int(r.height), int(r.height))
	C.display(win)			// redraw everything; TODO only if resize() was called?
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

// 27 february 2014

package ui

/*
This creates a class goAppDelegate that will be used as the delegate for /everything/. Specifically, it:
	- handles window close events (windowShouldClose:)
	- handles window resize events (windowDidResize:)
	- handles button click events (buttonClicked:)
	- handles the application-global Quit event (such as from the Dock) (applicationShouldTerminate)
*/

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
func appDelegate_windowShouldClose(win C.id) C.BOOL {
	sysData := getSysData(win)
	return toBOOL(sysData.close())
}

//export appDelegate_windowDidResize
func appDelegate_windowDidResize(win C.id) {
	s := getSysData(win)
	wincv := C.windowGetContentView(win) // we want the content view's size, not the window's
	r := C.frame(wincv)
	// (0,0) is the bottom left corner but this is handled in sysData.translateAllocationCoords()
	s.resizeWindow(int(r.width), int(r.height))
	C.display(win) // redraw everything
}

//export appDelegate_buttonClicked
func appDelegate_buttonClicked(button C.id) {
	sysData := getSysData(button)
	sysData.event()
}

// 1 march 2014

package ui

// #include "objc_darwin.h"
import "C"

/*
Cocoa doesn't provide a reliable way to get the preferred size of a control (you're supposed to use Interface Builder and have it set up autoresizing for you). The best we can do is call [control sizeToFit] (which is defined for NSControls and has a custom implementation for the other types here) and read the preferred size. Though this changes the size, we're immediately overriding the change on return from sysData.preferredSize(), so no harm done. (This is similar to what we are doing with GTK+, except GTK+ does not actually change the size.)
*/

// standard case: control immediately passed in
func controlPrefSize(control C.id, alternate C.BOOL) (width int, height int, yoff int) {
	r := C.controlPrefSize(control, alternate)
	return int(r.width), int(r.height), int(r.yoff)
}

// Labels have special yoff calculation
func labelPrefSize(control C.id, alternate C.BOOL) (width int, height int, yoff int) {
	r := C.labelPrefSize(control, alternate)
	return int(r.width), int(r.height), int(r.yoff)
}

// NSTableView is actually in a NSScrollView so we have to get it out first
func listboxPrefSize(control C.id, alternate C.BOOL) (width int, height int, yoff int) {
	r := C.listboxPrefSize(control, alternate)
	return int(r.width), int(r.height), int(r.yoff)
}

// and for type checking reasons, progress bars are separate despite responding to -[sizeToFit]
func pbarPrefSize(control C.id, alternate C.BOOL) (width int, height int, yoff int) {
	r := C.pbarPrefSize(control, alternate)
	return int(r.width), int(r.height), int(r.yoff)
}

// Areas know their own preferred size
func areaPrefSize(control C.id, alternate C.BOOL) (width int, height int, yoff int) {
	r := C.areaPrefSize(control, alternate)
	return int(r.width), int(r.height), int(r.yoff)
}

var prefsizefuncs = [nctypes]func(C.id, C.BOOL) (int, int, int){
	c_button:      controlPrefSize,
	c_checkbox:    controlPrefSize,
	c_combobox:    controlPrefSize,
	c_lineedit:    controlPrefSize,
	c_label:       labelPrefSize,
	c_listbox:     listboxPrefSize,
	c_progressbar: pbarPrefSize,
	c_area:        areaPrefSize,
}

func (s *sysData) preferredSize() (width int, height int, yoff int) {
	return prefsizefuncs[s.ctype](s.id, toBOOL(s.alternate))
}

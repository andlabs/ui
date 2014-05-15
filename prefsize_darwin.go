// 1 march 2014

package ui

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
// #include "prefsize_darwin.h"
import "C"

/*
Cocoa doesn't provide a reliable way to get the preferred size of a control (you're supposed to use Interface Builder and have it set up autoresizing for you). The best we can do is call [control sizeToFit] (which is defined for NSControls and has a custom implementation for the other types here) and read the preferred size. Though this changes the size, we're immediately overriding the change on return from sysData.preferredSize(), so no harm done. (This is similar to what we are doing with GTK+, except GTK+ does not actually change the size.)
*/

var (
	_sizeToFit = sel_getUid("sizeToFit")
	// _frame in sysdata_darwin.go
)

// standard case: control immediately passed in
func controlPrefSize(control C.id) (width int, height int) {
	r := C.controlPrefSize(control)
	return int(r.width), int(r.height)
}

// NSTableView is actually in a NSScrollView so we have to get it out first
func listboxPrefSize(control C.id) (width int, height int) {
	r := C.listboxPrefSize(control)
	return int(r.width), int(r.height)
}

// and for type checking reasons, progress bars are separate despite responding to -[sizeToFit]
func pbarPrefSize(control C.id) (width int, height int) {
	r := C.pbarPrefSize(control)
	return int(r.width), int(r.height)
}

var prefsizefuncs = [nctypes]func(C.id) (int, int){
	c_button:			controlPrefSize,
	c_checkbox:		controlPrefSize,
	c_combobox:		controlPrefSize,
	c_lineedit:		controlPrefSize,
	c_label:			controlPrefSize,
	c_listbox:			listboxPrefSize,
	c_progressbar:		pbarPrefSize,
}

func (s *sysData) preferredSize() (width int, height int) {
	return prefsizefuncs[s.ctype](s.id)
}

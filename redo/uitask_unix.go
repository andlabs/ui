// +build !windows,!darwin

// 7 july 2014

package ui

import (
	"unsafe"
)

// #cgo pkg-config: gtk+-3.0
// #include "gtk_unix.h"
// extern gboolean doissue(gpointer data);
import "C"

func uiinit() error {
	// TODO replace with the error-checking version
	C.gtk_init(nil, nil)
	return nil
}

func uimsgloop() {
	C.gtk_main()
}

func uistop() {
	C.gtk_main_quit()
}

func issue(f func()) {
	C.gdk_threads_add_idle(C.GSourceFunc(C.doissue), C.gpointer(unsafe.Pointer(&f)))
}

//export doissue
func doissue(data C.gpointer) C.gboolean {
	perform(unsafe.Pointer(data))
	return C.FALSE		// don't repeat
}

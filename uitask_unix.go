// +build !windows,!darwin

// 7 july 2014

package ui

import (
	"fmt"
	"unsafe"
)

// #cgo pkg-config: gtk+-3.0
// #cgo CFLAGS: --std=c99
// #include "gtk_unix.h"
// extern gboolean xdoissue(gpointer data);
import "C"

func uiinit() error {
	var err *C.GError = nil // redundant in Go, but let's explicitly assign it anyway

	// gtk_init_with_args() gives us error info (thanks chpe in irc.gimp.net/#gtk+)
	// don't worry about GTK+'s command-line arguments; they're also available as environment variables (thanks mclasen in irc.gimp.net/#gtk+)
	result := C.gtk_init_with_args(nil, nil, nil, nil, nil, &err)
	if result == C.FALSE {
		return fmt.Errorf("error actually initilaizing GTK+: %s", fromgstr(err.message))
	}
	return nil
}

func uimsgloop() {
	C.gtk_main()
}

func uistop() {
	C.gtk_main_quit()
}

func issue(f *func()) {
	C.gdk_threads_add_idle(C.GSourceFunc(C.xdoissue), C.gpointer(unsafe.Pointer(f)))
}

//export xdoissue
func xdoissue(data C.gpointer) C.gboolean {
	perform(unsafe.Pointer(data))
	return C.FALSE // don't repeat
}

//export doissue
func doissue(data unsafe.Pointer) {
	// for the modal queue functions
	perform(data)
}

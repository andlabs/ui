// +build !windows,!darwin,!plan9

// 16 february 2014

package ui

import (
	"fmt"
)

// #cgo pkg-config: gtk+-3.0
// #include "gtk_unix.h"
import "C"

var uitask chan func()

func uiinit() error {
	err := gtk_init()
	if err != nil {
		return fmt.Errorf("gtk_init() failed: %v", err)
	}

	// do this only on success, just to be safe
	uitask = make(chan func())
	return nil
}

func ui() {
	// thanks to tristan and Daniel_S in irc.gimp.net/#gtk
	// see our_idle_callback in callbacks_unix.go for details
	go func() {
		for {
			var f func()

			select {
			case f = <-uitask:
				// do nothing
			case <-Stop:
				f = func() {
					C.gtk_main_quit()
				}
			}
			done := make(chan struct{})
			gdk_threads_add_idle(&gtkIdleOp{
				what: f,
				done: done,
			})
			<-done
			close(done)
		}
	}()

	C.gtk_main()
}

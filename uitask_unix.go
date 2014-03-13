// +build !windows,!darwin,!plan9

// 16 february 2014
package ui

import (
	"fmt"
	"runtime"
)

var uitask chan func()

func ui(main func()) error {
	runtime.LockOSThread()

	uitask = make(chan func())
	if gtk_init() != true {
		return fmt.Errorf("gtk_init failed (reason unknown; TODO)")
	}

	// thanks to tristan and Daniel_S in irc.gimp.net/#gtk
	// see our_idle_callback in callbacks_unix.go for details
	go func() {
		for f := range uitask {
			done := make(chan struct{})
			gdk_threads_add_idle(&gtkIdleOp{
				what:	f,
				done:	done,
			})
			<-done
			close(done)
		}
	}()

	go func() {
		main()
		uitask <- gtk_main_quit
	}()

	gtk_main()
	return nil
}

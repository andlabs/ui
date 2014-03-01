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

	// thanks to tristan in irc.gimp.net/#gtk
	gdk_threads_add_idle(func() bool {
		select {
		case f := <-uitask:
			f()
		default:		// do not block
		}
		return true	// don't destroy the callback
	})

	go main()

	gtk_main()
	return nil
}

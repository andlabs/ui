// +build !windows,!darwin,!plan9

// 16 february 2014
package main

import (
	"fmt"
	"runtime"
)

var uitask chan func()

func ui(initDone chan error) {
	runtime.LockOSThread()

	uitask = make(chan func())
	if gtk_init() != true {
		initDone <- fmt.Errorf("gtk_init failed (reason unknown; TODO)")
		return
	}
	initDone <- nil

	// thanks to tristan in irc.gimp.net/#gtk
	gdk_threads_add_idle(func() bool {
		select {
		case f := <-uitask:
			f()
		default:		// do not block
		}
		return true	// don't destroy the callback
	})
	gtk_main()
}

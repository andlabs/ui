// +build !windows,!darwin,!plan9

// 16 february 2014

package ui

import (
	"fmt"
)

// #cgo pkg-config: gtk+-3.0
// #include "gtk_unix.h"
// /* this is called when we're done */
// static inline gboolean our_quit_callback(gpointer data)
// {
// 	gtk_main_quit();
// 	return FALSE;		/* remove from idle handler (not like it matters) */
// }
// /* I would call gdk_threads_add_idle() directly from ui() but cgo whines, so; trying to access our_quit_callback() in any way other than a call would cause _cgo_main.c to complain too */
// static inline void signalQuit(void)
// {
// 	gdk_threads_add_idle(our_quit_callback, NULL);
// }
import "C"

func uiinit() error {
	err := gtk_init()
	if err != nil {
		return fmt.Errorf("gtk_init() failed: %v", err)
	}

	return nil
}

func ui() {
	go func() {
		<-Stop
		C.signalQuit()
		// TODO wait for it to return?
	}()

	C.gtk_main()
}

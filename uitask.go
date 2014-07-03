// 11 february 2014

package ui

import (
	"runtime"
)

// Go sets up the UI environment.
// If initialization fails, Go returns an error and Ready is not pulsed.
// Otherwise, Go first runs start(), which should contain code to create the first Window, and then fires up the event loop, not returning to its caller until Stop is pulsed, at which point Go() will return nil.
// After Go() returns, you cannot call future ui functions/methods meaningfully.
// Pulsing Stop will cause Go() to return immediately; the programmer is responsible for cleaning up (for instance, hiding open Windows) beforehand.
//
// It is not safe to call ui.Go() in a goroutine. It must be called directly from main(). This means if your code calls other code-modal servers (such as http.ListenAndServe()), they must be run from goroutines. (This is due to limitations in various OSs, such as Mac OS X.)
//
// Go() does not process the command line for flags (that is, it does not call flag.Parse()), nor does package ui add any of the underlying toolkit's supported command-line flags.
// If you must, and if the toolkit also has environment variable equivalents to these flags (for instance, GTK+), use those instead.
func Go(start func()) error {
	runtime.LockOSThread()
	if err := uiinit(); err != nil {
		return err
	}
	start()
	ui()
	return nil
}

// Post issues a request to the given Window to do something on the main thread.
// Note the name of the function: there is no guarantee that the request will be handled immediately.
// Because this can be safely called from any goroutine, it is a package-level function, and not a method on Window.
// TODO garbage collection
func Post(w *Window, data interface{}) {
	uipost(w, data)
}

// TODO this needs to be replaced with a function
// Stop should be pulsed when you are ready for Go() to return.
// Pulsing Stop will cause Go() to return immediately; the programmer is responsible for cleaning up (for instance, hiding open Windows) beforehand.
// Do not pulse Stop more than once.
var Stop = make(chan struct{})

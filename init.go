// 11 february 2014

package ui

import (
	"runtime"
)

// Go sets up the UI environment and pulses Ready.
// If initialization fails, Go returns an error and Ready is not pulsed.
// Otherwise, Go does not return to its caller until Stop is pulsed, at which point Go() will return nil.
// After Go() returns, you cannot call future ui functions/methods meaningfully.
// Pulsing Stop will cause Go() to return immediately; the programmer is responsible for cleaning up (for instance, hiding open Windows) beforehand.
//
// It is not safe to call ui.Go() in a goroutine. It must be called directly from main(). This means if your code calls other code-modal servers (such as http.ListenAndServe()), they must be run from goroutines. (This is due to limitations in various OSs, such as Mac OS X.)
//
// Go() does not process the command line for flags (that is, it does not call flag.Parse()), nor does package ui add any of the underlying toolkit's supported command-line flags.
// If you must, and if the toolkit also has environment variable equivalents to these flags (for instance, GTK+), use those instead.
func Go() error {
	runtime.LockOSThread()
	if err := uiinit(); err != nil {
		return err
	}
	Ready <- struct{}{}
	close(Ready)
	ui()
	return nil
}

// Ready is pulsed when Go() is ready to begin accepting requests to the safe methods.
// Go() will wait for something to receive on Ready, then Ready will be closed.
var Ready = make(chan struct{})

// Stop should be pulsed when you are ready for Go() to return.
// Pulsing Stop will cause Go() to return immediately; the programmer is responsible for cleaning up (for instance, hiding open Windows) beforehand.
// Do not pulse Stop more than once.
var Stop = make(chan struct{})

// This function is a simple helper functionn that basically pushes the effect of a function call for later. This allows the selected safe Window methods to be safe.
// TODO make sure this acts sanely if called from uitask itself
func touitask(f func()) {
	done := make(chan struct{})
	defer close(done)
	go func() {		// to avoid locking uitask itself
		done2 := make(chan struct{})		// make the chain uitask <- f <- uitask to avoid deadlocks
		defer close(done2)
		uitask <- func() {
			f()
			done2 <- struct{}{}
		}
		done <- <-done2
	}()
	<-done
}

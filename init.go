// 11 february 2014

package ui

// Go sets up the UI environment and runs main in a goroutine.
// If initialization fails, Go returns an error and main is not called.
// Otherwise, Go does not return to its caller until main does, at which point it returns nil.
// After it returns, you cannot call future ui functions/methods meaningfully.
//
// It is not safe to call ui.Go() in a goroutine. It must be called directly from main().
//
// This model is undesirable, but Cocoa limitations require it.
//
// Go does not process the command line for flags (that is, it does not call flag.Parse()), nor does package ui add any of the underlying toolkit's supported command-line flags.
// If you must, and if the toolkit also has environment variable equivalents to these flags (for instance, GTK+), use those instead.
func Go(main func()) error {
	return ui(main)
}

// This function is a simple helper functionn that basically pushes the effect of a function call for later. This allows the selected safe Window methods to be safe.
// It's also currently used by the various dialog box functions on Windows to allow them to return instantly, rather than wait for the dialog box to finish (which both GTK+ and Mac OS X let you do). I consider this a race condition bug. TODO (also TODO document the /intended/ behavior)
func touitask(f func()) {
	go func() {		// to avoid locking uitask itself
		uitask <- f
	}()
}

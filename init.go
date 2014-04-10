// 11 february 2014

package ui

// Go sets up the UI environment and runs main in a goroutine.
// If initialization fails, Go returns an error and main is not called.
// Otherwise, Go does not return to its caller until main does, at which point it returns nil.
// After it returns, you cannot call future ui functions/methods meaningfully.
// (TODO ideally we would want to be able to call ui.MsgBoxError() to report failures to the user, but I would need to figure out how to do this on platforms other than Windows.)
//
// It is not safe to call ui.Go() in a goroutine. It must be called directly from main().
//
// This model is undesirable, but Cocoa limitations require it.
func Go(main func()) error {
	return ui(main)
}

// AppQuit is pulsed when the user decides to quit the program if their operating system provides a facility for quitting an entire application, rather than merely close all windows (for instance, Mac OS X via the Dock icon).
// You should assign one of your Windows's Closing to this variable so the user choosing to quit the application is treated the same as closing that window.
// If you do not respond to this signal, nothing will happen; regardless of whether or not you respond to this signal, the application will not quit.
// Do not merely check this channel alone; it is not guaranteed to be pulsed on all systems or in all conditions.
var AppQuit chan struct{}

func init() {
	// don't expose this in the documentation
	AppQuit = newEvent()
}

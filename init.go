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

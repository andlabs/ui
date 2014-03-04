// 11 february 2014

// ui is a simple library for building cross-platform GUI programs in Go.
//
// It targets Windows, Mac OS X, Linux, and other Unixes, and provides a
// thread-safe, channel-based API.
package ui

// Go sets up the UI environment and runs main in a goroutine.
// If initialization fails, Go returns an error and main is not called.
// Otherwise, Go does not return to its caller until (unless? TODO) the application loop exits, at which point it returns nil.
//
// This model is undesirable, but Cocoa limitations require it.
func Go(main func()) error {
	return ui(main)
}

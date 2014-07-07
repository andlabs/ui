// 7 july 2014

package ui

// Window represents a top-level window on screen that contains other Controls.
// Windows in package ui can only contain one control; the Stack and Grid layout Controls allow you to pack multiple Controls in a Window.
// Note that a Window is not itself a Control.
type Window interface {
	// SetControl creates a Request to the Window's child Control.
	SetControl(c Control) *Request

	// Title and SetTitle create Requests to get and set the Window's title, respectively.
	Title() *Request
	SetTitle(title string) *Request

	// Show and Hide create Requests to bring the Window on-screen and off-screen, respectively.
	Show() *Request
	Hide() *Request

	// Close creates a Request to close the Window.
	// Any Controls within the Window are destroyed, and the Window itself is also destroyed.
	// Attempting to use a Window after it has been closed results in undefined behavior.
	Close() *Request

	// OnClosing registers an event handler that is triggered when the user clicks the Window's close button.
	// On systems where whole applications own windows, OnClosing is also triggered when the user asks to close the application.
	// If this handler returns true, the Window is closed as defined by Close above.
	// If this handler returns false, the Window is not closed.
	OnClosing(func(c Doer) bool)

	// TODO SetSize (TODO remove?)
	// TODO Center
}

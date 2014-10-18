// 7 july 2014

package ui

// Window represents a top-level window on screen that contains other Controls.
// Windows in package ui can only contain one control; the Stack, Grid, and SimpleGrid layout Controls allow you to pack multiple Controls in a Window.
// Note that a Window is not itself a Control.
type Window interface {
	// Title and SetTitle get and set the Window's title, respectively.
	Title() string
	SetTitle(title string)

	// Show and Hide bring the Window on-screen and off-screen, respectively.
	Show()
	Hide()

	// Close closes the Window.
	// Any Controls within the Window are destroyed, and the Window itself is also destroyed.
	// Attempting to use a Window after it has been closed results in undefined behavior.
	// Close unconditionally closes the Window; it neither raises OnClosing nor checks for a return from OnClosing.
	Close()

	// OnClosing registers an event handler that is triggered when the user clicks the Window's close button.
	// On systems where whole applications own windows, OnClosing is also triggered when the user asks to close the application.
	// If this handler returns true, the Window is closed as defined by Close above.
	// If this handler returns false, the Window is not closed.
	OnClosing(func() bool)

	// Margined and SetMargined get and set whether the contents of the Window have a margin around them.
	// The size of the margin is platform-dependent.
	Margined() bool
	SetMargined(margined bool)

	windowDialog
}

// NewWindow creates a new Window with the given title text, size, and control.
func NewWindow(title string, width int, height int, control Control) Window {
	return newWindow(title, width, height, control)
}

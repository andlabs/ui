// 7 february 2014

package ui

import (
	// ...
)

// MsgBox displays an informational message box to the user with just an OK button.
// primaryText should be a short string describing the message, and will be displayed with additional emphasis on platforms that support it.
// Optionally, secondaryText can be used to show additional information.
// If you pass an empty string for secondaryText, neither additional information nor space for additional information will be shown.
// On platforms that allow for the message box window to have a title, os.Args[0] is used.
// 
// If parent is nil, the message box is modal to the entire application: the user cannot interact with any other window until this one is dismissed.
// Whether or not resizing Windows will still be allowed is implementation-defined; if the implementation does allow it, resizes will still work properly.
// Whether or not the message box stays above all other W+indows in the program is also implementation-defined.
// 
// If parent is not nil, the message box is modal to that Window only.
// Attempts to interact with parent will be blocked, but all other Windows in the application can still be used properly.
// The message box will also stay above parent.
// As with parent == nil, resizing is implementation-defined, but will work properly if allowed. [TODO verify]
// If parent has not yet been created, MsgBox() panics. [TODO check what happens if hidden]
func MsgBox(parent *Window, primaryText string, secondaryText string) {
	msgBox(parent, primaryText, secondaryText)
}

// MsgBoxError displays a message box to the user with just an OK button and an icon indicating an error.
// Otherwise, it behaves like MsgBox.
func MsgBoxError(parent *Window, primaryText string, secondaryText string) {
	msgBoxError(parent, primaryText, secondaryText)
}

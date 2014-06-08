// 7 february 2014

package ui

import (
	// ...
)

// sentinel (presently nil; may be a private instance if subtle bugs start showing up in user code)
var dialogWindow *Window

// MsgBox displays an informational message box to the user with just an OK button.
// primaryText should be a short string describing the message, and will be displayed with additional emphasis on platforms that support it.
// Optionally, secondaryText can be used to show additional information.
// If you pass an empty string for secondaryText, neither additional information nor space for additional information will be shown.
// On platforms that allow for the message box window to have a title, os.Args[0] is used.
// 
// The message box is modal to the entire application: the user cannot interact with any other window until this one is dismissed.
// Whether or not resizing Windows will still be allowed is implementation-defined; if the implementation does allow it, resizes will still work properly.
// Whether or not the message box stays above all other W+indows in the program is also implementation-defined.
func MsgBox(primaryText string, secondaryText string) {
	<-dialogWindow.msgBox(primaryText, secondaryText)
}

// MsgBox behaves like the package-scope MsgBox function, except the message box is modal to w only.
// Attempts to interact with w will be blocked, but all other Windows in the application can still be used properly.
// The message box will also stay above w.
// Whether w can be resized while the message box is displayed is implementation-defined, but will work properly if allowed.
// If w has not yet been created, MsgBox() panics.
// If w has not been shown yet or is currently hidden, what MsgBox does is implementation-defined.
// 
// On return, done will be a channel that is pulsed when the message box is dismissed.
func (w *Window) MsgBox(primaryText string, secondaryText string) (done chan struct{}) {
	if !w.created {
		panic("parent window passed to Window.MsgBox() before it was created")
	}
	return w.msgBox(primaryText, secondaryText)
}

// MsgBoxError displays a message box to the user with just an OK button and an icon indicating an error.
// Otherwise, it behaves like MsgBox.
func MsgBoxError(primaryText string, secondaryText string) {
	<-dialogWindow.msgBoxError(primaryText, secondaryText)
}

// MsgBoxError displays a message box to the user with just an OK button and an icon indicating an error.
// Otherwise, it behaves like Window.MsgBox.
func (w *Window) MsgBoxError(primaryText string, secondaryText string) (done chan struct{}) {
	if !w.created {
		panic("parent window passed to MsgBoxError() before it was created")
	}
	return w.msgBoxError(primaryText, secondaryText)
}

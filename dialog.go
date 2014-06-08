// 7 february 2014

package ui

import (
	// ...
)

// sentinel (not nil so programmer errors don't go undetected)
// this window is invalid and cannot be used directly
var dialogWindow = new(Window)

// MsgBox displays an informational message box to the user with just an OK button.
// primaryText should be a short string describing the message, and will be displayed with additional emphasis on platforms that support it.
// Optionally, secondaryText can be used to show additional information.
// If you pass an empty string for secondaryText, neither additional information nor space for additional information will be shown.
// On platforms that allow for the message box window to have a title, os.Args[0] is used.
// 
// See "On Dialogs" in the package overview for behavioral information.
func MsgBox(primaryText string, secondaryText string) {
	<-dialogWindow.msgBox(primaryText, secondaryText)
}

// MsgBox is the Window method version of the package-scope function MsgBox.
// See that function's documentation and "On Dialogs" in the package overview for more information.
func (w *Window) MsgBox(primaryText string, secondaryText string) (done chan struct{}) {
	if !w.created {
		panic("parent window passed to Window.MsgBox() before it was created")
	}
	return w.msgBox(primaryText, secondaryText)
}

// MsgBoxError displays a message box to the user with just an OK button and an icon indicating an error.
// Otherwise, it behaves like MsgBox.
// 
// See "On Dialogs" in the package overview for more information.
func MsgBoxError(primaryText string, secondaryText string) {
	<-dialogWindow.msgBoxError(primaryText, secondaryText)
}

// MsgBoxError is the Window method version of the package-scope function MsgBoxError.
// See that function's documentation and "On Dialogs" in the package overview for more information.
func (w *Window) MsgBoxError(primaryText string, secondaryText string) (done chan struct{}) {
	if !w.created {
		panic("parent window passed to MsgBoxError() before it was created")
	}
	return w.msgBoxError(primaryText, secondaryText)
}

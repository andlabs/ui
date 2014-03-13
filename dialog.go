// 7 february 2014

package ui

import (
	// ...
)

// MsgBox displays an informational message box to the user with just an OK button.
// primaryText should be a short string describing the message, and will be displayed with additional emphasis on platforms that support it.
// secondaryText can be used to provide more information.
// On platforms that allow for the message box window to have a title, os.Args[0] is used.
func MsgBox(primaryText string, secondaryText string) {
	msgBox(primaryText, secondaryText)
}

// MsgBoxError displays a message box to the user with just an OK button and an icon indicating an error.
// Otherwise, it behaves like MsgBox.
func MsgBoxError(primaryText string, secondaryText string) {
	msgBoxError(primaryText, secondaryText)
}

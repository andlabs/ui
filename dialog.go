// 7 february 2014
package main

import (
	"fmt"
)


// MsgBox displays an informational message box to the user with just an OK button.
func MsgBox(title string, textfmt string, args ...interface{}) {
	msgBox(title, fmt.Sprintf(textfmt, args...))
}

// MsgBoxError displays a message box to the user with just an OK button and an icon indicating an error.
func MsgBoxError(title string, textfmt string, args ...interface{}) {
	msgBoxError(title, fmt.Sprintf(textfmt, args...))
}

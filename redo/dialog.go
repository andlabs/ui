// 18 august 2014

package ui

type windowDialog interface {
	openFile(f func(filename string))
}

// OpenFile opens a dialog box that asks the user to choose a file.
// The dialog box is modal to win, which mut not be nil.
// Some time after the dialog box is closed, OpenFile runs f on the main thread, passing filename.
// filename is the selected filename, or an empty string if no file was chosen.
// OpenFile does not ensure that f remains alive; the programmer is responsible for this.
// If possible on a given system, OpenFile() will not dereference links; it will return the link file itself.
// Hidden files will not be hidden by OpenFile().
func OpenFile(win Window, f func(filename string)) {
	if win == nil {
		panic("Window passed to OpenFile() cannot be nil")
	}
	win.openFile(f)
}

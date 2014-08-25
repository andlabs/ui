// 18 august 2014

package ui

// OpenFile opens a dialog box that asks the user to choose a file.
// It returns the selected filename, or an empty string if no file was chosen.
// All events stop while OpenFile is executing. (TODO move to doc.go)
// If possible on a given system, OpenFile() will not dereference links; it will return the link file itself.
func OpenFile() (filename string) {
	return openFile()
}

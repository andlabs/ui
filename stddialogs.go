// 20 december 2015

package ui

// #include "ui.h"
import "C"

// MsgBoxError opens a modal dialog box with graphical error hinting and returns when the
// user acknowledges the message.
func MsgBoxError(w *Window, title string, description string) {
	ctitle := C.CString(title)
	defer freestr(ctitle)
	cdescription := C.CString(description)
	defer freestr(cdescription)
	C.uiMsgBoxError(w.w, ctitle, cdescription)
}

// OpenFile opens a modal allowing the user to select a path to an existing file.
func OpenFile(w *Window) string {
	cname := C.uiOpenFile(w.w)
	if cname == nil {
		return ""
	}
	defer C.uiFreeText(cname)
	return C.GoString(cname)
}


// SaveFile opens a modal allowing the user to select a path to a new or existing file.
func SaveFile(w *Window) string {
	cname := C.uiSaveFile(w.w)
	if cname == nil {
		return ""
	}
	defer C.uiFreeText(cname)
	return C.GoString(cname)
}
// MsgBox opens a generic modal dialog box and returns when the user acknowledges the
// message.
func MsgBox(w *Window, title string, description string) {
	ctitle := C.CString(title)
	defer freestr(ctitle)
	cdescription := C.CString(description)
	defer freestr(cdescription)
	C.uiMsgBox(w.w, ctitle, cdescription)
}

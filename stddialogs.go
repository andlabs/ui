// 20 december 2015

package ui

// #include "ui.h"
import "C"

// TODO OpenFile, SaveFile, MsgBox

// TODO
func MsgBoxError(w *Window, title string, description string) {
	ctitle := C.CString(title)
	cdescription := C.CString(description)
	C.uiMsgBoxError(w.w, ctitle, cdescription)
	freestr(ctitle)
	freestr(cdescription)
}

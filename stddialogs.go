// 20 december 2015

package ui

// #include "ui.h"
import "C"

// TODO
func MsgBoxError(w *Window, title string, description string) {
	ctitle := C.CString(title)
	cdescription := C.CString(description)
	C.uiMsgBoxError(w.w, ctitle, cdescription)
	freestr(ctitle)
	freestr(cdescription)
}

func OpenFile(w *Window) string {
	cname := C.uiOpenFile(w.w)
	name := C.GoString(cname)
	freestr(cname)
	return name
}

func SaveFile(w *Window) string {
	cname := C.uiSaveFile(w.w)
	name := C.GoString(cname)
	freestr(cname)
	return name
}

func MsgBox(w *Window, title string, description string) {
	ctitle := C.CString(title)
	cdescription := C.CString(description)
	C.uiMsgBox(w.w, ctitle, cdescription)
	freestr(ctitle)
	freestr(cdescription)
}

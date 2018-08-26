// 20 december 2015

package ui

// #include "pkgui.h"
import "C"

// TODO
func MsgBoxError(w *Window, title string, description string) {
	ctitle := C.CString(title)
	defer freestr(ctitle)
	cdescription := C.CString(description)
	defer freestr(cdescription)
	C.uiMsgBoxError(w.w, ctitle, cdescription)
}

func OpenFile(w *Window) string {
	cname := C.uiOpenFile(w.w)
	if cname == nil {
		return ""
	}
	defer C.uiFreeText(cname)
	return C.GoString(cname)
}

func SaveFile(w *Window) string {
	cname := C.uiSaveFile(w.w)
	if cname == nil {
		return ""
	}
	defer C.uiFreeText(cname)
	return C.GoString(cname)
}

func MsgBox(w *Window, title string, description string) {
	ctitle := C.CString(title)
	defer freestr(ctitle)
	cdescription := C.CString(description)
	defer freestr(cdescription)
	C.uiMsgBox(w.w, ctitle, cdescription)
}

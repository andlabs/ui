// +build !windows,!darwin,!plan9

// 17 february 2014

package ui

import (
	"unsafe"
)

// this file contains functions that wrap around complex pointer casts to satisfy GTK+'s dumb type aliasing system
// fromxxx() converts from GTK+ type to Go type
// toxxxx() converts from Go type to GTK+ type
// Listbox casts are stored in listbox_unix.go

// #include "gtk_unix.h"
import "C"

func fromgbool(b C.gboolean) bool {
	return b != C.FALSE
}

func togbool(b bool) C.gboolean {
	if b {
		return C.TRUE
	}
	return C.FALSE
}

func fromgstr(what *C.gchar) string {
	cstr := (*C.char)(unsafe.Pointer(what))
	return C.GoString(cstr)
}

func togstr(what *C.char) *C.gchar {
	return (*C.gchar)(unsafe.Pointer(what))
}

func fromgtkwindow(x *C.GtkWindow) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtkwindow(what *C.GtkWidget) *C.GtkWindow {
	return (*C.GtkWindow)(unsafe.Pointer(what))
}

func fromgtkcontainer(x *C.GtkContainer) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtkcontainer(what *C.GtkWidget) *C.GtkContainer {
	return (*C.GtkContainer)(unsafe.Pointer(what))
}

func fromgtklayout(x *C.GtkLayout) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtklayout(what *C.GtkWidget) *C.GtkLayout {
	return (*C.GtkLayout)(unsafe.Pointer(what))
}

func fromgtkbutton(x *C.GtkButton) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtkbutton(what *C.GtkWidget) *C.GtkButton {
	return (*C.GtkButton)(unsafe.Pointer(what))
}

func fromgtktogglebutton(x *C.GtkToggleButton) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtktogglebutton(what *C.GtkWidget) *C.GtkToggleButton {
	return (*C.GtkToggleButton)(unsafe.Pointer(what))
}

func fromgtkcombobox(x *C.GtkComboBoxText) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtkcombobox(what *C.GtkWidget) *C.GtkComboBoxText {
	return (*C.GtkComboBoxText)(unsafe.Pointer(what))
}

func fromgtkentry(x *C.GtkEntry) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtkentry(what *C.GtkWidget) *C.GtkEntry {
	return (*C.GtkEntry)(unsafe.Pointer(what))
}

func fromgtklabel(x *C.GtkLabel) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtklabel(what *C.GtkWidget) *C.GtkLabel {
	return (*C.GtkLabel)(unsafe.Pointer(what))
}

func fromgtkprogressbar(x *C.GtkProgressBar) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(x))
}

func togtkprogressbar(what *C.GtkWidget) *C.GtkProgressBar {
	return (*C.GtkProgressBar)(unsafe.Pointer(what))
}

// +build !windows,!darwin,!plan9

// 17 february 2014

//
package ui

import (
	"unsafe"
)

// this file contains functions that wrap around complex pointer casts to satisfy GTK+'s dumb type aliasing system
// fromxxx() converts from GTK+ type to Go type
// toxxxx() converts from Go type to GTK+ type
// Listbox casts are stored in listbox_unix.go

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
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

func fromgtkwidget(x *C.GtkWidget) *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(x))
}

func togtkwidget(what *gtkWidget) *C.GtkWidget {
	return (*C.GtkWidget)(unsafe.Pointer(what))
}

// TODO wrap in C.GoString()?
func fromgchar(what *C.gchar) *C.char {
	return (*C.char)(unsafe.Pointer(what))
}

func togchar(what *C.char) *C.gchar {
	return (*C.gchar)(unsafe.Pointer(what))
}

func fromgtkwindow(x *C.GtkWindow) *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(x))
}

func togtkwindow(what *gtkWidget) *C.GtkWindow {
	return (*C.GtkWindow)(unsafe.Pointer(what))
}

func fromgtkcontainer(x *C.GtkContainer) *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(x))
}

func togtkcontainer(what *gtkWidget) *C.GtkContainer {
	return (*C.GtkContainer)(unsafe.Pointer(what))
}

func fromgtkfixed(x *C.GtkFixed) *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(x))
}

func togtkfixed(what *gtkWidget) *C.GtkFixed {
	return (*C.GtkFixed)(unsafe.Pointer(what))
}

func fromgtkbutton(x *C.GtkButton) *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(x))
}

func togtkbutton(what *gtkWidget) *C.GtkButton {
	return (*C.GtkButton)(unsafe.Pointer(what))
}

func fromgtktogglebutton(x *C.GtkToggleButton) *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(x))
}

func togtktogglebutton(what *gtkWidget) *C.GtkToggleButton {
	return (*C.GtkToggleButton)(unsafe.Pointer(what))
}

func fromgtkcombobox(x *C.GtkComboBoxText) *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(x))
}

func togtkcombobox(what *gtkWidget) *C.GtkComboBoxText {
	return (*C.GtkComboBoxText)(unsafe.Pointer(what))
}

func fromgtkentry(x *C.GtkEntry) *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(x))
}

func togtkentry(what *gtkWidget) *C.GtkEntry {
	return (*C.GtkEntry)(unsafe.Pointer(what))
}

func fromgtklabel(x *C.GtkLabel) *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(x))
}

func togtklabel(what *gtkWidget) *C.GtkLabel {
	return (*C.GtkLabel)(unsafe.Pointer(what))
}

func fromgtkprogressbar(x *C.GtkProgressBar) *gtkWidget {
	return (*gtkWidget)(unsafe.Pointer(x))
}

func togtkprogressbar(what *gtkWidget) *C.GtkProgressBar {
	return (*C.GtkProgressBar)(unsafe.Pointer(what))
}

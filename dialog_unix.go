// +build !windows,!darwin,!plan9

// 7 february 2014

//
package ui

import (
	"unsafe"
)

// #cgo pkg-config: gtk+-3.0
// #include <stdlib.h>
// #include <gtk/gtk.h>
// /* because cgo seems to choke on ... */
// /* TODO does NULL parent make the box application-global? docs are unclear */
// /* TODO does secondary text appear in the titlebar or above the message? if the latter, will gtk_window_set_title() work? */
// GtkWidget *gtkNewMsgBox(GtkMessageType type, GtkButtonsType buttons, char *title, char *text) { GtkWidget *k = gtk_message_dialog_new(NULL, GTK_DIALOG_MODAL, type, buttons, "%s", (gchar *) title); gtk_message_dialog_format_secondary_text((GtkMessageDialog *) k, "%s", (gchar *) text); return k; }
import "C"

func _msgBox(text string, title string, msgtype C.GtkMessageType, buttons C.GtkButtonsType) (result C.gint) {
	ret := make(chan C.gint)
	defer close(ret)
	uitask <- func() {
		ctitle := C.CString(title)
		defer C.free(unsafe.Pointer(ctitle))
		ctext := C.CString(text)
		defer C.free(unsafe.Pointer(ctext))
		box := C.gtkNewMsgBox(msgtype, buttons, ctitle, ctext)
		response := C.gtk_dialog_run((*C.GtkDialog)(unsafe.Pointer(box)))
		C.gtk_widget_destroy(box)
		ret <- response
	}
	return <-ret
}

func msgBox(title string, text string) {
	// TODO add an icon?
	_msgBox(text, title, C.GtkMessageType(C.GTK_MESSAGE_OTHER), C.GtkButtonsType(C.GTK_BUTTONS_OK))
}

func msgBoxError(title string, text string) {
	_msgBox(text, title, C.GtkMessageType(C.GTK_MESSAGE_ERROR), C.GtkButtonsType(C.GTK_BUTTONS_OK))
}

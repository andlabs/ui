// +build !windows,!darwin

// 19 august 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// #include "modalqueue.h"
// /* because cgo doesn't like ... */
// GtkWidget *newOpenFileDialog(void)
// {
// 	return gtk_file_chooser_dialog_new(NULL,	/* default title */
// 		NULL,		/* no parent window */
// 		GTK_FILE_CHOOSER_ACTION_OPEN,
// 		GTK_STOCK_CANCEL, GTK_RESPONSE_CANCEL,
// 		GTK_STOCK_OPEN, GTK_RESPONSE_ACCEPT,
// 		NULL);
// }
import "C"

func openFile() string {
	widget := C.newOpenFileDialog()
	defer C.gtk_widget_destroy(widget)
	dialog  := (*C.GtkDialog)(unsafe.Pointer(widget))
	fc := (*C.GtkFileChooser)(unsafe.Pointer(widget))
	C.gtk_file_chooser_set_local_only(fc, C.FALSE)
	C.gtk_file_chooser_set_select_multiple(fc, C.FALSE)
	C.gtk_file_chooser_set_show_hidden(fc, C.TRUE)
	C.beginModal()
	response := C.gtk_dialog_run(dialog)
	C.endModal()
	if response != C.GTK_RESPONSE_ACCEPT {
		return ""
	}
	filename := C.gtk_file_chooser_get_filename(fc)
	if filename == nil {
		panic("[DEBUG TODO] chosen filename NULL")
	}
	defer C.g_free(C.gpointer(unsafe.Pointer(filename)))
	return fromgstr(filename)
}

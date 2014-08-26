// +build !windows,!darwin

// 19 august 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// #include "modalqueue.h"
// extern void our_openfile_response_callback(GtkDialog *, gint, gpointer);
// /* because cgo doesn't like ... */
// static inline GtkWidget *newOpenFileDialog(GtkWindow *parent)
// {
// 	return gtk_file_chooser_dialog_new(NULL,	/* default title */
// 		parent,
// 		GTK_FILE_CHOOSER_ACTION_OPEN,
// 		GTK_STOCK_CANCEL, GTK_RESPONSE_CANCEL,
// 		GTK_STOCK_OPEN, GTK_RESPONSE_ACCEPT,
// 		NULL);
// }
import "C"

func (w *window) openFile(f func(filename string)) {
	widget := C.newOpenFileDialog(w.window)
	window := (*C.GtkWindow)(unsafe.Pointer(widget))
	dialog  := (*C.GtkDialog)(unsafe.Pointer(widget))
	fc := (*C.GtkFileChooser)(unsafe.Pointer(widget))
	C.gtk_file_chooser_set_local_only(fc, C.FALSE)
	C.gtk_file_chooser_set_select_multiple(fc, C.FALSE)
	C.gtk_file_chooser_set_show_hidden(fc, C.TRUE)
	C.gtk_window_set_modal(window, C.TRUE)
	g_signal_connect(
		C.gpointer(unsafe.Pointer(dialog)),
		"response",
		C.GCallback(C.our_openfile_response_callback),
		C.gpointer(unsafe.Pointer(&f)))
	C.gtk_widget_show_all(widget)
}

//export our_openfile_response_callback
func our_openfile_response_callback(dialog *C.GtkDialog, response C.gint, data C.gpointer) {
	f := (*func(string))(unsafe.Pointer(data))
	if response != C.GTK_RESPONSE_ACCEPT {
		(*f)("")
		C.gtk_widget_destroy((*C.GtkWidget)(unsafe.Pointer(dialog)))
		return
	}
	filename := C.gtk_file_chooser_get_filename((*C.GtkFileChooser)(unsafe.Pointer(dialog)))
	if filename == nil {
		panic("[DEBUG TODO] chosen filename NULL")
	}
	realfilename := fromgstr(filename)
	C.g_free(C.gpointer(unsafe.Pointer(filename)))
	C.gtk_widget_destroy((*C.GtkWidget)(unsafe.Pointer(dialog)))
	(*f)(realfilename)
}

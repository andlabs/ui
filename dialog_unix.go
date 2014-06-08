// +build !windows,!darwin,!plan9

// 7 february 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// extern void our_dialog_response_callback(GtkDialog *, gint, gpointer);
// /* because cgo seems to choke on ... */
// /* parent will be NULL if there is no parent, so this is fine */
// static inline GtkWidget *gtkNewMsgBox(GtkWindow *parent, GtkMessageType type, GtkButtonsType buttons, char *title, char *text)
// {
// 	GtkWidget *k;
// 
// 	k = gtk_message_dialog_new(parent, GTK_DIALOG_MODAL, type, buttons, "%s", (gchar *) title);
// 	if (text != NULL)
// 		gtk_message_dialog_format_secondary_text((GtkMessageDialog *) k, "%s", (gchar *) text);
// 	return k;
// }
import "C"

// dialog performs the bookkeeping involved for having a GtkDialog behave the way we want.
type dialog struct {
	parent		*Window
	pwin			*C.GtkWindow
	hadgroup		C.gboolean
	prevgroup	*C.GtkWindowGroup
	newgroup		*C.GtkWindowGroup
	result		chan int
}

func mkdialog(parent *Window) *dialog {
	return &dialog{
		parent:	parent,
		result:	make(chan int),
	}
}

func (d *dialog) prepare() {
	// to implement parent, we need to put the GtkMessageDialog into a new window group along with parent
	// a GtkWindow can only be part of one group
	// so we use this to save the parent window group (if there is one) and store the new window group
	// after showing the message box, we restore the previous window group, so future parent == dialogWindow can work properly
	// thanks to pbor and mclasen in irc.gimp.net/#gtk+
	if d.parent != dialogWindow {
		d.pwin = togtkwindow(d.parent.sysData.widget)
		d.hadgroup = C.gtk_window_has_group(d.pwin)
		// we can't remove a window from the "default window group"; otherwise this throws up Gtk-CRITICAL warnings
		if d.hadgroup != C.FALSE {
			d.prevgroup = C.gtk_window_get_group(d.pwin)
			C.gtk_window_group_remove_window(d.prevgroup, d.pwin)
		}
		d.newgroup = C.gtk_window_group_new()
		C.gtk_window_group_add_window(d.newgroup, d.pwin)
	}
}

func (d *dialog) run(mk func() *C.GtkWidget) {
	d.prepare()
	box := mk()
	if d.parent == dialogWindow {
		go func() {
			res := make(chan C.gint)
			defer close(res)
			uitask <- func() {
				r := C.gtk_dialog_run((*C.GtkDialog)(unsafe.Pointer(box)))
				d.cleanup(box)
				res <- r
			}
			d.send(<-res)
		}()
		return
	}
	// otherwise just connect the delete signal
	g_signal_connect_pointer(box, "response", dialog_response_callback, unsafe.Pointer(d))
	C.gtk_widget_show_all(box)
}

//export our_dialog_response_callback
func our_dialog_response_callback(box *C.GtkDialog, res C.gint, data C.gpointer) {
	d := (*dialog)(unsafe.Pointer(data))
	d.cleanup((*C.GtkWidget)(unsafe.Pointer(box)))
	go d.send(res)		// send on another goroutine, like everything else
}

var dialog_response_callback = C.GCallback(C.our_dialog_response_callback)

func (d *dialog) cleanup(box *C.GtkWidget) {
	// have to explicitly close the dialog box, otherwise wacky things will happen
	C.gtk_widget_destroy(box)
	if d.parent != dialogWindow {
		C.gtk_window_group_remove_window(d.newgroup, d.pwin)
		C.g_object_unref(C.gpointer(unsafe.Pointer(d.newgroup)))		// free the group
		if d.prevgroup != nil {
			C.gtk_window_group_add_window(d.prevgroup, d.pwin)
		}		// otherwise it'll go back into the default group on its own
	}
}

func (d *dialog) send(res C.gint) {
	// this is where processing would go
	d.result <- int(res)
}

func _msgBox(parent *Window, primarytext string, secondarytext string, msgtype C.GtkMessageType, buttons C.GtkButtonsType) (result chan int) {
	result = make(chan int)
	d := mkdialog(parent)
	uitask <- func() {
		cprimarytext := C.CString(primarytext)
		defer C.free(unsafe.Pointer(cprimarytext))
		csecondarytext := (*C.char)(nil)
		if secondarytext != "" {
			csecondarytext = C.CString(secondarytext)
			defer C.free(unsafe.Pointer(csecondarytext))
		}
		d.run(func() *C.GtkWidget {
			return C.gtkNewMsgBox(d.pwin, msgtype, buttons, cprimarytext, csecondarytext)
		})
	}
	return d.result
}

func (w *Window) msgBox(primarytext string, secondarytext string) (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		<-_msgBox(w, primarytext, secondarytext, C.GtkMessageType(C.GTK_MESSAGE_OTHER), C.GtkButtonsType(C.GTK_BUTTONS_OK))
		done <- struct{}{}
	}()
	return done
}

func (w *Window) msgBoxError(primarytext string, secondarytext string) (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		<-_msgBox(w, primarytext, secondarytext, C.GtkMessageType(C.GTK_MESSAGE_ERROR), C.GtkButtonsType(C.GTK_BUTTONS_OK))
		done <- struct{}{}
	}()
	return done
}

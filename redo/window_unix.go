// +build !windows,!darwin

// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// extern gboolean windowClosing(GtkWidget *, GdkEvent *, gpointer);
// extern void windowResizing(GtkWidget *, GdkRectangle *, gpointer);
import "C"

type window struct {
	widget	*C.GtkWidget
	wc		*C.GtkContainer
	bin		*C.GtkBin
	window	*C.GtkWindow

	closing	*event

	*layout
}

func newWindow(title string, width int, height int, control Control) *window {
	widget := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	ctitle := togstr(title)
	defer freegstr(ctitle)
	w := &window{
		widget:		widget,
		wc:			(*C.GtkContainer)(unsafe.Pointer(widget)),
		bin:			(*C.GtkBin)(unsafe.Pointer(widget)),
		window:		(*C.GtkWindow)(unsafe.Pointer(widget)),
		closing:		newEvent(),
	}
	C.gtk_window_set_title(w.window, ctitle)
	g_signal_connect(
		C.gpointer(unsafe.Pointer(w.window)),
		"delete-event",
		C.GCallback(C.windowClosing),
		C.gpointer(unsafe.Pointer(w)))
	C.gtk_window_resize(w.window, C.gint(width), C.gint(height))
	w.layout = newLayout(control)
	C.gtk_container_add(w.wc, w.layout.layoutwidget)
	return w
}

func (w *window) Title() string {
	return fromgstr(C.gtk_window_get_title(w.window))
}

func (w *window) SetTitle(title string) {
	ctitle := togstr(title)
	defer freegstr(ctitle)
	C.gtk_window_set_title(w.window, ctitle)
}

func (w *window) Show() {
	C.gtk_widget_show_all(w.widget)
}

func (w *window) Hide() {
	C.gtk_widget_hide(w.widget)
}

func (w *window) Close() {
	C.gtk_widget_destroy(w.widget)
}

func (w *window) OnClosing(e func() bool) {
	w.closing.setbool(e)
}

//export windowClosing
func windowClosing(wid *C.GtkWidget, e *C.GdkEvent, data C.gpointer) C.gboolean {
	w := (*window)(unsafe.Pointer(data))
	close := w.closing.fire()
	if close {
		return C.GDK_EVENT_PROPAGATE		// will do gtk_widget_destroy(), which is what we want (thanks ebassi in irc.gimp.net/#gtk+)
	}
	return C.GDK_EVENT_STOP				// keeps window alive
}

//export windowResizing
func windowResizing(wid *C.GtkWidget, r *C.GdkRectangle, data C.gpointer) {
	w := (*window)(unsafe.Pointer(data))
	// the origin of the window's content area is always (0, 0)
	w.resize(0, 0, int(r.width), int(r.height))
}

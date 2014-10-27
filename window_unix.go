// +build !windows,!darwin

// 7 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// extern gboolean windowClosing(GtkWidget *, GdkEvent *, gpointer);
import "C"

type window struct {
	widget *C.GtkWidget
	wc     *C.GtkContainer
	bin    *C.GtkBin
	window *C.GtkWindow

	group *C.GtkWindowGroup

	closing *event

	child			Control
	container		*container
}

func newWindow(title string, width int, height int, control Control) *window {
	widget := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	ctitle := togstr(title)
	defer freegstr(ctitle)
	w := &window{
		widget:  widget,
		wc:      (*C.GtkContainer)(unsafe.Pointer(widget)),
		bin:     (*C.GtkBin)(unsafe.Pointer(widget)),
		window:  (*C.GtkWindow)(unsafe.Pointer(widget)),
		closing: newEvent(),
		child:	control,
	}
	C.gtk_window_set_title(w.window, ctitle)
	g_signal_connect(
		C.gpointer(unsafe.Pointer(w.window)),
		"delete-event",
		C.GCallback(C.windowClosing),
		C.gpointer(unsafe.Pointer(w)))
	C.gtk_window_resize(w.window, C.gint(width), C.gint(height))
	w.container = newContainer()
	w.child.setParent(w.container.parent())
	w.container.resize = w.child.resize
	C.gtk_container_add(w.wc, w.container.widget)
	// for dialogs; otherwise, they will be modal to all windows, not just this one
	w.group = C.gtk_window_group_new()
	C.gtk_window_group_add_window(w.group, w.window)
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

func (w *window) Margined() bool {
	return w.container.margined
}

func (w *window) SetMargined(margined bool) {
	w.container.margined = margined
}

//export windowClosing
func windowClosing(wid *C.GtkWidget, e *C.GdkEvent, data C.gpointer) C.gboolean {
	w := (*window)(unsafe.Pointer(data))
	close := w.closing.fire()
	if close {
		return C.GDK_EVENT_PROPAGATE // will do gtk_widget_destroy(), which is what we want (thanks ebassi in irc.gimp.net/#gtk+)
	}
	return C.GDK_EVENT_STOP // keeps window alive
}

// no need for windowResized; the child container takes care of that

// +build !windows,!darwin

// 7 july 2014

package ui

import (
	"unsafe"
"fmt"
)

// #include "gtk_unix.h"
// extern gboolean windowClosing(GtkWidget *, GdkEvent *, gpointer);
// extern void windowResizing(GtkWidget *, GdkRectangle *, gpointer);
import "C"

type window struct {
	widget	*C.GtkWidget
	container	*C.GtkContainer
	bin		*C.GtkBin
	window	*C.GtkWindow

	layoutw	*C.GtkWidget
	layoutc	*C.GtkContainer
	layout	*C.GtkLayout

	child		Control

	closing	*event

	spaced	bool
}

func newWindow(title string, width int, height int) *window {
	widget := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
	ctitle := togstr(title)
	defer freegstr(ctitle)
	layoutw := C.gtk_layout_new(nil, nil)
	w := &window{
		widget:		widget,
		container:		(*C.GtkContainer)(unsafe.Pointer(widget)),
		bin:			(*C.GtkBin)(unsafe.Pointer(widget)),
		window:		(*C.GtkWindow)(unsafe.Pointer(widget)),
		layoutw:		layoutw,
		layoutc:		(*C.GtkContainer)(unsafe.Pointer(layoutw)),
		layout:		(*C.GtkLayout)(unsafe.Pointer(layoutw)),
		closing:		newEvent(),
	}
	C.gtk_window_set_title(w.window, ctitle)
	g_signal_connect(
		C.gpointer(unsafe.Pointer(w.window)),
		"delete-event",
		C.GCallback(C.windowClosing),
		C.gpointer(unsafe.Pointer(w)))
	// we connect to the layout's size-allocate, not to the window's configure-event
	// this allows us to handle client-side decoration-based configurations (such as GTK+ on Wayland) properly
	// also see commitResize() in sizing_unix.go for additional notes
	// thanks to many people in irc.gimp.net/#gtk+ for help (including tristan for suggesting g_signal_connect_after())
	g_signal_connect_after(
		C.gpointer(unsafe.Pointer(layoutw)),
		"size-allocate",
		C.GCallback(C.windowResizing),
		C.gpointer(unsafe.Pointer(w)))
	C.gtk_window_resize(w.window, C.gint(width), C.gint(height))
	C.gtk_container_add(w.container, layoutw)
	return w
}

func (w *window) SetControl(control Control) {
	if w.child != nil {		// unparent existing control
		w.child.unparent()
	}
	control.unparent()
	control.parent(w)
	w.child = control
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
	w.doresize(int(r.width), int(r.height))
	fmt.Printf("new size %d x %d\n", r.width, r.height)
}

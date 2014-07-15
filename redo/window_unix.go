// +build !darwin

// 7 july 2014

package ui

import (
	"unsafe"
"fmt"
)

// #include "gtk_unix.h"
// extern gboolean windowClosing(GtkWidget *, GdkEvent *, gpointer);
// extern gboolean windowResizing(GtkWidget *, GdkEvent *, gpointer);
import "C"

type window struct {
	widget	*C.GtkWidget
	container	*C.GtkContainer
	bin		*C.GtkBin
	window	*C.GtkWindow

	layoutc	*C.GtkContainer
	layout	*C.GtkLayout

	closing	*event
}

func newWindow(title string, width int, height int) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			widget := C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
			ctitle := togstr(title)
			defer freegstr(ctitle)
			layoutw := C.gtk_layout_new(nil, nil)
			w := &window{
				widget:		widget,
				container:		(*C.GtkContainer)(unsafe.Pointer(widget)),
				bin:			(*C.GtkBin)(unsafe.Pointer(widget)),
				window:		(*C.GtkWindow)(unsafe.Pointer(widget)),
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
			g_signal_connect(
				C.gpointer(unsafe.Pointer(w.window)),
				"configure-event",
				C.GCallback(C.windowResizing),
				C.gpointer(unsafe.Pointer(w)))
			// TODO size
			C.gtk_container_add(w.container, layoutw)
			c <- w
		},
		resp:		c,
	}
}

func (w *window) SetControl(control Control) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			// TODO unparent
			// TODO reparent
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Title() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			c <- fromgstr(C.gtk_window_get_title(w.window))
		},
		resp:		c,
	}
}

func (w *window) SetTitle(title string) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			ctitle := togstr(title)
			defer freegstr(ctitle)
			C.gtk_window_set_title(w.window, ctitle)
			c <- struct{}{}
		},
		resp:		c,
	}
}


func (w *window) Show() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.gtk_widget_show_all(w.widget)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Hide() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.gtk_widget_hide(w.widget)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Close() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			C.gtk_widget_destroy(w.widget)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) OnClosing(e func(c Doer) bool) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			w.closing.setbool(e)
			c <- struct{}{}
		},
		resp:		c,
	}
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
func windowResizing(wid *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	w := (*window)(unsafe.Pointer(data))
	e := (*C.GdkEventConfigure)(unsafe.Pointer(event))
	_ = w // TODO
	// TODO this does not take CSD into account; my attempts at doing so so far have failed to work correctly in the face of rapid live resizing
	// TODO triggered twice on each resize or maximize for some reason???
	fmt.Printf("new size %d x %d\n", e.width, e.height)
	return C.GDK_EVENT_PROPAGATE		// let's be safe
}

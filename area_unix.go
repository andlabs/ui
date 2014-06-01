// +build !windows,!darwin,!plan9

// 14 march 2014

package ui

import (
	"fmt"
	"unsafe"
	"image"
)

// #include "gtk_unix.h"
// extern gboolean our_area_draw_callback(GtkWidget *, cairo_t *, gpointer);
// extern gboolean our_area_button_press_event_callback(GtkWidget *, GdkEvent *, gpointer);
// extern gboolean our_area_button_release_event_callback(GtkWidget *, GdkEvent *, gpointer);
// extern gboolean our_area_motion_notify_event_callback(GtkWidget *, GdkEvent *, gpointer);
// extern gboolean our_area_enterleave_notify_event_callback(GtkWidget *, GdkEvent *, gpointer);
// extern gboolean our_area_key_press_event_callback(GtkWidget *, GdkEvent *, gpointer);
// extern gboolean our_area_key_release_event_callback(GtkWidget *, GdkEvent *, gpointer);
// /* because cgo doesn't like ... */
// static inline void gtkGetDoubleClickSettings(GtkSettings *settings, gint *maxTime, gint *maxDistance)
// {
// 	g_object_get(settings,
// 		"gtk-double-click-time", maxTime,
// 		"gtk-double-click-distance", maxDistance,
// 		NULL);
// }
import "C"

func gtkAreaNew() *C.GtkWidget {
	drawingarea := C.gtk_drawing_area_new()
	// the Area's size will be set later
	// we need to explicitly subscribe to mouse events with GtkDrawingArea
	C.gtk_widget_add_events(drawingarea,
		C.GDK_BUTTON_PRESS_MASK | C.GDK_BUTTON_RELEASE_MASK | C.GDK_POINTER_MOTION_MASK | C.GDK_BUTTON_MOTION_MASK | C.GDK_ENTER_NOTIFY_MASK | C.GDK_LEAVE_NOTIFY_MASK)
	// and we need to allow focusing on a GtkDrawingArea to enable keyboard events
	C.gtk_widget_set_can_focus(drawingarea, C.TRUE)
	scrollarea := C.gtk_scrolled_window_new((*C.GtkAdjustment)(nil), (*C.GtkAdjustment)(nil))
	// need a viewport because GtkDrawingArea isn't natively scrollable
	C.gtk_scrolled_window_add_with_viewport((*C.GtkScrolledWindow)(unsafe.Pointer(scrollarea)), drawingarea)
	return scrollarea
}

func gtkAreaGetControl(scrollarea *C.GtkWidget) *C.GtkWidget {
	viewport := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(scrollarea)))
	control := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(viewport)))
	return control
}

//export our_area_draw_callback
func our_area_draw_callback(widget *C.GtkWidget, cr *C.cairo_t, data C.gpointer) C.gboolean {
	var x0, y0, x1, y1 C.double
	var maxwid, maxht C.gint

	s := (*sysData)(unsafe.Pointer(data))
	// thanks to desrt in irc.gimp.net/#gtk+
	// TODO these are in "user coordinates"; is that what we want?
	C.cairo_clip_extents(cr, &x0, &y0, &x1, &y1)
	// we do not need to clear the cliprect; GtkDrawingArea did it for us beforehand
	cliprect := image.Rect(int(x0), int(y0), int(x1), int(y1))
	// the cliprect can actually fall outside the size of the Area; clip it by intersecting the two rectangles
	C.gtk_widget_get_size_request(widget, &maxwid, &maxht)
	cliprect = image.Rect(0, 0, int(maxwid), int(maxht)).Intersect(cliprect)
	if cliprect.Empty() {			// no intersection; nothing to paint
		return C.FALSE			// signals handled without stopping the event chain (thanks to desrt again)
	}
	i := s.handler.Paint(cliprect)
	surface := C.cairo_image_surface_create(
		C.CAIRO_FORMAT_ARGB32,			// alpha-premultiplied; native byte order
		C.int(i.Rect.Dx()),
		C.int(i.Rect.Dy()))
	if status := C.cairo_surface_status(surface); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("cairo_create_image_surface() failed: %s\n",
			C.GoString(C.cairo_status_to_string(status))))
	}
	// the flush and mark_dirty calls are required; see the cairo docs and https://git.gnome.org/browse/gtk+/tree/gdk/gdkcairo.c#n232 (thanks desrt in irc.gimp.net/#gtk+)
	C.cairo_surface_flush(surface)
	toARGB(i, uintptr(unsafe.Pointer(C.cairo_image_surface_get_data(surface))),
		int(C.cairo_image_surface_get_stride(surface)))
	C.cairo_surface_mark_dirty(surface)
	C.cairo_set_source_surface(cr,
		surface,
		0, 0)			// origin of the surface
	// that just set the brush that cairo uses: we have to actually draw now
	// (via https://developer.gnome.org/gtkmm-tutorial/stable/sec-draw-images.html.en)
	// TODO see above about user coordinates; if we do change to device coordinates the following line will need to change or be added to
	C.cairo_rectangle(cr, x0, y0, x1, y1)		// breaking the nrom here since we have the coordinates as a C double already
	C.cairo_fill(cr)
	C.cairo_surface_destroy(surface)		// free surface
	return C.FALSE		// signals handled without stopping the event chain (thanks to desrt again)
}

var area_draw_callback = C.GCallback(C.our_area_draw_callback)

func translateModifiers(state C.guint, window *C.GdkWindow) C.guint {
	// GDK doesn't initialize the modifier flags fully; we have to explicitly tell it to (thanks to Daniel_S and daniels (two different people) in irc.gimp.net/#gtk+)
	C.gdk_keymap_add_virtual_modifiers(
		C.gdk_keymap_get_for_display(C.gdk_window_get_display(window)),
		(*C.GdkModifierType)(unsafe.Pointer(&state)))
	return state
}

func makeModifiers(state C.guint) (m Modifiers) {
	if (state & C.GDK_CONTROL_MASK) != 0 {
		m |= Ctrl
	}
	if (state & C.GDK_META_MASK) != 0 {		// TODO get equivalent for Alt
		m |= Alt
	}
	if (state & C.GDK_SHIFT_MASK) != 0 {
		m |= Shift
	}
	if (state & C.GDK_SUPER_MASK) != 0 {
		m |= Super
	}
	return m
}

// shared code for finishing up and sending a mouse event
func finishMouseEvent(widget *C.GtkWidget, data C.gpointer, me MouseEvent, mb uint, x C.gdouble, y C.gdouble, state C.guint, gdkwindow *C.GdkWindow) {
	var areawidth, areaheight C.gint

	// on GTK+, mouse buttons 4-7 are for scrolling; if we got here, that's a mistake (and see the TODOs on return values below)
	if mb >= 4 && mb <= 7 {
		return
	}
	s := (*sysData)(unsafe.Pointer(data))
	state = translateModifiers(state, gdkwindow)
	me.Modifiers = makeModifiers(state)
	// the mb != # checks exclude the Up/Down button from Held
	if mb != 1 && (state & C.GDK_BUTTON1_MASK) != 0 {
		me.Held = append(me.Held, 1)
	}
	if mb != 2 && (state & C.GDK_BUTTON2_MASK) != 0 {
		me.Held = append(me.Held, 2)
	}
	if mb != 3 && (state & C.GDK_BUTTON3_MASK) != 0 {
		me.Held = append(me.Held, 3)
	}
	// don't check GDK_BUTTON4_MASK or GDK_BUTTON5_MASK because those are for the scrolling buttons mentioned above; there doesn't seem to be a way to detect higher buttons... (TODO)
	me.Pos = image.Pt(int(x), int(y))
	C.gtk_widget_get_size_request(widget, &areawidth, &areaheight)
	if !me.Pos.In(image.Rect(0, 0, int(areawidth), int(areaheight))) {		// outside the actual Area; no event
		return
	}
	// and finally, if the button ID >= 8, continue counting from 4, as above and as in the MouseEvent spec
	if me.Down >= 8 {
		me.Down -= 4
	}
	if me.Up >= 8 {
		me.Up -= 4
	}
	repaint := s.handler.Mouse(me)
	if repaint {
		C.gtk_widget_queue_draw(widget)
	}
}

//export our_area_button_press_event_callback
func our_area_button_press_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	// clicking doesn't automatically transfer keyboard focus; we must do so manually (thanks tristan in irc.gimp.net/#gtk+)
	C.gtk_widget_grab_focus(widget)
	e := (*C.GdkEventButton)(unsafe.Pointer(event))
	me := MouseEvent{
		// GDK button ID == our button ID with some exceptions taken care of by finishMouseEvent()
		Down:	uint(e.button),
	}

	var maxTime C.gint
	var maxDistance C.gint

	if e._type != C.GDK_BUTTON_PRESS {
		// ignore GDK's generated double-clicks and beyond; we handled those ourselves below
		return C.FALSE		// TODO really false?
	}
	s := (*sysData)(unsafe.Pointer(data))
	// e.time is unsigned and in milliseconds
	// maxTime is also milliseconds; despite being gint, it is only allowed to be positive
	// maxDistance is also only allowed to be positive
	settings := C.gtk_widget_get_settings(widget)
	C.gtkGetDoubleClickSettings(settings, &maxTime, &maxDistance)
	me.Count = s.clickCounter.click(me.Down, int(e.x), int(e.y),
		uintptr(e.time), uintptr(maxTime),
		int(maxDistance), int(maxDistance))

	finishMouseEvent(widget, data, me, me.Down, e.x, e.y, e.state, e.window)
	return C.FALSE			// TODO really false?
}

var area_button_press_event_callback = C.GCallback(C.our_area_button_press_event_callback)

//export our_area_button_release_event_callback
func our_area_button_release_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	e := (*C.GdkEventButton)(unsafe.Pointer(event))
	me := MouseEvent{
		// GDK button ID == our button ID with some exceptions taken care of by finishMouseEvent()
		Up:		uint(e.button),
	}
	finishMouseEvent(widget, data, me, me.Up, e.x, e.y, e.state, e.window)
	return C.FALSE			// TODO really false?
}

var area_button_release_event_callback = C.GCallback(C.our_area_button_release_event_callback)

//export our_area_motion_notify_event_callback
func our_area_motion_notify_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	e := (*C.GdkEventMotion)(unsafe.Pointer(event))
	me := MouseEvent{}
	finishMouseEvent(widget, data, me, 0, e.x, e.y, e.state, e.window)
	return C.FALSE			// TODO really false?
}

var area_motion_notify_event_callback = C.GCallback(C.our_area_motion_notify_event_callback)

// we want switching away from the control to reset the double-click counter, like with WM_ACTIVATE on Windows
// according to tristan in irc.gimp.net/#gtk+, doing this on enter-notify-event and leave-notify-event is correct (and it seems to be true in my own tests; plus the events DO get sent when switching programs with the keyboard (just pointing that out))
// differentiating between enter-notify-event and leave-notify-event is unimportant

//export our_area_enterleave_notify_event_callback
func our_area_enterleave_notify_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	s := (*sysData)(unsafe.Pointer(data))
	s.clickCounter.reset()
	return C.FALSE		// TODO really false?
}

var area_enterleave_notify_event_callback = C.GCallback(C.our_area_enterleave_notify_event_callback)

// shared code for doing a key event
func doKeyEvent(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer, up bool) bool {
	var ke KeyEvent

	e := (*C.GdkEventKey)(unsafe.Pointer(event))
	s := (*sysData)(unsafe.Pointer(data))
	keyval := e.keyval
	// get modifiers now in case a modifier was pressed
	state := translateModifiers(e.state, e.window)
	ke.Modifiers = makeModifiers(state)
	if extkey, ok := extkeys[keyval]; ok {
		ke.ExtKey = extkey
	} else if mod, ok := modonlykeys[keyval]; ok {
		ke.Modifier = mod
		// don't include the modifier in ke.Modifiers
		ke.Modifiers &^= mod
	} else if xke, ok := fromScancode(uintptr(e.hardware_keycode) - 8); ok {
		// see events_notdarwin.go for details of the above map lookup
		// one of these will be nonzero
		ke.Key = xke.Key
		ke.ExtKey = xke.ExtKey
	} else {		// no match
		// TODO really stop here? [or should we handle modifiers?]
		return false		// pretend unhandled
	}
	ke.Up = up
	handled, repaint := s.handler.Key(ke)
	if repaint {
		C.gtk_widget_queue_draw(widget)
	}
	return handled
}

//export our_area_key_press_event_callback
func our_area_key_press_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
/*
	fmt.Printf("PRESS %#v\n", e)
	fmt.Printf("this (%d/GDK_KEY_%s):\n", e.keyval,
		C.GoString((*C.char)(unsafe.Pointer(
			C.gdk_keyval_name(e.keyval)))))
	pk(e.keyval, e.window)
	fmt.Printf("%d/GDK_KEY_A:\n", C.GDK_KEY_A)
	pk(C.GDK_KEY_A, e.window)
	fmt.Printf("%d/GDK_KEY_a:\n", C.GDK_KEY_a)
	pk(C.GDK_KEY_a, e.window)
*/
	ret := doKeyEvent(widget, event, data, false)
	_ = ret
	return C.FALSE			// TODO really false? should probably return !ret (since true indicates stop processing)
}

var area_key_press_event_callback = C.GCallback(C.our_area_key_press_event_callback)

//export our_area_key_release_event_callback
func our_area_key_release_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	ret := doKeyEvent(widget, event, data, true)
	_ = ret
	return C.FALSE			// TODO really false? should probably return !ret (since true indicates stop processing)
}

var area_key_release_event_callback = C.GCallback(C.our_area_key_release_event_callback)

var extkeys = map[C.guint]ExtKey{
	C.GDK_KEY_Escape:			Escape,
	C.GDK_KEY_Insert:			Insert,
	C.GDK_KEY_Delete:			Delete,
	C.GDK_KEY_Home:			Home,
	C.GDK_KEY_End:			End,
	C.GDK_KEY_Page_Up:		PageUp,
	C.GDK_KEY_Page_Down:		PageDown,
	C.GDK_KEY_Up:			Up,
	C.GDK_KEY_Down:			Down,
	C.GDK_KEY_Left:			Left,
	C.GDK_KEY_Right:			Right,
	C.GDK_KEY_F1:			F1,
	C.GDK_KEY_F2:			F2,
	C.GDK_KEY_F3:			F3,
	C.GDK_KEY_F4:			F4,
	C.GDK_KEY_F5:			F5,
	C.GDK_KEY_F6:			F6,
	C.GDK_KEY_F7:			F7,
	C.GDK_KEY_F8:			F8,
	C.GDK_KEY_F9:			F9,
	C.GDK_KEY_F10:			F10,
	C.GDK_KEY_F11:			F11,
	C.GDK_KEY_F12:			F12,
	// numpad numeric keys and . are handled in events_notdarwin.go
	C.GDK_KEY_KP_Enter:		NEnter,
	C.GDK_KEY_KP_Add:		NAdd,
	C.GDK_KEY_KP_Subtract:		NSubtract,
	C.GDK_KEY_KP_Multiply:		NMultiply,
	C.GDK_KEY_KP_Divide:		NDivide,
}

// sanity check
func init() {
	included := make([]bool, _nextkeys)
	for _, v := range extkeys {
		included[v] = true
	}
	for i := 1; i < int(_nextkeys); i++ {
		if i >= int(N0) && i <= int(N9) {		// skip numpad numbers and .
			continue
		}
		if i == int(NDot) {
			continue
		}
		if !included[i] {
			panic(fmt.Errorf("error: not all ExtKeys defined on Unix (missing %d)", i))
		}
	}
}

var modonlykeys =  map[C.guint]Modifiers{
	C.GDK_KEY_Control_L:	Ctrl,
	C.GDK_KEY_Control_R:	Ctrl,
	C.GDK_KEY_Alt_L:		Alt,
	C.GDK_KEY_Alt_R:		Alt,
	C.GDK_KEY_Meta_L:		Alt,
	C.GDK_KEY_Meta_R:	Alt,
	C.GDK_KEY_Shift_L:		Shift,
	C.GDK_KEY_Shift_R:		Shift,
	C.GDK_KEY_Super_L:	Super,
	C.GDK_KEY_Super_R:	Super,
}

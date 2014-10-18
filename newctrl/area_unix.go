// +build !windows,!darwin,!plan9

// 14 march 2014

package ui

import (
	"fmt"
	"image"
	"unsafe"
)

// #include "gtk_unix.h"
// extern gboolean our_area_get_child_position_callback(GtkOverlay *, GtkWidget *, GdkRectangle *, gpointer);
// extern void our_area_textfield_populate_popup_callback(GtkEntry *, GtkMenu *, gpointer);
// extern gboolean our_area_textfield_focus_out_event_callback(GtkWidget *, GdkEvent *, gpointer);
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

type area struct {
	*areabase

	*scroller
	drawingarea *C.GtkDrawingArea

	clickCounter *clickCounter

	textfieldw    *C.GtkWidget
	textfield     *C.GtkEntry
	textfieldx    int
	textfieldy    int
	textfielddone *event
	inmenu        bool
}

func newArea(ab *areabase) Area {
	widget := C.gtk_drawing_area_new()
	// the Area's size will be set later
	// we need to explicitly subscribe to mouse events with GtkDrawingArea
	C.gtk_widget_add_events(widget,
		C.GDK_BUTTON_PRESS_MASK|C.GDK_BUTTON_RELEASE_MASK|C.GDK_POINTER_MOTION_MASK|C.GDK_BUTTON_MOTION_MASK|C.GDK_ENTER_NOTIFY_MASK|C.GDK_LEAVE_NOTIFY_MASK)
	// and we need to allow focusing on a GtkDrawingArea to enable keyboard events
	C.gtk_widget_set_can_focus(widget, C.TRUE)
	textfieldw := C.gtk_entry_new()
	a := &area{
		areabase:      ab,
		drawingarea:   (*C.GtkDrawingArea)(unsafe.Pointer(widget)),
		scroller:      newScroller(widget, false, false, true), // not natively scrollable; no border; have an overlay for OpenTextFieldAt()
		clickCounter:  new(clickCounter),
		textfieldw:    textfieldw,
		textfield:     (*C.GtkEntry)(unsafe.Pointer(textfieldw)),
		textfielddone: newEvent(),
	}
	a.fpreferredSize = a.xpreferredSize
	for _, c := range areaCallbacks {
		g_signal_connect(
			C.gpointer(unsafe.Pointer(a.drawingarea)),
			c.name,
			c.callback,
			C.gpointer(unsafe.Pointer(a)))
	}
	a.SetSize(a.width, a.height)
	C.gtk_overlay_add_overlay(a.scroller.overlayoverlay, a.textfieldw)
	g_signal_connect(
		C.gpointer(unsafe.Pointer(a.scroller.overlayoverlay)),
		"get-child-position",
		area_get_child_position_callback,
		C.gpointer(unsafe.Pointer(a)))
	// this part is important
	// entering the context menu is considered focusing out
	// so we connect to populate-popup to mark that we're entering the context menu (thanks slaf in irc.gimp.net/#gtk+)
	// and we have to connect_after to focus-out-event so that it runs after the populate-popup
	g_signal_connect(
		C.gpointer(unsafe.Pointer(a.textfield)),
		"populate-popup",
		area_textfield_populate_popup_callback,
		C.gpointer(unsafe.Pointer(a)))
	g_signal_connect_after(
		C.gpointer(unsafe.Pointer(a.textfield)),
		"focus-out-event",
		area_textfield_focus_out_event_callback,
		C.gpointer(unsafe.Pointer(a)))
	// the widget shows up initially
	C.gtk_widget_set_no_show_all(a.textfieldw, C.TRUE)
	return a
}

func (a *area) SetSize(width, height int) {
	a.width = width
	a.height = height
	C.gtk_widget_set_size_request(a.widget, C.gint(a.width), C.gint(a.height))
}

func (a *area) Repaint(r image.Rectangle) {
	r = image.Rect(0, 0, a.width, a.height).Intersect(r)
	if r.Empty() {
		return
	}
	C.gtk_widget_queue_draw_area(a.widget, C.gint(r.Min.X), C.gint(r.Min.Y), C.gint(r.Dx()), C.gint(r.Dy()))
}

func (a *area) RepaintAll() {
	C.gtk_widget_queue_draw(a.widget)
}

func (a *area) OpenTextFieldAt(x, y int) {
	if x < 0 || x >= a.width || y < 0 || y >= a.height {
		panic(fmt.Errorf("point (%d,%d) outside Area in Area.OpenTextFieldAt()", x, y))
	}
	a.textfieldx = x
	a.textfieldy = y
	a.inmenu = false // to start
	// we disabled this for the initial Area show; we don't need to anymore
	C.gtk_widget_set_no_show_all(a.textfieldw, C.FALSE)
	C.gtk_widget_show_all(a.textfieldw)
	C.gtk_widget_grab_focus(a.textfieldw)
}

func (a *area) TextFieldText() string {
	return fromgstr(C.gtk_entry_get_text(a.textfield))
}

func (a *area) SetTextFieldText(text string) {
	ctext := togstr(text)
	defer freegstr(ctext)
	C.gtk_entry_set_text(a.textfield, ctext)
}

func (a *area) OnTextFieldDismissed(f func()) {
	a.textfielddone.set(f)
}

//export our_area_get_child_position_callback
func our_area_get_child_position_callback(overlay *C.GtkOverlay, widget *C.GtkWidget, rect *C.GdkRectangle, data C.gpointer) C.gboolean {
	var nat C.GtkRequisition

	a := (*area)(unsafe.Pointer(data))
	rect.x = C.int(a.textfieldx)
	rect.y = C.int(a.textfieldy)
	C.gtk_widget_get_preferred_size(a.textfieldw, nil, &nat)
	rect.width = C.int(nat.width)
	rect.height = C.int(nat.height)
	return C.TRUE
}

var area_get_child_position_callback = C.GCallback(C.our_area_get_child_position_callback)

//export our_area_textfield_populate_popup_callback
func our_area_textfield_populate_popup_callback(entry *C.GtkEntry, menu *C.GtkMenu, data C.gpointer) {
	a := (*area)(unsafe.Pointer(data))
	a.inmenu = true
}

var area_textfield_populate_popup_callback = C.GCallback(C.our_area_textfield_populate_popup_callback)

//export our_area_textfield_focus_out_event_callback
func our_area_textfield_focus_out_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	a := (*area)(unsafe.Pointer(data))
	if !a.inmenu {
		C.gtk_widget_hide(a.textfieldw)
		a.textfielddone.fire()
	}
	a.inmenu = false // for next time
	return continueEventChain
}

var area_textfield_focus_out_event_callback = C.GCallback(C.our_area_textfield_focus_out_event_callback)

var areaCallbacks = []struct {
	name     string
	callback C.GCallback
}{
	{"draw", area_draw_callback},
	{"button-press-event", area_button_press_event_callback},
	{"button-release-event", area_button_release_event_callback},
	{"motion-notify-event", area_motion_notify_event_callback},
	{"enter-notify-event", area_enterleave_notify_event_callback},
	{"leave-notify-event", area_enterleave_notify_event_callback},
	{"key-press-event", area_key_press_event_callback},
	{"key-release-event", area_key_release_event_callback},
}

//export our_area_draw_callback
func our_area_draw_callback(widget *C.GtkWidget, cr *C.cairo_t, data C.gpointer) C.gboolean {
	var x0, y0, x1, y1 C.double

	a := (*area)(unsafe.Pointer(data))
	// thanks to desrt in irc.gimp.net/#gtk+
	// these are in user coordinates, which match what coordinates we want by default, even out of a draw event handler (thanks johncc3, mclasen, and Company in irc.gimp.net/#gtk+)
	C.cairo_clip_extents(cr, &x0, &y0, &x1, &y1)
	// we do not need to clear the cliprect; GtkDrawingArea did it for us beforehand
	cliprect := image.Rect(int(x0), int(y0), int(x1), int(y1))
	// the cliprect can actually fall outside the size of the Area; clip it by intersecting the two rectangles
	cliprect = image.Rect(0, 0, a.width, a.height).Intersect(cliprect)
	if cliprect.Empty() { // no intersection; nothing to paint
		return C.FALSE // signals handled without stopping the event chain (thanks to desrt again)
	}
	i := a.handler.Paint(cliprect)
	surface := C.cairo_image_surface_create(
		C.CAIRO_FORMAT_ARGB32, // alpha-premultiplied; native byte order
		C.int(i.Rect.Dx()),
		C.int(i.Rect.Dy()))
	if status := C.cairo_surface_status(surface); status != C.CAIRO_STATUS_SUCCESS {
		panic(fmt.Errorf("cairo_create_image_surface() failed: %s\n",
			C.GoString(C.cairo_status_to_string(status))))
	}
	// the flush and mark_dirty calls are required; see the cairo docs and https://git.gnome.org/browse/gtk+/tree/gdk/gdkcairo.c#n232 (thanks desrt in irc.gimp.net/#gtk+)
	C.cairo_surface_flush(surface)
	toARGB(i, uintptr(unsafe.Pointer(C.cairo_image_surface_get_data(surface))),
		int(C.cairo_image_surface_get_stride(surface)), false) // not NRGBA
	C.cairo_surface_mark_dirty(surface)
	C.cairo_set_source_surface(cr,
		surface,
		x0, y0) // point on cairo_t where we want to draw (thanks Company in irc.gimp.net/#gtk+)
	// that just set the brush that cairo uses: we have to actually draw now
	// (via https://developer.gnome.org/gtkmm-tutorial/stable/sec-draw-images.html.en)
	C.cairo_rectangle(cr, x0, y0, x1, y1) // breaking the norm here since we have the coordinates as a C double already
	C.cairo_fill(cr)
	C.cairo_surface_destroy(surface) // free surface
	return C.FALSE                   // signals handled without stopping the event chain (thanks to desrt again)
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
	if (state & C.GDK_META_MASK) != 0 {
		m |= Alt
	}
	if (state & C.GDK_MOD1_MASK) != 0 { // GTK+ itself requires this to be Alt (just read through gtkaccelgroup.c)
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

	// on GTK+, mouse buttons 4-7 are for scrolling; if we got here, that's a mistake
	if mb >= 4 && mb <= 7 {
		return
	}
	a := (*area)(unsafe.Pointer(data))
	state = translateModifiers(state, gdkwindow)
	me.Modifiers = makeModifiers(state)
	// the mb != # checks exclude the Up/Down button from Held
	if mb != 1 && (state&C.GDK_BUTTON1_MASK) != 0 {
		me.Held = append(me.Held, 1)
	}
	if mb != 2 && (state&C.GDK_BUTTON2_MASK) != 0 {
		me.Held = append(me.Held, 2)
	}
	if mb != 3 && (state&C.GDK_BUTTON3_MASK) != 0 {
		me.Held = append(me.Held, 3)
	}
	// don't check GDK_BUTTON4_MASK or GDK_BUTTON5_MASK because those are for the scrolling buttons mentioned above
	// GDK expressly does not support any more buttons in the GdkModifierType; see https://git.gnome.org/browse/gtk+/tree/gdk/x11/gdkdevice-xi2.c#n763 (thanks mclasen in irc.gimp.net/#gtk+)
	me.Pos = image.Pt(int(x), int(y))
	C.gtk_widget_get_size_request(widget, &areawidth, &areaheight)
	if !me.Pos.In(image.Rect(0, 0, int(areawidth), int(areaheight))) { // outside the actual Area; no event
		return
	}
	// and finally, if the button ID >= 8, continue counting from 4, as above and as in the MouseEvent spec
	if me.Down >= 8 {
		me.Down -= 4
	}
	if me.Up >= 8 {
		me.Up -= 4
	}
	a.handler.Mouse(me)
}

// convenience name to make our intent clear
const continueEventChain C.gboolean = C.FALSE
const stopEventChain C.gboolean = C.TRUE

// checking for a mouse click that makes the program/window active is meaningless on GTK+: it's a property of the window manager/X11, and it's the WM that decides if the program should become active or not
// however, one thing is certain: the click event will ALWAYS be sent (to the window that the X11 decides to send it to)
// I assume the same is true for Wayland
// thanks Chipzz in irc.gimp.net/#gtk+

//export our_area_button_press_event_callback
func our_area_button_press_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	// clicking doesn't automatically transfer keyboard focus; we must do so manually (thanks tristan in irc.gimp.net/#gtk+)
	C.gtk_widget_grab_focus(widget)
	e := (*C.GdkEventButton)(unsafe.Pointer(event))
	me := MouseEvent{
		// GDK button ID == our button ID with some exceptions taken care of by finishMouseEvent()
		Down: uint(e.button),
	}

	var maxTime C.gint
	var maxDistance C.gint

	if e._type != C.GDK_BUTTON_PRESS {
		// ignore GDK's generated double-clicks and beyond; we handled those ourselves below
		return continueEventChain
	}
	a := (*area)(unsafe.Pointer(data))
	// e.time is unsigned and in milliseconds
	// maxTime is also milliseconds; despite being gint, it is only allowed to be positive
	// maxDistance is also only allowed to be positive
	settings := C.gtk_widget_get_settings(widget)
	C.gtkGetDoubleClickSettings(settings, &maxTime, &maxDistance)
	me.Count = a.clickCounter.click(me.Down, int(e.x), int(e.y),
		uintptr(e.time), uintptr(maxTime),
		int(maxDistance), int(maxDistance))

	finishMouseEvent(widget, data, me, me.Down, e.x, e.y, e.state, e.window)
	return continueEventChain
}

var area_button_press_event_callback = C.GCallback(C.our_area_button_press_event_callback)

//export our_area_button_release_event_callback
func our_area_button_release_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	e := (*C.GdkEventButton)(unsafe.Pointer(event))
	me := MouseEvent{
		// GDK button ID == our button ID with some exceptions taken care of by finishMouseEvent()
		Up: uint(e.button),
	}
	finishMouseEvent(widget, data, me, me.Up, e.x, e.y, e.state, e.window)
	return continueEventChain
}

var area_button_release_event_callback = C.GCallback(C.our_area_button_release_event_callback)

//export our_area_motion_notify_event_callback
func our_area_motion_notify_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	e := (*C.GdkEventMotion)(unsafe.Pointer(event))
	me := MouseEvent{}
	finishMouseEvent(widget, data, me, 0, e.x, e.y, e.state, e.window)
	return continueEventChain
}

var area_motion_notify_event_callback = C.GCallback(C.our_area_motion_notify_event_callback)

// we want switching away from the control to reset the double-click counter, like with WM_ACTIVATE on Windows
// according to tristan in irc.gimp.net/#gtk+, doing this on enter-notify-event and leave-notify-event is correct (and it seems to be true in my own tests; plus the events DO get sent when switching programs with the keyboard (just pointing that out))
// differentiating between enter-notify-event and leave-notify-event is unimportant

//export our_area_enterleave_notify_event_callback
func our_area_enterleave_notify_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	a := (*area)(unsafe.Pointer(data))
	a.clickCounter.reset()
	return continueEventChain
}

var area_enterleave_notify_event_callback = C.GCallback(C.our_area_enterleave_notify_event_callback)

// shared code for doing a key event
func doKeyEvent(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer, up bool) bool {
	var ke KeyEvent

	e := (*C.GdkEventKey)(unsafe.Pointer(event))
	a := (*area)(unsafe.Pointer(data))
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
	} else { // no match
		return false
	}
	ke.Up = up
	return a.handler.Key(ke)
}

//export our_area_key_press_event_callback
func our_area_key_press_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	if doKeyEvent(widget, event, data, false) == true {
		return stopEventChain
	}
	return continueEventChain
}

var area_key_press_event_callback = C.GCallback(C.our_area_key_press_event_callback)

//export our_area_key_release_event_callback
func our_area_key_release_event_callback(widget *C.GtkWidget, event *C.GdkEvent, data C.gpointer) C.gboolean {
	if doKeyEvent(widget, event, data, true) == true {
		return stopEventChain
	}
	return continueEventChain
}

var area_key_release_event_callback = C.GCallback(C.our_area_key_release_event_callback)

var extkeys = map[C.guint]ExtKey{
	C.GDK_KEY_Escape:    Escape,
	C.GDK_KEY_Insert:    Insert,
	C.GDK_KEY_Delete:    Delete,
	C.GDK_KEY_Home:      Home,
	C.GDK_KEY_End:       End,
	C.GDK_KEY_Page_Up:   PageUp,
	C.GDK_KEY_Page_Down: PageDown,
	C.GDK_KEY_Up:        Up,
	C.GDK_KEY_Down:      Down,
	C.GDK_KEY_Left:      Left,
	C.GDK_KEY_Right:     Right,
	C.GDK_KEY_F1:        F1,
	C.GDK_KEY_F2:        F2,
	C.GDK_KEY_F3:        F3,
	C.GDK_KEY_F4:        F4,
	C.GDK_KEY_F5:        F5,
	C.GDK_KEY_F6:        F6,
	C.GDK_KEY_F7:        F7,
	C.GDK_KEY_F8:        F8,
	C.GDK_KEY_F9:        F9,
	C.GDK_KEY_F10:       F10,
	C.GDK_KEY_F11:       F11,
	C.GDK_KEY_F12:       F12,
	// numpad numeric keys and . are handled in events_notdarwin.go
	C.GDK_KEY_KP_Enter:    NEnter,
	C.GDK_KEY_KP_Add:      NAdd,
	C.GDK_KEY_KP_Subtract: NSubtract,
	C.GDK_KEY_KP_Multiply: NMultiply,
	C.GDK_KEY_KP_Divide:   NDivide,
}

// sanity check
func init() {
	included := make([]bool, _nextkeys)
	for _, v := range extkeys {
		included[v] = true
	}
	for i := 1; i < int(_nextkeys); i++ {
		if i >= int(N0) && i <= int(N9) { // skip numpad numbers and .
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

var modonlykeys = map[C.guint]Modifiers{
	C.GDK_KEY_Control_L: Ctrl,
	C.GDK_KEY_Control_R: Ctrl,
	C.GDK_KEY_Alt_L:     Alt,
	C.GDK_KEY_Alt_R:     Alt,
	C.GDK_KEY_Meta_L:    Alt,
	C.GDK_KEY_Meta_R:    Alt,
	C.GDK_KEY_Shift_L:   Shift,
	C.GDK_KEY_Shift_R:   Shift,
	C.GDK_KEY_Super_L:   Super,
	C.GDK_KEY_Super_R:   Super,
}

func (a *area) xpreferredSize(d *sizing) (width, height int) {
	// the preferred size of an Area is its size
	return a.width, a.height
}

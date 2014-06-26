// +build !windows,!darwin,!plan9
// this is manual but either this or the opposite (listing all valid systems) really are the only ways to do it; proposals for a 'unix' tag were rejected (https://code.google.com/p/go/issues/detail?id=6325)

// 16 february 2014

package ui

import (
	"fmt"
	"unsafe"
)

// #include "gtk_unix.h"
// static inline void gtkSetComboBoxArbitrarilyResizeable(GtkWidget *w)
// {
// 	/* we can safely assume that the cell renderers of a GtkComboBoxText are a single GtkCellRendererText by default */
// 	/* thanks mclasen and tristan in irc.gimp.net/#gtk+ */
// 	GtkCellLayout *cbox = (GtkCellLayout *) w;
// 	GList *renderer;
//
// 	renderer = gtk_cell_layout_get_cells(cbox);
// 	g_object_set(renderer->data,
// 		"ellipsize-set", TRUE,		/* with no ellipsizing or wrapping, the combobox won't resize */
// 		"ellipsize", PANGO_ELLIPSIZE_END,
// 		"width-chars", 0,
// 		NULL);
// 	g_list_free(renderer);
// }
import "C"

var gtkStyles = []byte(`
/*
this first one is for ProgressBar so we can arbitrarily resize it
min-horizontal-bar-width is a style property; we do it through CSS
thanks tristan in irc.gimp.net/#gtk+
*/
* {
	-GtkProgressBar-min-horizontal-bar-width: 1;
}
`)

func gtk_init() error {
	var err *C.GError = nil // redundant in Go, but let's explicitly assign it anyway

	// gtk_init_with_args() gives us error info (thanks chpe in irc.gimp.net/#gtk+)
	// don't worry about GTK+'s command-line arguments; they're also available as environment variables (thanks mclasen in irc.gimp.net/#gtk+)
	result := C.gtk_init_with_args(nil, nil, nil, nil, nil, &err)
	if result == C.FALSE {
		return fmt.Errorf("error actually initilaizing GTK+: %s", fromgstr(err.message))
	}

	// now apply our custom program-global styles
	provider := C.gtk_css_provider_new()
	if C.gtk_css_provider_load_from_data(provider, (*C.gchar)(unsafe.Pointer(&gtkStyles[0])), C.gssize(len(gtkStyles)), &err) == C.FALSE {
		return fmt.Errorf("error applying package ui's custom program-global styles to GTK+: %v", fromgstr(err.message))
	}
	// GDK (at least as far back as GTK+ 3.4, but officially documented as of 3.10) merges all screens into one big one, so we don't need to worry about multimonitor
	// thanks to baedert and mclasen in irc.gimp.net/#gtk+
	C.gtk_style_context_add_provider_for_screen(C.gdk_screen_get_default(),
		(*C.GtkStyleProvider)(unsafe.Pointer(provider)),
		C.GTK_STYLE_PROVIDER_PRIORITY_USER)

	return nil
}

func gtk_main_quit() {
	C.gtk_main_quit()
}

func gtk_window_new() *C.GtkWidget {
	return C.gtk_window_new(C.GTK_WINDOW_TOPLEVEL)
}

func gtk_widget_show(widget *C.GtkWidget) {
	C.gtk_widget_show_all(widget)
}

func gtk_widget_hide(widget *C.GtkWidget) {
	C.gtk_widget_hide(widget)
}

func gtk_window_set_title(window *C.GtkWidget, title string) {
	ctitle := C.CString(title)
	defer C.free(unsafe.Pointer(ctitle))
	C.gtk_window_set_title(togtkwindow(window), togstr(ctitle))
}

func gtk_window_get_title(window *C.GtkWidget) string {
	return fromgstr(C.gtk_window_get_title(togtkwindow(window)))
}

func gtk_window_resize(window *C.GtkWidget, width int, height int) {
	C.gtk_window_resize(togtkwindow(window), C.gint(width), C.gint(height))
}

func gtk_window_get_size(window *C.GtkWidget) (int, int) {
	var width, height C.gint

	C.gtk_window_get_size(togtkwindow(window), &width, &height)
	return int(width), int(height)
}

// this should allow us to resize the window arbitrarily
// thanks to Company in irc.gimp.net/#gtk+
func gtkNewWindowLayout() *C.GtkWidget {
	layout := C.gtk_layout_new(nil, nil)
	scrollarea := C.gtk_scrolled_window_new((*C.GtkAdjustment)(nil), (*C.GtkAdjustment)(nil))
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(scrollarea)), layout)
	// never show scrollbars; we're just doing this to allow arbitrary resizes
	C.gtk_scrolled_window_set_policy((*C.GtkScrolledWindow)(unsafe.Pointer(scrollarea)),
		C.GTK_POLICY_NEVER, C.GTK_POLICY_NEVER)
	return scrollarea
}

func gtk_container_add(container *C.GtkWidget, widget *C.GtkWidget) {
	C.gtk_container_add(togtkcontainer(container), widget)
}

func gtkAddWidgetToLayout(container *C.GtkWidget, widget *C.GtkWidget) {
	layout := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(container)))
	C.gtk_container_add((*C.GtkContainer)(unsafe.Pointer(layout)), widget)
}

func gtkMoveWidgetInLayout(container *C.GtkWidget, widget *C.GtkWidget, x int, y int) {
	layout := C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(container)))
	C.gtk_layout_move((*C.GtkLayout)(unsafe.Pointer(layout)), widget,
		C.gint(x), C.gint(y))
}

func gtk_widget_set_size_request(widget *C.GtkWidget, width int, height int) {
	C.gtk_widget_set_size_request(widget, C.gint(width), C.gint(height))
}

func gtk_button_new() *C.GtkWidget {
	return C.gtk_button_new()
}

func gtk_button_set_label(button *C.GtkWidget, label string) {
	clabel := C.CString(label)
	defer C.free(unsafe.Pointer(clabel))
	C.gtk_button_set_label(togtkbutton(button), togstr(clabel))
}

func gtk_button_get_label(button *C.GtkWidget) string {
	return fromgstr(C.gtk_button_get_label(togtkbutton(button)))
}

func gtk_check_button_new() *C.GtkWidget {
	return C.gtk_check_button_new()
}

func gtk_toggle_button_get_active(widget *C.GtkWidget) bool {
	return fromgbool(C.gtk_toggle_button_get_active(togtktogglebutton(widget)))
}

func gtk_combo_box_text_new() *C.GtkWidget {
	w := C.gtk_combo_box_text_new()
	C.gtkSetComboBoxArbitrarilyResizeable(w)
	return w
}

func gtk_combo_box_text_new_with_entry() *C.GtkWidget {
	w := C.gtk_combo_box_text_new_with_entry()
	// we need to set the entry's width-chars to achieve the arbitrary sizing we want
	// thanks to tristan in irc.gimp.net/#gtk+
	// in my tests, this is sufficient to get the effect we want; setting the cell renderer isn't necessary like it is with an entry-less combobox
	entry := togtkentry(C.gtk_bin_get_child((*C.GtkBin)(unsafe.Pointer(w))))
	C.gtk_entry_set_width_chars(entry, 0)
	return w
}

func gtk_combo_box_text_get_active_text(widget *C.GtkWidget) string {
	return fromgstr(C.gtk_combo_box_text_get_active_text(togtkcombobox(widget)))
}

func gtk_combo_box_text_append_text(widget *C.GtkWidget, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_combo_box_text_append_text(togtkcombobox(widget), togstr(ctext))
}

func gtk_combo_box_text_insert_text(widget *C.GtkWidget, index int, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_combo_box_text_insert_text(togtkcombobox(widget), C.gint(index), togstr(ctext))
}

func gtk_combo_box_get_active(widget *C.GtkWidget) int {
	cb := (*C.GtkComboBox)(unsafe.Pointer(widget))
	return int(C.gtk_combo_box_get_active(cb))
}

func gtk_combo_box_text_remove(widget *C.GtkWidget, index int) {
	C.gtk_combo_box_text_remove(togtkcombobox(widget), C.gint(index))
}

func gtkComboBoxLen(widget *C.GtkWidget) int {
	cb := (*C.GtkComboBox)(unsafe.Pointer(widget))
	model := C.gtk_combo_box_get_model(cb)
	// this is the same as with a Listbox so
	return gtkTreeModelListLen(model)
}

func gtk_entry_new() *C.GtkWidget {
	e := C.gtk_entry_new()
	// allows the GtkEntry to be resized with the window smaller than what it thinks the size should be
	// thanks to Company in irc.gimp.net/#gtk+
	C.gtk_entry_set_width_chars(togtkentry(e), 0)
	return e
}

func gtkPasswordEntryNew() *C.GtkWidget {
	e := gtk_entry_new()
	C.gtk_entry_set_visibility(togtkentry(e), C.FALSE)
	return e
}

func gtk_entry_set_text(widget *C.GtkWidget, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_entry_set_text(togtkentry(widget), togstr(ctext))
}

func gtk_entry_get_text(widget *C.GtkWidget) string {
	return fromgstr(C.gtk_entry_get_text(togtkentry(widget)))
}

var _emptystring = [1]C.gchar{0}
var emptystring = &_emptystring[0]

func gtk_misc_set_alignment(widget *C.GtkWidget, x float64, y float64) {
	C.gtk_misc_set_alignment((*C.GtkMisc)(unsafe.Pointer(widget)), C.gfloat(x), C.gfloat(y))
}

func gtk_label_new() *C.GtkWidget {
	label := C.gtk_label_new(emptystring)
	C.gtk_label_set_line_wrap(togtklabel(label), C.FALSE) // turn off line wrap
	// don't call gtk_label_set_line_wrap_mode(); there's no "wrap none" value there anyway
	C.gtk_label_set_ellipsize(togtklabel(label), C.PANGO_ELLIPSIZE_NONE) // turn off ellipsizing; this + line wrapping above will guarantee cutoff as documented
	// there's a function gtk_label_set_justify() that indicates GTK_JUSTIFY_LEFT is the default
	// but this actually is NOT the control justification, just the multi-line justification
	// so we need to do THIS instead
	// this will valign to the center, which is what the HIG says (https://developer.gnome.org/hig-book/3.4/design-text-labels.html.en) for all controls that label to the side; thankfully this means we don't need to do any extra positioning magic
	// this will also valign to the top
	// thanks to mclasen in irc.gimp.net/#gtk+
	gtk_misc_set_alignment(label, 0, 0.5)
	return label
}

func gtk_label_new_standalone() *C.GtkWidget {
	label := gtk_label_new()
	// this will valign to the top
	gtk_misc_set_alignment(label, 0, 0)
	return label
}

func gtk_label_set_text(widget *C.GtkWidget, text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.gtk_label_set_text(togtklabel(widget), togstr(ctext))
}

func gtk_label_get_text(widget *C.GtkWidget) string {
	return fromgstr(C.gtk_label_get_text(togtklabel(widget)))
}

func gtk_widget_get_preferred_size(widget *C.GtkWidget) (minWidth int, minHeight int, natWidth int, natHeight int) {
	var minimum, natural C.GtkRequisition

	C.gtk_widget_get_preferred_size(widget, &minimum, &natural)
	return int(minimum.width), int(minimum.height),
		int(natural.width), int(natural.height)
}

func gtk_progress_bar_new() *C.GtkWidget {
	return C.gtk_progress_bar_new()
}

func gtk_progress_bar_set_fraction(w *C.GtkWidget, percent int) {
	p := C.gdouble(percent) / 100
	C.gtk_progress_bar_set_fraction(togtkprogressbar(w), p)
}

func gtk_progress_bar_pulse(w *C.GtkWidget) {
	C.gtk_progress_bar_pulse(togtkprogressbar(w))
}

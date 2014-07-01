// +build !windows,!darwin,!plan9

// 16 february 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
// extern gboolean our_pulse_callback(gpointer);
import "C"

type sysData struct {
	cSysData

	widget       *C.GtkWidget
	container    *C.GtkWidget // for moving
	pulseTimer	C.guint	// for indeterminate progress bars
	clickCounter clickCounter // for Areas
	// we probably don't need to save these, but we'll do so for sysData.preferredSize() just in case
	areawidth  int
	areaheight int
}

type classData struct {
	make     func() *C.GtkWidget
	makeAlt  func() *C.GtkWidget
	setText  func(widget *C.GtkWidget, text string)
	text     func(widget *C.GtkWidget) string
	append   func(widget *C.GtkWidget, text string)
	insert   func(widget *C.GtkWidget, index int, text string)
	selected func(widget *C.GtkWidget) int
	selMulti func(widget *C.GtkWidget) []int
	smtexts  func(widget *C.GtkWidget) []string
	delete   func(widget *C.GtkWidget, index int)
	len      func(widget *C.GtkWidget) int
	// ...
	signals   callbackMap
	child     func(widget *C.GtkWidget) *C.GtkWidget
	childsigs callbackMap
}

var classTypes = [nctypes]*classData{
	c_window: &classData{
		make:    gtk_window_new,
		setText: gtk_window_set_title,
		text:    gtk_window_get_title,
		signals: callbackMap{
			"delete-event":    window_delete_event_callback,
			"configure-event": window_configure_event_callback,
		},
	},
	c_button: &classData{
		make:    gtk_button_new,
		setText: gtk_button_set_label,
		text:    gtk_button_get_label,
		signals: callbackMap{
			"clicked": button_clicked_callback,
		},
	},
	c_checkbox: &classData{
		make:    gtk_check_button_new,
		setText: gtk_button_set_label,
		text:    gtk_button_get_label,
	},
	c_combobox: &classData{
		make:     gtk_combo_box_text_new,
		makeAlt:  gtk_combo_box_text_new_with_entry,
		text:     gtk_combo_box_text_get_active_text,
		append:   gtk_combo_box_text_append_text,
		insert:   gtk_combo_box_text_insert_text,
		selected: gtk_combo_box_get_active,
		delete:   gtk_combo_box_text_remove,
		len:      gtkComboBoxLen,
	},
	c_lineedit: &classData{
		make:    gtk_entry_new,
		makeAlt: gtkPasswordEntryNew,
		setText: gtk_entry_set_text,
		text:    gtk_entry_get_text,
	},
	c_label: &classData{
		make:    gtk_label_new,
		makeAlt:	gtk_label_new_standalone,
		setText: gtk_label_set_text,
		text:    gtk_label_get_text,
	},
	c_listbox: &classData{
		make:     gListboxNewSingle,
		makeAlt:  gListboxNewMulti,
		text:     gListboxText,
		append:   gListboxAppend,
		insert:   gListboxInsert,
		selMulti: gListboxSelectedMulti,
		smtexts:  gListboxSelMultiTexts,
		delete:   gListboxDelete,
		len:      gListboxLen,
	},
	c_progressbar: &classData{
		make: gtk_progress_bar_new,
	},
	c_area: &classData{
		make:  gtkAreaNew,
		child: gtkAreaGetControl,
		childsigs: callbackMap{
			"draw":                 area_draw_callback,
			"button-press-event":   area_button_press_event_callback,
			"button-release-event": area_button_release_event_callback,
			"motion-notify-event":  area_motion_notify_event_callback,
			"enter-notify-event":   area_enterleave_notify_event_callback,
			"leave-notify-event":   area_enterleave_notify_event_callback,
			"key-press-event":      area_key_press_event_callback,
			"key-release-event":    area_key_release_event_callback,
		},
	},
}

func (s *sysData) make(window *sysData) error {
	ct := classTypes[s.ctype]
	if s.alternate {
		s.widget = ct.makeAlt()
	} else {
		s.widget = ct.make()
	}
	if window == nil {
		fixed := gtkNewWindowLayout()
		gtk_container_add(s.widget, fixed)
		for signame, sigfunc := range ct.signals {
			g_signal_connect(s.widget, signame, sigfunc, s)
		}
		s.container = fixed
	} else {
		s.container = window.container
		gtkAddWidgetToLayout(s.container, s.widget)
		for signame, sigfunc := range ct.signals {
			g_signal_connect(s.widget, signame, sigfunc, s)
		}
		if ct.child != nil {
			child := ct.child(s.widget)
			for signame, sigfunc := range ct.childsigs {
				g_signal_connect(child, signame, sigfunc, s)
			}
		}
	}
	return nil
}

// see sysData.center()
func (s *sysData) resetposition() {
	C.gtk_window_set_position(togtkwindow(s.widget), C.GTK_WIN_POS_NONE)
}

// used for Windows; nothing special needed elsewhere
func (s *sysData) firstShow() error {
	s.show()
	return nil
}

func (s *sysData) show() {
	gtk_widget_show(s.widget)
	s.resetposition()
}

func (s *sysData) hide() {
	gtk_widget_hide(s.widget)
	s.resetposition()
}

func (s *sysData) setText(text string) {
	classTypes[s.ctype].setText(s.widget, text)
}

func (s *sysData) setRect(x int, y int, width int, height int, winheight int) error {
	gtkMoveWidgetInLayout(s.container, s.widget, x, y)
	gtk_widget_set_size_request(s.widget, width, height)
	return nil
}

func (s *sysData) isChecked() bool {
	return gtk_toggle_button_get_active(s.widget)
}

func (s *sysData) text() string {
	return classTypes[s.ctype].text(s.widget)
}

func (s *sysData) append(what string) {
	classTypes[s.ctype].append(s.widget, what)
}

func (s *sysData) insertBefore(what string, before int) {
	classTypes[s.ctype].insert(s.widget, before, what)
}

func (s *sysData) selectedIndex() int {
	return classTypes[s.ctype].selected(s.widget)
}

func (s *sysData) selectedIndices() []int {
	return classTypes[s.ctype].selMulti(s.widget)
}

func (s *sysData) selectedTexts() []string {
	return classTypes[s.ctype].smtexts(s.widget)
}

func (s *sysData) setWindowSize(width int, height int) error {
	// does not take window geometry into account (and cannot, since the window manager won't give that info away)
	// thanks to TingPing in irc.gimp.net/#gtk+
	gtk_window_resize(s.widget, width, height)
	return nil
}

func (s *sysData) delete(index int) {
	classTypes[s.ctype].delete(s.widget, index)
}

// With GTK+, we must manually pulse the indeterminate progressbar ourselves.
// To ensure pulsing runs on the main lop and doesn't cause any other weird racy things, we'll use g_timeout_add().
// Zenity 3.4 does this too (https://git.gnome.org/browse/zenity/tree/src/progress.c?id=3.4.0).
// The following is Zenity 3.4's pulse rate.
const pulseRate = 100	// in milliseconds

//export our_pulse_callback
func our_pulse_callback(data C.gpointer) C.gboolean {
	// TODO this can be called when closing the window
	s := (*sysData)(unsafe.Pointer(data))
	gtk_progress_bar_pulse(s.widget)
	return C.TRUE		// continue processing
}

func (s *sysData) setProgress(percent int) {
	if s.pulseTimer != 0 {		// kill current timer
		// TODO only if not indeterminate already?
		C.g_source_remove(s.pulseTimer)
		s.pulseTimer = 0
	}
	if percent == -1 {
		gtk_progress_bar_pulse(s.widget)		// start animation now
		s.pulseTimer = C.g_timeout_add(pulseRate,
			C.GSourceFunc(C.our_pulse_callback),
			C.gpointer(unsafe.Pointer(s)))
		return
	}
	gtk_progress_bar_set_fraction(s.widget, percent)
}

func (s *sysData) len() int {
	return classTypes[s.ctype].len(s.widget)
}

func (s *sysData) setAreaSize(width int, height int) {
	c := gtkAreaGetControl(s.widget)
	gtk_widget_set_size_request(c, width, height)
	s.areawidth = width // for sysData.preferredSize()
	s.areaheight = height
	C.gtk_widget_queue_draw(c)
}

// TODO should this be made safe? (TODO move to area.go)
func (s *sysData) repaintAll() {
	c := gtkAreaGetControl(s.widget)
	C.gtk_widget_queue_draw(c)
}

func (s *sysData) center() {
	if C.gtk_widget_get_visible(s.widget) == C.FALSE {
		// hint to the WM to make it centered when it is shown again
		// thanks to Jasper in irc.gimp.net/#gtk+
		C.gtk_window_set_position(togtkwindow(s.widget), C.GTK_WIN_POS_CENTER)
	} else {
		var width, height C.gint

		s.resetposition()
		//we should be able to use gravity to simplify this, but it doesn't take effect immediately, and adding show calls does nothing (thanks Jasper in irc.gimp.net/#gtk+)
		C.gtk_window_get_size(togtkwindow(s.widget), &width, &height)
		C.gtk_window_move(togtkwindow(s.widget),
			(C.gdk_screen_width() / 2) - (width / 2),
			(C.gdk_screen_height() / 2) - (width / 2))
	}
}

func (s *sysData) setChecked(checked bool) {
	gtk_toggle_button_set_active(s.widget, checked)
}

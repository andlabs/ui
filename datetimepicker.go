// 12 december 2015

package ui

import (
	"time"
	"unsafe"
)

// #include <time.h>
// #include "ui.h"
// static inline struct tm *allocTimeStruct(void)
// {
// 	/* TODO handle error */
// 	return (struct tm *) malloc(sizeof (struct tm));
// }
// extern void doDateTimePickerChanged(uiDateTimePicker *, void *);
// static inline void realuiDateTimePickerOnChanged(uiDateTimePicker *d)
// {
// 	uiDateTimePickerOnChanged(d, doDateTimePickerOnChanged, NULL);
// }
import "C"

// DateTimePicker is a Control that represents a field where the user
// can enter a date and/or a time.
type DateTimePicker struct {
	c	*C.uiControl
	d	*C.uiDateTimePicker

	onChanged	func(*DateTimePicker)
}

// NewDateTimePicker creates a new DateTimePicker that shows
// both a date and a time.
func NewDateTimePicker() *DateTimePicker {
	d := new(DateTimePicker)

	d.d = C.uiNewDateTimePicker()
	d.c = (*C.uiControl)(unsafe.Pointer(d.d))

	C.realuiDateTimePickerOnChanged(d.d)

	return d
}

// NewDatePicker creates a new DateTimePicker that shows
// only a date.
func NewDatePicker() *DateTimePicker {
	d := new(DateTimePicker)

	d.d = C.uiNewDatePicker()
	d.c = (*C.uiControl)(unsafe.Pointer(d.d))

	C.realuiDateTimePickerOnChanged(d.d)

	return d
}

// NewTimePicker creates a new DateTimePicker that shows
// only a time.
func NewTimePicker() *DateTimePicker {
	d := new(DateTimePicker)

	d.d = C.uiNewTimePicker()
	d.c = (*C.uiControl)(unsafe.Pointer(d.d))

	C.realuiDateTimePickerOnChanged(d.d)

	return d
}

// Destroy destroys the DateTimePicker.
func (d *DateTimePicker) Destroy() {
	C.uiControlDestroy(d.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (d *DateTimePicker) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(d.c))
}

// Handle returns the OS-level handle associated with this DateTimePicker.
// On Windows this is an HWND of a standard Windows API
// DATETIMEPICK_CLASS class (as provided by Common Controls
// version 6).
// On GTK+ this is a pointer to a libui-internal class.
// On OS X this is a pointer to a NSDatePicker.
func (d *DateTimePicker) Handle() uintptr {
	return uintptr(C.uiControlHandle(d.c))
}

// Show shows the DateTimePicker.
func (d *DateTimePicker) Show() {
	C.uiControlShow(d.c)
}

// Hide hides the DateTimePicker.
func (d *DateTimePicker) Hide() {
	C.uiControlHide(d.c)
}

// Enable enables the DateTimePicker.
func (d *DateTimePicker) Enable() {
	C.uiControlEnable(d.c)
}

// Disable disables the DateTimePicker.
func (d *DateTimePicker) Disable() {
	C.uiControlDisable(d.c)
}

// Time returns the time stored in the uiDateTimePicker.
// The time is assumed to be local time.
func (d *DateTimePicker) Time() time.Time {
	tm := C.allocTimeStruct()
	defer C.free(unsafe.Pointer(tm))
	C.uiDateTimePickerTime(d.d, tm)
	return time.Date(
		int(tm.tm_year + 1900),
		time.Month(tm.tm_mon + 1),
		int(tm.tm_mday),
		int(tm.tm_hour),
		int(tm.tm_min),
		int(tm.tm_sec),
		0, time.Local)
}

// SetTime sets the time in the DateTimePicker to t.
// t's components are read as-is; no time zone manipulations
// are done.
func (d *DateTimePicker) SetTime(t time.Time) {
	tm := C.allocTimeStruct()
	defer C.free(unsafe.Pointer(tm))
	year, mon, mday := t.Date()
	tm.tm_year = C.int(year - 1900)
	tm.tm_mon = C.int(mon - 1)
	tm.tm_mday = C.int(mday)
	hour, min, sec := t.Time()
	tm.tm_hour = C.int(hour)
	tm.tm_min = C.int(min)
	tm.tm_sec = C.int(sec)
	tm.tm_isdst = -1
	C.uiDateTimePickerSetTime(d.d, tm)
}

// OnChanged registers f to be run when the user changes the time in the DateTimePicker.
// Only one function can be registered at a time.
func (d *DateTimePicker) OnChanged(f func(*DateTimePicker)) {
	d.onChanged = f
}

//export doDateTimePickerOnChanged
func doDateTimePickerOnChanged(dd *C.uiDateTimePicker, data unsafe.Pointer) {
	d := dateTimePickers[dd]
	if d.onChanged != nil {
		d.onChanged(d)
	}
}

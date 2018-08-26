// 12 december 2015

package ui

import (
	"time"
	"unsafe"
)

// #include "pkgui.h"
import "C"

// DateTimePicker is a Control that represents a field where the user
// can enter a date and/or a time.
type DateTimePicker struct {
	ControlBase
	d	*C.uiDateTimePicker
	onChanged	func(*DateTimePicker)
}

func finishNewDateTimePicker(dd *C.uiDateTimePicker) *DateTimePicker {
	d := new(DateTimePicker)

	d.d = dd

	C.pkguiDateTimePickerOnChanged(d.d)

	d.ControlBase = NewControlBase(d, uintptr(unsafe.Pointer(d.d)))
	return d
}

// NewDateTimePicker creates a new DateTimePicker that shows
// both a date and a time.
func NewDateTimePicker() *DateTimePicker {
	return finishNewDateTimePicker(C.uiNewDateTimePicker())
}

// NewDatePicker creates a new DateTimePicker that shows
// only a date.
func NewDatePicker() *DateTimePicker {
	return finishNewDateTimePicker(C.uiNewDatePicker())
}

// NewTimePicker creates a new DateTimePicker that shows
// only a time.
func NewTimePicker() *DateTimePicker {
	return finishNewDateTimePicker(C.uiNewTimePicker())
}

// Time returns the time stored in the uiDateTimePicker.
// The time is assumed to be local time.
func (d *DateTimePicker) Time() time.Time {
	tm := C.pkguiAllocTime()
	defer C.pkguiFreeTime(tm)
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
// t's components are read as-is via t.Date() and t.Clock();
// no time zone manipulations are done.
func (d *DateTimePicker) SetTime(t time.Time) {
	tm := C.pkguiAllocTime()
	defer C.pkguiFreeTime(tm)
	year, mon, mday := t.Date()
	tm.tm_year = C.int(year - 1900)
	tm.tm_mon = C.int(mon - 1)
	tm.tm_mday = C.int(mday)
	hour, min, sec := t.Clock()
	tm.tm_hour = C.int(hour)
	tm.tm_min = C.int(min)
	tm.tm_sec = C.int(sec)
	tm.tm_isdst = -1
	C.uiDateTimePickerSetTime(d.d, tm)
}

// OnChanged registers f to be run when the user changes the time
// in the DateTimePicker. Only one function can be registered at a
// time.
func (d *DateTimePicker) OnChanged(f func(*DateTimePicker)) {
	d.onChanged = f
}

//export pkguiDoDateTimePickerOnChanged
func pkguiDoDateTimePickerOnChanged(dd *C.uiDateTimePicker, data unsafe.Pointer) {
	d := ControlFromLibui(uintptr(unsafe.Pointer(dd))).(*DateTimePicker)
	if d.onChanged != nil {
		d.onChanged(d)
	}
}

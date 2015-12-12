// +build !windows,!darwin

// 4 november 2014

package ui

import (
	"fmt"
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type progressbar struct {
	*controlSingleWidget
	pbar		*C.GtkProgressBar
}

func newProgressBar() ProgressBar {
	widget := C.gtk_progress_bar_new();
	p := &progressbar{
		controlSingleWidget:	newControlSingleWidget(widget),
		pbar:				(*C.GtkProgressBar)(unsafe.Pointer(widget)),
	}
	return p
}

func (p *progressbar) Percent() int {
	return int(C.gtk_progress_bar_get_fraction(p.pbar) * 100)
}

func (p *progressbar) SetPercent(percent int) {
	if percent < 0 || percent > 100 {
		panic(fmt.Errorf("given ProgressBar percentage %d out of range", percent))
	}
	C.gtk_progress_bar_set_fraction(p.pbar, C.gdouble(percent) / 100)
}

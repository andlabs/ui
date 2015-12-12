// 4 november 2014

package ui

import (
	"fmt"
)

// #include "winapi_windows.h"
import "C"

type progressbar struct {
	*controlSingleHWND
}

func newProgressBar() ProgressBar {
	hwnd := C.newControl(C.xPROGRESS_CLASS,
		C.PBS_SMOOTH,
		0)
	p := &progressbar{
		controlSingleHWND:		newControlSingleHWND(hwnd),
	}
	p.fpreferredSize = p.xpreferredSize
	p.fnTabStops = func() int {
		// progress bars are not tab stops
		return 0
	}
	return p
}

func (p *progressbar) Percent() int {
	return int(C.SendMessageW(p.hwnd, C.PBM_GETPOS, 0, 0))
}

func (p *progressbar) SetPercent(percent int) {
	if percent < 0 || percent > 100 {
		panic(fmt.Errorf("given ProgressBar percentage %d out of range", percent))
	}
	// TODO circumvent aero
	C.SendMessageW(p.hwnd, C.PBM_SETPOS, C.WPARAM(percent), 0)
}

const (
	// via http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx
	// this is the double-width option
	progressbarWidth = 237
	progressbarHeight = 8
)

func (p *progressbar) xpreferredSize(d *sizing) (width, height int) {
	return fromdlgunitsX(progressbarWidth, d), fromdlgunitsY(progressbarHeight, d)
}

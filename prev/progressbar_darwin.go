// 4 november 2014

package ui

import (
	"fmt"
)

// #include "objc_darwin.h"
import "C"

type progressbar struct {
	*controlSingleObject
}

func newProgressBar() ProgressBar {
	return &progressbar{
		controlSingleObject:		newControlSingleObject(C.newProgressBar()),
	}
}

func (p *progressbar) Percent() int {
	return int(C.progressbarPercent(p.id))
}

func (p *progressbar) SetPercent(percent int) {
	if percent < 0 || percent > 100 {
		panic(fmt.Errorf("given ProgressBar percentage %d out of range", percent))
	}
	C.progressbarSetPercent(p.id, C.intmax_t(percent))
}

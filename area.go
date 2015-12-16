// 16 december 2015

package ui

// #include "ui.h"
import "C"

// no need to lock this; only the GUI thread can access it
var areas = make(map[*C.uiArea]*Area)

// TODO.
type Area struct {
	Control
}

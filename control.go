// 11 february 2014

//
package ui

import (
// ...
)

// A Control represents an UI control. Note that Control contains unexported members; this has the consequence that you can't build custom controls that interface directly with the system-specific code (fo rinstance, to import an unsupported control), or at least not without some hackery. If you want to make your own controls, embed Area and provide its necessities.
type Control interface {
	make(window *sysData) error
	setRect(x int, y int, width int, height int, winheight int) error
	preferredSize() (width int, height int, err error)
}

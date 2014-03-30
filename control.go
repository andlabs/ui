// 11 february 2014

package ui

import (
	// ...
)

// A Control represents an UI control. Note that Control contains unexported members; this has the consequence that you can't build custom controls that interface directly with the system-specific code (fo rinstance, to import an unsupported control), or at least not without some hackery. If you want to make your own controls, create an Area and provide an AreaHandler that does what you need.
type Control interface {
	make(window *sysData) error
	setRect(x int, y int, width int, height int, rr *[]resizerequest)
	preferredSize() (width int, height int)
}

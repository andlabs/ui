// 28 october 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

// interface builder notes
// - the tops of the alignment rects should be identical
// - spinner properties: auto repeat
// - http://stackoverflow.com/questions/702829/integrate-nsstepper-with-nstextfield we'll need to bind the int value :S
// 	- TODO experiment with a dummy project
// - http://juliuspaintings.co.uk/cgi-bin/paint_css/animatedPaint/059-NSStepper-NSTextField.pl
// - http://www.youtube.com/watch?v=ZZSHU-O7HVo
// - http://andrehoffmann.wordpress.com/tag/nsstepper/ ?
// TODO
// - proper spacing between edit and spinner: Interface Builder isn't clear; NSDatePicker doesn't spill the beans

type spinbox struct {
	textfield		C.id
	stepper		C.id
	objcspinbox	C.id
}

func newSpinbox() Spinbox {
	s := new(spinbox)
	s.textfield = C.newTextField()
	s.stepper = C.newSpinboxStepper()
	s.objcspinbox = C.spinboxSetup(s.textfield, s.stepper, unsafe.Pointer(s))
	return s
}

func (s *spinbox) setParent(p *controlParent) {
	C.parent(s.textfield, p.id)
	C.parent(s.stepper, p.id)
}

func (s *spinbox) preferredSize(d *sizing) (width, height int) {
	// TODO
	return 20, 20
}

func (s *spinbox) resize(x int, y int, width int, height int, d *sizing) {
	// TODO
	C.moveControl(s.textfield, C.intptr_t(x), C.intptr_t(y), C.intptr_t(width - 20), C.intptr_t(height))
	C.moveControl(s.stepper, C.intptr_t(x + width - 15), C.intptr_t(y), C.intptr_t(15), C.intptr_t(height))
}

func (s *spinbox) nTabStops() int {
	// TODO does the stepper count?
	return 1
}

func (s *spinbox) containerShow() {
	// only provided for the Windows backend
}

func (s *spinbox) containerHide() {
	// only provided for the Windows backend
}

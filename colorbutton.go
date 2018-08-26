// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// ColorButton is a Control that represents a button that the user can
// click to select a color.
type ColorButton struct {
	ControlBase
	b	*C.uiColorButton
	onChanged		func(*ColorButton)
}

// NewColorButton creates a new ColorButton.
func NewColorButton() *ColorButton {
	b := new(ColorButton)

	b.b = C.uiNewColorButton()

	C.pkguiColorButtonOnChanged(b.b)

	b.ControlBase = NewControlBase(b, uintptr(unsafe.Pointer(b.b)))
	return b
}

// Color returns the color currently selected in the ColorButton.
// Colors are not alpha-premultiplied.
// TODO rename b or bl
func (b *ColorButton) Color() (r, g, bl, a float64) {
	c := C.pkguiAllocColorDoubles()
	defer C.pkguiFreeColorDoubles(c)
	C.uiColorButtonColor(b.b, c.r, c.g, c.b, c.a)
	return float64(*(c.r)), float64(*(c.g)), float64(*(c.b)), float64(*(c.a))
}

// SetColor sets the currently selected color in the ColorButton.
// Colors are not alpha-premultiplied.
// TODO rename b or bl
func (b *ColorButton) SetColor(r, g, bl, a float64) {
	C.uiColorButtonSetColor(b.b, C.double(r), C.double(g), C.double(bl), C.double(a))
}

// OnChanged registers f to be run when the user changes the
// currently selected color in the ColorButton. Only one function
// can be registered at a time.
func (b *ColorButton) OnChanged(f func(*ColorButton)) {
	b.onChanged = f
}

//export pkguiDoColorButtonOnChanged
func pkguiDoColorButtonOnChanged(bb *C.uiColorButton, data unsafe.Pointer) {
	b := ControlFromLibui(uintptr(unsafe.Pointer(bb))).(*ColorButton)
	if b.onChanged != nil {
		b.onChanged(b)
	}
}

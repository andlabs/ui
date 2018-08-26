// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// FontButton is a Control that represents a button that the user can
// click to select a font.
type FontButton struct {
	ControlBase
	b	*C.uiFontButton
	onChanged		func(*FontButton)
}

// NewFontButton creates a new FontButton.
func NewFontButton() *FontButton {
	b := new(FontButton)

	b.b = C.uiNewFontButton()

	C.pkguiFontButtonOnChanged(b.b)

	b.ControlBase = NewControlBase(b, uintptr(unsafe.Pointer(b.b)))
	return b
}

// Font returns the font currently selected in the FontButton.
func (b *FontButton) Font() *FontDescriptor {
	cfd := C.pkguiNewFontDescriptor()
	defer C.pkguiFreeFontDescriptor(cfd)
	C.uiFontButtonFont(b.b, cfd)
	defer C.uiFreeFontButtonFont(cfd)
	fd := &FontDescriptor{}
	fd.fromLibui(cfd)
	return fd
}

// OnChanged registers f to be run when the user changes the
// currently selected font in the FontButton. Only one function can
// be registered at a time.
func (b *FontButton) OnChanged(f func(*FontButton)) {
	b.onChanged = f
}

//export pkguiDoFontButtonOnChanged
func pkguiDoFontButtonOnChanged(bb *C.uiFontButton, data unsafe.Pointer) {
	b := ControlFromLibui(uintptr(unsafe.Pointer(bb))).(*FontButton)
	if b.onChanged != nil {
		b.onChanged(b)
	}
}

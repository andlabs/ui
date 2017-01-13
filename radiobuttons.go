// 12 december 2015

package ui

import "unsafe"

// #include "ui.h"
// extern void doRadioButtonsOnSelected(uiRadioButtons *, void *);
// static inline void realuiRadioButtonsOnSelected(uiRadioButtons *r)
// {
// 	uiRadioButtonsOnSelected(r, doRadioButtonsOnSelected, NULL);
// }
import "C"

// no need to lock this; only the GUI thread can access it
var radioButtons = make(map[*C.uiRadioButtons]*RadioButtons)

// RadioButtons is a Control that represents a set of checkable
// buttons from which exactly one may be chosen by the user.
type RadioButtons struct {
	c *C.uiControl
	r *C.uiRadioButtons

	onSelected func(*RadioButtons)
}

// NewRadioButtons creates a new RadioButtons.
func NewRadioButtons() *RadioButtons {
	r := new(RadioButtons)

	r.r = C.uiNewRadioButtons()
	r.c = (*C.uiControl)(unsafe.Pointer(r.r))

	C.realuiRadioButtonsOnSelected(r.r)
	radioButtons[r.r] = r

	return r
}

// Destroy destroys the RadioButtons.
func (r *RadioButtons) Destroy() {
	delete(radioButtons, r.r)
	C.uiControlDestroy(r.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Window. This is only used by package ui itself and should
// not be called by programs.
func (r *RadioButtons) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(r.c))
}

// Handle returns the OS-level handle associated with this RadioButtons.
// On Windows this is an HWND of a libui-internal class; its
// child windows are instances of the standard Windows API
// BUTTON class (as provided by Common Controls version 6).
// On GTK+ this is a pointer to a GtkBox containing GtkRadioButtons.
// On OS X this is a pointer to a NSView with each radio button as a NSButton subview.
func (r *RadioButtons) Handle() uintptr {
	return uintptr(C.uiControlHandle(r.c))
}

// Show shows the RadioButtons.
func (r *RadioButtons) Show() {
	C.uiControlShow(r.c)
}

// Hide hides the RadioButtons.
func (r *RadioButtons) Hide() {
	C.uiControlHide(r.c)
}

// Enable enables the RadioButtons.
func (r *RadioButtons) Enable() {
	C.uiControlEnable(r.c)
}

// Disable disables the RadioButtons.
func (r *RadioButtons) Disable() {
	C.uiControlDisable(r.c)
}

// Append adds the named button to the end of the RadioButtons.
// If this button is the first button, it is automatically selected.
func (r *RadioButtons) Append(text string) {
	ctext := C.CString(text)
	C.uiRadioButtonsAppend(r.r, ctext)
	freestr(ctext)
}

// Selected returns the index of the currently selected item in the
// RadioButtons.
func (r *RadioButtons) Selected() int {
	return int(C.uiRadioButtonsSelected(r.r))
}

// SetSelected sets the currently select item in the RadioButtons
// to index.
func (r *RadioButtons) SetSelected(index int) {
	C.uiRadioButtonsSetSelected(r.r, C.int(index))
}

// OnSelected registers f to be run when the user selects an item in
// the RadioButtons. Only one function can be registered at a time.
func (r *RadioButtons) OnSelected(f func(*RadioButtons)) {
	r.onSelected = f
}

//export doRadioButtonsOnSelected
func doRadioButtonsOnSelected(rr *C.uiRadioButtons, data unsafe.Pointer) {
	r := radioButtons[rr]
	if r.onSelected != nil {
		r.onSelected(r)
	}
}

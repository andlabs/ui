// 16 july 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type button struct {
	*controlbase
	clicked		*event
}

func finishNewButton(id C.id, text string) *button {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	b := &button{
		controlbase:	newControl(id),
		clicked:		newEvent(),
	}
	C.buttonSetText(b.id, ctext)
	C.buttonSetDelegate(b.id, unsafe.Pointer(b))
	return b
}

func newButton(text string) *button {
	return finishNewButton(C.newButton(), text)
}

func (b *button) OnClicked(e func()) {
	b.clicked.set(e)
}

//export buttonClicked
func buttonClicked(xb unsafe.Pointer) {
	b := (*button)(unsafe.Pointer(xb))
	b.clicked.fire()
	println("button clicked")
}

func (b *button) Text() string {
	return C.GoString(C.buttonText(b.id))
}

func (b *button) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.buttonSetText(b.id, ctext)
}

type checkbox struct {
	*button
}

func newCheckbox(text string) *checkbox {
	return &checkbox{
		button:	finishNewButton(C.newCheckbox(), text),
	}
}

// we don't need to define our own event here; we can just reuse Button's
// (it's all target-action anyway)

func (c *checkbox) Checked() bool {
	return fromBOOL(C.checkboxChecked(c.id))
}

func (c *checkbox) SetChecked(checked bool) {
	C.checkboxSetChecked(c.id, toBOOL(checked))
}

type textField struct {
	*controlbase
}

func finishNewTextField(id C.id) *textField {
	return &textField{
		controlbase:	newControl(id),
	}
}

func newTextField() *textField {
	return finishNewTextField(C.newTextField())
}

func newPasswordField() *textField {
	return finishNewTextField(C.newPasswordField())
}

func (t *textField) Text() string {
	return C.GoString(C.textFieldText(t.id))
}

func (t *textField) SetText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.textFieldSetText(t.id, ctext)
}

// cheap trick
type label struct {
	*textField
	standalone			bool
	supercommitResize		func(c *allocation, d *sizing)
}

func finishNewLabel(text string, standalone bool) *label {
	l := &label{
		textField:		finishNewTextField(C.newLabel()),
		standalone:	standalone,
	}
	l.SetText(text)
	l.supercommitResize = l.fcommitResize
	l.fcommitResize = l.labelcommitResize
	return l
}

func newLabel(text string) Label {
	return finishNewLabel(text, false)
}

func newStandaloneLabel(text string) Label {
	return finishNewLabel(text, true)
}

func (l *label) labelcommitResize(c *allocation, d *sizing) {
	if !l.standalone && c.neighbor != nil {
		c.neighbor.getAuxResizeInfo(d)
		if d.neighborAlign.baseline != 0 {		// no adjustment needed if the given control has no baseline
			// in order for the baseline value to be correct, the label MUST BE AT THE HEIGHT THAT OS X WANTS IT TO BE!
			// otherwise, the baseline calculation will be relative to the bottom of the control, and everything will be wrong
			origsize := C.controlPrefSize(l.id)
			c.height = int(origsize.height)
			newrect := C.struct_xrect{
				x:		C.intptr_t(c.x),
				y:		C.intptr_t(c.y),
				width:	C.intptr_t(c.width),
				height:	C.intptr_t(c.height),
			}
			ourAlign := C.alignmentInfo(l.id, newrect)
			// we need to find the exact Y positions of the baselines
			// fortunately, this is easy now that (x,y) is the bottom-left corner
			thisbasey := ourAlign.rect.y + ourAlign.baseline
			neighborbasey := d.neighborAlign.rect.y + d.neighborAlign.baseline
			// now the amount we have to move the label down by is easy to find
			yoff := neighborbasey - thisbasey
			// and we just add that
			c.y += int(yoff)
		}
		// TODO if there's no baseline, the alignment should be to the top /of the alignment rect/, not the frame
	}
	l.supercommitResize(c, d)
}

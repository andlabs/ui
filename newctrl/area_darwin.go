// 29 march 2014

package ui

import (
	"fmt"
	"image"
	"unsafe"
)

//// #include <HIToolbox/Events.h>
// #include "objc_darwin.h"
import "C"

type area struct {
	*areabase

	*scroller
	textfield     C.id
	textfielddone *event
}

func newArea(ab *areabase) Area {
	a := &area{
		areabase:      ab,
		textfielddone: newEvent(),
	}
	id := C.newArea(unsafe.Pointer(a))
	a.scroller = newScroller(id, false) // no border on Area
	a.fpreferredSize = a.xpreferredSize
	a.SetSize(a.width, a.height)
	a.textfield = C.newTextField()
	C.areaSetTextField(a.id, a.textfield)
	return a
}

func (a *area) SetSize(width, height int) {
	a.width = width
	a.height = height
	// set the frame size to set the area's effective size on the Cocoa side
	C.moveControl(a.id, 0, 0, C.intptr_t(a.width), C.intptr_t(a.height))
}

func (a *area) Repaint(r image.Rectangle) {
	var s C.struct_xrect

	r = image.Rect(0, 0, a.width, a.height).Intersect(r)
	if r.Empty() {
		return
	}
	s.x = C.intptr_t(r.Min.X)
	s.y = C.intptr_t(r.Min.Y)
	s.width = C.intptr_t(r.Dx())
	s.height = C.intptr_t(r.Dy())
	C.areaRepaint(a.id, s)
}

func (a *area) RepaintAll() {
	C.areaRepaintAll(a.id)
}

func (a *area) OpenTextFieldAt(x, y int) {
	if x < 0 || x >= a.width || y < 0 || y >= a.height {
		panic(fmt.Errorf("point (%d,%d) outside Area in Area.OpenTextFieldAt()", x, y))
	}
	C.areaTextFieldOpen(a.id, a.textfield, C.intptr_t(x), C.intptr_t(y))
}

func (a *area) TextFieldText() string {
	return C.GoString(C.textfieldText(a.textfield))
}

func (a *area) SetTextFieldText(text string) {
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	C.textfieldSetText(a.textfield, ctext)
}

func (a *area) OnTextFieldDismissed(f func()) {
	a.textfielddone.set(f)
}

//export areaTextFieldDismissed
func areaTextFieldDismissed(data unsafe.Pointer) {
	a := (*area)(unsafe.Pointer(data))
	C.controlSetHidden(a.textfield, C.YES)
	a.textfielddone.fire()
}

//export areaView_drawRect
func areaView_drawRect(self C.id, rect C.struct_xrect, data unsafe.Pointer) {
	a := (*area)(data)
	// no need to clear the clip rect; the NSScrollView does that for us (see the setDrawsBackground: call in objc_darwin.m)
	// rectangles in Cocoa are origin/size, not point0/point1; if we don't watch for this, weird things will happen when scrolling
	cliprect := image.Rect(int(rect.x), int(rect.y), int(rect.x+rect.width), int(rect.y+rect.height))
	cliprect = image.Rect(0, 0, int(a.width), int(a.height)).Intersect(cliprect)
	if cliprect.Empty() { // no intersection; nothing to paint
		return
	}
	i := a.handler.Paint(cliprect)
	success := C.drawImage(
		unsafe.Pointer(pixelData(i)), C.intptr_t(i.Rect.Dx()), C.intptr_t(i.Rect.Dy()), C.intptr_t(i.Stride),
		C.intptr_t(cliprect.Min.X), C.intptr_t(cliprect.Min.Y))
	if success == C.NO {
		panic("error drawing into Area (exactly what is unknown)")
	}
}

func parseModifiers(e C.id) (m Modifiers) {
	mods := C.modifierFlags(e)
	if (mods & C.cNSControlKeyMask) != 0 {
		m |= Ctrl
	}
	if (mods & C.cNSAlternateKeyMask) != 0 {
		m |= Alt
	}
	if (mods & C.cNSShiftKeyMask) != 0 {
		m |= Shift
	}
	if (mods & C.cNSCommandKeyMask) != 0 {
		m |= Super
	}
	return m
}

func areaMouseEvent(self C.id, e C.id, click bool, up bool, data unsafe.Pointer) {
	var me MouseEvent

	a := (*area)(data)
	xp := C.getTranslatedEventPoint(self, e)
	me.Pos = image.Pt(int(xp.x), int(xp.y))
	// for the most part, Cocoa won't geenerate an event outside the Area... except when dragging outside the Area, so check for this
	if !me.Pos.In(image.Rect(0, 0, int(a.width), int(a.height))) {
		return
	}
	me.Modifiers = parseModifiers(e)
	which := uint(C.buttonNumber(e)) + 1
	if which == 3 { // swap middle and right button numbers
		which = 2
	} else if which == 2 {
		which = 3
	}
	if click && up {
		me.Up = which
	} else if click {
		me.Down = which
		// this already works the way we want it to so nothing special needed like with Windows and GTK+
		me.Count = uint(C.clickCount(e))
	} else {
		which = 0 // reset for Held processing below
	}
	// the docs do say don't use this for tracking (mouseMoved:) since it returns the state now, and mouse move events work by tracking, but as far as I can tell dragging the mouse over the inactive window does not generate an event on Mac OS X, so :/ (tracking doesn't touch dragging anyway except during mouseEntered: and mouseExited:, which we don't handle, and the only other tracking message, cursorChanged:, we also don't handle (yet...? need to figure out if this is how to set custom cursors or not), so)
	held := C.pressedMouseButtons()
	if which != 1 && (held&1) != 0 { // button 1
		me.Held = append(me.Held, 1)
	}
	if which != 2 && (held&4) != 0 { // button 2; mind the swap
		me.Held = append(me.Held, 2)
	}
	if which != 3 && (held&2) != 0 { // button 3
		me.Held = append(me.Held, 3)
	}
	held >>= 3
	for i := uint(4); held != 0; i++ {
		if which != i && (held&1) != 0 {
			me.Held = append(me.Held, i)
		}
		held >>= 1
	}
	a.handler.Mouse(me)
}

//export areaView_mouseMoved_mouseDragged
func areaView_mouseMoved_mouseDragged(self C.id, e C.id, data unsafe.Pointer) {
	// for moving, this is handled by the tracking rect stuff above
	// for dragging, if multiple buttons are held, only one of their xxxMouseDragged: messages will be sent, so this is OK to do
	areaMouseEvent(self, e, false, false, data)
}

//export areaView_mouseDown
func areaView_mouseDown(self C.id, e C.id, data unsafe.Pointer) {
	// no need to manually set focus; Mac OS X has already done that for us by this point since we set our view to be a first responder
	areaMouseEvent(self, e, true, false, data)
}

//export areaView_mouseUp
func areaView_mouseUp(self C.id, e C.id, data unsafe.Pointer) {
	areaMouseEvent(self, e, true, true, data)
}

func sendKeyEvent(self C.id, ke KeyEvent, data unsafe.Pointer) C.BOOL {
	a := (*area)(data)
	handled := a.handler.Key(ke)
	return toBOOL(handled)
}

func areaKeyEvent(self C.id, e C.id, up bool, data unsafe.Pointer) C.BOOL {
	var ke KeyEvent

	keyCode := uintptr(C.keyCode(e))
	ke, ok := fromKeycode(keyCode)
	if !ok {
		// no such key; modifiers by themselves are handled by -[self flagsChanged:]
		return C.NO
	}
	// either ke.Key or ke.ExtKey will be set at this point
	ke.Modifiers = parseModifiers(e)
	ke.Up = up
	return sendKeyEvent(self, ke, data)
}

//export areaView_keyDown
func areaView_keyDown(self C.id, e C.id, data unsafe.Pointer) C.BOOL {
	return areaKeyEvent(self, e, false, data)
}

//export areaView_keyUp
func areaView_keyUp(self C.id, e C.id, data unsafe.Pointer) C.BOOL {
	return areaKeyEvent(self, e, true, data)
}

//export areaView_flagsChanged
func areaView_flagsChanged(self C.id, e C.id, data unsafe.Pointer) C.BOOL {
	var ke KeyEvent

	// Mac OS X sends this event on both key up and key down.
	// Fortunately -[e keyCode] IS valid here, so we can simply map from key code to Modifiers, get the value of [e modifierFlags], and check if the respective bit is set or not — that will give us the up/down state
	keyCode := uintptr(C.keyCode(e))
	mod, ok := keycodeModifiers[keyCode] // comma-ok form to avoid adding entries
	if !ok {                             // unknown modifier; ignore
		return C.NO
	}
	ke.Modifiers = parseModifiers(e)
	ke.Up = (ke.Modifiers & mod) == 0
	ke.Modifier = mod
	// don't include the modifier in ke.Modifiers
	ke.Modifiers &^= mod
	return sendKeyEvent(self, ke, data)
}

func (a *area) xpreferredSize(d *sizing) (width, height int) {
	// the preferred size of an Area is its size
	return a.width, a.height
}

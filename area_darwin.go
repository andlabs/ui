// 29 march 2014

package ui

import (
	"unsafe"
	"image"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
//// #include <HIToolbox/Events.h>
// #include "objc_darwin.h"
import "C"

func makeArea(parentWindow C.id, alternate bool, s *sysData) C.id {
	area := C.makeArea()
	area = makeScrollView(area)
	addControl(parentWindow, area)
	return area
}

func areaInScrollView(scrollview C.id) C.id {
	return getScrollViewContent(scrollview)
}

//export areaView_drawRect
func areaView_drawRect(self C.id, rect C.struct_xrect) {
	s := getSysData(self)
	// TODO clear clip rect
	// rectangles in Cocoa are origin/size, not point0/point1; if we don't watch for this, weird things will happen when scrolling
	// TODO change names EVERYWHERE ELSE to match
	cliprect := image.Rect(int(rect.x), int(rect.y), int(rect.x + rect.width), int(rect.y + rect.height))
	max := C.frame(self)
	cliprect = image.Rect(0, 0, int(max.width), int(max.height)).Intersect(cliprect)
	if cliprect.Empty() {			// no intersection; nothing to paint
		return
	}
	i := s.handler.Paint(cliprect)
	C.drawImage(
		unsafe.Pointer(pixelData(i)), C.intptr_t(i.Rect.Dx()), C.intptr_t(i.Rect.Dy()), C.intptr_t(i.Stride),
		C.intptr_t(cliprect.Min.X), C.intptr_t(cliprect.Min.Y))
}

func parseModifiers(e C.id) (m Modifiers) {
	const (
		_NSShiftKeyMask = 1 << 17
		_NSControlKeyMask = 1 << 18
		_NSAlternateKeyMask = 1 << 19
		_NSCommandKeyMask = 1 << 20
	)

	mods := uintptr(C.modifierFlags(e))
	if (mods & _NSControlKeyMask) != 0 {
		m |= Ctrl
	}
	if (mods & _NSAlternateKeyMask) != 0 {
		m |= Alt
	}
	if (mods & _NSShiftKeyMask) != 0 {
		m |= Shift
	}
	if (mods & _NSCommandKeyMask) != 0 {
		m |= Super
	}
	return m
}

// TODO pressing both buttons 1 and 3 simultaneously gets turned into button 2; see if we can turn that off for our NSView only
func areaMouseEvent(self C.id, e C.id, click bool, up bool) {
	var me MouseEvent

	s := getSysData(self)
	xp := C.getTranslatedEventPoint(self, e)
	me.Pos = image.Pt(int(xp.x), int(xp.y))
	// for the most part, Cocoa won't geenerate an event outside the Area... except when dragging outside the Area, so check for this
	max := C.frame(self)
	if !me.Pos.In(image.Rect(0, 0, int(max.width), int(max.height))) {
		return
	}
	me.Modifiers = parseModifiers(e)
	which := uint(C.buttonNumber(e)) + 1
	if which == 3 {		// swap middle and right button numbers
		which = 2
	} else if which == 2 {
		which = 3
	}
	if click && up {
		me.Up = which
	} else if click {
		me.Down = which
		me.Count = uint(C.clickCount(e))
	} else {
		which = 0			// reset for Held processing below
	}
	// the docs do say don't use this for tracking (mouseMoved:) since it returns the state now, and mouse move events work by tracking, but as far as I can tell dragging the mouse over the inactive window does not generate an event on Mac OS X, so :/ (tracking doesn't touch dragging anyway so)
	held := C.pressedMouseButtons()
	if which != 1 && (held & 1) != 0 {		// button 1
		me.Held = append(me.Held, 1)
	}
	if which != 2 && (held & 4) != 0 {		// button 2; mind the swap
		me.Held = append(me.Held, 2)
	}
	if which != 3 && (held & 2) != 0 {		// button 3
		me.Held = append(me.Held, 3)
	}
	held >>= 3
	for i := uint(4); held != 0; i++ {
		if which != i && (held & 1) != 0 {
			me.Held = append(me.Held, i)
		}
		held >>= 1
	}
	repaint := s.handler.Mouse(me)
	if repaint {
		C.display(self)
	}
}

//export areaView_mouseMoved_mouseDragged
func areaView_mouseMoved_mouseDragged(self C.id, e C.id) {
	// for moving, this is handled by the tracking rect stuff above
	// for dragging, if multiple buttons are held, only one of their xxxMouseDragged: messages will be sent, so this is OK to do
	areaMouseEvent(self, e, false, false)
}

//export areaView_mouseDown
func areaView_mouseDown(self C.id, e C.id) {
	areaMouseEvent(self, e, true, false)
}

//export areaView_mouseUp
func areaView_mouseUp(self C.id, e C.id) {
	areaMouseEvent(self, e, true, true)
}

func sendKeyEvent(self C.id, e C.id, ke KeyEvent) bool {
	s := getSysData(self)
	handled, repaint := s.handler.Key(ke)
	if repaint {
		C.display(self)
	}
	return handled
}

func areaKeyEvent(self C.id, e C.id, up bool) bool {
	var ke KeyEvent

	keyCode := uintptr(C.keyCode(e))
	ke, ok := fromKeycode(keyCode)
	if !ok {
		// no such key; modifiers by themselves are handled by -[self flagsChanged:]
		return false
	}
	// either ke.Key or ke.ExtKey will be set at this point
	ke.Modifiers = parseModifiers(e)
	ke.Up = up
	return sendKeyEvent(self, e, ke)
}

//export areaView_keyDown
func areaView_keyDown(self C.id, e C.id) C.BOOL {
	return toBOOL(areaKeyEvent(self, e, false))
}

//export areaView_keyUp
func areaView_keyUp(self C.id, e C.id) C.BOOL {
	return toBOOL(areaKeyEvent(self, e, true))
}

//export areaView_flagsChanged
func areaView_flagsChanged(self C.id, e C.id) C.BOOL {
	var ke KeyEvent

	// Mac OS X sends this event on both key up and key down.
	// Fortunately -[e keyCode] IS valid here, so we can simply map from key code to Modifiers, get the value of [e modifierFlags], the respective bit is set or not — that will give us the up/down state
	keyCode := uintptr(C.keyCode(e))
	mod, ok := keycodeModifiers[keyCode]		// comma-ok form to avoid adding entries
	if !ok {		// unknown modifier; ignore
		return C.NO
	}
	ke.Modifiers = parseModifiers(e)
	ke.Up = (ke.Modifiers & mod) == 0
	return toBOOL(sendKeyEvent(self, e, ke))
}

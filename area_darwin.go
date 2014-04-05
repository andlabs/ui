// 29 march 2014

package ui

import (
	"fmt"
	"unsafe"
	"image"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include <stdlib.h>
//// #include <HIToolbox/Events.h>
// #include "objc_darwin.h"
// extern void areaView_drawRect(id, struct xrect);
// extern BOOL areaView_isFlipped_acceptsFirstResponder(id, SEL);
// extern void areaView_mouseMoved(id, SEL, id);
// extern void areaView_mouseDown_mouseDragged(id, SEL, id);
// extern void areaView_mouseUp(id, SEL, id);
// extern void areaView_keyDown(id, SEL, id);
// extern void areaView_keyUp(id, SEL, id);
// extern void areaView_flagsChanged(id, SEL, id);
import "C"

const (
	__goArea = "goArea"
)

var (
	_goArea C.id
	_NSView = objc_getClass("NSView")

	_drawRect = sel_getUid("drawRect:")
)

var goAreaSels = []selector{
	selector{"isFlipped", uintptr(C.areaView_isFlipped_acceptsFirstResponder), sel_bool,
		"ensuring that an Area's coordinate system has (0,0) at the top-left corner"},
	selector{"acceptsFirstResponder", uintptr(C.areaView_isFlipped_acceptsFirstResponder), sel_bool,
		"registering that Areas are to receive events"},
	selector{"mouseMoved:", uintptr(C.areaView_mouseMoved), sel_void_id,
		"handling mouse movements in Area"},
	selector{"mouseDown:", uintptr(C.areaView_mouseDown_mouseDragged), sel_void_id,
		"handling mouse button 1 presses in Area"},
	selector{"mouseDragged:", uintptr(C.areaView_mouseDown_mouseDragged), sel_void_id,
		"handling mouse button 1 dragging in Area"},
	selector{"mouseUp:", uintptr(C.areaView_mouseUp), sel_void_id,
		"handling mouse button 1 releases in Area"},
	selector{"rightMouseDown:", uintptr(C.areaView_mouseDown_mouseDragged), sel_void_id,
		"handling mouse button 3 presses in Area"},
	selector{"rightMouseDragged:", uintptr(C.areaView_mouseDown_mouseDragged), sel_void_id,
		"handling mouse button 3 dragging in Area"},
	selector{"rightMouseUp:", uintptr(C.areaView_mouseUp), sel_void_id,
		"handling mouse button 3 releases in Area"},
	selector{"otherMouseDown:", uintptr(C.areaView_mouseDown_mouseDragged), sel_void_id,
		"handling mouse button 2 (and 4 and higher) presses in Area"},
	selector{"otherMouseDragged:", uintptr(C.areaView_mouseDown_mouseDragged), sel_void_id,
		"handling mouse button 2 (and 4 and higher) dragging in Area"},
	selector{"otherMouseUp:", uintptr(C.areaView_mouseUp), sel_void_id,
		"handling mouse button 2 (and 4 and higher) releases in Area"},
	selector{"keyDown:", uintptr(C.areaView_keyDown), sel_void_id,
		"handling key presses in Area"},
	selector{"keyUp:", uintptr(C.areaView_keyUp), sel_void_id,
		"handling key releases in Area"},
	selector{"flagsChanged:", uintptr(C.areaView_flagsChanged), sel_void_id,
		"handling Modifiers presses and releases in Area"},
}

func mkAreaClass() error {
	id, class, err := makeClass(__goArea, _NSView, goAreaSels,
		"the implementation of Area on Mac OS X")
	if err != nil {
		return err
	}
	// addAreaViewDrawMethod() is in bleh_darwin.m
	ok := C.addAreaViewDrawMethod(class)
	if ok != C.BOOL(C.YES) {
		return fmt.Errorf("unable to add selector drawRect: to class %s (needed for actually drawing Areas; reason unknown)", __goArea)
	}
	_goArea = id
	return nil
}

var (
	_drawAtPoint = sel_getUid("drawAtPoint:")
)

//export areaView_drawRect
func areaView_drawRect(self C.id, rect C.struct_xrect) {
	s := getSysData(self)
	// TODO clear clip rect
	// rectangles in Cocoa are origin/size, not point0/point1; if we don't watch for this, weird things will happen when scrolling
	// TODO change names EVERYWHERE ELSE to match
	cliprect := image.Rect(int(rect.x), int(rect.y), int(rect.x + rect.width), int(rect.y + rect.height))
	max := C.objc_msgSend_stret_rect_noargs(self, _frame)
	cliprect = image.Rect(0, 0, int(max.width), int(max.height)).Intersect(cliprect)
	if cliprect.Empty() {			// no intersection; nothing to paint
		return
	}
	i := s.handler.Paint(cliprect)
	C.drawImage(
		unsafe.Pointer(&i.Pix[0]), C.int64_t(i.Rect.Dx()), C.int64_t(i.Rect.Dy()), C.int64_t(i.Stride),
		C.int64_t(cliprect.Min.X), C.int64_t(cliprect.Min.Y))
}

//export areaView_isFlipped_acceptsFirstResponder
func areaView_isFlipped_acceptsFirstResponder(self C.id, sel C.SEL) C.BOOL {
	// yes use the same function for both methods since they're just going to return YES anyway
	// isFlipped gives us a coordinate system with (0,0) at the top-left
	// acceptsFirstResponder lets us respond to events
	return C.BOOL(C.YES)
}

var (
	_NSEvent = objc_getClass("NSEvent")

	_modifierFlags = sel_getUid("modifierFlags")
	_buttonNumber = sel_getUid("buttonNumber")
	_clickCount = sel_getUid("clickCount")
	_pressedMouseButtons = sel_getUid("pressedMouseButtons")
)

func parseModifiers(e C.id) (m Modifiers) {
	const (
		_NSShiftKeyMask = 1 << 17
		_NSControlKeyMask = 1 << 18
		_NSAlternateKeyMask = 1 << 19
		_NSCommandKeyMask = 1 << 20
	)

	mods := uintptr(C.objc_msgSend_uintret_noargs(e, _modifierFlags))
	if (mods & _NSShiftKeyMask) != 0 {
		m |= Shift
	}
	if (mods & _NSControlKeyMask) != 0 {
		// TODO
	}
	if (mods & _NSAlternateKeyMask) != 0 {
		m |= Alt
	}
	if (mods & _NSCommandKeyMask) != 0 {
		m |= Ctrl
	}
	return m
}

func areaMouseEvent(self C.id, e C.id, click bool, up bool) {
	var me MouseEvent

	s := getSysData(self)
	xp := C.getTranslatedEventPoint(self, e)
	me.Pos = image.Pt(int(xp.x), int(xp.y))
	me.Modifiers = parseModifiers(e)
	which := uint(C.objc_msgSend_intret_noargs(e, _buttonNumber)) + 1
	if which == 3 {		// swap middle and right button numbers
		which = 2
	} else if which == 2 {
		which = 3
	}
	if click && up {
		me.Up = which
	} else if click {
		me.Down = which
		me.Count = uint(C.objc_msgSend_intret_noargs(e, _clickCount))
	} else {
		which = 0			// reset for Held processing below
	}
	// pressedMouseButtons is a class method; calling objc_msgSend() directly with e as an argument on these will panic (in real Objective-C the compiler can detect [e pressedMouseButtons])
	held := C.objc_msgSend_uintret_noargs(_NSEvent, _pressedMouseButtons)
	if which != 1 && (held & 1) != 0 {		// button 1
		me.Held = append(me.Held, 1)
	}
	if which != 2 && (held & 4) != 0 {		// button 2; mind the swap
		me.Held = append(me.Held, 2)
	}
	if which != 3 && (held & 2) != 0 {		// button 3
		me.Held = append(me.Held, 3)
	}
	// TODO remove this part?
	held >>= 3
	for i := uint(4); held != 0; i++ {
		if which != i && (held & 1) != 0 {
			me.Held = append(me.Held, i)
		}
		held >>= 1
	}
	repaint := s.handler.Mouse(me)
	if repaint {
		C.objc_msgSend_noargs(self, _display)
	}
}

//export areaView_mouseMoved
func areaView_mouseMoved(self C.id, sel C.SEL, e C.id) {
	// TODO not triggered?
	areaMouseEvent(self, e, false, false)
}

//export areaView_mouseDown_mouseDragged
func areaView_mouseDown_mouseDragged(self C.id, sel C.SEL, e C.id) {
	areaMouseEvent(self, e, true, false)
}

//export areaView_mouseUp
func areaView_mouseUp(self C.id, sel C.SEL, e C.id) {
	areaMouseEvent(self, e, true, true)
}

var (
	_keyCode = sel_getUid("keyCode")
)

func handleKeyEvent(self C.id) {
	// TODO
}

func sendKeyEvent(self C.id, ke KeyEvent) {
	s := getSysData(self)
	handled, repaint := s.handler.Key(ke)
	if repaint {
		C.objc_msgSend_noargs(self, _display)
	}
	if !handled {
		handleKeyEvent(self)
	}
}

func areaKeyEvent(self C.id, e C.id, up bool) {
	var ke KeyEvent

	keyCode := uintptr(C.objc_msgSend_ushortret_noargs(e, _keyCode))
	ke, ok := fromKeycode(keyCode)
	if !ok {
		// no such key; modifiers by themselves are handled by -[self flagsChanged:]
		handleKeyEvent(self)
		return
	}
	// either ke.Key or ke.ExtKey will be set at this point
	ke.Modifiers = parseModifiers(e)
	ke.Up = up
	sendKeyEvent(self, ke)
}

//export areaView_keyDown
func areaView_keyDown(self C.id, sel C.SEL, e C.id) {
	areaKeyEvent(self, e, false)
}

//export areaView_keyUp
func areaView_keyUp(self C.id, sel C.SEL, e C.id) {
	areaKeyEvent(self, e, true)
}

//export areaView_flagsChanged
func areaView_flagsChanged(self C.id, sel C.SEL, e C.id) {
	var ke KeyEvent

	// Mac OS X sends this event on both key up and key down.
	// Fortunately -[e keyCode] IS valid here, so we can simply map from key code to Modifiers, get the value of [e modifierFlags], the respective bit is set or not — that will give us the up/down state
	keyCode := uintptr(C.objc_msgSend_ushortret_noargs(e, _keyCode))
	mod, ok := keycodeModifiers[keyCode]		// comma-ok form to avoid adding entries
	if !ok {		// unknown modifier; ignore
		handleKeyEvent(self)
		return
	}
	ke.Modifiers = parseModifiers(e)
	ke.Up = (ke.Modifiers & mod) == 0
	sendKeyEvent(self, ke)
}

// TODO combine these with the listbox functions?

func newAreaScrollView(area C.id) C.id {
	scrollview := C.objc_msgSend_noargs(_NSScrollView, _alloc)
	scrollview = initWithDummyFrame(scrollview)
	C.objc_msgSend_bool(scrollview, _setHasHorizontalScroller, C.BOOL(C.YES))
	C.objc_msgSend_bool(scrollview, _setHasVerticalScroller, C.BOOL(C.YES))
	C.objc_msgSend_bool(scrollview, _setAutohidesScrollers, C.BOOL(C.YES))
	C.objc_msgSend_id(scrollview, _setDocumentView, area)
	return scrollview
}

func areaInScrollView(scrollview C.id) C.id {
	return C.objc_msgSend_noargs(scrollview, _documentView)
}

func makeArea(parentWindow C.id, alternate bool) C.id {
	area := C.objc_msgSend_noargs(_goArea, _alloc)
	area = initWithDummyFrame(area)
	// TODO others?
	area = newAreaScrollView(area)
	addControl(parentWindow, area)
	return area
}

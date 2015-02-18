// 24 march 2014

package ui

import (
	"fmt"
	"image"
	"syscall"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type area struct {
	*areabase

	*controlSingleHWND

	clickCounter *clickCounter

	textfield     C.HWND
	textfielddone *event
}

func makeAreaWindowClass() error {
	var errmsg *C.char

	err := C.makeAreaWindowClass(&errmsg)
	if err != 0 || errmsg != nil {
		return fmt.Errorf("%s: %v", C.GoString(errmsg), syscall.Errno(err))
	}
	return nil
}

func newArea(ab *areabase) Area {
	a := &area{
		areabase:      ab,
		clickCounter:  new(clickCounter),
		textfielddone: newEvent(),
	}
	a.controlSingleHWND = newControlSingleHWND(C.newArea(unsafe.Pointer(a)))
	a.fpreferredSize = a.xpreferredSize
	a.SetSize(a.width, a.height)
	a.textfield = C.newAreaTextField(a.hwnd, unsafe.Pointer(a))
	C.controlSetControlFont(a.textfield)
	return a
}

func (a *area) SetSize(width, height int) {
	a.width = width
	a.height = height
	C.SendMessageW(a.hwnd, C.msgAreaSizeChanged, 0, 0)
}

func (a *area) Repaint(r image.Rectangle) {
	var hscroll, vscroll C.int
	var rect C.RECT

	C.SendMessageW(a.hwnd, C.msgAreaGetScroll, C.WPARAM(uintptr(unsafe.Pointer(&hscroll))), C.LPARAM(uintptr(unsafe.Pointer(&vscroll))))
	r = r.Add(image.Pt(int(hscroll), int(vscroll))) // adjust by scroll position
	r = image.Rect(0, 0, a.width, a.height).Intersect(r)
	if r.Empty() {
		return
	}
	rect.left = C.LONG(r.Min.X)
	rect.top = C.LONG(r.Min.Y)
	rect.right = C.LONG(r.Max.X)
	rect.bottom = C.LONG(r.Max.Y)
	C.SendMessageW(a.hwnd, C.msgAreaRepaint, 0, C.LPARAM(uintptr(unsafe.Pointer(&rect))))
}

func (a *area) RepaintAll() {
	C.SendMessageW(a.hwnd, C.msgAreaRepaintAll, 0, 0)
}

func (a *area) OpenTextFieldAt(x, y int) {
	if x < 0 || x >= a.width || y < 0 || y >= a.height {
		panic(fmt.Errorf("point (%d,%d) outside Area in Area.OpenTextFieldAt()", x, y))
	}
	C.areaOpenTextField(a.hwnd, a.textfield, C.int(x), C.int(y), textfieldWidth, textfieldHeight)
}

func (a *area) TextFieldText() string {
	return getWindowText(a.textfield)
}

func (a *area) SetTextFieldText(text string) {
	t := toUTF16(text)
	C.setWindowText(a.textfield, t)
}

func (a *area) OnTextFieldDismissed(f func()) {
	a.textfielddone.set(f)
}

//export areaTextFieldDone
func areaTextFieldDone(data unsafe.Pointer) {
	a := (*area)(data)
	C.areaMarkTextFieldDone(a.hwnd)
	a.textfielddone.fire()
}

//export doPaint
func doPaint(xrect *C.RECT, hscroll C.int, vscroll C.int, data unsafe.Pointer, dx *C.intptr_t, dy *C.intptr_t) unsafe.Pointer {
	a := (*area)(data)
	// both Windows RECT and Go image.Rect are point..point, so the following is correct
	cliprect := image.Rect(int(xrect.left), int(xrect.top), int(xrect.right), int(xrect.bottom))
	cliprect = cliprect.Add(image.Pt(int(hscroll), int(vscroll))) // adjust by scroll position
	// make sure the cliprect doesn't fall outside the size of the Area
	cliprect = cliprect.Intersect(image.Rect(0, 0, a.width, a.height))
	if !cliprect.Empty() { // we have an update rect
		i := a.handler.Paint(cliprect)
		*dx = C.intptr_t(i.Rect.Dx())
		*dy = C.intptr_t(i.Rect.Dy())
		return unsafe.Pointer(i)
	}
	return nil
}

//export dotoARGB
func dotoARGB(img unsafe.Pointer, ppvBits unsafe.Pointer, toNRGBA C.BOOL) {
	i := (*image.RGBA)(unsafe.Pointer(img))
	t := toNRGBA != C.FALSE
	// the bitmap Windows gives us has a stride == width
	// TODO use GetObject() and get the stride from the resultant BITMAP to be *absolutely* sure
	toARGB(i, uintptr(ppvBits), i.Rect.Dx()*4, t)
}

//export areaWidthLONG
func areaWidthLONG(data unsafe.Pointer) C.LONG {
	a := (*area)(data)
	return C.LONG(a.width)
}

//export areaHeightLONG
func areaHeightLONG(data unsafe.Pointer) C.LONG {
	a := (*area)(data)
	return C.LONG(a.height)
}

func getModifiers() (m Modifiers) {
	down := func(x C.int) bool {
		// GetKeyState() gets the key state at the time of the message, so this is what we want
		return (C.GetKeyState(x) & 0x80) != 0
	}

	if down(C.VK_CONTROL) {
		m |= Ctrl
	}
	if down(C.VK_MENU) {
		m |= Alt
	}
	if down(C.VK_SHIFT) {
		m |= Shift
	}
	if down(C.VK_LWIN) || down(C.VK_RWIN) {
		m |= Super
	}
	return m
}

//export finishAreaMouseEvent
func finishAreaMouseEvent(data unsafe.Pointer, cbutton C.DWORD, up C.BOOL, heldButtons C.uintptr_t, xpos C.int, ypos C.int) {
	var me MouseEvent

	a := (*area)(data)
	button := uint(cbutton)
	me.Pos = image.Pt(int(xpos), int(ypos))
	if !me.Pos.In(image.Rect(0, 0, a.width, a.height)) { // outside the actual Area; no event
		return
	}
	if up != C.FALSE {
		me.Up = button
	} else if button != 0 { // don't run the click counter if the mouse was only moved
		me.Down = button
		// this returns a LONG, which is int32, but we don't need to worry about the signedness because for the same bit widths and two's complement arithmetic, s1-s2 == u1-u2 if bits(s1)==bits(s2) and bits(u1)==bits(u2) (and Windows requires two's complement: http://blogs.msdn.com/b/oldnewthing/archive/2005/05/27/422551.aspx)
		// signedness isn't much of an issue for these calls anyway because http://stackoverflow.com/questions/24022225/what-are-the-sign-extension-rules-for-calling-windows-api-functions-stdcall-t and that we're only using unsigned values (think back to how you (didn't) handle signedness in assembly language) AND because of the above AND because the statistics below (time interval and width/height) really don't make sense if negative
		time := C.GetMessageTime()
		maxTime := C.GetDoubleClickTime()
		// ignore zero returns and errors; MSDN says zero will be returned on error but that GetLastError() is meaningless
		xdist := C.GetSystemMetrics(C.SM_CXDOUBLECLK)
		ydist := C.GetSystemMetrics(C.SM_CYDOUBLECLK)
		me.Count = a.clickCounter.click(button, me.Pos.X, me.Pos.Y,
			uintptr(time), uintptr(maxTime), int(xdist/2), int(ydist/2))
	}
	// though wparam will contain control and shift state, let's use just one function to get modifiers for both keyboard and mouse events; it'll work the same anyway since we have to do this for alt and windows key (super)
	me.Modifiers = getModifiers()
	if button != 1 && (heldButtons&C.MK_LBUTTON) != 0 {
		me.Held = append(me.Held, 1)
	}
	if button != 2 && (heldButtons&C.MK_MBUTTON) != 0 {
		me.Held = append(me.Held, 2)
	}
	if button != 3 && (heldButtons&C.MK_RBUTTON) != 0 {
		me.Held = append(me.Held, 3)
	}
	if button != 4 && (heldButtons&C.MK_XBUTTON1) != 0 {
		me.Held = append(me.Held, 4)
	}
	if button != 5 && (heldButtons&C.MK_XBUTTON2) != 0 {
		me.Held = append(me.Held, 5)
	}
	a.handler.Mouse(me)
}

//export areaKeyEvent
func areaKeyEvent(data unsafe.Pointer, up C.BOOL, wParam C.WPARAM, lParam C.LPARAM) C.BOOL {
	var ke KeyEvent

	a := (*area)(data)
	lp := uint32(lParam) // to be safe
	// the numeric keypad keys when Num Lock is off are considered left-hand keys as the separate navigation buttons were added later
	// the numeric keypad enter, however, is a right-hand key because it has the same virtual-key code as the typewriter enter
	righthand := (lp & 0x01000000) != 0

	scancode := byte((lp >> 16) & 0xFF)
	ke.Modifiers = getModifiers()
	if extkey, ok := numpadextkeys[wParam]; ok && !righthand {
		// the above is special handling for numpad keys to ignore the state of Num Lock and Shift; see http://blogs.msdn.com/b/oldnewthing/archive/2004/09/06/226045.aspx and https://github.com/glfw/glfw/blob/master/src/win32_window.c#L152
		ke.ExtKey = extkey
	} else if wParam == C.VK_RETURN && righthand {
		ke.ExtKey = NEnter
	} else if extkey, ok := extkeys[wParam]; ok {
		ke.ExtKey = extkey
	} else if mod, ok := modonlykeys[wParam]; ok {
		ke.Modifier = mod
		// don't include the modifier in ke.Modifiers
		ke.Modifiers &^= mod
	} else if xke, ok := fromScancode(uintptr(scancode)); ok {
		// one of these will be nonzero
		ke.Key = xke.Key
		ke.ExtKey = xke.ExtKey
	} else if ke.Modifiers == 0 {
		// no key, extkey, or modifiers; do nothing
		return C.FALSE
	}
	ke.Up = up != C.FALSE
	handled := a.handler.Key(ke)
	if handled {
		return C.TRUE
	}
	return C.FALSE
}

// all mappings come from GLFW - https://github.com/glfw/glfw/blob/master/src/win32_window.c#L152
var numpadextkeys = map[C.WPARAM]ExtKey{
	C.VK_HOME:   N7,
	C.VK_UP:     N8,
	C.VK_PRIOR:  N9,
	C.VK_LEFT:   N4,
	C.VK_CLEAR:  N5,
	C.VK_RIGHT:  N6,
	C.VK_END:    N1,
	C.VK_DOWN:   N2,
	C.VK_NEXT:   N3,
	C.VK_INSERT: N0,
	C.VK_DELETE: NDot,
}

var extkeys = map[C.WPARAM]ExtKey{
	C.VK_ESCAPE: Escape,
	C.VK_INSERT: Insert,
	C.VK_DELETE: Delete,
	C.VK_HOME:   Home,
	C.VK_END:    End,
	C.VK_PRIOR:  PageUp,
	C.VK_NEXT:   PageDown,
	C.VK_UP:     Up,
	C.VK_DOWN:   Down,
	C.VK_LEFT:   Left,
	C.VK_RIGHT:  Right,
	C.VK_F1:     F1,
	C.VK_F2:     F2,
	C.VK_F3:     F3,
	C.VK_F4:     F4,
	C.VK_F5:     F5,
	C.VK_F6:     F6,
	C.VK_F7:     F7,
	C.VK_F8:     F8,
	C.VK_F9:     F9,
	C.VK_F10:    F10,
	C.VK_F11:    F11,
	C.VK_F12:    F12,
	// numpad numeric keys and . are handled in events_notdarwin.go
	// numpad enter is handled in code above
	C.VK_ADD:      NAdd,
	C.VK_SUBTRACT: NSubtract,
	C.VK_MULTIPLY: NMultiply,
	C.VK_DIVIDE:   NDivide,
}

// sanity check
func init() {
	included := make([]bool, _nextkeys)
	for _, v := range extkeys {
		included[v] = true
	}
	for i := 1; i < int(_nextkeys); i++ {
		if i >= int(N0) && i <= int(N9) { // skip numpad numbers, ., and enter
			continue
		}
		if i == int(NDot) || i == int(NEnter) {
			continue
		}
		if !included[i] {
			panic(fmt.Errorf("error: not all ExtKeys defined on Windows (missing %d)", i))
		}
	}
}

var modonlykeys = map[C.WPARAM]Modifiers{
	// even if the separate left/right aren't necessary, have them here anyway, just to be safe
	C.VK_CONTROL:  Ctrl,
	C.VK_LCONTROL: Ctrl,
	C.VK_RCONTROL: Ctrl,
	C.VK_MENU:     Alt,
	C.VK_LMENU:    Alt,
	C.VK_RMENU:    Alt,
	C.VK_SHIFT:    Shift,
	C.VK_LSHIFT:   Shift,
	C.VK_RSHIFT:   Shift,
	// there's no combined Windows key virtual-key code as there is with the others
	C.VK_LWIN: Super,
	C.VK_RWIN: Super,
}

//export areaResetClickCounter
func areaResetClickCounter(data unsafe.Pointer) {
	a := (*area)(data)
	a.clickCounter.reset()
}

func (a *area) xpreferredSize(d *sizing) (width, height int) {
	// the preferred size of an Area is its size
	return a.width, a.height
}

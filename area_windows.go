// 24 march 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
	"sync"
	"image"
)

const (
	areastyle = _WS_HSCROLL | _WS_VSCROLL | controlstyle
	areaxstyle = 0 | controlxstyle
)

const (
	areaWndClassFormat = "gouiarea%X"
)

var (
	areaWndClassNum uintptr
	areaWndClassNumLock sync.Mutex
)

func getScrollPos(hwnd _HWND) (xpos int32, ypos int32) {
	var si _SCROLLINFO

	si.cbSize = uint32(unsafe.Sizeof(si))
	si.fMask = _SIF_POS | _SIF_TRACKPOS
	r1, _, err := _getScrollInfo.Call(
		uintptr(hwnd),
		uintptr(_SB_HORZ),
		uintptr(unsafe.Pointer(&si)))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error getting horizontal scroll position for Area: %v", err))
	}
	xpos = si.nPos
	si.cbSize = uint32(unsafe.Sizeof(si))			// MSDN example code reinitializes this each time, so we'll do it too just to be safe
	si.fMask = _SIF_POS | _SIF_TRACKPOS
	r1, _, err = _getScrollInfo.Call(
		uintptr(hwnd),
		uintptr(_SB_VERT),
		uintptr(unsafe.Pointer(&si)))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error getting vertical scroll position for Area: %v", err))
	}
	ypos = si.nPos
	return xpos, ypos
}

var (
	_getUpdateRect = user32.NewProc("GetUpdateRect")
	_beginPaint = user32.NewProc("BeginPaint")
	_endPaint = user32.NewProc("EndPaint")
	_fillRect = user32.NewProc("FillRect")
	_gdipCreateBitmapFromScan0 = gdiplus.NewProc("GdipCreateBitmapFromScan0")
	_gdipCreateFromHDC = gdiplus.NewProc("GdipCreateFromHDC")
	_gdipDrawImageI = gdiplus.NewProc("GdipDrawImageI")
	_gdipDeleteGraphics = gdiplus.NewProc("GdipDeleteGraphics")
	_gdipDisposeImage = gdiplus.NewProc("GdipDisposeImage")
)

const (
	areaBackgroundBrush = _HBRUSH(_COLOR_BTNFACE + 1)

	// from winuser.h
	_WM_PAINT = 0x000F
)

func paintArea(s *sysData) {
	const (
		// from gdipluspixelformats.h
		_PixelFormatGDI = 0x00020000
		_PixelFormatAlpha = 0x00040000
		_PixelFormatCanonical = 0x00200000
		_PixelFormat32bppARGB = (10 | (32 << 8) | _PixelFormatAlpha | _PixelFormatGDI | _PixelFormatCanonical)
	)

	var xrect _RECT
	var ps _PAINTSTRUCT

	r1, _, _ := _getUpdateRect.Call(
		uintptr(s.hwnd),
		uintptr(unsafe.Pointer(&xrect)),
		uintptr(_TRUE))		// erase the update rect with the background color
	if r1 == 0 {			// no update rect; do nothing
		return
	}

	hscroll, vscroll := getScrollPos(s.hwnd)

	cliprect := image.Rect(int(xrect.Left), int(xrect.Top), int(xrect.Right), int(xrect.Bottom))
	cliprect = cliprect.Add(image.Pt(int(hscroll), int(vscroll)))			// adjust by scroll position
	// make sure the cliprect doesn't fall outside the size of the Area
	cliprect = cliprect.Intersect(image.Rect(0, 0, s.areawidth, s.areaheight))
	if cliprect.Empty() {		// still no update rect
		return
	}

	r1, _, err := _beginPaint.Call(
		uintptr(s.hwnd),
		uintptr(unsafe.Pointer(&ps)))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error beginning Area repaint: %v", err))
	}
	hdc := _HANDLE(r1)

	// Windows won't necessarily erase the update rect for us; we need to do so ourselves
	// thanks to the people at http://stackoverflow.com/questions/23001890/winapi-getupdaterect-with-brepaint-true-inside-wm-paint-doesnt-clear-the-pai
	// TODO this whole thing is inefficient, as explained in the page; we probably don't need _getUpdateRect
	if ps.fErase != 0 {		// if Windows didn't
		r1, _, err := _fillRect.Call(
			uintptr(hdc),
			uintptr(unsafe.Pointer(&xrect)),
			uintptr(areaBackgroundBrush))
		if r1 == 0 {		// failure
			panic(fmt.Errorf("error manually clearing Area background: %v", err))
		}
	}

	i := s.handler.Paint(cliprect)
	// the pixels are arranged in RGBA order, but GDI+ requires BGRA
	// we don't have a choice but to convert it ourselves
	// TODO make realbits a part of sysData to conserve memory
	realbits := make([]byte, 4 * i.Rect.Dx() * i.Rect.Dy())
	p := pixelDataPos(i)
	q := 0
	for y := i.Rect.Min.Y; y < i.Rect.Max.Y; y++ {
		nextp := p + i.Stride
		for x := i.Rect.Min.X; x < i.Rect.Max.X; x++ {
			realbits[q + 0] = byte(i.Pix[p + 2])		// B
			realbits[q + 1] = byte(i.Pix[p + 1])		// G
			realbits[q + 2] = byte(i.Pix[p + 0])		// R
			realbits[q + 3] = byte(i.Pix[p + 3])		// A
			p += 4
			q += 4
		}
		p = nextp
	}

	var bitmap, graphics uintptr

	r1, _, err = _gdipCreateBitmapFromScan0.Call(
		uintptr(i.Rect.Dx()),
		uintptr(i.Rect.Dy()),
		uintptr(i.Rect.Dx() * 4),			// got rid of extra stride
		uintptr(_PixelFormat32bppARGB),
		uintptr(unsafe.Pointer(&realbits[0])),
		uintptr(unsafe.Pointer(&bitmap)))
	if r1 != 0 {			// failure
		panic(fmt.Errorf("error creating GDI+ bitmap to blit (GDI+ error code %d; Windows last error %v)", r1, err))
	}
	r1, _, err = _gdipCreateFromHDC.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(&graphics)))
	if r1 != 0 {			// failure
		panic(fmt.Errorf("error creating GDI+ graphics context to blit to (GDI+ error code %d; Windows last error %v)", r1, err))
	}
	r1, _, err = _gdipDrawImageI.Call(
		graphics,
		bitmap,
		uintptr(xrect.Left),			// cliprect is adjusted; use original
		uintptr(xrect.Top))
	if r1 != 0 {			// failure
		panic(fmt.Errorf("error blitting GDI+ bitmap (GDI+ error code %d; Windows last error %v)", r1, err))
	}
	r1, _, err = _gdipDeleteGraphics.Call(graphics)
	if r1 != 0 {			// failure
		panic(fmt.Errorf("error freeing GDI+ graphics context to blit to (GDI+ error code %d; Windows last error %v)", r1, err))
	}
	// TODO this is the destructor of Image (Bitmap's base class); I don't see a specific destructor for Bitmap itself so
	r1, _, err = _gdipDisposeImage.Call(bitmap)
	if r1 != 0 {			// failure
		panic(fmt.Errorf("error freeing GDI+ bitmap to blit (GDI+ error code %d; Windows last error %v)", r1, err))
	}

	// return value always nonzero according to MSDN
	_endPaint.Call(
		uintptr(s.hwnd),
		uintptr(unsafe.Pointer(&ps)))
}

func getAreaControlSize(hwnd _HWND) (width int, height int) {
	var rect _RECT

	r1, _, err := _getClientRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error getting size of actual Area control: %v", err))
	}
	return int(rect.Right - rect.Left),
		int(rect.Bottom - rect.Top)
}

func scrollArea(s *sysData, wparam _WPARAM, which uintptr) {
	var si _SCROLLINFO

	cwid, cht := getAreaControlSize(s.hwnd)
	pagesize := int32(cwid)
	maxsize := int32(s.areawidth)
	if which == uintptr(_SB_VERT) {
		pagesize = int32(cht)
		maxsize = int32(s.areaheight)
	}

	si.cbSize = uint32(unsafe.Sizeof(si))
	si.fMask = _SIF_POS | _SIF_TRACKPOS
	r1, _, err := _getScrollInfo.Call(
		uintptr(s.hwnd),
		which,
		uintptr(unsafe.Pointer(&si)))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error getting current scroll position for scrolling: %v", err))
	}

	newpos := si.nPos
	switch wparam & 0xFFFF {
	case _SB_LEFT:			// also _SB_TOP but Go won't let me
		newpos = 0
	case _SB_RIGHT:		// also _SB_BOTTOM
		// see comment in adjustAreaScrollbars() below
		newpos = maxsize - pagesize
	case _SB_LINELEFT:		// also _SB_LINEUP
		newpos--
	case _SB_LINERIGHT:	// also _SB_LINEDOWN
		newpos++
	case _SB_PAGELEFT:		// also _SB_PAGEUP
		newpos -= pagesize
	case _SB_PAGERIGHT:	// also _SB_PAGEDOWN
		newpos += pagesize
	case _SB_THUMBPOSITION:
		// TODO is this the same as SB_THUMBTRACK instead? MSDN says use of thumb pos is only for that one
		// do nothing; newpos already has the thumb's position
	case _SB_THUMBTRACK:
		newpos = si.nTrackPos
	}		// otherwise just keep the current position (that's what MSDN example code says, anyway)

	// make sure we're not out of range
	if newpos < 0 {
		newpos = 0
	}
	if newpos > (maxsize - pagesize) {
		newpos = maxsize - pagesize
	}

	// TODO is this the right thing to do for SB_THUMBTRACK? or will it conflict?
	if newpos == si.nPos {		// no change; no scrolling
		return
	}

	delta := -(newpos - si.nPos)	// negative because ScrollWindowEx() scrolls in the opposite direction
	dx := delta
	dy := int32(0)
	if which == uintptr(_SB_VERT) {
		dx = int32(0)
		dy = delta
	}
	r1, _, err = _scrollWindowEx.Call(
		uintptr(s.hwnd),
		uintptr(dx),
		uintptr(dy),
		uintptr(0),			// these four change what is scrolled and record info about the scroll; we're scrolling the whole client area and don't care about the returned information here
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(_SW_INVALIDATE | _SW_ERASE))					// mark the remaining rect as needing redraw and erase...
	if r1 == _ERROR {		// failure
		panic(fmt.Errorf("error scrolling Area: %v", err))
	}
	// ...but don't redraw the window yet; we need to apply our scroll changes

	// we actually have to commit the change back to the scrollbar; otherwise the scroll position will merely reset itself
	si.cbSize = uint32(unsafe.Sizeof(si))
	si.fMask = _SIF_POS
	si.nPos = newpos
	_setScrollInfo.Call(
		uintptr(s.hwnd),
		which,
		uintptr(unsafe.Pointer(&si)))

	// NOW redraw it
	r1, _, err = _updateWindow.Call(uintptr(s.hwnd))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error updating Area after scrolling: %v", err))
	}

	// TODO in some cases wine will show a thumb one pixel away from the advance arrow button if going to the end; the values are correct though... weirdness in wine or something I never noticed about Windows?
}

func adjustAreaScrollbars(s *sysData) {
	var si _SCROLLINFO

	cwid, cht := getAreaControlSize(s.hwnd)

	// the trick is we want a page to be the width/height of the visible area
	// so the scroll range would go from [0..image_dimension - control_dimension]
	// but judging from the sample code on MSDN, we don't need to do this; the scrollbar will do it for us
	// we DO need to handle it when scrolling, though, since the thumb can only go up to this upper limit

	// have to do horizontal and vertical separately
	si.cbSize = uint32(unsafe.Sizeof(si))
	si.fMask = _SIF_RANGE | _SIF_PAGE
	si.nMin = 0
	si.nMax = int32(s.areawidth)
	si.nPage = uint32(cwid)
	_setScrollInfo.Call(
		uintptr(s.hwnd),
		uintptr(_SB_HORZ),
		uintptr(unsafe.Pointer(&si)),
		uintptr(_TRUE))			// redraw the scroll bar

	si.cbSize = uint32(unsafe.Sizeof(si))			// MSDN sample code does this a second time; let's do it too to be safe
	si.fMask = _SIF_RANGE | _SIF_PAGE
	si.nMin = 0
	si.nMax = int32(s.areaheight)
	si.nPage = uint32(cht)
	_setScrollInfo.Call(
		uintptr(s.hwnd),
		uintptr(_SB_VERT),
		uintptr(unsafe.Pointer(&si)),
		uintptr(_TRUE))			// redraw the scroll bar
}

var (
	_invalidateRect = user32.NewProc("InvalidateRect")
)

func repaintArea(s *sysData) {
	r1, _, err := _invalidateRect.Call(
		uintptr(s.hwnd),
		uintptr(0),			// the whole area
		uintptr(_TRUE))		// have Windows erase if possible
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error flagging Area as needing repainting after event (last error: %v)", err))
	}
	r1, _, err = _updateWindow.Call(uintptr(s.hwnd))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error repainting Area after event: %v", err))
	}
}

var (
	_getKeyState = user32.NewProc("GetKeyState")
)

func getModifiers() (m Modifiers) {
	down := func(x uintptr) bool {
		r1, _, _ := _getKeyState.Call(x)
		return (r1 & 0x80) != 0
	}

	if down(_VK_CONTROL) {
		m |= Ctrl
	}
	if down(_VK_MENU) {
		m |= Alt
	}
	if down(_VK_SHIFT) {
		m |= Shift
	}
	// TODO windows key (super)
	return m
}

// TODO populate me.Held
func areaMouseEvent(s *sysData, button uint, up bool, count uint, wparam _WPARAM, lparam _LPARAM) {
	var me MouseEvent

	xpos, ypos := getScrollPos(s.hwnd)		// mouse coordinates are relative to control; make them relative to Area
	xpos += lparam._X()
	ypos += lparam._Y()
	me.Pos = image.Pt(int(xpos), int(ypos))
	if !me.Pos.In(image.Rect(0, 0, s.areawidth, s.areaheight)) {		// outside the actual Area; no event
		return
	}
	if up {
		me.Up = button
	} else {
		me.Down = button
		me.Count = count
	}
	// though wparam will contain control and shift state, let's use just one function to get modifiers for both keyboard and mouse events; it'll work the same anyway since we have to do this for alt and windows key (super)
	me.Modifiers = getModifiers()
	if button != 1 && (wparam & _MK_LBUTTON) != 0 {
		me.Held = append(me.Held, 1)
	}
	if button != 2 && (wparam & _MK_MBUTTON) != 0 {
		me.Held = append(me.Held, 2)
	}
	if button != 3 && (wparam & _MK_RBUTTON) != 0 {
		me.Held = append(me.Held, 3)
	}
	// TODO XBUTTONs?
	repaint := s.handler.Mouse(me)
	if repaint {
		repaintArea(s)
	}
}

func areaKeyEvent(s *sysData, up bool, wparam _WPARAM, lparam _LPARAM) bool {
	var ke KeyEvent

	scancode := byte((lparam >> 16) & 0xFF)
	ke.Modifiers = getModifiers()
	if wparam == _VK_RETURN && (lparam & 0x01000000) != 0 {
		// the above is special handling for numpad enter
		// bit 24 of LPARAM (0x01000000) indicates right-hand keys
		ke.ExtKey = NEnter
	} else if extkey, ok := extkeys[wparam]; ok {
		ke.ExtKey = extkey
	} else if xke, ok := fromScancode(uintptr(scancode)); ok {
		// one of these will be nonzero
		ke.Key = xke.Key
		ke.ExtKey = xke.ExtKey
	} else if ke.Modifiers == 0 {
		// no key, extkey, or modifiers; do nothing but mark not handled
		return false
	}
	ke.Up = up
	handled, repaint := s.handler.Key(ke)
	if repaint {
		repaintArea(s)
	}
	return handled
}

var extkeys = map[_WPARAM]ExtKey{
	_VK_ESCAPE:		Escape,
	_VK_INSERT:		Insert,
	_VK_DELETE:		Delete,
	_VK_HOME:		Home,
	_VK_END:			End,
	_VK_PRIOR:		PageUp,
	_VK_NEXT:		PageDown,
	_VK_UP:			Up,
	_VK_DOWN:		Down,
	_VK_LEFT:		Left,
	_VK_RIGHT:		Right,
	_VK_F1:			F1,
	_VK_F2:			F2,
	_VK_F3:			F3,
	_VK_F4:			F4,
	_VK_F5:			F5,
	_VK_F6:			F6,
	_VK_F7:			F7,
	_VK_F8:			F8,
	_VK_F9:			F9,
	_VK_F10:			F10,
	_VK_F11:			F11,
	_VK_F12:			F12,
	// numpad numeric keys and . are handled in events_notdarwin.go
	// numpad enter is handled in code above
	_VK_ADD:		NAdd,
	_VK_SUBTRACT:	NSubtract,
	_VK_MULTIPLY:	NMultiply,
	_VK_DIVIDE:		NDivide,
}

// sanity check
func init() {
	included := make([]bool, _nextkeys)
	for _, v := range extkeys {
		included[v] = true
	}
	for i := 1; i < int(_nextkeys); i++ {
		if i >= int(N0) && i <= int(N9) {		// skip numpad numbers, ., and enter
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

var (
	_setFocus = user32.NewProc("SetFocus")
)

func areaWndProc(s *sysData) func(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
	return func(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
		const (
			_MA_ACTIVATE = 1
		)

		defwndproc := func() _LRESULT {
			r1, _, _ := defWindowProc.Call(
				uintptr(hwnd),
				uintptr(uMsg),
				uintptr(wParam),
				uintptr(lParam))
			return _LRESULT(r1)
		}

		switch uMsg {
		case _WM_PAINT:
			paintArea(s)
			return 0
		case _WM_HSCROLL:
			// TODO make this unnecessary
			if s != nil && s.hwnd != 0 {			// this message can be sent before s is assigned properly
				scrollArea(s, wParam, _SB_HORZ)
			}
			return 0
		case _WM_VSCROLL:
			// TODO make this unnecessary
			if s != nil && s.hwnd != 0 {			// this message can be sent before s is assigned properly
				scrollArea(s, wParam, _SB_VERT)
				return 0
			}
			return defwndproc()
		case _WM_SIZE:
			// TODO make this unnecessary
			if s != nil && s.hwnd != 0 {			// this message can be sent before s is assigned properly
				adjustAreaScrollbars(s)
				return 0
			}
			return defwndproc()
		case _WM_MOUSEACTIVATE:
			// register our window for keyboard input
			// (see http://www.catch22.net/tuts/custom-controls)
			r1, _, err := _setFocus.Call(uintptr(s.hwnd))
			if r1 == 0 {		// failure
				panic(fmt.Errorf("error giving Area keyboard focus: %v", err))
				return _MA_ACTIVATE		// TODO eat the click?
			}
			return defwndproc()
		case _WM_MOUSEMOVE:
			areaMouseEvent(s, 0, false, 0, wParam, lParam)
			return 0
		case _WM_LBUTTONDOWN:
			areaMouseEvent(s, 1, false, 1, wParam, lParam)
			return 0
		case _WM_LBUTTONDBLCLK:
			areaMouseEvent(s, 1, false, 2, wParam, lParam)
			return 0
		case _WM_LBUTTONUP:
			areaMouseEvent(s, 1, true, 0, wParam, lParam)
			return 0
		case _WM_MBUTTONDOWN:
			areaMouseEvent(s, 2, false, 1, wParam, lParam)
			return 0
		case _WM_MBUTTONDBLCLK:
			areaMouseEvent(s, 2, false, 2, wParam, lParam)
			return 0
		case _WM_MBUTTONUP:
			areaMouseEvent(s, 2, true, 0, wParam, lParam)
			return 0
		case _WM_RBUTTONDOWN:
			areaMouseEvent(s, 3, false, 1, wParam, lParam)
			return 0
		case _WM_RBUTTONDBLCLK:
			areaMouseEvent(s, 3, false, 2, wParam, lParam)
			return 0
		case _WM_RBUTTONUP:
			areaMouseEvent(s, 3, true, 0, wParam, lParam)
			return 0
		// TODO XBUTTONs?
		case _WM_KEYDOWN:
			areaKeyEvent(s, false, wParam, lParam)
			return 0
		case _WM_KEYUP:
			areaKeyEvent(s, true, wParam, lParam)
			return 0
		// Alt+[anything] and F10 send these instead
		case _WM_SYSKEYDOWN:
			handled := areaKeyEvent(s, false, wParam, lParam)
			if handled {
				return 0
			}
			return defwndproc()
		case _WM_SYSKEYUP:
			handled := areaKeyEvent(s, true, wParam, lParam)
			if handled {
				return 0
			}
			return defwndproc()
		case msgSetAreaSize:
			s.areawidth = int(wParam)		// see setAreaSize() in sysdata_windows.go
			s.areaheight = int(lParam)
			adjustAreaScrollbars(s)
			repaintArea(s)					// this calls for an update
			return 0
		default:
			return defwndproc()
		}
		panic(fmt.Sprintf("areaWndProc message %d did not return: internal bug in ui library", uMsg))
	}
}

func registerAreaWndClass(s *sysData) (newClassName string, err error) {
	const (
		// from winuser.h
		_CS_DBLCLKS = 0x0008
	)

	areaWndClassNumLock.Lock()
	newClassName = fmt.Sprintf(areaWndClassFormat, areaWndClassNum)
	areaWndClassNum++
	areaWndClassNumLock.Unlock()

	wc := &_WNDCLASS{
		style:			_CS_DBLCLKS,		// needed to be able to register double-clicks
		lpszClassName:	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(newClassName))),
		lpfnWndProc:		syscall.NewCallback(areaWndProc(s)),
		hInstance:		hInstance,
		hIcon:			icon,
		hCursor:			cursor,
		hbrBackground:	areaBackgroundBrush,
	}

	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_registerClass,
		p:		[]uintptr{uintptr(unsafe.Pointer(wc))},
		ret:		ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		return "", r.err
	}
	return newClassName, nil
}

type _PAINTSTRUCT struct {
	hdc			_HANDLE
	fErase		int32		// originally BOOL
	rcPaint		_RECT
	fRestore		int32		// originally BOOL
	fIncUpdate	int32		// originally BOOL
	rgbReserved	[32]byte
}

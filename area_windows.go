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
	_alphaBlend = msimg32.NewProc("AlphaBlend")
	_beginPaint = user32.NewProc("BeginPaint")
	_bitBlt = gdi32.NewProc("BitBlt")
	_createCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	_createCompatibleDC = gdi32.NewProc("CreateCompatibleDC")
	_createDIBSection = gdi32.NewProc("CreateDIBSection")
	_deleteDC = gdi32.NewProc("DeleteDC")
	_deleteObject = gdi32.NewProc("DeleteObject")
	_endPaint = user32.NewProc("EndPaint")
	_fillRect = user32.NewProc("FillRect")
	_getUpdateRect = user32.NewProc("GetUpdateRect")
	// _selectObject in prefsize_windows.go
)

const (
	areaBackgroundBrush = _HBRUSH(_COLOR_BTNFACE + 1)

	// from winuser.h
	_WM_PAINT = 0x000F
)

func paintArea(s *sysData) {
	const (
		// from wingdi.h
		_BI_RGB = 0
		_DIB_RGB_COLORS = 0
		_SRCCOPY = 0x00CC0020
		_AC_SRC_OVER = 0x00
		_AC_SRC_ALPHA = 0x01
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

	// TODO don't do the above, but always draw the background color?

	r1, _, err := _beginPaint.Call(
		uintptr(s.hwnd),
		uintptr(unsafe.Pointer(&ps)))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error beginning Area repaint: %v", err))
	}
	hdc := _HANDLE(r1)

	// very big thanks to Ninjifox for suggesting this technique and helping me go through it

	// first let's create the destination image, which we fill with the windows background color
	// this is how we fake drawing the background; see also http://msdn.microsoft.com/en-us/library/ms969905.aspx
	r1, _, err = _createCompatibleDC.Call(uintptr(hdc))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error creating off-screen rendering DC: %v", err))
	}
	rdc := _HANDLE(r1)
	// the bitmap has to be compatible with the window
	// if we create a bitmap compatible with the DC we just created, it'll be monochrome
	// thanks to David Heffernan in http://stackoverflow.com/questions/23033636/winapi-gdi-fillrectcolor-btnface-fills-with-strange-grid-like-brush-on-window
	r1, _, err = _createCompatibleBitmap.Call(
		uintptr(hdc),
		uintptr(xrect.Right - xrect.Left),
		uintptr(xrect.Bottom - xrect.Top))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error creating off-screen rendering bitmap: %v", err))
	}
	rbitmap := _HANDLE(r1)
	r1, _, err = _selectObject.Call(
		uintptr(rdc),
		uintptr(rbitmap))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error connecting off-screen rendering bitmap to off-screen rendering DC: %v", err))
	}
	prevrbitmap := _HANDLE(r1)
	rrect := _RECT{
		Left:		0,
		Right:	xrect.Right - xrect.Left,
		Top:		0,
		Bottom:	xrect.Bottom - xrect.Top,
	}
	r1, _, err = _fillRect.Call(
		uintptr(rdc),
		uintptr(unsafe.Pointer(&rrect)),
		uintptr(areaBackgroundBrush))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error filling off-screen rendering bitmap with the system background color: %v", err))
	}

	i := s.handler.Paint(cliprect)
	// don't convert to BRGA just yet; see below

	// now we need to shove realbits into a bitmap
	// technically bitmaps don't know about alpha; they just ignore the alpha byte
	// AlphaBlend(), however, sees it - see http://msdn.microsoft.com/en-us/library/windows/desktop/dd183352%28v=vs.85%29.aspx
	bi := _BITMAPINFO{}
	bi.bmiHeader.biSize = uint32(unsafe.Sizeof(bi.bmiHeader))
	bi.bmiHeader.biWidth = int32(i.Rect.Dx())
	bi.bmiHeader.biHeight = -int32(i.Rect.Dy())		// negative height to force top-down drawing
	bi.bmiHeader.biPlanes = 1
	bi.bmiHeader.biBitCount = 32
	bi.bmiHeader.biCompression = _BI_RGB
	bi.bmiHeader.biSizeImage = uint32(i.Rect.Dx() * i.Rect.Dy() * 4)
	// this is all we need, but because this confused me at first, I will say the two pixels-per-meter fields are unused (see http://blogs.msdn.com/b/oldnewthing/archive/2013/05/15/10418646.aspx and page 581 of Charles Petzold's Programming Windows, Fifth Edition)
	ppvBits := uintptr(0)		// now for the trouble: CreateDIBSection() allocates the memory for us...
	r1, _, err = _createDIBSection.Call(
		uintptr(0),		// TODO is this safe? Ninjifox does it
		uintptr(unsafe.Pointer(&bi)),
		uintptr(_DIB_RGB_COLORS),
		uintptr(unsafe.Pointer(&ppvBits)),
		uintptr(0),		// we're not dealing with hSection or dwOffset
		uintptr(0))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error creating HBITMAP for image returned by AreaHandler.Paint(): %v", err))
	}
	ibitmap := _HANDLE(r1)

	// now we have to do TWO MORE things before we can finally do alpha blending
	// first, we need to load the bitmap memory, because Windows makes it for us
	// the pixels are arranged in RGBA order, but GDI requires BGRA
	// this turns out to be just ARGB in little endian; let's convert into this memory
	// the bitmap Windows gives us has a stride == width
	toARGB(i, ppvBits, i.Rect.Dx() * 4)

	// the second thing is... make a device context for the bitmap :|
	// Ninjifox just makes another compatible DC; we'll do the same
	r1, _, err = _createCompatibleDC.Call(uintptr(hdc))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error creating HDC for image returned by AreaHandler.Paint(): %v", err))
	}
	idc := _HANDLE(r1)
	r1, _, err = _selectObject.Call(
		uintptr(idc),
		uintptr(ibitmap))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error connecting HBITMAP for image returned by AreaHandler.Paint() to its HDC: %v", err))
	}
	previbitmap := _HANDLE(r1)

	// AND FINALLY WE CAN DO THE ALPHA BLENDING!!!!!!111
	blendfunc := _BLENDFUNCTION{
		BlendOp:				_AC_SRC_OVER,
		BlendFlags:			0,
		SourceConstantAlpha:	255,					// only use per-pixel alphas
		AlphaFormat:			_AC_SRC_ALPHA,		// premultiplied
	}
	r1, _, err = _alphaBlend.Call(
		uintptr(rdc),	// destination
		uintptr(0),		// origin and size
		uintptr(0),
		uintptr(i.Rect.Dx()),
		uintptr(i.Rect.Dy()),
		uintptr(idc),	// source image
		uintptr(0),
		uintptr(0),
		uintptr(i.Rect.Dx()),
		uintptr(i.Rect.Dy()),
		blendfunc.arg())
	if r1 == _FALSE {		// failure
		panic(fmt.Errorf("error alpha-blending image returned by AreaHandler.Paint() onto background: %v", err))
	}

	// and finally we can just blit that into the window
	r1, _, err = _bitBlt.Call(
		uintptr(hdc),
		uintptr(xrect.Left),
		uintptr(xrect.Top),
		uintptr(xrect.Right - xrect.Left),
		uintptr(xrect.Bottom - xrect.Top),
		uintptr(rdc),
		uintptr(0),			// from the rdc's origin
		uintptr(0),
		uintptr(_SRCCOPY))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error blitting Area image to Area: %v", err))
	}

	// now to clean up
	r1, _, err = _selectObject.Call(
		uintptr(idc),
		uintptr(previbitmap))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error reverting HDC for image returned by AreaHandler.Paint() to original HBITMAP: %v", err))
	}
	r1, _, err = _selectObject.Call(
		uintptr(rdc),
		uintptr(prevrbitmap))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error reverting HDC for off-screen rendering to original HBITMAP: %v", err))
	}
	r1, _, err = _deleteObject.Call(uintptr(ibitmap))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error deleting HBITMAP for image returned by AreaHandler.Paint(): %v", err))
	}
	r1, _, err = _deleteObject.Call(uintptr(rbitmap))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error deleting HBITMAP for off-screen rendering: %v", err))
	}
	r1, _, err = _deleteDC.Call(uintptr(idc))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error deleting HDC for image returned by AreaHandler.Paint(): %v", err))
	}
	r1, _, err = _deleteDC.Call(uintptr(rdc))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error deleting HDC for off-screen rendering: %v", err))
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
		// raymond chen says to just set the newpos to the SCROLLINFO nPos for this message; see http://blogs.msdn.com/b/oldnewthing/archive/2003/07/31/54601.aspx and http://blogs.msdn.com/b/oldnewthing/archive/2003/08/05/54602.aspx
		// do nothing here; newpos already has nPos
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
	si.nMax = int32(s.areawidth - 1)		// the max point is inclusive, so we have to pass in the last valid value, not the first invalid one (see http://blogs.msdn.com/b/oldnewthing/archive/2003/07/31/54601.aspx); if we don't, we get weird things like the scrollbar sometimes showing one extra scroll position at the end that you can never scroll to
	si.nPage = uint32(cwid)
	_setScrollInfo.Call(
		uintptr(s.hwnd),
		uintptr(_SB_HORZ),
		uintptr(unsafe.Pointer(&si)),
		uintptr(_TRUE))			// redraw the scroll bar

	si.cbSize = uint32(unsafe.Sizeof(si))			// MSDN sample code does this a second time; let's do it too to be safe
	si.fMask = _SIF_RANGE | _SIF_PAGE
	si.nMin = 0
	si.nMax = int32(s.areaheight - 1)
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
		case _WM_ERASEBKGND:
			// don't draw a background; we'll do so when painting
			// this is to make things flicker-free; see http://msdn.microsoft.com/en-us/library/ms969905.aspx
			return 1
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
			// transfer keyboard focus to our Area on an activating click
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
		_CS_HREDRAW = 0x0002
		_CS_VREDRAW = 0x0001

		// from winuser.h
		_CS_DBLCLKS = 0x0008		// needed to be able to register double-clicks
	)

	areaWndClassNumLock.Lock()
	newClassName = fmt.Sprintf(areaWndClassFormat, areaWndClassNum)
	areaWndClassNum++
	areaWndClassNumLock.Unlock()

	wc := &_WNDCLASS{
		style:			_CS_DBLCLKS | _CS_HREDRAW | _CS_VREDRAW,
		lpszClassName:	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(newClassName))),
		lpfnWndProc:		syscall.NewCallback(areaWndProc(s)),
		hInstance:		hInstance,
		hIcon:			icon,
		hCursor:			cursor,
		hbrBackground:	_HBRUSH(_NULL),		// no brush; we handle WM_ERASEBKGND
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

type _BITMAPINFO struct {
	bmiHeader	_BITMAPINFOHEADER
	bmiColors	[32]uintptr	// we don't use it; make it an arbitrary number that wouldn't cause issues
}

type _BITMAPINFOHEADER struct {
	biSize			uint32
	biWidth			int32
	biHeight			int32
	biPlanes			uint16
	biBitCount		uint16
	biCompression		uint32
	biSizeImage		uint32
	biXPelsPerMeter	int32
	biYPelsPerMeter	int32
	biClrUsed			uint32
	biClrImportant		uint32
}

type _BLENDFUNCTION struct {
	BlendOp				byte
	BlendFlags			byte
	SourceConstantAlpha	byte
	AlphaFormat			byte
}

// AlphaBlend() takes a BLENDFUNCTION value
func (b _BLENDFUNCTION) arg() (x uintptr) {
	// little endian
	x = uintptr(b.AlphaFormat) << 24
	x |= uintptr(b.SourceConstantAlpha) << 16
	x |= uintptr(b.BlendFlags) << 8
	x |= uintptr(b.BlendOp)
	return x
}

type _PAINTSTRUCT struct {
	hdc			_HANDLE
	fErase		int32		// originally BOOL
	rcPaint		_RECT
	fRestore		int32		// originally BOOL
	fIncUpdate	int32		// originally BOOL
	rgbReserved	[32]byte
}

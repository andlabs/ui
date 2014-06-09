// 24 february 2014

package ui

import (
	"fmt"
//	"syscall"
	"unsafe"
)

// For Windows, Microsoft just hands you a list of preferred control sizes as part of the MSDN documentation and tells you to roll with it.
// These sizes are given in "dialog units", which are independent of the font in use.
// We need to convert these into standard pixels, which requires we get the device context of the OS window.
// References:
// - http://msdn.microsoft.com/en-us/library/windows/desktop/aa511279.aspx#controlsizing for control sizes
// - http://msdn.microsoft.com/en-us/library/ms645502%28VS.85%29.aspx - the calculation needed
// - http://support.microsoft.com/kb/125681 - to get the base X and Y
// (thanks to http://stackoverflow.com/questions/58620/default-button-size)
// For push buttons, date/time pickers, links (which we don't use), toolbars, and rebars (another type of toolbar), Common Controls version 6 provides convenient methods to use instead, falling back to the old way if it fails.

// As we are left with incomplete data, an arbitrary size will be chosen
const (
	defaultWidth = 100		// 2 * preferred width of buttons
)

type dlgunits struct {
	width	int
	height	int
	longest	bool		// TODO actually use this
	getsize	uintptr
}

var stdDlgSizes = [nctypes]dlgunits{
	c_button:		dlgunits{
		width:	50,
		height:	14,
		getsize:	_BCM_GETIDEALSIZE,
	},
	c_checkbox:	dlgunits{
		// widtdh is not defined here so assume longest
		longest:	true,
		height:	10,
	},
	c_combobox:	dlgunits{
		// technically the height of a combobox has to include the drop-down list (this is a historical accident: originally comboboxes weren't drop-down)
		// but since we're forcing Common Controls version 6, we can take advantage of one of its mechanisms to automatically fix this mistake (bad practice but whatever)
		// see also: http://blogs.msdn.com/b/oldnewthing/archive/2006/03/10/548537.aspx
		// note that the Microsoft guidelines pages don't take the list size into account
		longest:	true,
		height:	12,		// from http://msdn.microsoft.com/en-us/library/windows/desktop/bb226818%28v=vs.85%29.aspx; the page linked above says 14
	},
	c_lineedit:	dlgunits{
		longest:	true,
		height:	14,
	},
	c_label:		dlgunits{
		longest:	true,
		height:	8,
	},
	c_listbox:		dlgunits{
		longest:	true,
		// height is not clearly defined here ("an integral number of items (3 items minimum)") so just use a three-line edit control
		height:	14 + 10 + 10,
	},
	c_progressbar:		dlgunits{
		width:	237,		// the first reference says 107 also works; TODO decide which to use
		height:	8,
	},
	// TODO area
}

var (
	_selectObject = gdi32.NewProc("SelectObject")
	_getTextExtentPoint32 = gdi32.NewProc("GetTextExtentPoint32W")
	_getTextMetrics = gdi32.NewProc("GetTextMetricsW")
	_getWindowDC = user32.NewProc("GetWindowDC")
	_releaseDC = user32.NewProc("ReleaseDC")
)

// This function runs on uitask; call the functions directly.
func (s *sysData) preferredSize() (width int, height int) {
	// the preferred size of an Area is its size
	if s.ctype == c_area {
		return s.areawidth, s.areaheight
	}

	if msg := stdDlgSizes[s.ctype].getsize; msg != 0 {
		var size _SIZE

		r1, _, _ := _sendMessage.Call(
			uintptr(s.hwnd),
			msg,
			uintptr(0),
			uintptr(unsafe.Pointer(&size)))
		if r1 != uintptr(_FALSE) {		// success
			return int(size.cx), int(size.cy)
		}
		// otherwise the message approach failed, so fall back to the regular approach
		println("message failed; falling back")
	}

	var dc _HANDLE
	var tm _TEXTMETRICS
	var baseX, baseY int

	r1, _, err := _getWindowDC.Call(uintptr(s.hwnd))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error getting DC for preferred size calculations: %v", err))
	}
	dc = _HANDLE(r1)
	r1, _, err = _selectObject.Call(
		uintptr(dc),
		uintptr(controlFont))
	if r1 == 0  {		// failure
		panic(fmt.Errorf("error loading control font into device context for preferred size calculation: %v", err))
	}
	r1, _, err = _getTextMetrics.Call(
		uintptr(dc),
		uintptr(unsafe.Pointer(&tm)))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error getting text metrics for preferred size calculations: %v", err))
	}
	baseX = int(tm.tmAveCharWidth)		// TODO not optimal; third reference has better way
	baseY = int(tm.tmHeight)
	r1, _, err = _releaseDC.Call(
		uintptr(s.hwnd),
		uintptr(dc))
	if r1 == 0 {		// failure
		panic(fmt.Errorf("error releasing DC for preferred size calculations: %v", err))
	}

	// now that we have the conversion factors...
	width = stdDlgSizes[s.ctype].width
	if width == 0 {
		width = defaultWidth
	}
	height = stdDlgSizes[s.ctype].height
	width = muldiv(width, baseX, 4)		// equivalent to right of rect
	height = muldiv(height, baseY, 8)		// equivalent to bottom of rect
	return width, height
}

var (
	_mulDiv = kernel32.NewProc("MulDiv")
)

func muldiv(ma int, mb int, div int) int {
	// div will not be 0 in the usages above
	// we also ignore overflow; that isn't likely to happen for our use case anytime soon
	r1, _, _ := _mulDiv.Call(
		uintptr(int32(ma)),
		uintptr(int32(mb)),
		uintptr(int32(div)))
	return int(int32(r1))
}

type _SIZE struct {
	cx	int32		// originally LONG
	cy	int32
}

type _TEXTMETRICS struct {
	tmHeight				int32
	tmAscent				int32
	tmDescent			int32
	tmInternalLeading		int32
	tmExternalLeading		int32
	tmAveCharWidth		int32
	tmMaxCharWidth		int32
	tmWeight				int32
	tmOverhang			int32
	tmDigitizedAspectX		int32
	tmDigitizedAspectY		int32
	tmFirstChar			uint16
	tmLastChar			uint16
	tmDefaultChar			uint16
	tmBreakChar			uint16
	tmItalic				byte
	tmUnderlined			byte
	tmStruckOut			byte
	tmPitchAndFamily		byte
	tmCharSet			byte
}

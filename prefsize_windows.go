// 24 february 2014
package ui

import (
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

// As we are left with incomplete data, an arbitrary size will be chosen
const (
	defaultWidth = 100		// 2 * preferred width of buttons
)

type dlgunits struct {
	width	int
	height	int
	longest	bool		// TODO actually use this
}

var stdDlgSizes = [nctypes]dlgunits{
	c_button:		dlgunits{
		width:	50,
		height:	14,
	},
	c_checkbox:	dlgunits{
		// widtdh is not defined here so assume longest
		longest:	true,
		height:	10,
	},
	c_combobox:	dlgunits{
		longest:	true,
		height:	14,
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
}

var (
	_getTextMetrics = gdi32.NewProc("GetTextMetricsW")
	_getWindowDC = user32.NewProc("GetWindowDC")
	_releaseDC = user32.NewProc("ReleaseDC")
)

// This function runs on uitask; call the functions directly.
func (s *sysData) preferredSize() (width int, height int) {
	var dc _HANDLE
	var tm _TEXTMETRICS
	var baseX, baseY int

	// TODO use GetDC() and not GetWindowDC()?
	r1, _, err := _getWindowDC.Call(uintptr(s.hwnd))
	if r1 == 0 {		// failure
		panic(err)		// TODO return it instead
	}
	dc = _HANDLE(r1)
	r1, _, err = _getTextMetrics.Call(
		uintptr(dc),
		uintptr(unsafe.Pointer(&tm)))
	if r1 == 0 {		// failure
		panic(err)		// TODO return it instead
	}
	baseX = int(tm.tmAveCharWidth)		// TODO not optimal; third reference has better way
	baseY = int(tm.tmHeight)
	r1, _, err = _releaseDC.Call(
		uintptr(s.hwnd),
		uintptr(dc))
	if r1 == 0 {		// failure
		panic(err)		// TODO return it instead
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

// attempts to mimic the behavior of kernel32.MulDiv()
// caling it directly would be better (TODO)
// alternatively TODO make sure the rounding is correct
func muldiv(ma int, mb int, div int) int {
	xa := int64(ma) * int64(mb)
	xa /= int64(div)
	return int(xa)
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

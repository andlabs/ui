// 4 august 2014

package ui

// #include "winapi_windows.h"
import "C"

type sizing struct {
	sizingbase

	// for size calculations
	baseX           C.int
	baseY           C.int
	internalLeading C.LONG // for Label; see Label.commitResize() for details

	// for the actual resizing
	// possibly the HDWP
}

// For Windows, Microsoft just hands you a list of preferred control sizes as part of the MSDN documentation and tells you to roll with it.
// These sizes are given in "dialog units", which are independent of the font in use.
// We need to convert these into standard pixels, which requires we get the device context of the OS window.
// References:
// - http://msdn.microsoft.com/en-us/library/ms645502%28VS.85%29.aspx - the calculation needed
// - http://support.microsoft.com/kb/125681 - to get the base X and Y
// (thanks to http://stackoverflow.com/questions/58620/default-button-size)
// In my tests (see https://github.com/andlabs/windlgunits), the GetTextExtentPoint32() option for getting the base X produces much more accurate results than the tmAveCharWidth option when tested against the sample values given in http://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing, but can be off by a pixel in either direction (probably due to rounding errors).

// note on MulDiv():
// div will not be 0 in the usages below
// we also ignore overflow; that isn't likely to happen for our use case anytime soon

func fromdlgunitsX(du int, d *sizing) int {
	return int(C.MulDiv(C.int(du), d.baseX, 4))
}

func fromdlgunitsY(du int, d *sizing) int {
	return int(C.MulDiv(C.int(du), d.baseY, 8))
}

const (
	// shared by multiple containers
	marginDialogUnits  = 7
	paddingDialogUnits = 4
)

func beginResize(hwnd C.HWND) (d *sizing) {
	var baseX, baseY C.int
	var internalLeading C.LONG

	d = new(sizing)

	C.calculateBaseUnits(hwnd, &baseX, &baseY, &internalLeading)
	d.baseX = baseX
	d.baseY = baseY
	d.internalLeading = internalLeading

	d.xpadding = fromdlgunitsX(paddingDialogUnits, d)
	d.ypadding = fromdlgunitsY(paddingDialogUnits, d)

	return d
}

func marginRectDLU(r *C.RECT, top int, bottom int, left int, right int, d *sizing) {
	r.left += C.LONG(fromdlgunitsX(left, d))
	r.top += C.LONG(fromdlgunitsY(top, d))
	r.right -= C.LONG(fromdlgunitsX(right, d))
	r.bottom -= C.LONG(fromdlgunitsY(bottom, d))
}

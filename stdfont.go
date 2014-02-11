// 10 february 2014
package main

import (
//	"syscall"
	"unsafe"
)

const (
	SPI_GETNONCLIENTMETRICS = 0x0029
	LF_FACESIZE = 32		// from wingdi.h
)

type LOGFONT struct {
	lfHeight			int32
	lfWidth			int32
	lfEscapement		int32
	lfOrientation		int32
	lfWeight			int32
	lfItalic			byte
	lfUnderline		byte
	lfStrikeOut		byte
	lfCharSet			byte
	lfOutPrecision		byte
	lfClipPrecision		byte
	lfQuality			byte
	lfPitchAndFamily	byte
	lfFaceName		[LF_FACESIZE]uint16
}

type NONCLIENTMETRICS struct {
	cbSize			uint32
	iBorderWidth		int
	iScrollWidth		int
	iScrollHeight		int
	iCaptionWidth		int
	iCaptionHeight		int
	lfCaptionFont		LOGFONT
	iSmCaptionWidth	int
	iSmCaptionHeight	int
	lfSmCaptionFont	LOGFONT
	iMenuWidth		int
	iMenuHeight		int
	lfMenuFont		LOGFONT
	lfStatusFont		LOGFONT
	lfMessageFont		LOGFONT
}

var (
	systemParametersInfo = user32.NewProc("SystemParametersInfoW")
	createFontIndirect = gdi32.NewProc("CreateFontIndirectW")
)

// TODO adorn errors with which step failed?
func getStandardWindowFont() (hfont HANDLE, err error) {
	var ncm NONCLIENTMETRICS

	ncm.cbSize = uint32(unsafe.Sizeof(ncm))
	r1, _, err := systemParametersInfo.Call(
		uintptr(SPI_GETNONCLIENTMETRICS),
		uintptr(unsafe.Sizeof(ncm)),
		uintptr(unsafe.Pointer(&ncm)),
		0)
	if r1 == 0 {		// failure
		return NULL, err
	}
	// TODO does this specify an error?
	r1, _, err = createFontIndirect.Call(uintptr(unsafe.Pointer(&ncm.lfMessageFont)))
	if r1 == 0 {		// failure
		return NULL, err
	}
	return HANDLE(r1), nil
}

// 10 february 2014

package ui

import (
	"fmt"
//	"syscall"
	"unsafe"
)

var (
	controlFont		_HANDLE
	titleFont			_HANDLE
	smallTitleFont		_HANDLE
	menubarFont		_HANDLE
	statusbarFont		_HANDLE
)

const (
	_SPI_GETNONCLIENTMETRICS = 0x0029
	_LF_FACESIZE = 32		// from wingdi.h
)

type _LOGFONT struct {
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
	lfFaceName		[_LF_FACESIZE]uint16
}

type _NONCLIENTMETRICS struct {
	cbSize			uint32
	iBorderWidth		int32		// originally int
	iScrollWidth		int32		// originally int
	iScrollHeight		int32		// originally int
	iCaptionWidth		int32		// originally int
	iCaptionHeight		int32		// originally int
	lfCaptionFont		_LOGFONT
	iSmCaptionWidth	int32		// originally int
	iSmCaptionHeight	int32		// originally int
	lfSmCaptionFont	_LOGFONT
	iMenuWidth		int32		// originally int
	iMenuHeight		int32		// originally int
	lfMenuFont		_LOGFONT
	lfStatusFont		_LOGFONT
	lfMessageFont		_LOGFONT
}

var (
	_systemParametersInfo = user32.NewProc("SystemParametersInfoW")
	_createFontIndirect = gdi32.NewProc("CreateFontIndirectW")
)

// TODO the lfMessageFont doesn't seem like the right one for controls but that's all I could find for what people actually use; also I need to return the other ones and check HWND types to make sure I apply the right font to the right thing...
func getStandardWindowFonts() (err error) {
	var ncm _NONCLIENTMETRICS

	ncm.cbSize = uint32(unsafe.Sizeof(ncm))
	r1, _, err := _systemParametersInfo.Call(
		uintptr(_SPI_GETNONCLIENTMETRICS),
		uintptr(unsafe.Sizeof(ncm)),
		uintptr(unsafe.Pointer(&ncm)),
		0)
	if r1 == 0 {		// failure
		return fmt.Errorf("error getting system parameters: %v", err)
	}

	getfont := func(which *_LOGFONT, what string) (_HANDLE, error) {
		// TODO does this specify an error?
		r1, _, err = _createFontIndirect.Call(uintptr(unsafe.Pointer(which)))
		if r1 == 0 {		// failure
			return _NULL, fmt.Errorf("error getting %s font", what, err)
		}
		return _HANDLE(r1), nil
	}

	controlFont, err = getfont(&ncm.lfMessageFont, "control")
	if err != nil {
		return err
	}
	titleFont, err = getfont(&ncm.lfCaptionFont, "titlebar")
	if err != nil {
		return err
	}
	smallTitleFont, err = getfont(&ncm.lfSmCaptionFont, "small titlebar")
	if err != nil {
		return err
	}
	menubarFont, err = getfont(&ncm.lfMenuFont, "menubar")
	if err != nil {
		return err
	}
	statusbarFont, err = getfont(&ncm.lfStatusFont, "statusbar")
	if err != nil {
		return err
	}
	return nil		// all good
}

// 8 february 2014

package ui

import (
//	"syscall"
	"unsafe"
)

var (
	_adjustWindowRectEx = user32.NewProc("AdjustWindowRectEx")
	_createWindowEx = user32.NewProc("CreateWindowExW")
	_getClientRect = user32.NewProc("GetClientRect")
	_moveWindow = user32.NewProc("MoveWindow")
	_setWindowLong = user32.NewProc("SetWindowLongW")
	_setWindowPos = user32.NewProc("SetWindowPos")
	_setWindowText = user32.NewProc("SetWindowTextW")
	_showWindow = user32.NewProc("ShowWindow")
)

type _MINMAXINFO struct {
	ptReserved		_POINT
	ptMaxSize		_POINT
	ptMaxPosition		_POINT
	ptMinTrackSize		_POINT
	ptMaxTrackSize	_POINT
}

func (l _LPARAM) MINMAXINFO() *_MINMAXINFO {
	return (*_MINMAXINFO)(unsafe.Pointer(l))
}

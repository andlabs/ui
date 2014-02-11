// 11 february 2014
package main

import (
	"syscall"
	"unsafe"
)

type sysData struct {
	cSysData

	hwnd	_HWND
	cid		_HMENU
}

type classData struct {
	name	uintptr
	style		uint32
	xstyle	uint32
}

//const controlstyle = _WS_CHILD | _WS_VISIBLE | _WS_TABSTOP
//const controlxstyle = 0

var classTypes = [nctypes]*classData{
	c_window:	&classData{
		name:	uintptr(unsafe.Pointer(windowclass)),
		style:	xxxx,
		xstyle:	xxxx,
	},
//	c_button:		&classData{
//		name:	uintptr(unsafe.Pointer("BUTTON"))
//		style:	_BS_PUSHBUTTON | controlstyle,
//		xstyle:	0 | controlxstyle,
//	},
}

func (s *sysData) make() (err error) {
	
}

func (s *sysData) show() (err error) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		os_showWindow,
		p:		[]uintptr{uintptr(s.hwnd, _SW_SHOW},
		ret:		ret,
	}
	r := <-ret
	close(ret)
	return r.err
}

func (s *sysData) hide() (err error) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_showWindow,
		p:		[]uintptr{uintptr(s.hwnd), _SW_HIDE},
		ret:		ret,
	}
	r := <-ret
	close(ret)
	return r.err
}

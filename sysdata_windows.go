// 11 february 2014
package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

type sysData struct {
	cSysData

	hwnd			_HWND
	cid				_HMENU
	shownAlready		bool
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
		name:	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(stdWndClass))),
		style:	_WS_OVERLAPPEDWINDOW,
		xstyle:	0,
	},
//	c_button:		&classData{
//		name:	uintptr(unsafe.Pointer("BUTTON"))
//		style:	_BS_PUSHBUTTON | controlstyle,
//		xstyle:	0 | controlxstyle,
//	},
}

func (s *sysData) make() (err error) {
	sysDatasLock.Lock()
	defer sysDatasLock.Unlock()

	ret := make(chan uiret)
	defer close(ret)
	ct := classTypes[s.ctype]
	uitask <- &uimsg{
		call:		_createWindowEx,	
		p:		[]uintptr{
			uintptr(ct.xstyle),
			ct.name,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s.text))),
			uintptr(ct.style),
			uintptr(_CW_USEDEFAULT),		// TODO
			uintptr(_CW_USEDEFAULT),
			uintptr(_CW_USEDEFAULT),
			uintptr(_CW_USEDEFAULT),
			uintptr(_NULL),					// TODO parent
			uintptr(s.cid),
			uintptr(hInstance),
			uintptr(_NULL),
		},
		ret:	ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		return r.err
	}
	s.hwnd = _HWND(r.ret)
	addSysData(s.hwnd, s)
	// TODO parent
	return nil
}

var (
	_updateWindow = user32.NewProc("UpdateWindow")
)

// if the object is a window, we need to do the following the first time
// 	ShowWindow(hwnd, nCmdShow);
// 	UpdateWindow(hwnd);
// otherwise we go ahead and show the object normally with SW_SHOW
func (s *sysData) show() (err error) {
	if s.ctype != c_window {		// don't do the init ShowWindow/UpdateWindow chain on non-windows
		s.shownAlready = true
	}
	show := uintptr(_SW_SHOW)
	if !s.shownAlready {
		show = uintptr(nCmdShow)
	}
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_showWindow,
		p:		[]uintptr{uintptr(s.hwnd), show},
		ret:		ret,
	}
	r := <-ret
	if r.err != nil {
		return r.err
	}
	if !s.shownAlready {
		uitask <- &uimsg{
			call:		_updateWindow,
			p:		[]uintptr{uintptr(s.hwnd)},
			ret:		ret,
		}
		r = <-ret
		if r.ret == 0 {		// failure
			return fmt.Errorf("error updating window for the first time: %v", r.err)
		}
		s.shownAlready = true
	}
	return nil
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

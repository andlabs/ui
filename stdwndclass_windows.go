// 8 february 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
	"sync"
)

const (
	stdWndClassFormat = "gouiwnd%X"
)

var (
	curWndClassNum uintptr
	curWndClassNumLock sync.Mutex
)

var (
	_defWindowProc = user32.NewProc("DefWindowProcW")
)

func defWindowProc(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
	r1, _, _ := _defWindowProc.Call(
		uintptr(hwnd),
		uintptr(uMsg),
		uintptr(wParam),
		uintptr(lParam))
	return _LRESULT(r1)
}

// don't worry about error returns from GetWindowLongPtr()/SetWindowLongPtr()
// see comments of http://blogs.msdn.com/b/oldnewthing/archive/2014/02/03/10496248.aspx

func getWindowLongPtr(hwnd _HWND, what uintptr) uintptr {
	r1, _, _ := _getWindowLongPtr.Call(
		uintptr(hwnd),
		what)
	return r1
}

func setWindowLongPtr(hwnd _HWND, what uintptr, value uintptr) {
	_setWindowLongPtr.Call(
		uintptr(hwnd),
		what,
		value)
}

// we can store a pointer in extra space provided by Windows
// we'll store sysData there
// see http://blogs.msdn.com/b/oldnewthing/archive/2005/03/03/384285.aspx
func getSysData(hwnd _HWND) *sysData {
	return (*sysData)(unsafe.Pointer(getWindowLongPtr(hwnd, negConst(_GWLP_USERDATA))))
}

func storeSysData(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
	// we can get the lpParam from CreateWindowEx() in WM_NCCREATE and WM_CREATE
	// we can freely skip any messages that come prior
	// see http://blogs.msdn.com/b/oldnewthing/archive/2005/04/22/410773.aspx and http://blogs.msdn.com/b/oldnewthing/archive/2014/02/03/10496248.aspx (note the date on the latter one!)
	if uMsg == _WM_NCCREATE {
		// the lpParam to CreateWindowEx() is the first uintptr of the CREATESTRUCT
		// so rather than create that whole structure, we'll just grab the uintptr at the address pointed to by lParam
		cs := (*uintptr)(unsafe.Pointer(lParam))
		saddr := *cs
		setWindowLongPtr(hwnd, negConst(_GWLP_USERDATA), saddr)
		// don't set s; we return here
	}
	// TODO is this correct for WM_NCCREATE? I think the above link does it but I'm not entirely sure...
	return defWindowProc(hwnd, uMsg, wParam, lParam)
}

func stdWndProc(unused *sysData) func(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
	return func(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
		s := getSysData(hwnd)
		if s == nil {		// not yet saved
			return storeSysData(hwnd, uMsg, wParam, lParam)
		}

		switch uMsg {
		case _WM_COMMAND:
			id := _HMENU(wParam.LOWORD())
			s.childrenLock.Lock()
			ss := s.children[id]
			s.childrenLock.Unlock()
			switch ss.ctype {
			case c_button:
				if wParam.HIWORD() == _BN_CLICKED {
					ss.signal()
				}
			}
			return 0
		case _WM_GETMINMAXINFO:
			mm := lParam.MINMAXINFO()
			// ... minimum size
			_ = mm
			return 0
		case _WM_SIZE:
			if s.resize != nil {
				var r _RECT

				r1, _, err := _getClientRect.Call(
					uintptr(hwnd),
					uintptr(unsafe.Pointer(&r)))
				if r1 == 0 {
					panic("GetClientRect failed: " + err.Error())
				}
				// top-left corner is (0,0) so no need for winheight
				s.doResize(int(r.left), int(r.top), int(r.right), int(r.bottom), 0)
				// TODO use the Defer movement functions here?
				// TODO redraw window and all children here?
			}
			return 0
		case _WM_CLOSE:
			s.signal()
			return 0
		default:
			return defWindowProc(hwnd, uMsg, wParam, lParam)
		}
		panic(fmt.Sprintf("stdWndProc message %d did not return: internal bug in ui library", uMsg))
	}
}

type _WNDCLASS struct {
	style				uint32
	lpfnWndProc		uintptr
	cbClsExtra		int32		// originally int
	cbWndExtra		int32		// originally int
	hInstance			_HANDLE
	hIcon			_HANDLE
	hCursor			_HANDLE
	hbrBackground	_HBRUSH
	lpszMenuName	*uint16
	lpszClassName		uintptr
}

var (
	icon, cursor _HANDLE
)

var (
	_registerClass = user32.NewProc("RegisterClassW")
)

func registerStdWndClass(s *sysData) (newClassName string, err error) {
	curWndClassNumLock.Lock()
	newClassName = fmt.Sprintf(stdWndClassFormat, curWndClassNum)
	curWndClassNum++
	curWndClassNumLock.Unlock()

	wc := &_WNDCLASS{
		lpszClassName:	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(newClassName))),
		lpfnWndProc:		syscall.NewCallback(stdWndProc(s)),
		hInstance:		hInstance,
		hIcon:			icon,
		hCursor:			cursor,
		hbrBackground:	_HBRUSH(_COLOR_BTNFACE + 1),
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

// no need to use/recreate MAKEINTRESOURCE() here as the Windows constant generator already took care of that because Microsoft's headers do already
func initWndClassInfo() (err error) {
	r1, _, err := user32.NewProc("LoadIconW").Call(
		uintptr(_NULL),
		uintptr(_IDI_APPLICATION))
	if r1 == 0 {		// failure
		return fmt.Errorf("error getting window icon: %v", err)
	}
	icon = _HANDLE(r1)

	r1, _, err = user32.NewProc("LoadCursorW").Call(
		uintptr(_NULL),
		uintptr(_IDC_ARROW))
	if r1 == 0 {		// failure
		return fmt.Errorf("error getting window cursor: %v", err)
	}
	cursor = _HANDLE(r1)

	return nil
}

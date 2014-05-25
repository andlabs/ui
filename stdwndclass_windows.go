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
	defWindowProc = user32.NewProc("DefWindowProcW")
)

func stdWndProc(s *sysData) func(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
	return func(hwnd _HWND, uMsg uint32, wParam _WPARAM, lParam _LPARAM) _LRESULT {
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
			r1, _, _ := defWindowProc.Call(
				uintptr(hwnd),
				uintptr(uMsg),
				uintptr(wParam),
				uintptr(lParam))
			return _LRESULT(r1)
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

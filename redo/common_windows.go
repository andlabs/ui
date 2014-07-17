// 12 july 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

//export xpanic
func xpanic(msg *C.char, lasterr C.DWORD) {
	panic(fmt.Errorf("%s: %s", C.GoString(msg), syscall.Errno(lasterr))
}

//export xmissedmsg
func xmissedmsg(purpose *C.char, f *C.char, uMsg C.UINT) {
	panic(fmt.Errorf("%s window procedure message %d does not return a value (bug in %s)", C.GoString(purpose), uMsg, C.GoString(f)))
}

func getWindowText(hwnd uintptr) string {
	// WM_GETTEXTLENGTH and WM_GETTEXT return the count /without/ the terminating null character
	// but WM_GETTEXT expects the buffer size handed to it to /include/ the terminating null character
	n := f_SendMessageW(hwnd, c_WM_GETTEXTLENGTH, 0, 0)
	buf := make([]uint16, int(n + 1))
	if f_SendMessageW(hwnd, c_WM_GETTEXT,
		t_WPARAM(n + 1), t_LPARAM(uintptr(unsafe.Pointer(&buf[0])))) != n {
		panic(fmt.Errorf("WM_GETTEXT did not copy exactly %d characters out", n))
	}
	return syscall.UTF16ToString(buf)
}

func setWindowText(hwnd uintptr, text string, errors []t_LRESULT) {
	res := f_SendMessageW(hwnd, c_WM_SETTEXT,
		0, t_LPARAM(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text)))))
	for _, err := range errors {
		if res == err {
			panic(fmt.Errorf("WM_SETTEXT failed; error code %d", res))
		}
	}
}

func updateWindow(hwnd uintptr, caller string) {
	res, err := f_UpdateWindow(hwnd)
	if res == 0 {
		panic(fmt.Errorf("error calling UpdateWindow() from %s: %v", caller, err))
	}
}

func storelpParam(hwnd uintptr, lParam t_LPARAM) {
	var cs *s_CREATESTRUCTW

	cs = (*s_CREATESTRUCTW)(unsafe.Pointer(uintptr(lParam)))
	f_SetWindowLongPtrW(hwnd, c_GWLP_USERDATA, cs.lpCreateParams)
}

func (w t_WPARAM) HIWORD() uint16 {
	u := uintptr(w) & 0xFFFF0000
	return uint16(u >> 16)
}

func (w t_WPARAM) LOWORD() uint16 {
	u := uintptr(w) & 0x0000FFFF
	return uint16(u)
}

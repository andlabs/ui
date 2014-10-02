// 12 july 2014

package ui

import (
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

//export xpanic
func xpanic(msg *C.char, lasterr C.DWORD) {
	panic(fmt.Errorf("%s: %s", C.GoString(msg), syscall.Errno(lasterr)))
}

//export xpanichresult
func xpanichresult(msg *C.char, hresult C.HRESULT) {
	panic(fmt.Errorf("%s; HRESULT: 0x%X", C.GoString(msg), hresult))
}

//export xpaniccomdlg
func xpaniccomdlg(msg *C.char, err C.DWORD) {
	panic(fmt.Errorf("%s; comdlg32.dll extended error: 0x%X", C.GoString(msg), err))
}

//export xmissedmsg
func xmissedmsg(purpose *C.char, f *C.char, uMsg C.UINT) {
	panic(fmt.Errorf("%s window procedure message %d does not return a value (bug in %s)", C.GoString(purpose), uMsg, C.GoString(f)))
}

func toUTF16(s string) C.LPWSTR {
	return C.LPWSTR(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

func getWindowText(hwnd C.HWND) string {
	// WM_GETTEXTLENGTH and WM_GETTEXT return the count /without/ the terminating null character
	// but WM_GETTEXT expects the buffer size handed to it to /include/ the terminating null character
	n := C.getWindowTextLen(hwnd)
	buf := make([]uint16, int(n+1))
	C.getWindowText(hwnd, C.WPARAM(n),
		C.LPWSTR(unsafe.Pointer(&buf[0])))
	return syscall.UTF16ToString(buf)
}

func wstrToString(wstr *C.WCHAR) string {
	n := C.wcslen((*C.wchar_t)(unsafe.Pointer(wstr)))
	xbuf := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(wstr)),
		Len:  int(n + 1),
		Cap:  int(n + 1),
	}
	buf := (*[]uint16)(unsafe.Pointer(xbuf))
	return syscall.UTF16ToString(*buf)
}

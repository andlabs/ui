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
	panic(fmt.Errorf("%s: %s", C.GoString(msg), syscall.Errno(lasterr)))
}

//export xmissedmsg
func xmissedmsg(purpose *C.char, f *C.char, uMsg C.UINT) {
	panic(fmt.Errorf("%s window procedure message %d does not return a value (bug in %s)", C.GoString(purpose), uMsg, C.GoString(f)))
}

func toUTF16(s string) C.LPCWSTR {
	return C.LPCWSTR(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

func getWindowText(hwnd C.HWND) string {
	// WM_GETTEXTLENGTH and WM_GETTEXT return the count /without/ the terminating null character
	// but WM_GETTEXT expects the buffer size handed to it to /include/ the terminating null character
	n := C.getWindowTextLen(hwnd)
	buf := make([]uint16, int(n + 1))
	C.getWindowText(hwnd, C.WPARAM(n),
		C.LPCWSTR(unsafe.Pointer(&buf[0])))
	return syscall.UTF16ToString(buf)
}

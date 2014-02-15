// 7 february 2014
package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

// MessageBox button types.
const (
	_MB_ABORTRETRYIGNORE = 0x00000002
	_MB_CANCELTRYCONTINUE = 0x00000006
	_MB_HELP = 0x00004000
	_MB_OK = 0x00000000
	_MB_OKCANCEL = 0x00000001
	_MB_RETRYCANCEL = 0x00000005
	_MB_YESNO = 0x00000004
	_MB_YESNOCANCEL = 0x00000003
)

// MessageBox icon types.
const (
	_MB_ICONEXCLAMATION = 0x00000030
	_MB_ICONWARNING = 0x00000030
	_MB_ICONINFORMATION = 0x00000040
	_MB_ICONASTERISK = 0x00000040
	_MB_ICONQUESTION = 0x00000020
	_MB_ICONSTOP = 0x00000010
	_MB_ICONERROR = 0x00000010
	_MB_ICONHAND = 0x00000010
)

// MessageBox default button types.
const (
	_MB_DEFBUTTON1 = 0x00000000
	_MB_DEFBUTTON2 = 0x00000100
	_MB_DEFBUTTON3 = 0x00000200
	_MB_DEFBUTTON4 = 0x00000300
)

// MessageBox modality types.
const (
	_MB_APPLMODAL = 0x00000000
	_MB_SYSTEMMODAL = 0x00001000
	_MB_TASKMODAL = 0x00002000
)

// MessageBox miscellaneous types.
const (
	_MB_DEFAULT_DESKTOP_ONLY = 0x00020000
	_MB_RIGHT = 0x00080000
	_MB_RTLREADING = 0x00100000
	_MB_SETFOREGROUND = 0x00010000
	_MB_TOPMOST = 0x00040000
	_MB_SERVICE_NOTIFICATION = 0x00200000
)

// MessageBox return values.
const (
	_IDABORT = 3
	_IDCANCEL = 2
	_IDCONTINUE = 11
	_IDIGNORE = 5
	_IDNO = 7
	_IDOK = 1
	_IDRETRY = 4
	_IDTRYAGAIN = 10
	_IDYES = 6
)

var (
	_messageBox = user32.NewProc("MessageBoxW")
)

func msgBox(lpText string, lpCaption string, uType uint32) (result int) {
	r1, _, err := _messageBox.Call(
		uintptr(_NULL),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpText))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpCaption))),
		uintptr(uType))
	if r1 == 0 {		// failure
		panic(fmt.Sprintf("error displaying message box to user: %v\nstyle: 0x%08X\ntitle: %q\ntext:\n%s", err, uType, lpCaption, lpText))
	}
	return int(r1)
}

// MsgBox displays an informational message box to the user with just an OK button.
func MsgBox(title string, textfmt string, args ...interface{}) {
	// TODO add an icon?
	msgBox(fmt.Sprintf(textfmt, args...), title, _MB_OK)
}

// MsgBoxError displays a message box to the user with just an OK button and an icon indicating an error.
func MsgBoxError(title string, textfmt string, args ...interface{}) {
	msgBox(fmt.Sprintf(textfmt, args...), title, _MB_OK | _MB_ICONERROR)
}

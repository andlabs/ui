// 7 february 2014

package ui

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// TODO change what the default window titles are?

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

func _msgBox(primarytext string, secondarytext string, uType uint32) (result int) {
	// http://msdn.microsoft.com/en-us/library/windows/desktop/aa511267.aspx says "Use task dialogs whenever appropriate to achieve a consistent look and layout. Task dialogs require Windows Vista® or later, so they aren't suitable for earlier versions of Windows. If you must use a message box, separate the main instruction from the supplemental instruction with two line breaks."
	text := primarytext
	if secondarytext != "" {
		text += "\n\n" + secondarytext
	}
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_messageBox,
		p:		[]uintptr{
			uintptr(_NULL),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(os.Args[0]))),
			uintptr(uType),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		panic(fmt.Sprintf("error displaying message box to user: %v\nstyle: 0x%08X\ntitle: %q\ntext:\n%s", r.err, uType, os.Args[0], text))
	}
	return int(r.ret)
}

func msgBox(primarytext string, secondarytext string) {
	_msgBox(primarytext, secondarytext, _MB_OK)
}

func msgBoxError(primarytext string, secondarytext string) {
	_msgBox(primarytext, secondarytext, _MB_OK | _MB_ICONERROR)
}

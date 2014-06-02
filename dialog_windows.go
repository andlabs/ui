// 7 february 2014

package ui

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// TODO change what the default window titles are?

var (
	_messageBox = user32.NewProc("MessageBoxW")
)

func _msgBox(primarytext string, secondarytext string, uType uint32) (result int) {
	// http://msdn.microsoft.com/en-us/library/windows/desktop/aa511267.aspx says "Use task dialogs whenever appropriate to achieve a consistent look and layout. Task dialogs require Windows Vista® or later, so they aren't suitable for earlier versions of Windows. If you must use a message box, separate the main instruction from the supplemental instruction with two line breaks."
	text := primarytext
	if secondarytext != "" {
		text += "\n\n" + secondarytext
	}
	uType |= _MB_TASKMODAL		// make modal to every window in the program (they're all windows of the uitask, which is a single thread)
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

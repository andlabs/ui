// 25 february 2014
package ui

import (
	"fmt"
//	"syscall"
	"unsafe"
)

// pretty much every constant here except _WM_USER is from commctrl.h
// TODO for all: filter out constants not available in Windows 2000

// InitCommonControlsEx constants.
const (
	_ICC_LISTVIEW_CLASSES = 0x00000001
	_ICC_TREEVIEW_CLASSES = 0x00000002
	_ICC_BAR_CLASSES = 0x00000004
	_ICC_TAB_CLASSES = 0x00000008
	_ICC_UPDOWN_CLASS = 0x00000010
	_ICC_PROGRESS_CLASS = 0x00000020
	_ICC_HOTKEY_CLASS = 0x00000040
	_ICC_ANIMATE_CLASS = 0x00000080
	_ICC_WIN95_CLASSES = 0x000000FF
	_ICC_DATE_CLASSES = 0x00000100
	_ICC_USEREX_CLASSES = 0x00000200
	_ICC_COOL_CLASSES = 0x00000400
	_ICC_INTERNET_CLASSES = 0x00000800
	_ICC_PAGESCROLLER_CLASS = 0x00001000
	_ICC_NATIVEFNTCTL_CLASS = 0x00002000
	_ICC_STANDARD_CLASSES = 0x00004000
	_ICC_LINK_CLASS = 0x00008000
)

var (
	_initCommonControlsEx = comctl32.NewProc("InitCommonControlsEx")
)

func initCommonControls() (err error) {
	var icc struct {
		dwSize	uint32
		dwICC	uint32
	}

	icc.dwSize = uint32(unsafe.Sizeof(icc))
	icc.dwICC = _ICC_PROGRESS_CLASS
	r1, _, err := _initCommonControlsEx.Call(uintptr(unsafe.Pointer(&icc)))
	if r1 == _FALSE {		// failure
		// TODO does it set GetLastError()?
		return fmt.Errorf("error initializing Common Controls (comctl32.dll): %v", err)
	}
	return nil
}

// Common Controls class names.
const (
	_PROGRESS_CLASS = "msctls_progress32"
)

// Shared Common Controls styles.
const (
	_WM_USER = 0x0400
	_CCM_FIRST = 0x2000
	_CCM_SETBKCOLOR = (_CCM_FIRST + 1)
)

// Progress Bar styles.
const (
	_PBS_SMOOTH = 0x01
	_PBS_VERTICAL = 0x04
)

// Progress Bar messages.
const (
	_PBM_SETRANGE = (_WM_USER + 1)
	_PBM_SETPOS = (_WM_USER + 2)
	_PBM_DELTAPOS = (_WM_USER + 3)
	_PBM_SETSTEP = (_WM_USER + 4)
	_PBM_STEPIT = (_WM_USER + 5)
	_PBM_SETRANGE32 = (_WM_USER + 6)
	_PBM_GETRANGE = (_WM_USER + 7)
	_PBM_GETPOS = (_WM_USER + 8)
	_PBM_SETBARCOLOR = (_WM_USER + 9)
	_PBM_SETBKCOLOR = _CCM_SETBKCOLOR
)

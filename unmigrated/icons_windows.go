// 8 february 2014
package ui

import (
//	"syscall"
//	"unsafe"
)

// Predefined icon resource IDs.
const (
	IDI_APPLICATION = 32512
	IDI_ASTERISK = 32516
	IDI_ERROR = 32513
	IDI_EXCLAMATION = 32515
	IDI_HAND = 32513
	IDI_INFORMATION = 32516
	IDI_QUESTION = 32514
	IDI_SHIELD = 32518
	IDI_WARNING = 32515
	IDI_WINLOGO = 32517
)

var (
	loadIcon = user32.NewProc("LoadIconW")
)

func LoadIcon_ResourceID(hInstance HANDLE, lpIconName uint16) (icon HANDLE, err error) {
	r1, _, err := loadIcon.Call(
		uintptr(hInstance),
		MAKEINTRESOURCE(lpIconName))
	if r1 == 0 {		// failure
		return NULL, err
	}
	return HANDLE(r1), nil
}

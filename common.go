// 7 february 2014
package main

import (
	"syscall"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
)

type HWND uintptr

const (
	NULL = 0
)


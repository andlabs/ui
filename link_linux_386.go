// +build !windows
// +build !darwin

// 11 december 2015

package ui

// #cgo LDFLAGS: ${SRCDIR}/libui_linux_386.a -lm -ldl
// #cgo pkg-config: gtk+-3.0
import "C"

func ensureMainThread() {
	// do nothing; GTK+ doesn't care which thread we're on so long as we don't change it after starting
}

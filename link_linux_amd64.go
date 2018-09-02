// +build !windows
// +build !darwin

// 11 december 2015

package ui

// #cgo LDFLAGS: ${SRCDIR}/libui_linux_amd64.a -lm -ldl
// #cgo pkg-config: gtk+-3.0
import "C"

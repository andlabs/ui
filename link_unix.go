// +build !windows
// +build !darwin

// 11 december 2015

package ui

// #cgo LDFLAGS: -L${SRCDIR} -lui -Wl,-rpath=$ORIGIN
import "C"

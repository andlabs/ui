// +build !windows
// +build !darwin

// 11 december 2015

package ui

// #cgo LDFLAGS: -L${SRCDIR} -lui -Wl,-rpath=$ORIGIN
import "C"

func ensureMainThread() {
	// do nothing; GTK+ doesn't care which thread we're on so long as we don't change it after starting
}

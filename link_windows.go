// 13 december 2015

package ui

// #cgo LDFLAGS: -L${SRCDIR} -lui
import "C"

func ensureMainThread() {
	// do nothing; Windows doesn't care which thread we're on so long as we don't change it after starting
}

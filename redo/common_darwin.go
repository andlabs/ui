// 20 july 2014

package ui

// #include "objc_darwin.h"
import "C"

func fromBOOL(b C.BOOL) bool {
	if b != C.NO {
		return true
	}
	return false
}

func toBOOL(b bool) C.BOOL {
	if b == true {
		return C.YES
	}
	return C.NO
}

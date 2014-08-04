// 4 august 2014

package ui

import (
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type container struct {
	view		C.id
	*sizer
}

func newContainer(child Control) *container {
	c := &container{
		sizer:	new(sizer),
	}
	c.view = C.newContainerView(unsafe.Pointer(c))
	c.child = child
	c.child.setParent(&controlParent{c.view})
	return c
}

//export containerResized
func containerResized(data unsafe.Pointer, width C.intptr_t, height C.intptr_t) {
	c := (*container)(unsafe.Pointer(data))
	c.resize(0, 0, int(width), int(height))
}

// 11 february 2014
//package ui
package main

import (
	"syscall"
	"unsafe"
	"runtime"
)

var uitask chan *uimsg

type uimsg struct {
	call		*syscall.LazyProc
	p		[]uintptr
	ret		chan uiret
}

type uiret struct {
	ret		uintptr
	err		error
}

func ui(initDone chan error) {
	runtime.LockOSThread()

	// initialize hInstance
	// initialize nCmdShow
	// initialize the common window class
	uitask = make(chan *uimsg)
	initDone <- nil

	for m := range uitask {
		r1, _, err := m.msg.Call(m.p...)
		m.ret <- uiret{
			ret:	r1,
			err:	err,
		}
	}
}

// +build !windows,!darwin,!plan9

// 16 february 2014
//package ui
package main

import (
	"fmt"
	"runtime"
)

var uitask chan func()

func ui(initDone chan error) {
	runtime.LockOSThread()

	uitask = make(chan func())
	if gtk_init() != true {
		initDone <- fmt.Errorf("gtk_init failed (reason unknown; TODO)")
		return
	}
	initDone <- nil

	gtk_main()
}

func uistep() {
	select {
	case f := <-uitask:
		f()
	default:		// do not block
	}
}

// temporary
func MsgBox(string, string, ...interface{}) {}
func MsgBoxError(title string, text string, args ...interface{}) {panic(title+"\n"+fmt.Sprintf(text,args...))}

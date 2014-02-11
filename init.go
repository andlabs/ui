// 11 february 2014
//package ui
package main

// TODO this will be system-defined
func initpanic(err error) {
	panic("failure during init: " + err.Error())
}

func init() {
	initDone := make(chan error)
	go ui(initDone)
	err := <-initDone
	if err != nil {
		initpanic(err)
	}
}

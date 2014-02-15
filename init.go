// 11 february 2014
//package ui
package main

import (
	"os"
)

func init() {
	initDone := make(chan error)
	go ui(initDone)
	err := <-initDone
	if err != nil {
		// TODO provide copying instructions? will need to be system-specific
		MsgBoxError("UI Library Init Failure",
			"A failure occured during UI library initialization:\n%v\n" +
			"Please report this to the application developer or on http://github.com/andlabs/ui.",
			err)
		os.Exit(1)
	}
}

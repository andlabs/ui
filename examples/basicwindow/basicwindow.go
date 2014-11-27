package main

import (
	"github.com/andlabs/ui"
	"log"
)

func main() {
	// This runs the code that displays our GUI.
	// All code that interfaces with package ui (except event handlers) must be run from within a ui.Do() call.
	go ui.Do(gui)

	err := ui.Go()
	if err != nil {
		log.Print(err)
	}
}

func gui() {
	// All windows must have a control inside.
	// ui.Space() creates a control that is just a blank space for us to use.
	newControl := ui.Space()

	// Then we create a window.
	w := ui.NewWindow("Window", 280, 350, newControl)

	// We tell package ui to destroy our window and shut down cleanly when the user closes the window by clicking the X button in the titlebar.
	w.OnClosing(func() bool {
		// This informs package ui to shut down cleanly when it can.
		ui.Stop()
		// And this informs package ui that we want to hide AND destroy the window.
		return true
	})

	// And finally, we need to show the window.
	w.Show()
}

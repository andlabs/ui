package main

import (
	"github.com/andlabs/ui"
	"log"
)

func main() {

	go ui.Do(gui)

	err := ui.Go()
	if err != nil {
		log.Print(err)
	}
}

func gui() {

	// Here we create a new space
	newControl := ui.Space()

    // Then we create a window
	w := ui.NewWindow("Window", 280, 350, newControl)
	w.OnClosing(func() bool {
		ui.Stop()
		return true
	})

	w.Show()
}

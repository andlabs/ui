// 12 march 2014

/*
Package ui is a simple package which provides a way to write portable GUI programs quickly and easily. It aims to run on as many systems as Go itself, but presently has support for Windows, Mac OS X, and other Unix systems using the Windows API, Cocoa, and GTK+ 3, respectively. It also aims to be Go-like: events are transmitted via channels, and the library is fully safe for concurrent use.

To use the library, place your main program code in another function and call Go(), passing that function as a parameter. (This is necessary due to threading restrictions on some environments, such as Cocoa.) Once in the function you pass to Go(), you can safely use the rest of the library. When this function returns, so does Go(), and package functions become unavailable.

Building GUIs is as simple as creating a Window, populating it with Controls, and then calling Open() on the Window. A Window only has one Control: you pack multiple Controls into a Window by arranging them in layouts (Layouts are also Controls). There are presently two Layouts, Stack and Grid, each with different semantics on sizing and placement. See their documentation.

Once a Window is open, you cannot make layout or event channel changes.

Once your Window is open, you can begin to handle events. Handling events is simple: because all events are channels exposed as exported members of the Window and Control types, simply select on them. Event channels are initialized by default. However, before you Open a Window, you can freely reassign event channels, such that multiple events trigger the same channel, making event logic more compact. You may also choose not to handle events; events are sent asynchronously so the GUI loop is not initerrupted.

Here is a simple, complete program that asks the user for their name and greets them after clicking a button.
	package main

	import (
		"github.com/andlabs/ui"
	)

	func myMain() {
		w := ui.NewWindow("Hello", 400, 100)
		ui.AppQuit = w.Closing		// treat quitting the application like closing the main window
		nameField := ui.NewLineEdit("Enter Your Name Here")
		button := ui.NewButton("Click Here For a Greeting")
		w.Open(ui.NewVerticalStack(nameField, button))

		for {
			select {
			case <-w.Closing:		// user tries to close the window or quit the program
				return
			case <-button.Clicked:	// user clicked the button
				ui.MsgBox("Hello, " + nameField.Text() + "!", "")
			}
		}
	}

	func main() {
		err := ui.Go(myMain)
		if err != nil {
			panic(err)
		}
	}
*/
package ui

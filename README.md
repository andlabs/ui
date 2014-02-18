# Native UI library for Go
### THIS PACKAGE IS UNSTABLE AND PRELIMINARY. In fact, it is presently compiled as package `main` for ease of cross-platform testing and debugging. Once major issues are dealt with and the Mac OS X build working, I will likely move to packge `ui` and move `main()` to a test.

This is a simple library for building cross-platform GUI programs in Go. It targets Windows and all Unix variants (except Mac OS X until further notice) and provides a thread-safe, channel-based API.

There is documentation, but due to the note above, you won't be able to see it just yet. Refer to `main.go` for an example.

## Future Readme Contents
This is a simple library for building cross-platform GUI programs in Go. It targets Windows, Mac OS X, Linux, and other Unixes, and provides a thread-safe, channel-based API. The API itself is minimal; it aims to provide only what is necessary for GUI program design. That being said, suggestions are welcome. Layout is done using various layout managers, and some effort is taken to conform to the target platform's UI guidelines. Otherwise, the library uses native toolkits.

ui aims to run on all supported versions of supported platforms. To be more precise, the system requirements are:

* Windows: Windows 2000 or newer. The Windows backend uses package `syscall` and calls Windows DLLs directly, so does not rely on cgo.
* Mac OS X: Mac OS X 10.6 (Snow Leopard) or newer. Objective-C dispatch is done by interfacing with libobjc directly, and thus this uses cgo.
* Other Unixes: The Unix backend uses GTK+, and thus cgo. It requires GTK+ 3.4 or newer; for Ubuntu this means 12.04 LTS (Precise Pangolin) at minimum. Check your distribution.

ui itself has no outside Go package dependencies; it is entirely self-contained.

To install, simply `go get` this package. On Mac OS X, make sure you have the Apple development headers. On other Unixes, make sure you have the GTK+ development files (for Ubuntu, `libgtk-3-dev` is sufficient).

Package documentation is available at http://godoc.org/github.com/andlabs/ui.

The following is an example program to illustrate what programming with ui is like:
```go
package main

import (
	"fmt"
	"github.com/andlabs/ui"
)

func main() {
	w := ui.NewWindow("Main Window", 320, 240)
	w.Closing = ui.Event()
	b := ui.NewButton("Click Me")
	b2 := ui.NewButton("Or Me")
	s2 := ui.NewStack(Horizontal, b, b2)
	c := ui.NewCheckbox("Check Me")
	cb1 := ui.NewCombobox(true, "You can edit me!", "Yes you can!", "Yes you will!")
	cb2 := ui.NewCombobox(false, "You can't edit me!", "No you can't!", "No you won't!")
	e := ui.NewLineEdit("Enter text here too")
	l := ui.NewLabel("This is a label")
	b3 := ui.NewButton("List Info")
	s3 := ui.NewStack(ui.Horizontal, l, b3)
	s0 := ui.NewStack(ui.Vertical, s2, c, cb1, cb2, e, s3)
	lb1 := ui.NewListbox(true, "Select One", "Or More", "To Continue")
	lb2 := ui.NewListbox(false, "Select", "Only", "One", "Please")
	i := 0
	doAdjustments := func() {
		cb1.Append("append")
		cb2.InsertBefore(fmt.Sprintf("before %d", i), 1)
		lb1.InsertBefore(fmt.Sprintf("%d", i), 2)
		lb2.Append("Please")
		i++
	}
	doAdjustments()
	s1 := ui.NewStack(ui.Vertical, lb2, lb1)
	s := ui.NewStack(ui.Horizontal, s1, s0)
	err := w.Open(s)
	if err != nil {
		panic(err)
	}

mainloop:
	for {
		select {
		case <-w.Closing:
			break mainloop
		case <-b.Clicked:
			err = w.SetTitle(fmt.Sprintf("%v | %s | %s | %s",
				c.Checked(),
				cb1.Selection(),
				cb2.Selection(),
				e.Text()))
			if err != nil {
				panic(err)
			}
			doAdjustments()
		case <-b2.Clicked:
			cb1.Delete(1)
			cb2.Delete(2)
			lb1.Delete(3)
			lb2.Delete(4)
		case <-b3.Clicked:
			MsgBox("List Info",
				"cb1: %d %q\ncb2: %d %q\nlb1: %d %q\nlb2: %d %q",
				cb1.SelectedIndex(), cb1.Selection(),
				cb2.SelectedIndex(), cb2.Selection(),
				lb1.SelectedIndices(), lb1.Selection(),
				lb2.SelectedIndices(), lb2.Selection())
		}
	}
	w.Hide()
}
```

## Contributing
Contributions are welcome. File issues, pull requests, approach me on IRC (pietro10 in #go-nuts; andlabs elsewhere), etc. Even suggestions are welcome: while I'm mainly drawing from my own GUI programming experience, everyone is different.

If you want to dive in, read implementation.md: this is a description of how the library works. (Feel free to suggest improvements to this as well.) The other .md files in this repository contain various development notes.

Please suggest documentation improvements as well.

# Go UI package planning
Pietro Gagliardi
http://github.com/andlabs

## Goals
- Simple, easy to use GUI library for hard-coding GUI layouts
- Go-like: uses Go's concurrency features, interfaces, etc. and behaves like other Go libraries
- Portable; runs on all OSs Go supports and uses native toolkits (wherever possible)
- Minimal: only support what's absolutely necessary (for instance, only events that we will actually use in a program); if functionality can be done cleanly in an existing thing, use that (for instnaces, if adjustable sliding dividers are ever added, they can be made part of `Stack` instead of their own thing)
- Lightweight and fast
- Error-safe
- Correct: uses APIs properly and conforms to system-specific UI design guidelines

## Layouts
Layouts control positioning and sizing. Layouts are controls, so they can be added recursively. The layout types are:
* `Stack`: a stack of controls, all sized alike, with padding between controls and spacing around the whole set. Controls can be arranged horizontally or vertically. (Analogues: Qt's `QBoxLayout`)
>* TODO change the name?
* `RadioSet`: like `Stack` but for radio buttons: only has radio buttons and handles exclusivity automatically (this is also the only way to add radio buttons)
* `Grid`: a grid of controls; they size themselves. Spacing is handled like `Stack`. (Analogues: Qt's `QGridLayout`)
* `Form`: a set of label-control pairs arranged to resemble options on a dialog form. Sizing, positioning, and spacing are handled in an OS-dependent way. (Analogues: Qt's `QFormLayout`)

## Windows
There's only one (maybe two, if I choose to add floating toolboxes) window type. You can add one control to the content area of a window.

In the case of dialogue boxes, you can call a function, say `RunDaialogue()` , that runs the dialogue modal, and adds standard OK/Cancel/Apply buttons for you.

## An example
``` go
package main

import (
	"github.com/andlabs/ui"
)

func main() {
	win := ui.NewWindow("Hello")
	form := ui.NewForm()
	name := ui.NewLineEntry()
	form.Append("Enter your name:", name)
	button := ui.NewButton("Click Me")
	form.Append("", button)
	win.SetControl(form)
	
	events, err := win.RunDialogue(ui.OkCancel)
	if err != nil {
		panic(err)
	}
	done := false
	for !done {
		select {
		case event := <-events:
			switch event {
			case ui.Ok:
				ui.MsgBox("Hi", "Hello, " + name.Text(), ui.Ok)
			case ui.Cancel:
				done = true
			}
		case <-button.Click:
			ui.MsgBox("Hi", "You clicked me!", ui.Ok)
		}
	}
	window.Close()
}
```

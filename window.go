// 11 february 2014

package ui

import (
	"fmt"
)

// Window represents an on-screen window.
type Window struct {
	created    bool
	sysData    *sysData
	initTitle  string
	initWidth  int
	initHeight int
	shownOnce  bool
	spaced	bool
	handler	WindowHandler
}

// WindowHandler represents an event handler for a Window and all its child Controls.
// 
// When an event on a Window or one of its child Controls comes in, the respect Window's handler's Event() method is called. The method call occurs on the main thread, and thus any call to any package ui method can be performed.
// 
// Each Event() call takes two parameters: the event ID and a data argument. For most events, the data argument is a pointer to the Control that triggered the event.
// 
// For Closing, the data argument is a pointer to a bool variable. If, after returning from Event, the value of this variable is true, the Window is closed; if false, the Window is not closed. The default value on entry to the function is [TODO].
// 
// For any event >= CustomEvent, the data argument is the argument passed to the Window's SendEvent() method.
type WindowHandler interface {
	Event(e Event, data interface{})
}

// Event represents an event; see WindowHandler for details.
// All event values >= CustomEvent are available for program use.
type Event int
const (
	Closing Event = iota			// Window close
	Clicked					// Button click
	Dismissed					// Dialog closed
	CustomEvent = 5000		// very high number; higher than the package would ever need, anyway
)

// NewWindow allocates a new Window with the given title and size. The window is not created until a call to Create() or Open().
func NewWindow(title string, width int, height int, handler WindowHandler) *Window {
	return &Window{
		sysData:    mksysdata(c_window),
		initTitle:  title,
		initWidth:  width,
		initHeight: height,
		handler:		handler,
	}
}

// SetTitle sets the window's title.
func (w *Window) SetTitle(title string) {
	if w.created {
		w.sysData.setText(title)
		return
	}
	w.initTitle = title
}

// SetSize sets the window's size.
func (w *Window) SetSize(width int, height int) (err error) {
	if w.created {
		err := w.sysData.setWindowSize(width, height)
		if err != nil {
			return fmt.Errorf("error setting window size: %v", err)
		}
		return nil
	}
	w.initWidth = width
	w.initHeight = height
	return nil
}

// SetSpaced sets whether the Window's child control takes padding and spacing into account.
// That is, with w.SetSpaced(true), w's child will have a margin around the window frame and will have sub-controls separated by an implementation-defined amount.
// Currently, only Stack and Grid explicitly understand this property.
// This property is visible recursively throughout the widget hierarchy of the Window.
// This property cannot be set after the Window has been created.
func (w *Window) SetSpaced(spaced bool) {
	if w.created {
		panic(fmt.Errorf("Window.SetSpaced() called after window created"))
	}
	w.spaced = spaced
}

// Open creates the Window with Create and then shows the Window with Show. As with Create, you cannot call Open more than once per window.
func (w *Window) Open(control Control) {
	done := make(chan struct{})
	defer close(done)
	touitask(func() {
		w.Create(control)
		w.Show()
		done <- struct{}{}
	})
	<-done
}

// Create creates the Window, setting its control to the given control. It does not show the window. This can only be called once per window, and finalizes all initialization of the control.
func (w *Window) Create(control Control) {
	done := make(chan struct{})
	defer close(done)
	touitask(func() {
		if w.created {
			panic("window already open")
		}
		w.sysData.spaced = w.spaced
		w.sysData.winhandler = w.handler
		w.sysData.close = func(b *bool) {
			w.sysData.winhandler.Event(Closing, b)
		}
		err := w.sysData.make(nil)
		if err != nil {
			panic(fmt.Errorf("error opening window: %v", err))
		}
		if control != nil {
			w.sysData.allocate = control.allocate
			err = control.make(w.sysData)
			if err != nil {
				panic(fmt.Errorf("error adding window's control: %v", err))
			}
		}
		err = w.sysData.setWindowSize(w.initWidth, w.initHeight)
		if err != nil {
			panic(fmt.Errorf("error setting window size (in Window.Open()): %v", err))
		}
		w.sysData.setText(w.initTitle)
		w.created = true
		done <- struct{}{}
	})
	<-done
}

// Show shows the window.
func (w *Window) Show() {
	if !w.shownOnce {
		w.shownOnce = true
		err := w.sysData.firstShow()
		if err != nil {
			panic(fmt.Errorf("error showing window for the first time: %v", err))
		}
		return
	}
	w.sysData.show()
}

// Hide hides the window.
func (w *Window) Hide() {
	w.sysData.hide()
}

// Center centers the Window on-screen.
// The concept of "screen" in the case of a multi-monitor setup is implementation-defined.
// It presently panics if the Window has not been created.
func (w *Window) Center() {
	if !w.created {
		panic("attempt to center Window before it has been created")
	}
	w.sysData.center()
}

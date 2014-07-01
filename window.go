// 11 february 2014

package ui

import (
	"fmt"
)

// Window represents an on-screen window.
type Window struct {
	// Closing is called when the Close button is pressed by the user, or when the application needs to be quit (should the underlying system provide a concept of application distinct from window).
	// Return true to allow the window to be closed; false otherwise.
	// You cannot change this field after the Window has been created.
	// [TODO close vs. hide]
	// If Closing is nil, a default which rejects the close will be used.
	Closing		func() bool

	created    bool
	sysData    *sysData
	initTitle  string
	initWidth  int
	initHeight int
	shownOnce  bool
	spaced	bool
}

// NewWindow allocates a new Window with the given title and size. The window is not created until a call to Create() or Open().
func NewWindow(title string, width int, height int) *Window {
	return &Window{
		sysData:    mksysdata(c_window),
		initTitle:  title,
		initWidth:  width,
		initHeight: height,
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
	w.create(control, true)
}

// Create creates the Window, setting its control to the given control. It does not show the window. This can only be called once per window, and finalizes all initialization of the control.
func (w *Window) Create(control Control) {
	w.create(control, false)
}

func (w *Window) create(control Control, show bool) {
	touitask(func() {
		if w.created {
			panic("window already open")
		}
		w.sysData.spaced = w.spaced
		w.sysData.close = w.Closing
		if w.sysData.close == nil {
			w.sysData.close = func() bool {
				return false
			}
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
		if show {
			w.Show()
		}
	})
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

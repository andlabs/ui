// 11 february 2014

package ui

import (
	"fmt"
	"sync"
)

// Window represents an on-screen window.
type Window struct {
	// Closing gets a message when the user clicks the window's close button.
	// You cannot change it once the Window has been created.
	// If you do not respond to this signal, nothing will happen; regardless of whether you handle the signal or not, the window will not be closed.
	Closing		chan struct{}

	lock			sync.Mutex
	created		bool
	sysData		*sysData
	initTitle		string
	initWidth		int
	initHeight		int
}

// NewWindow creates a new window with the given title and size. The window is not constructed at the OS level until a call to Open().
func NewWindow(title string, width int, height int) *Window {
	return &Window{
		sysData:		mksysdata(c_window),
		initTitle:		title,
		initWidth:		width,
		initHeight:	height,
		Closing:		newEvent(),
	}
}

// SetTitle sets the window's title.
func (w *Window) SetTitle(title string) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.created {
		w.sysData.setText(title)
		return
	}
	w.initTitle = title
}

// SetSize sets the window's size.
func (w *Window) SetSize(width int, height int) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

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

// Open opens the window, setting its control to the given control, and then shows the window. This can only be called once per window, and finalizes all initialization of the control.
// TODO rename?
func (w *Window) Open(control Control) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.created {
		// TODO return an error instead?
		panic("window already open")
	}
	w.sysData.event = w.Closing
	err = w.sysData.make(nil)
	if err != nil {
		return fmt.Errorf("error opening window: %v", err)
	}
	if control != nil {
		w.sysData.resize = control.setRect
		err = control.make(w.sysData)
		if err != nil {
			return fmt.Errorf("error adding window's control: %v", err)
		}
	}
	err = w.sysData.setWindowSize(w.initWidth, w.initHeight)
	if err != nil {
		return fmt.Errorf("error setting window size (in Window.Open()): %v", err)
	}
	w.sysData.setText(w.initTitle)
	// TODO separate showing?
	err = w.sysData.firstShow()
	if err != nil {
		return fmt.Errorf("error showing window (in Window.Open()): %v", err)
	}
	w.created = true
	return nil
}

// Show shows the window.
func (w *Window) Show() {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.sysData.show()
}

// Hide hides the window.
func (w *Window) Hide() {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.sysData.hide()
}

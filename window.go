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
	Closing chan struct{}

	lock       sync.Mutex
	created    bool
	sysData    *sysData
	initTitle  string
	initWidth  int
	initHeight int
	shownOnce  bool
}

// NewWindow allocates a new Window with the given title and size. The window is not created until a call to Create() or Open().
func NewWindow(title string, width int, height int) *Window {
	return &Window{
		sysData:    mksysdata(c_window),
		initTitle:  title,
		initWidth:  width,
		initHeight: height,
		Closing:    newEvent(),
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

// Open creates the Window with Create and then shows the Window with Show. As with Create, you cannot call Open more than once per window.
func (w *Window) Open(control Control) {
	w.Create(control)
	w.Show()
}

// Create creates the Window, setting its control to the given control. It does not show the window. This can only be called once per window, and finalizes all initialization of the control.
func (w *Window) Create(control Control) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.created {
		panic("window already open")
	}
	w.sysData.event = w.Closing
	err := w.sysData.make(nil)
	if err != nil {
		panic(fmt.Errorf("error opening window: %v", err))
	}
	if control != nil {
		w.sysData.resize = control.setRect
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
}

// Show shows the window.
func (w *Window) Show() {
	w.lock.Lock()
	defer w.lock.Unlock()

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
	w.lock.Lock()
	defer w.lock.Unlock()

	w.sysData.hide()
}

// Center centers the Window on-screen.
// The concept of "screen" in the case of a multi-monitor setup is implementation-defined.
// It presently panics if the Window has not been created.
func (w *Window) Center() {
	w.lock.Lock()
	defer w.lock.Unlock()

	if !w.created {
		panic("attempt to center Window before it has been created")
	}
	w.sysData.center()
}

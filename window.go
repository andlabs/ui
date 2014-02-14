// 11 february 2014
//package ui
package main

import (
	"sync"
)

// TODO adorn errors in each stage with which stage failed?

// Window represents an on-screen window.
type Window struct {
	// If this channel is non-nil, the event loop will receive on this when the user clicks the window's close button.
	// This channel can only be set before initially opening the window.
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
	}
}

// SetTitle sets the window's title.
func (w *Window) SetTitle(title string) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.created {
		return w.sysData.setText(title)
	}
	w.initTitle = title
	return nil
}

// SetSize sets the window's size.
func (w *Window) SetSize(width int, height int) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.created {
		panic("TODO")
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
	err = w.sysData.make(w.initTitle, w.initWidth, w.initHeight, nil)
	if err != nil {
		return err
	}
	if control != nil {
		w.sysData.resize = control.setRect
		err = control.apply(w.sysData)
		if err != nil {
			return err
		}
	}
	w.created = true
	return w.sysData.show()
}

// Show shows the window.
func (w *Window) Show() (err error) {
	return w.sysData.show()
}

// Hide hides the window.
func (w *Window) Hide() (err error) {
	return w.sysData.hide()
}

// These satisfy the Control interface, allowing a window to own a control. As a consequence, Windows are themselves Controls. THIS IS UNDOCUMENTED AND UNSUPPORTED. I can make it supported in the future, but for now, no. You shouldn't be depending on the internals of the library to develop your programs: if the documentation is incomplete and/or wrong, get the person responsible to fix it, as the documentation, not the implementation, is your contract to what you can or cannot do. Don't worry, this package is in good company: Go itself was designed spec-first for this reason.
// If I decide not to support windows as controls, a better way to deal with controls would be in order. Perhaps separate interfaces...? Making Windows Controls seems the cleanest option for now (and with correct usage of the library costs nothing).
func (w *Window) apply(window *sysData) error {
	panic("Window.apply() should never be called")
}
func (w *Window) setRect(x int, y int, width int, height int) error {
	panic("Window.setRect() should never be called")
}

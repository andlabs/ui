// 11 february 2014
//package ui
package main

import (
	"fmt"
	"sync"
)

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
		err = w.sysData.setText(title)
		if err != nil {
			return fmt.Errorf("error setting window title: %v", err)
		}
		return nil
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
		return fmt.Errorf("error opening window: %v", err)
	}
	if control != nil {
		w.sysData.resize = control.setRect
		err = control.make(w.sysData)
		if err != nil {
			return fmt.Errorf("error adding window's control: %v", err)
		}
	}
	// TODO resize window to apply control sizes
	// TODO separate showing?
	err = w.sysData.show()
	if err != nil {
		return fmt.Errorf("error showing window (in Window.Open()): %v", err)
	}
	w.created = true
	return nil
}

// Show shows the window.
func (w *Window) Show() (err error) {
	err = w.sysData.show()
	if err != nil {
		return fmt.Errorf("error showing window: %v", err)
	}
	return nil
}

// Hide hides the window.
func (w *Window) Hide() (err error) {
	err = w.sysData.hide()
	if err != nil {
		return fmt.Errorf("error hiding window: %v", err)
	}
	return nil
}

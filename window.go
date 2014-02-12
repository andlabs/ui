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
	Closing	chan struct{}

	lock		sync.Mutex
	created	bool
	control	Control
	sysData	*sysData
}

// NewWindow creates a new window with the given title. The window is not constructed at the OS level until a call to Open().
func NewWindow(title string) *Window {
	return &Window{
		sysData:	&sysData{
			cSysData:		cSysData{
				ctype:	c_window,
				text:		title,
			},
		},
	}
}

// SetControl sets the window's central control to control.
func (w *Window) SetControl(control Control) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.control = control
	err = w.control.unapply()
	if err != nil {
		return err
	}
	w.control.setParent(w)
	if w.created {
		err = w.control.apply()
		if err != nil {
			return err
		}
	}
	return nil
}

// Open opens the window. If the OS window has not been created yet, this function will.
func (w *Window) Open() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// If the window has already been created, show it.
	if !w.created {
		w.sysData.closing = w.Closing
		err = w.sysData.make()
		if err != nil {
			return err
		}
		if w.control != nil {
			err = w.control.apply()
			if err != nil {
				return err
			}
		}
	}
	return w.sysData.show()
}

// Close closes the window. The window is not destroyed; it is merely hidden.
// TODO don't send on w.Closing
func (w *Window) Close() (err error) {
	return w.sysData.hide()
}

// These satisfy the Control interface, allowing a window to own a control. As a consequence, Windows are themselves Controls. THIS IS UNDOCUMENTED AND UNSUPPORTED. I can make it supported in the future, but for now, no. You shouldn't be depending on the internals of the library to develop your programs: if the documentation is incomplete and/or wrong, get the person responsible to fix it, as the documentation, not the implementation, is your contract to what you can or cannot do. Don't worry, this package is in good company: Go itself was designed spec-first for this reason.
// If I decide not to support windows as controls, a better way to deal with controls would be in order. Perhaps separate interfaces...? Making Windows Controls seems the cleanest option for now (and with correct usage of the library costs nothing).
func (w *Window) apply() error {
	panic("Window.apply() should never be called")
}
func (w *Window) unapply() error {
	panic("Window.unapply() should never be called")
}
func (w *Window) setParent(c Control) {
	panic("Window.setParent() should never be called")
}

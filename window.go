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
	control		Control
	sysData		*sysData
	initTitle		string
	initWidth		int
	initHeight		int
}

// NewWindow creates a new window with the given title and size. The window is not constructed at the OS level until a call to Open().
func NewWindow(title string, width int, height int) *Window {
	return &Window{
		sysData:		&sysData{
			cSysData:		cSysData{
				ctype:	c_window,
			},
		},
		initTitle:		title,
		initWidth:		width,
		initHeight:	height,
	}
}

// SetControl sets the window's central control to control. This function cannot be called once the window has been opened.
func (w *Window) SetControl(control Control) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.created {
		panic("cannot set window control after window has been opened")
	}
	w.control = control
	w.control.setParent(w)
	return nil
}

// SetTitle sets the window's title.
func (w *Window) SetTitle(title string) (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.created {
		panic("TODO")
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

// Open opens the window. If the OS window has not been created yet, this function will.
func (w *Window) Open() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// If the window has already been created, show it.
	if !w.created {
		w.sysData.closing = w.Closing
		err = w.sysData.make(w.initTitle, w.initWidth, w.initHeight)
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
func (w *Window) setParent(c Control) {
	panic("Window.setParent() should never be called")
}
func (w *Window) parentWindow() *Window {
	panic("Window.parentWindow() should never be called")
}

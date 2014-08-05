// +build !windows,!darwin

// 25 july 2014

package ui

import (
	"unsafe"
)

// #include "gtk_unix.h"
import "C"

type tab struct {
	_widget		*C.GtkWidget
	notebook		*C.GtkNotebook

	tabs			[]*container
}

func newTab() Tab {
	widget := C.gtk_notebook_new()
	t := &tab{
		_widget:		widget,
		notebook:		(*C.GtkNotebook)(unsafe.Pointer(widget)),
	}
	// there are no scrolling arrows by default; add them in case there are too many tabs
	C.gtk_notebook_set_scrollable(t.notebook, C.TRUE)
	return t
}

func (t *tab) Append(name string, control Control) {
	c := newContainer(control)
	t.tabs = append(t.tabs, c)
	cname := togstr(name)
	defer freegstr(cname)
	tab := C.gtk_notebook_append_page(t.notebook,
		// TODO figure out how to keep this private
		c.layoutwidget,
		C.gtk_label_new(cname))
	if tab == -1 {
		panic("gtk_notebook_append_page() failed")
	}
}

func (t *tab) widget() *C.GtkWidget {
	return t._widget
}

func (t *tab) setParent(p *controlParent) {
	basesetParent(t, p)
}

func (t *tab) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(t, x, y, width, height, d)
}

func (t *tab) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(t, d)
}

// no need to override Control.commitResize() as only prepared the tabbed control; its children will be reallocated when that one is resized
func (t *tab) commitResize(a *allocation, d *sizing) {
	basecommitResize(t, a, d)
}

func (t *tab) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(t, d)
}

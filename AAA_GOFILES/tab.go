// 12 december 2015

package ui

import (
	"unsafe"
)

// #include "ui.h"
import "C"

// Tab is a Control that holds tabbed pages of Controls. Each tab
// has a label. The user can click on the tabs themselves to switch
// pages. Individual pages can also have margins.
type Tab struct {
	c	*C.uiControl
	t	*C.uiTab

	children	[]Control
}

// NewTab creates a new Tab.
func NewTab() *Tab {
	t := new(Tab)

	t.t = C.uiNewTab()
	t.c = (*C.uiControl)(unsafe.Pointer(t.t))

	return t
}

// Destroy destroys the Tab. If the Tab has pages,
// Destroy calls Destroy on the pages's Controls as well.
func (t *Tab) Destroy() {
	for len(t.children) != 0 {
		c := t.children[0]
		t.Delete(0)
		c.Destroy()
	}
	C.uiControlDestroy(t.c)
}

// LibuiControl returns the libui uiControl pointer that backs
// the Tab. This is only used by package ui itself and should
// not be called by programs.
func (t *Tab) LibuiControl() uintptr {
	return uintptr(unsafe.Pointer(t.c))
}

// Handle returns the OS-level handle associated with this Tab.
// On Windows this is an HWND of a standard Windows API
// WC_TABCONTROL class (as provided by Common Controls
// version 6). The pages are not children of this window and there
// currently is no way to directly access them.
// On GTK+ this is a pointer to a GtkNotebook.
// On OS X this is a pointer to a NSTabView.
func (t *Tab) Handle() uintptr {
	return uintptr(C.uiControlHandle(t.c))
}

// Show shows the Tab.
func (t *Tab) Show() {
	C.uiControlShow(t.c)
}

// Hide hides the Tab.
func (t *Tab) Hide() {
	C.uiControlHide(t.c)
}

// Enable enables the Tab.
func (t *Tab) Enable() {
	C.uiControlEnable(t.c)
}

// Disable disables the Tab.
func (t *Tab) Disable() {
	C.uiControlDisable(t.c)
}

// Append adds the given page to the end of the Tab.
func (t *Tab) Append(name string, child Control) {
	t.InsertAt(name, len(t.children), child)
}

// InsertAt adds the given page to the Tab such that it is the
// nth page of the Tab (starting at 0).
func (t *Tab) InsertAt(name string, n int, child Control) {
	c := (*C.uiControl)(nil)
	if child != nil {
		c = touiControl(child.LibuiControl())
	}
	cname := C.CString(name)
	// TODO why is this uintmax_t and not intmax_t
	C.uiTabInsertAt(t.t, cname, C.uintmax_t(n), c)
	freestr(cname)
	ch := make([]Control, len(t.children) + 1)
	// and insert into t.children at the right place
	copy(ch[:n], t.children[:n])
	ch[n] = child
	copy(ch[n + 1:], t.children[n:])
	t.children = ch
}

// Delete deletes the nth page of the Tab.
func (t *Tab) Delete(n int) {
	t.children = append(t.children[:n], t.children[n + 1:]...)
	C.uiTabDelete(t.t, C.uintmax_t(n))
}

// NumPages returns the number of pages in the Tab.
func (t *Tab) NumPages() int {
	return len(t.children)
}

// Margined returns whether page n (starting at 0) of the Tab
// has margins around its child.
func (t *Tab) Margined(n int) bool {
	return tobool(C.uiTabMargined(t.t, C.uintmax_t(n)))
}

// SetMargined controls whether page n (starting at 0) of the Tab
// has margins around its child. The size of the margins are
// determined by the OS and its best practices.
func (t *Tab) SetMargined(n int, margined bool) {
	C.uiTabSetMargined(t.t, C.uintmax_t(n), frombool(margined))
}

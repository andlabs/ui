package ui

import (
	"unsafe"
)

// #include "ui.h"
// extern void doMenuItemOnClicked(uiMenuItem *, uiWindow *, void *);
// static inline void realuiMenuItemOnClicked(uiMenuItem *i)
// {
// 	uiMenuItemOnClicked(i, doMenuItemOnClicked, NULL);
// }
import "C"

var menuItems = make(map[*C.uiMenuItem]*MenuItem)

// Menu represents a menu on the menu bar located at the top
// of a Window. Currently, every window has the same menu bar
// and thus the same menus and menu items.
type Menu struct {
	m *C.uiMenu
}

// MenuItem represents an item in a Menu.
type MenuItem struct {
	i *C.uiMenuItem

	onClicked        func(*MenuItem, *Window)
	onClickedAllowed bool
}

// NewMenu creates a new Menu.
func NewMenu(name string) *Menu {
	m := new(Menu)

	cname := C.CString(name)
	m.m = C.uiNewMenu(cname)
	freestr(cname)

	return m
}

// AppendSeparator appends a separator to the Menu.
func (m *Menu) AppendSeparator() {
	C.uiMenuAppendSeparator(m.m)
}

// AppendItem appends a new item to the Menu.
func (m *Menu) AppendItem(name string) *MenuItem {
	cname := C.CString(name)
	i := C.uiMenuAppendItem(m.m, cname)
	freestr(cname)

	return newMenuItem(i, true)
}

// AppendCheckItem appends a check item to the Menu.
func (m *Menu) AppendCheckItem(name string) *MenuItem {
	cname := C.CString(name)
	i := C.uiMenuAppendCheckItem(m.m, cname)
	freestr(cname)

	return newMenuItem(i, true)
}

// AppendQuitItem appends a quit item to the Menu.
func (m *Menu) AppendQuitItem() *MenuItem {
	i := C.uiMenuAppendQuitItem(m.m)
	return newMenuItem(i, false)
}

// AppendPreferencesItem appends a preferences item to the Menu.
func (m *Menu) AppendPreferencesItem() *MenuItem {
	i := C.uiMenuAppendPreferencesItem(m.m)
	return newMenuItem(i, true)
}

// AppendAboutItem appends an about item to the Menu.
func (m *Menu) AppendAboutItem() *MenuItem {
	i := C.uiMenuAppendAboutItem(m.m)
	return newMenuItem(i, true)
}

// Enable enables the MenuItem.
func (i *MenuItem) Enable() {
	C.uiMenuItemEnable(i.i)
}

// Disable disables the MenuItem.
func (i *MenuItem) Disable() {
	C.uiMenuItemDisable(i.i)
}

// Checked returns the check status of a menu item.
func (i *MenuItem) Checked() bool {
	return tobool(C.uiMenuItemChecked(i.i))
}

// SetChecked sets the check status of a menu item.
func (i *MenuItem) SetChecked(checked bool) {
	C.uiMenuItemSetChecked(i.i, frombool(checked))
}

// OnClicked registers f to be run when the user clicks the MenuItem.
// Only one function can be registered at a time.
func (i *MenuItem) OnClicked(f func(*MenuItem, *Window)) {
	i.onClicked = f
}

//export doMenuItemOnClicked
func doMenuItemOnClicked(i *C.uiMenuItem, w *C.uiWindow, data unsafe.Pointer) {
	ii := menuItems[i]
	ww := windows[w]

	if ii.onClicked != nil {
		ii.onClicked(ii, ww)
	}
}

func newMenuItem(i *C.uiMenuItem, onClickedAllowed bool) *MenuItem {
	if onClickedAllowed {
		C.realuiMenuItemOnClicked(i)
	}

	result := &MenuItem{i: i, onClickedAllowed: onClickedAllowed}
	menuItems[i] = result

	return result
}

// 11 february 2014
package main

import (
	"runtime"
)

// The sysData type contains all system data. It provides the system-specific underlying implementation. It is guaranteed to have the following by embedding:
type cSysData struct {
	ctype	int
	event	chan struct{}
	resize	func(x int, y int, width int, height int) error
	alternate	bool		// editable for Combobox, multi-select for listbox
}
func (c *cSysData) make(initText string, window *sysData) error {
	panic(runtime.GOOS + " sysData does not define make()")
}
func (c *cSysData) show() error {
	panic(runtime.GOOS + " sysData does not define show()")
}
func (c *cSysData) hide() error {
	panic(runtime.GOOS + " sysData does not define hide()")
}
func (c *cSysData) setText(text string) error {
	panic(runtime.GOOS + " sysData does not define setText()")
}
func (c *cSysData) setRect(x int, y int, width int, height int) error {
	panic(runtime.GOOS + " sysData does not define setRect()")
}
func (c *cSysData) isChecked() (bool, error) {
	panic(runtime.GOOS + " sysData does not define isChecked()")
}
func (c *cSysData) text() (string, error) {
	panic(runtime.GOOS + " sysData does not define text()")
}
func (c *cSysData) append(string) error {
	panic(runtime.GOOS + " sysData does not define append()")
}
func (c *cSysData) insertBefore(string, int) error {
	panic(runtime.GOOS + " sysData does not define insertBefore()")
}
func (c *cSysData) selectedIndex() (int, error) {
	panic(runtime.GOOS + " sysData does not define selectedIndex()")
}
func (c *cSysData) selectedIndices() ([]int, error) {
	panic(runtime.GOOS + " sysData does not define selectedIndices()")
}
func (c *cSysData) selectedTexts() ([]string, error) {
	panic(runtime.GOOS + " sysData does not define selectedIndex()")
}
func (c *cSysData) setWindowSize(int, int) error {
	panic(runtime.GOOS + " sysData does not define setWindowSize()")
}

const (
	c_window = iota
	c_button
	c_checkbox
	c_combobox
	c_lineedit
	c_label
	c_listbox
	nctypes
)

func mksysdata(ctype int) *sysData {
	return &sysData{
		cSysData:		cSysData{
			ctype:	ctype,
		},
	}
}

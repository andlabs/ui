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
}
func (c *cSysData) make(initText string, initWidth int, initHeight int, window *sysData) error {
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

const (
	c_window = iota
	c_button
	c_checkbox
	nctypes
)

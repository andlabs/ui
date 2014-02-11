// 11 february 2014
package main

import (
	"runtime"
)

// The sysData type contains all system data. It provides the system-specific underlying implementation. It is guaranteed to have the following by embedding:
type cSysData struct {
	ctype	int
	text		string
}
func (c *cSysData) make() error {
	panic(runtime.GOOS + " sysData does not define make()")
}
func (c *cSysData) show() error {
	panic(runtime.GOOS + " sysData does not define show()")
}
func (c *cSysData) show() error {
	panic(runtime.GOOS + " sysData does not define hide()")
}

const (
	c_window = iota
	c_button
)

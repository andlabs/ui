// 11 february 2014
package ui

import (
	"runtime"
)

const eventbufsiz = 100		// suggested by skelterjohn

// Event returns a new channel suitable for listening for events.
func Event() chan struct{} {
	return make(chan struct{}, eventbufsiz)
}

// The sysData type contains all system data. It provides the system-specific underlying implementation. It is guaranteed to have the following by embedding:
type cSysData struct {
	ctype	int
	event	chan struct{}
	resize	func(x int, y int, width int, height int, winheight int) error
	alternate	bool		// editable for Combobox, multi-select for listbox, password for lineedit
}
func (c *cSysData) make(initText string, window *sysData) error {
	panic(runtime.GOOS + " sysData does not define make()")
}
func (c *cSysData) firstShow() error {
	panic(runtime.GOOS + " sysData does not define firstShow()")
}
func (c *cSysData) show() {
	panic(runtime.GOOS + " sysData does not define show()")
}
func (c *cSysData) hide() {
	panic(runtime.GOOS + " sysData does not define hide()")
}
func (c *cSysData) setText(text string) {
	panic(runtime.GOOS + " sysData does not define setText()")
}
func (c *cSysData) setRect(x int, y int, width int, height int, winheight int) error {
	panic(runtime.GOOS + " sysData does not define setRect()")
}
func (c *cSysData) isChecked() bool {
	panic(runtime.GOOS + " sysData does not define isChecked()")
}
func (c *cSysData) text() string {
	panic(runtime.GOOS + " sysData does not define text()")
}
func (c *cSysData) append(string) {
	panic(runtime.GOOS + " sysData does not define append()")
}
func (c *cSysData) insertBefore(string, int) {
	panic(runtime.GOOS + " sysData does not define insertBefore()")
}
func (c *cSysData) selectedIndex() int {
	panic(runtime.GOOS + " sysData does not define selectedIndex()")
}
func (c *cSysData) selectedIndices() []int {
	panic(runtime.GOOS + " sysData does not define selectedIndices()")
}
func (c *cSysData) selectedTexts() []string {
	panic(runtime.GOOS + " sysData does not define selectedIndex()")
}
func (c *cSysData) setWindowSize(int, int) error {
	panic(runtime.GOOS + " sysData does not define setWindowSize()")
}
func (c *cSysData) delete(int) error {
	panic(runtime.GOOS + " sysData does not define delete()")
}
func (c *cSysData) preferredSize() (int, int) {
	panic(runtime.GOOS + " sysData does not define preferredSize()")
}
func (c *cSysData) setProgress(int) {
	panic(runtime.GOOS + " sysData does not define setProgress()")
}
func (c *cSysData) len() int {
	panic(runtime.GOOS + " sysData does not define len()")
}

// signal sends the event signal. This raise is done asynchronously to avoid deadlocking the UI task.
// Thanks skelterjohn for this techinque: if we can't queue any more events, drop them
func (s *cSysData) signal() {
	if s.event != nil {
		go func() {
			select {
			case s.event <- struct{}{}:
			default:
			}
		}()
	}
}

const (
	c_window = iota
	c_button
	c_checkbox
	c_combobox
	c_lineedit
	c_label
	c_listbox
	c_progressbar
	nctypes
)

func mksysdata(ctype int) *sysData {
	return &sysData{
		cSysData:		cSysData{
			ctype:	ctype,
		},
	}
}

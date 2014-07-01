// 11 february 2014

package ui

// The sysData type contains all system data. It provides the system-specific underlying implementation. It is guaranteed to have the following by embedding:
type cSysData struct {
	ctype     int
	allocate    func(x int, y int, width int, height int, d *sysSizeData) []*allocation
	spaced	bool
	alternate bool        // editable for Combobox, multi-select for listbox, password for lineedit
	handler   AreaHandler // for Areas; TODO rename to areahandler
	winhandler	WindowHandler	// for Windows
	close	func() bool	// provided by each Window
	event	func()		// provided by each control
}

// this interface is used to make sure all sysDatas are synced
var _xSysData interface {
	sysDataSizingFunctions
	make(window *sysData) error
	firstShow() error
	show()
	hide()
	setText(text string)
	setRect(x int, y int, width int, height int, winheight int) error
	isChecked() bool
	text() string
	append(string)
	insertBefore(string, int)
	selectedIndex() int
	selectedIndices() []int
	selectedTexts() []string
	setWindowSize(int, int) error
	setProgress(int)
	len() int
	setAreaSize(int, int)
	repaintAll()
	center()
	setChecked(bool)
} = &sysData{} // this line will error if there's an inconsistency

const (
	c_window = iota
	c_button
	c_checkbox
	c_combobox
	c_lineedit
	c_label
	c_listbox
	c_progressbar
	c_area
	nctypes
)

func mksysdata(ctype int) *sysData {
	s := &sysData{
		cSysData: cSysData{
			ctype: ctype,
		},
	}
	return s
}

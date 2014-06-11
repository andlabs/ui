// 11 february 2014

package ui

const eventbufsiz = 100 // suggested by skelterjohn

// newEvent returns a new channel suitable for listening for events.
func newEvent() chan struct{} {
	return make(chan struct{}, eventbufsiz)
}

// The sysData type contains all system data. It provides the system-specific underlying implementation. It is guaranteed to have the following by embedding:
type cSysData struct {
	ctype     int
	event     chan struct{}
	resize    func(x int, y int, width int, height int, rr *[]resizerequest)
	resizes   []resizerequest
	alternate bool        // editable for Combobox, multi-select for listbox, password for lineedit
	handler   AreaHandler // for Areas
}

// this interface is used to make sure all sysDatas are synced
var _xSysData interface {
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
	delete(int)
	preferredSize() (int, int)
	setProgress(int)
	len() int
	setAreaSize(int, int)
	repaintAll()
	center()
} = &sysData{} // this line will error if there's an inconsistency

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
	c_area
	nctypes
)

func mksysdata(ctype int) *sysData {
	s := &sysData{
		cSysData: cSysData{
			ctype: ctype,
		},
	}
	if ctype == c_window { // make resizes non-nil so it can be passed in
		s.resizes = make([]resizerequest, 0, 0)
	}
	return s
}

type resizerequest struct {
	sysData *sysData
	x       int
	y       int
	width   int
	height  int
}

func (s *sysData) doResize(x int, y int, width int, height int, winheight int) {
	if s.resize != nil {
		s.resizes = s.resizes[0:0] // set len to 0 without changing cap
		s.resize(x, y, width, height, &s.resizes)
		for _, s := range s.resizes {
			err := s.sysData.setRect(s.x, s.y, s.width, s.height, winheight)
			if err != nil {
				panic("child resize failed: " + err.Error())
			}
		}
	}
}

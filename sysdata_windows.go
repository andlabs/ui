// 11 february 2014

package ui

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

type sysData struct {
	cSysData

	hwnd         _HWND
	children     map[_HMENU]*sysData
	nextChildID  _HMENU
	childrenLock sync.Mutex
	isMarquee    bool // for sysData.setProgress()
	// unlike with GTK+ and Mac OS X, we're responsible for sizing Area properly ourselves
	areawidth    int
	areaheight   int
	clickCounter clickCounter
	lastfocus    _HWND
}

type classData struct {
	name             *uint16
	style            uint32
	xstyle           uint32
	altStyle         uint32
	storeSysData     bool
	doNotLoadFont    bool
	appendMsg        uintptr
	insertBeforeMsg  uintptr
	deleteMsg        uintptr
	selectedIndexMsg uintptr
	selectedIndexErr uintptr
	addSpaceErr      uintptr
	lenMsg           uintptr
}

const controlstyle = _WS_CHILD | _WS_VISIBLE | _WS_TABSTOP
const controlxstyle = 0

var classTypes = [nctypes]*classData{
	c_window: &classData{
		name:          stdWndClass,
		style:         _WS_OVERLAPPEDWINDOW,
		xstyle:        0,
		storeSysData:  true,
		doNotLoadFont: true,
	},
	c_button: &classData{
		name:   toUTF16("BUTTON"),
		style:  _BS_PUSHBUTTON | controlstyle,
		xstyle: 0 | controlxstyle,
	},
	c_checkbox: &classData{
		name: toUTF16("BUTTON"),
		// don't use BS_AUTOCHECKBOX because http://blogs.msdn.com/b/oldnewthing/archive/2014/05/22/10527522.aspx
		style:  _BS_CHECKBOX | controlstyle,
		xstyle: 0 | controlxstyle,
	},
	c_combobox: &classData{
		name:             toUTF16("COMBOBOX"),
		style:            _CBS_DROPDOWNLIST | _WS_VSCROLL | controlstyle,
		xstyle:           0 | controlxstyle,
		altStyle:         _CBS_DROPDOWN | _CBS_AUTOHSCROLL | _WS_VSCROLL | controlstyle,
		appendMsg:        _CB_ADDSTRING,
		insertBeforeMsg:  _CB_INSERTSTRING,
		deleteMsg:        _CB_DELETESTRING,
		selectedIndexMsg: _CB_GETCURSEL,
		selectedIndexErr: negConst(_CB_ERR),
		addSpaceErr:      negConst(_CB_ERRSPACE),
		lenMsg:           _CB_GETCOUNT,
	},
	c_lineedit: &classData{
		name: toUTF16("EDIT"),
		// WS_EX_CLIENTEDGE without WS_BORDER will apply visual styles
		// thanks to MindChild in irc.efnet.net/#winprog
		style:    _ES_AUTOHSCROLL | controlstyle,
		xstyle:   _WS_EX_CLIENTEDGE | controlxstyle,
		altStyle: _ES_PASSWORD | _ES_AUTOHSCROLL | controlstyle,
	},
	c_label: &classData{
		name: toUTF16("STATIC"),
		// SS_NOPREFIX avoids accelerator translation; SS_LEFTNOWORDWRAP clips text past the end
		// controls are vertically aligned to the top by default (thanks Xeek in irc.freenode.net/#winapi)
		// also note that tab stops are remove dfor labels
		style:  (_SS_NOPREFIX | _SS_LEFTNOWORDWRAP | controlstyle) &^ _WS_TABSTOP,
		xstyle: 0 | controlxstyle,
		// MAKE SURE THIS IS THE SAME
		altStyle:		(_SS_NOPREFIX | _SS_LEFTNOWORDWRAP | controlstyle) &^ _WS_TABSTOP,
	},
	c_listbox: &classData{
		name: toUTF16("LISTBOX"),
		// we don't use LBS_STANDARD because it sorts (and has WS_BORDER; see above)
		// LBS_NOINTEGRALHEIGHT gives us exactly the size we want
		// LBS_MULTISEL sounds like it does what we want but it actually doesn't; instead, it toggles item selection regardless of modifier state, which doesn't work like anything else (see http://msdn.microsoft.com/en-us/library/windows/desktop/bb775149%28v=vs.85%29.aspx and http://msdn.microsoft.com/en-us/library/windows/desktop/aa511485.aspx)
		style:            _LBS_NOTIFY | _LBS_NOINTEGRALHEIGHT | _WS_VSCROLL | controlstyle,
		xstyle:           _WS_EX_CLIENTEDGE | controlxstyle,
		altStyle:         _LBS_EXTENDEDSEL | _LBS_NOTIFY | _LBS_NOINTEGRALHEIGHT | _WS_VSCROLL | controlstyle,
		appendMsg:        _LB_ADDSTRING,
		insertBeforeMsg:  _LB_INSERTSTRING,
		deleteMsg:        _LB_DELETESTRING,
		selectedIndexMsg: _LB_GETCURSEL,
		selectedIndexErr: negConst(_LB_ERR),
		addSpaceErr:      negConst(_LB_ERRSPACE),
		lenMsg:           _LB_GETCOUNT,
	},
	c_progressbar: &classData{
		name:          toUTF16(x_PROGRESS_CLASS),
		// note that tab stops are disabled for progress bars
		style:         (_PBS_SMOOTH | controlstyle) &^ _WS_TABSTOP,
		xstyle:        0 | controlxstyle,
		doNotLoadFont: true,
	},
	c_area: &classData{
		name:          areaWndClass,
		style:         areastyle,
		xstyle:        areaxstyle,
		storeSysData:  true,
		doNotLoadFont: true,
	},
}

func (s *sysData) addChild(child *sysData) _HMENU {
	s.childrenLock.Lock()
	defer s.childrenLock.Unlock()
	s.nextChildID++ // start at 1
	if s.children == nil {
		s.children = map[_HMENU]*sysData{}
	}
	s.children[s.nextChildID] = child
	return s.nextChildID
}

func (s *sysData) delChild(id _HMENU) {
	s.childrenLock.Lock()
	defer s.childrenLock.Unlock()
	delete(s.children, id)
}

var (
	_blankString = toUTF16("")
	blankString  = utf16ToArg(_blankString)
)

func (s *sysData) make(window *sysData) (err error) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		ct := classTypes[s.ctype]
		cid := _HMENU(0)
		pwin := uintptr(_NULL)
		if window != nil { // this is a child control
			cid = window.addChild(s)
			pwin = uintptr(window.hwnd)
		}
		style := uintptr(ct.style)
		if s.alternate {
			style = uintptr(ct.altStyle)
		}
		lpParam := uintptr(_NULL)
		if ct.storeSysData {
			lpParam = uintptr(unsafe.Pointer(s))
		}
		r1, _, err := _createWindowEx.Call(
			uintptr(ct.xstyle),
			utf16ToArg(ct.name),
			blankString, // we set the window text later
			style,
			negConst(_CW_USEDEFAULT),
			negConst(_CW_USEDEFAULT),
			negConst(_CW_USEDEFAULT),
			negConst(_CW_USEDEFAULT),
			pwin,
			uintptr(cid),
			uintptr(hInstance),
			lpParam)
		if r1 == 0 { // failure
			if window != nil {
				window.delChild(cid)
			}
			panic(fmt.Errorf("error actually creating window/control: %v", err))
		}
		if !ct.storeSysData { // regular control; store s.hwnd ourselves
			s.hwnd = _HWND(r1)
		} else if s.hwnd != _HWND(r1) { // we store sysData in storeSysData(); sanity check
			panic(fmt.Errorf("hwnd mismatch creating window/control: storeSysData() stored 0x%X but CreateWindowEx() returned 0x%X", s.hwnd, r1))
		}
		if !ct.doNotLoadFont {
			_sendMessage.Call(
				uintptr(s.hwnd),
				uintptr(_WM_SETFONT),
				uintptr(_WPARAM(controlFont)),
				uintptr(_LPARAM(_TRUE)))
		}
		ret <- struct{}{}
	}
	<-ret
	return nil
}

var (
	_updateWindow = user32.NewProc("UpdateWindow")
)

// if the object is a window, we need to do the following the first time
// 	ShowWindow(hwnd, nCmdShow);
// 	UpdateWindow(hwnd);
func (s *sysData) firstShow() error {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		_showWindow.Call(
			uintptr(s.hwnd),
			uintptr(nCmdShow))
		r1, _, err := _updateWindow.Call(uintptr(s.hwnd))
		if r1 == 0 { // failure
			panic(fmt.Errorf("error updating window for the first time: %v", err))
		}
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) show() {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		_showWindow.Call(
			uintptr(s.hwnd),
			uintptr(_SW_SHOW))
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) hide() {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		_showWindow.Call(
			uintptr(s.hwnd),
			uintptr(_SW_HIDE))
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) setText(text string) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		ptext := toUTF16(text)
		r1, _, err := _setWindowText.Call(
			uintptr(s.hwnd),
			utf16ToArg(ptext))
		if r1 == 0 { // failure
			panic(fmt.Errorf("error setting window/control text: %v", err))
		}
		ret <- struct{}{}
	}
	<-ret
}

// runs on uitask
func (s *sysData) setRect(x int, y int, width int, height int, winheight int) error {
	r1, _, err := _moveWindow.Call(
		uintptr(s.hwnd),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(_TRUE))
	if r1 == 0 { // failure
		return fmt.Errorf("error setting window/control rect: %v", err)
	}
	return nil
}

func (s *sysData) isChecked() bool {
	ret := make(chan bool)
	defer close(ret)
	uitask <- func() {
		r1, _, _ := _sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(_BM_GETCHECK),
			uintptr(0),
			uintptr(0))
		ret <- r1 == _BST_CHECKED
	}
	return <-ret
}

func (s *sysData) text() (str string) {
	ret := make(chan string)
	defer close(ret)
	uitask <- func() {
		var tc []uint16

		r1, _, _ := _sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(_WM_GETTEXTLENGTH),
			uintptr(0),
			uintptr(0))
		length := r1 + 1 // terminating null
		tc = make([]uint16, length)
		_sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(_WM_GETTEXT),
			uintptr(_WPARAM(length)),
			uintptr(_LPARAM(unsafe.Pointer(&tc[0]))))
		ret <- syscall.UTF16ToString(tc)
	}
	return <-ret
}

func (s *sysData) append(what string) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		pwhat := toUTF16(what)
		r1, _, err := _sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].appendMsg),
			uintptr(_WPARAM(0)),
			utf16ToLPARAM(pwhat))
		if r1 == uintptr(classTypes[s.ctype].addSpaceErr) {
			panic(fmt.Errorf("out of space adding item to combobox/listbox (last error: %v)", err))
		} else if r1 == uintptr(classTypes[s.ctype].selectedIndexErr) {
			panic(fmt.Errorf("failed to add item to combobox/listbox (last error: %v)", err))
		}
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) insertBefore(what string, index int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		pwhat := toUTF16(what)
		r1, _, err := _sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].insertBeforeMsg),
			uintptr(_WPARAM(index)),
			utf16ToLPARAM(pwhat))
		if r1 == uintptr(classTypes[s.ctype].addSpaceErr) {
			panic(fmt.Errorf("out of space adding item to combobox/listbox (last error: %v)", err))
		} else if r1 == uintptr(classTypes[s.ctype].selectedIndexErr) {
			panic(fmt.Errorf("failed to add item to combobox/listbox (last error: %v)", err))
		}
		ret <- struct{}{}
	}
	<-ret
}

// runs on uitask
func (s *sysData) doSelectedIndex() int {
	r1, _, _ := _sendMessage.Call(
		uintptr(s.hwnd),
		uintptr(classTypes[s.ctype].selectedIndexMsg),
		uintptr(_WPARAM(0)),
		uintptr(_LPARAM(0)))
	if r1 == uintptr(classTypes[s.ctype].selectedIndexErr) { // no selection or manually entered text (apparently, for the latter)
		return -1
	}
	return int(r1)
}

func (s *sysData) selectedIndex() int {
	ret := make(chan int)
	defer close(ret)
	uitask <- func() {
		ret <- s.doSelectedIndex()
	}
	return <-ret
}

// runs on uitask
func (s *sysData) doSelectedIndices() []int {
	if !s.alternate { // single-selection list box; use single-selection method
		index := s.doSelectedIndex()
		if index == -1 {
			return nil
		}
		return []int{index}
	}

	r1, _, err := _sendMessage.Call(
		uintptr(s.hwnd),
		uintptr(_LB_GETSELCOUNT),
		uintptr(0),
		uintptr(0))
	if r1 == negConst(_LB_ERR) {
		panic(fmt.Errorf("error: LB_ERR from LB_GETSELCOUNT in what we know is a multi-selection listbox: %v", err))
	}
	if r1 == 0 { // nothing selected
		return nil
	}
	indices := make([]int, r1)
	r1, _, err = _sendMessage.Call(
		uintptr(s.hwnd),
		uintptr(_LB_GETSELITEMS),
		uintptr(_WPARAM(r1)),
		uintptr(_LPARAM(unsafe.Pointer(&indices[0]))))
	if r1 == negConst(_LB_ERR) {
		panic(fmt.Errorf("error: LB_ERR from LB_GETSELITEMS in what we know is a multi-selection listbox: %v", err))
	}
	return indices
}

func (s *sysData) selectedIndices() []int {
	ret := make(chan []int)
	defer close(ret)
	uitask <- func() {
		ret <- s.doSelectedIndices()
	}
	return <-ret
}

func (s *sysData) selectedTexts() []string {
	ret := make(chan []string)
	defer close(ret)
	uitask <- func() {
		indices := s.doSelectedIndices()
		strings := make([]string, len(indices))
		for i, v := range indices {
			r1, _, err := _sendMessage.Call(
				uintptr(s.hwnd),
				uintptr(_LB_GETTEXTLEN),
				uintptr(_WPARAM(v)),
				uintptr(0))
			if r1 == negConst(_LB_ERR) {
				panic(fmt.Errorf("error: LB_ERR from LB_GETTEXTLEN in what we know is a valid listbox index (came from LB_GETSELITEMS): %v", err))
			}
			str := make([]uint16, r1)
			r1, _, err = _sendMessage.Call(
				uintptr(s.hwnd),
				uintptr(_LB_GETTEXT),
				uintptr(_WPARAM(v)),
				uintptr(_LPARAM(unsafe.Pointer(&str[0]))))
			if r1 == negConst(_LB_ERR) {
				panic(fmt.Errorf("error: LB_ERR from LB_GETTEXT in what we know is a valid listbox index (came from LB_GETSELITEMS): %v", err))
			}
			strings[i] = syscall.UTF16ToString(str)
		}
		ret <- strings
	}
	return <-ret
}

func (s *sysData) setWindowSize(width int, height int) error {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		var rect _RECT

		r1, _, err := _getClientRect.Call(
			uintptr(s.hwnd),
			uintptr(unsafe.Pointer(&rect)))
		if r1 == 0 {
			panic(fmt.Errorf("error getting upper-left of window for resize: %v", err))
		}
		// TODO AdjustWindowRect() on the result
		// 0 because (0,0) is top-left so no winheight
		err = s.setRect(int(rect.left), int(rect.top), width, height, 0)
		if err != nil {
			panic(fmt.Errorf("error actually resizing window: %v", err))
		}
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) delete(index int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		r1, _, err := _sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].deleteMsg),
			uintptr(_WPARAM(index)),
			uintptr(0))
		if r1 == uintptr(classTypes[s.ctype].selectedIndexErr) {
			panic(fmt.Errorf("failed to delete item from combobox/listbox (last error: %v)", err))
		}
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) setIndeterminate() {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		r1, _, err := _setWindowLongPtr.Call(
			uintptr(s.hwnd),
			negConst(_GWL_STYLE),
			uintptr(classTypes[s.ctype].style | _PBS_MARQUEE))
		if r1 == 0 {
			panic(fmt.Errorf("error setting progress bar style to enter indeterminate mode: %v", err))
		}
		_sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(_PBM_SETMARQUEE),
			uintptr(_WPARAM(_TRUE)),
			uintptr(0))
		s.isMarquee = true
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) setProgress(percent int) {
	if percent == -1 {
		s.setIndeterminate()
		return
	}
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		if s.isMarquee {
			// turn off marquee before switching back
			_sendMessage.Call(
				uintptr(s.hwnd),
				uintptr(_PBM_SETMARQUEE),
				uintptr(_WPARAM(_FALSE)),
				uintptr(0))
			r1, _, err := _setWindowLongPtr.Call(
				uintptr(s.hwnd),
				negConst(_GWL_STYLE),
				uintptr(classTypes[s.ctype].style))
			if r1 == 0 {
				panic(fmt.Errorf("error setting progress bar style to leave indeterminate mode (percent %d): %v", percent, err))
			}
			s.isMarquee = false
		}
		send := func(msg uintptr, n int, l _LPARAM) {
			_sendMessage.Call(
				uintptr(s.hwnd),
				msg,
				uintptr(_WPARAM(n)),
				uintptr(l))
		}
		// Windows 7 has a non-disableable slowly-animating progress bar increment
		// there isn't one for decrement, so we'll work around by going one higher and then lower again
		// for the case where percent == 100, we need to increase the range temporarily
		// sources: http://social.msdn.microsoft.com/Forums/en-US/61350dc7-6584-4c4e-91b0-69d642c03dae/progressbar-disable-smooth-animation http://stackoverflow.com/questions/2217688/windows-7-aero-theme-progress-bar-bug http://discuss.joelonsoftware.com/default.asp?dotnet.12.600456.2 http://stackoverflow.com/questions/22469876/progressbar-lag-when-setting-position-with-pbm-setpos http://stackoverflow.com/questions/6128287/tprogressbar-never-fills-up-all-the-way-seems-to-be-updating-too-fast
		if percent == 100 {
			send(_PBM_SETRANGE32, 0, 101)
		}
		send(_PBM_SETPOS, percent+1, 0)
		send(_PBM_SETPOS, percent, 0)
		if percent == 100 {
			send(_PBM_SETRANGE32, 0, 100)
		}
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) len() int {
	ret := make(chan int)
	defer close(ret)
	uitask <- func() {
		r1, _, err := _sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].lenMsg),
			uintptr(_WPARAM(0)),
			uintptr(_LPARAM(0)))
		if r1 == uintptr(classTypes[s.ctype].selectedIndexErr) {
			panic(fmt.Errorf("unexpected error return from sysData.len(); GetLastError() says %v", err))
		}
		ret <- int(r1)
	}
	return <-ret
}

func (s *sysData) setAreaSize(width int, height int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		_sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(msgSetAreaSize),
			uintptr(width), // WPARAM is UINT_PTR on Windows XP and newer at least, so we're good with this
			uintptr(height))
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) repaintAll() {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		_sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(msgRepaintAll),
			uintptr(0),
			uintptr(0))
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) center() {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		var ws _RECT

		r1, _, err := _getWindowRect.Call(
			uintptr(s.hwnd),
			uintptr(unsafe.Pointer(&ws)))
		if r1 == 0 {
			panic(fmt.Errorf("error getting window rect for sysData.center(): %v", err))
		}
		// TODO should this be using the monitor functions instead? http://blogs.msdn.com/b/oldnewthing/archive/2005/05/05/414910.aspx
		// error returns from GetSystemMetrics() is meaningless because the return value, 0, is still valid
		dw, _, _ := _getSystemMetrics.Call(uintptr(_SM_CXFULLSCREEN))
		dh, _, _ := _getSystemMetrics.Call(uintptr(_SM_CYFULLSCREEN))
		ww := ws.right - ws.left
		wh := ws.bottom - ws.top
		wx := (int32(dw) / 2) - (ww / 2)
		wy := (int32(dh) / 2) - (wh / 2)
		s.setRect(int(wx), int(wy), int(ww), int(wh), 0)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) setChecked(checked bool) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		c := uintptr(_BST_CHECKED)
		if !checked {
			c = uintptr(_BST_UNCHECKED)
		}
		_sendMessage.Call(
			uintptr(s.hwnd),
			uintptr(_BM_SETCHECK),
			c,
			uintptr(0))
		ret <- struct{}{}
	}
	<-ret
}

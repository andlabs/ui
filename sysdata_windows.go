// 11 february 2014

//
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
	shownAlready bool
}

type classData struct {
	name             string
	style            uint32
	xstyle           uint32
	mkid             bool
	altStyle         uint32
	font             *_HANDLE
	appendMsg        uintptr
	insertBeforeMsg  uintptr
	deleteMsg        uintptr
	selectedIndexMsg uintptr
	selectedIndexErr int
	addSpaceErr      int
}

const controlstyle = _WS_CHILD | _WS_VISIBLE | _WS_TABSTOP
const controlxstyle = 0

var classTypes = [nctypes]*classData{
	c_window: &classData{
		style:  _WS_OVERLAPPEDWINDOW,
		xstyle: 0,
	},
	c_button: &classData{
		name:   "BUTTON",
		style:  _BS_PUSHBUTTON | controlstyle,
		xstyle: 0 | controlxstyle,
		font:   &controlFont,
	},
	c_checkbox: &classData{
		name:   "BUTTON",
		style:  _BS_AUTOCHECKBOX | controlstyle,
		xstyle: 0 | controlxstyle,
		font:   &controlFont,
	},
	c_combobox: &classData{
		name:             "COMBOBOX",
		style:            _CBS_DROPDOWNLIST | _WS_VSCROLL | controlstyle,
		xstyle:           0 | controlxstyle,
		altStyle:         _CBS_DROPDOWN | _CBS_AUTOHSCROLL | _WS_VSCROLL | controlstyle,
		font:             &controlFont,
		appendMsg:        _CB_ADDSTRING,
		insertBeforeMsg:  _CB_INSERTSTRING,
		deleteMsg:        _CB_DELETESTRING,
		selectedIndexMsg: _CB_GETCURSEL,
		selectedIndexErr: _CB_ERR,
		addSpaceErr:      _CB_ERRSPACE,
	},
	c_lineedit: &classData{
		name:     "EDIT",
		style:    _ES_AUTOHSCROLL | _WS_BORDER | controlstyle,
		xstyle:   0 | controlxstyle,
		altStyle: _ES_PASSWORD | _ES_AUTOHSCROLL | _WS_BORDER | controlstyle,
		font:     &controlFont,
	},
	c_label: &classData{
		name:   "STATIC",
		style:  _SS_NOPREFIX | controlstyle,
		xstyle: 0 | controlxstyle,
		font:   &controlFont,
	},
	c_listbox: &classData{
		name: "LISTBOX",
		// TODO also _WS_HSCROLL?
		style:            _WS_VSCROLL | controlstyle,
		xstyle:           0 | controlxstyle,
		altStyle:         _LBS_EXTENDEDSEL | _WS_VSCROLL | controlstyle,
		font:             &controlFont,
		appendMsg:        _LB_ADDSTRING,
		insertBeforeMsg:  _LB_INSERTSTRING,
		deleteMsg:        _LB_DELETESTRING,
		selectedIndexMsg: _LB_GETCURSEL,
		selectedIndexErr: _LB_ERR,
		addSpaceErr:      _LB_ERRSPACE,
	},
	c_progressbar: &classData{
		name:   _PROGRESS_CLASS,
		style:  _PBS_SMOOTH | controlstyle,
		xstyle: 0 | controlxstyle,
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

func (s *sysData) make(initText string, window *sysData) (err error) {
	ret := make(chan uiret)
	defer close(ret)
	ct := classTypes[s.ctype]
	classname := ct.name
	cid := _HMENU(0)
	pwin := uintptr(_NULL)
	if window != nil { // this is a child control
		cid = window.addChild(s)
		pwin = uintptr(window.hwnd)
	} else { // need a new class
		n, err := registerStdWndClass(s)
		if err != nil {
			return fmt.Errorf("error creating window class for new window: %v", err)
		}
		classname = n
	}
	style := uintptr(ct.style)
	if s.alternate {
		style = uintptr(ct.altStyle)
	}
	uitask <- &uimsg{
		call: _createWindowEx,
		p: []uintptr{
			uintptr(ct.xstyle),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(classname))),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(initText))),
			style,
			uintptr(_CW_USEDEFAULT),
			uintptr(_CW_USEDEFAULT),
			uintptr(_CW_USEDEFAULT),
			uintptr(_CW_USEDEFAULT),
			pwin,
			uintptr(cid),
			uintptr(hInstance),
			uintptr(_NULL),
		},
		ret: ret,
	}
	r := <-ret
	if r.ret == 0 { // failure
		if window != nil {
			window.delChild(cid)
		}
		return fmt.Errorf("error actually creating window/control: %v", r.err)
	}
	s.hwnd = _HWND(r.ret)
	if ct.font != nil {
		uitask <- &uimsg{
			call: _sendMessage,
			p: []uintptr{
				uintptr(s.hwnd),
				uintptr(_WM_SETFONT),
				uintptr(_WPARAM(*ct.font)),
				uintptr(_LPARAM(_TRUE)),
			},
			ret: ret,
		}
		<-ret
	}
	return nil
}

var (
	_updateWindow = user32.NewProc("UpdateWindow")
)

// if the object is a window, we need to do the following the first time
// 	ShowWindow(hwnd, nCmdShow);
// 	UpdateWindow(hwnd);
// otherwise we go ahead and show the object normally with SW_SHOW
func (s *sysData) show() (err error) {
	if s.ctype != c_window { // don't do the init ShowWindow/UpdateWindow chain on non-windows
		s.shownAlready = true
	}
	show := uintptr(_SW_SHOW)
	if !s.shownAlready {
		show = uintptr(nCmdShow)
	}
	ret := make(chan uiret)
	defer close(ret)
	// TODO figure out how to handle error
	uitask <- &uimsg{
		call: _showWindow,
		p:    []uintptr{uintptr(s.hwnd), show},
		ret:  ret,
	}
	<-ret
	if !s.shownAlready {
		uitask <- &uimsg{
			call: _updateWindow,
			p:    []uintptr{uintptr(s.hwnd)},
			ret:  ret,
		}
		r := <-ret
		if r.ret == 0 { // failure
			return fmt.Errorf("error updating window for the first time: %v", r.err)
		}
		s.shownAlready = true
	}
	return nil
}

func (s *sysData) hide() (err error) {
	ret := make(chan uiret)
	defer close(ret)
	// TODO figure out how to handle error
	uitask <- &uimsg{
		call: _showWindow,
		p:    []uintptr{uintptr(s.hwnd), _SW_HIDE},
		ret:  ret,
	}
	<-ret
	return nil
}

func (s *sysData) setText(text string) error {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call: _setWindowText,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		},
		ret: ret,
	}
	r := <-ret
	if r.ret == 0 { // failure
		return r.err
	}
	return nil
}

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
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call: _sendMessage,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(_BM_GETCHECK),
			uintptr(0),
			uintptr(0),
		},
		ret: ret,
	}
	r := <-ret
	return r.ret == _BST_CHECKED
}

func (s *sysData) text() (str string) {
	var tc []uint16

	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call: _sendMessage,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(_WM_GETTEXTLENGTH),
			uintptr(0),
			uintptr(0),
		},
		ret: ret,
	}
	r := <-ret
	length := r.ret + 1 // terminating null
	tc = make([]uint16, length)
	uitask <- &uimsg{
		call: _sendMessage,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(_WM_GETTEXT),
			uintptr(_WPARAM(length)),
			uintptr(_LPARAM(unsafe.Pointer(&tc[0]))),
		},
		ret: ret,
	}
	<-ret
	return syscall.UTF16ToString(tc)
}

func (s *sysData) append(what string) (err error) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call: _sendMessage,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].appendMsg),
			uintptr(_WPARAM(0)),
			uintptr(_LPARAM(unsafe.Pointer(syscall.StringToUTF16Ptr(what)))),
		},
		ret: ret,
	}
	r := <-ret
	if r.ret == uintptr(classTypes[s.ctype].addSpaceErr) {
		return fmt.Errorf("out of space adding item to combobox/listbox (last error: %v)", r.err)
	} else if r.ret == uintptr(classTypes[s.ctype].selectedIndexErr) {
		return fmt.Errorf("failed to add item to combobox/listbox (last error: %v)", r.err)
	}
	return nil
}

func (s *sysData) insertBefore(what string, index int) (err error) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call: _sendMessage,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].insertBeforeMsg),
			uintptr(_WPARAM(index)),
			uintptr(_LPARAM(unsafe.Pointer(syscall.StringToUTF16Ptr(what)))),
		},
		ret: ret,
	}
	r := <-ret
	if r.ret == uintptr(classTypes[s.ctype].addSpaceErr) {
		return fmt.Errorf("out of space adding item to combobox/listbox (last error: %v)", r.err)
	} else if r.ret == uintptr(classTypes[s.ctype].selectedIndexErr) {
		return fmt.Errorf("failed to add item to combobox/listbox (last error: %v)", r.err)
	}
	return nil
}

func (s *sysData) selectedIndex() int {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call: _sendMessage,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].selectedIndexMsg),
			uintptr(_WPARAM(0)),
			uintptr(_LPARAM(0)),
		},
		ret: ret,
	}
	r := <-ret
	if r.ret == uintptr(classTypes[s.ctype].selectedIndexErr) { // no selection or manually entered text (apparently, for the latter)
		return -1
	}
	return int(r.ret)
}

func (s *sysData) selectedIndices() []int {
	if !s.alternate { // single-selection list box; use single-selection method
		index := s.selectedIndex()
		if index == -1 {
			return nil
		}
		return []int{index}
	}

	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call: _sendMessage,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(_LB_GETSELCOUNT),
			uintptr(0),
			uintptr(0),
		},
		ret: ret,
	}
	r := <-ret
	if r.ret == uintptr(_LB_ERR) {
		panic("UI library internal error: LB_ERR from LB_GETSELCOUNT in what we know is a multi-selection listbox")
	}
	if r.ret == 0 { // nothing selected
		return nil
	}
	indices := make([]int, r.ret)
	uitask <- &uimsg{
		call: _sendMessage,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(_LB_GETSELITEMS),
			uintptr(_WPARAM(r.ret)),
			uintptr(_LPARAM(unsafe.Pointer(&indices[0]))),
		},
		ret: ret,
	}
	r = <-ret
	if r.ret == uintptr(_LB_ERR) {
		panic("UI library internal error: LB_ERR from LB_GETSELITEMS in what we know is a multi-selection listbox")
	}
	return indices
}

func (s *sysData) selectedTexts() []string {
	indices := s.selectedIndices()
	ret := make(chan uiret)
	defer close(ret)
	strings := make([]string, len(indices))
	for i, v := range indices {
		uitask <- &uimsg{
			call: _sendMessage,
			p: []uintptr{
				uintptr(s.hwnd),
				uintptr(_LB_GETTEXTLEN),
				uintptr(_WPARAM(v)),
				uintptr(0),
			},
			ret: ret,
		}
		r := <-ret
		if r.ret == uintptr(_LB_ERR) {
			panic("UI library internal error: LB_ERR from LB_GETTEXTLEN in what we know is a valid listbox index (came from LB_GETSELITEMS")
		}
		str := make([]uint16, r.ret)
		uitask <- &uimsg{
			call: _sendMessage,
			p: []uintptr{
				uintptr(s.hwnd),
				uintptr(_LB_GETTEXT),
				uintptr(_WPARAM(v)),
				uintptr(_LPARAM(unsafe.Pointer(&str[0]))),
			},
			ret: ret,
		}
		<-ret
		if r.ret == uintptr(_LB_ERR) {
			panic("UI library internal error: LB_ERR from LB_GETTEXT in what we know is a valid listbox index (came from LB_GETSELITEMS")
		}
		strings[i] = syscall.UTF16ToString(str)
	}
	return strings
}

func (s *sysData) setWindowSize(width int, height int) error {
	var rect _RECT

	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call: _getClientRect,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(unsafe.Pointer(&rect)),
		},
		ret: ret,
	}
	r := <-ret
	if r.ret == 0 {
		return fmt.Errorf("error getting upper-left of window for resize: %v", r.err)
	}
	// 0 because (0,0) is top-left so no winheight
	err := s.setRect(int(rect.Left), int(rect.Top), width, height, 0)
	if err != nil {
		return fmt.Errorf("error actually resizing window: %v", err)
	}
	return nil
}

func (s *sysData) delete(index int) (err error) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call: _sendMessage,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].deleteMsg),
			uintptr(_WPARAM(index)),
			uintptr(0),
		},
		ret: ret,
	}
	r := <-ret
	if r.ret == uintptr(classTypes[s.ctype].selectedIndexErr) {
		return fmt.Errorf("failed to delete item from combobox/listbox (last error: %v)", r.err)
	}
	return nil
}

func (s *sysData) setProgress(percent int) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call: _sendMessage,
		p: []uintptr{
			uintptr(s.hwnd),
			uintptr(_PBM_SETPOS),
			uintptr(_WPARAM(percent)),
			uintptr(0),
		},
		ret: ret,
	}
	<-ret
}

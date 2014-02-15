// 11 february 2014
package main

import (
	"fmt"
	"syscall"
	"unsafe"
	"sync"
)

type sysData struct {
	cSysData

	hwnd			_HWND
	children			map[_HMENU]*sysData
	nextChildID		_HMENU
	childrenLock		sync.Mutex
	shownAlready		bool
}

type classData struct {
	name			string
	style				uint32
	xstyle			uint32
	mkid				bool
	altStyle			uint32
	appendMsg		uintptr
	insertBeforeMsg	uintptr
	deleteMsg			uintptr
	selectedIndexMsg	uintptr
	selectedIndexErr	int
}

const controlstyle = _WS_CHILD | _WS_VISIBLE | _WS_TABSTOP
const controlxstyle = 0

var classTypes = [nctypes]*classData{
	c_window:			&classData{
		style:			_WS_OVERLAPPEDWINDOW,
		xstyle:			0,
	},
	c_button:		&classData{
		name:			"BUTTON",
		style:			_BS_PUSHBUTTON | controlstyle,
		xstyle:			0 | controlxstyle,
	},
	c_checkbox:	&classData{
		name:			"BUTTON",
		style:			_BS_AUTOCHECKBOX | controlstyle,
		xstyle:			0 | controlxstyle,
	},
	c_combobox:	&classData{
		name:			"COMBOBOX",
		style:			_CBS_DROPDOWNLIST | controlstyle,
		xstyle:			0 | controlxstyle,
		altStyle:			_CBS_DROPDOWN | _CBS_AUTOHSCROLL | controlstyle,
		appendMsg:		_CB_ADDSTRING,
		insertBeforeMsg:	_CB_INSERTSTRING,
		deleteMsg:		_CB_DELETESTRING,
		selectedIndexMsg:	_CB_GETCURSEL,
		selectedIndexErr:	_CB_ERR,
	},
	c_lineedit:	&classData{
		name:			"EDIT",
		style:			_ES_AUTOHSCROLL | _WS_BORDER | controlstyle,
		xstyle:			0 | controlxstyle,
	},
	c_label:		&classData{
		name:			"STATIC",
		style:			_SS_NOPREFIX | controlstyle,
		xstyle:			0 | controlxstyle,
	},
	c_listbox:		&classData{
		name:			"LISTBOX",
		// TODO also _WS_HSCROLL?
		style:			_WS_VSCROLL | controlstyle,
		xstyle:			0 | controlxstyle,
		altStyle:			_LBS_EXTENDEDSEL | _WS_VSCROLL | controlstyle,
		appendMsg:		_LB_ADDSTRING,
		insertBeforeMsg:	_LB_INSERTSTRING,
		deleteMsg:		_LB_DELETESTRING,
		selectedIndexMsg:	_LB_GETCURSEL,
		selectedIndexErr:	_LB_ERR,
	},
}

func (s *sysData) addChild(child *sysData) _HMENU {
	s.childrenLock.Lock()
	defer s.childrenLock.Unlock()
	s.nextChildID++		// start at 1
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

// TODO adorn error messages with what stage failed?
func (s *sysData) make(initText string, initWidth int, initHeight int, window *sysData) (err error) {
	ret := make(chan uiret)
	defer close(ret)
	ct := classTypes[s.ctype]
	classname := ct.name
	cid := _HMENU(0)
	pwin := uintptr(_NULL)
	if window != nil {		// this is a child control
		cid = window.addChild(s)
		pwin = uintptr(window.hwnd)
	} else {				// need a new class
		n, err := registerStdWndClass(s)
		if err != nil {
			return err
		}
		classname = n
	}
	style := uintptr(ct.style)
	if s.alternate {
		style = uintptr(ct.altStyle)
	}
	uitask <- &uimsg{
		call:		_createWindowEx,	
		p:		[]uintptr{
			uintptr(ct.xstyle),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(classname))),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(initText))),
			style,
			uintptr(_CW_USEDEFAULT),		// TODO
			uintptr(_CW_USEDEFAULT),
			uintptr(initWidth),
			uintptr(initHeight),
			pwin,
			uintptr(cid),
			uintptr(hInstance),
			uintptr(_NULL),
		},
		ret:	ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		if window != nil {
			window.delChild(cid)
		}
		return r.err
	}
	s.hwnd = _HWND(r.ret)
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
	if s.ctype != c_window {		// don't do the init ShowWindow/UpdateWindow chain on non-windows
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
		call:		_showWindow,
		p:		[]uintptr{uintptr(s.hwnd), show},
		ret:		ret,
	}
	<-ret
	if !s.shownAlready {
		uitask <- &uimsg{
			call:		_updateWindow,
			p:		[]uintptr{uintptr(s.hwnd)},
			ret:		ret,
		}
		r := <-ret
		if r.ret == 0 {		// failure
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
		call:		_showWindow,
		p:		[]uintptr{uintptr(s.hwnd), _SW_HIDE},
		ret:		ret,
	}
	<-ret
	return nil
}

func (s *sysData) setText(text string) error {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_setWindowText,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		return r.err
	}
	return nil
}

func (s *sysData) setRect(x int, y int, width int, height int) error {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_moveWindow,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(x),
			uintptr(y),
			uintptr(width),
			uintptr(height),
			uintptr(_TRUE),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		return r.err
	}
	return nil
}

// TODO figure out how to handle error
func (s *sysData) isChecked() (bool, error) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(_BM_GETCHECK),
			uintptr(0),
			uintptr(0),
		},
		ret:		ret,
	}
	r := <-ret
	return r.ret == _BST_CHECKED, nil
}

// TODO adorn error messages with which part failed
func (s *sysData) text() (str string, err error) {
	var tc []uint16

	ret := make(chan uiret)
	defer close(ret)
	// TODO figure out how to handle errors
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(_WM_GETTEXTLENGTH),
			uintptr(0),
			uintptr(0),
		},
		ret:		ret,
	}
	r := <-ret
	length := r.ret + 1		// terminating null
	tc = make([]uint16, length)
	// TODO figure out how to handle errors
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(_WM_GETTEXT),
			uintptr(_WPARAM(length)),
			uintptr(_LPARAM(unsafe.Pointer(&tc[0]))),
		},
		ret:		ret,
	}
	<-ret
	// TODO check character count
	return syscall.UTF16ToString(tc), nil
}

// TODO figure out how to handle errors
func (s *sysData) append(what string) (err error) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].appendMsg),
			uintptr(_WPARAM(0)),
			uintptr(_LPARAM(unsafe.Pointer(syscall.StringToUTF16Ptr(what)))),
		},
		ret:		ret,
	}
	<-ret
	// TODO error handling
	return nil
}

// TODO figure out how to handle errors
func (s *sysData) insertBefore(what string, index int) (err error) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].insertBeforeMsg),
			uintptr(_WPARAM(index)),
			uintptr(_LPARAM(unsafe.Pointer(syscall.StringToUTF16Ptr(what)))),
		},
		ret:		ret,
	}
	<-ret
	// TODO error handling
	return nil
}

// TODO handle actual errors
// TODO differentiate between nothing selected and custom text entered for a Combobox
func (s *sysData) selectedIndex() (int, error) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].selectedIndexMsg),
			uintptr(_WPARAM(0)),
			uintptr(_LPARAM(0)),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == uintptr(classTypes[s.ctype].selectedIndexErr) {
		return -1, nil
	}
	return int(r.ret), nil
}

// TODO handle actual errors
func (s *sysData) selectedIndices() ([]int, error) {
	if !s.alternate {		// single-selection list box; use single-selection method
		index, err := s.selectedIndex()
		if err != nil {
			return nil, fmt.Errorf("error getting indices of single-selection list box: %v", err)
		}
		return []int{index}, nil
	}

	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(_LB_GETSELCOUNT),
			uintptr(0),
			uintptr(0),
		},
		ret:		ret,
	}
	r := <-ret
	// TODO handle errors
	indices := make([]int, r.ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(_LB_GETSELITEMS),
			uintptr(_WPARAM(r.ret)),
			uintptr(_LPARAM(unsafe.Pointer(&indices[0]))),
		},
		ret:		ret,
	}
	<-ret
	// TODO handle errors
	return indices, nil
}

func (s *sysData) selectedTexts() ([]string, error) {
	indices, err := s.selectedIndices()
	if err != nil {
		return nil, fmt.Errorf("error getting selected indices for selected texts: %v", err)
	}
	ret := make(chan uiret)
	defer close(ret)
	strings := make([]string, len(indices))
	for i, v := range indices {
		uitask <- &uimsg{
			call:		_sendMessage,
			p:		[]uintptr{
				uintptr(s.hwnd),
				uintptr(_LB_GETTEXTLEN),
				uintptr(_WPARAM(v)),
				uintptr(0),
			},
			ret:		ret,
		}
		r := <-ret
		// TODO handle errors
		str := make([]uint16, r.ret)
		uitask <- &uimsg{
			call:		_sendMessage,
			p:		[]uintptr{
				uintptr(s.hwnd),
				uintptr(_LB_GETTEXT),
				uintptr(_WPARAM(v)),
				uintptr(_LPARAM(unsafe.Pointer(&str[0]))),
			},
			ret:		ret,
		}
		<-ret
		// TODO handle errors
		strings[i] = syscall.UTF16ToString(str)
	}
	return strings, nil
}

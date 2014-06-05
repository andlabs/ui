// 11 february 2014

package ui

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
	isMarquee		bool			// for sysData.setProgress()
	// unlike with GTK+ and Mac OS X, we're responsible for sizing Area properly ourselves
	areawidth			int
	areaheight		int
	clickCounter		clickCounter
}

type classData struct {
	name			*uint16
	style				uint32
	xstyle			uint32
	altStyle			uint32
	storeSysData		bool
	doNotLoadFont	bool
	appendMsg		uintptr
	insertBeforeMsg	uintptr
	deleteMsg			uintptr
	selectedIndexMsg	uintptr
	selectedIndexErr	uintptr
	addSpaceErr		uintptr
	lenMsg			uintptr
}

const controlstyle = _WS_CHILD | _WS_VISIBLE | _WS_TABSTOP
const controlxstyle = 0

var classTypes = [nctypes]*classData{
	c_window:		&classData{
		name:			stdWndClass,
		style:			_WS_OVERLAPPEDWINDOW,
		xstyle:			0,
		storeSysData:		true,
		doNotLoadFont:	true,
	},
	c_button:			&classData{
		name:			toUTF16("BUTTON"),
		style:			_BS_PUSHBUTTON | controlstyle,
		xstyle:			0 | controlxstyle,
	},
	c_checkbox:		&classData{
		name:			toUTF16("BUTTON"),
		// don't use BS_AUTOCHECKBOX because http://blogs.msdn.com/b/oldnewthing/archive/2014/05/22/10527522.aspx
		style:			_BS_CHECKBOX | controlstyle,
		xstyle:			0 | controlxstyle,
	},
	c_combobox:		&classData{
		name:			toUTF16("COMBOBOX"),
		style:			_CBS_DROPDOWNLIST | _WS_VSCROLL | controlstyle,
		xstyle:			0 | controlxstyle,
		altStyle:			_CBS_DROPDOWN | _CBS_AUTOHSCROLL | _WS_VSCROLL | controlstyle,
		appendMsg:		_CB_ADDSTRING,
		insertBeforeMsg:	_CB_INSERTSTRING,
		deleteMsg:		_CB_DELETESTRING,
		selectedIndexMsg:	_CB_GETCURSEL,
		selectedIndexErr:	negConst(_CB_ERR),
		addSpaceErr:		negConst(_CB_ERRSPACE),
		lenMsg:			_CB_GETCOUNT,
	},
	c_lineedit:		&classData{
		name:			toUTF16("EDIT"),
		// WS_EX_CLIENTEDGE without WS_BORDER will apply visual styles
		// thanks to MindChild in irc.efnet.net/#winprog
		style:			_ES_AUTOHSCROLL | controlstyle,
		xstyle:			_WS_EX_CLIENTEDGE | controlxstyle,
		altStyle:			_ES_PASSWORD | _ES_AUTOHSCROLL | controlstyle,
	},
	c_label:			&classData{
		name:			toUTF16("STATIC"),
		// SS_NOPREFIX avoids accelerator translation; SS_LEFTNOWORDWRAP clips text past the end
		// TODO find out if the default behavior is not to ellipsize
		style:			_SS_NOPREFIX | _SS_LEFTNOWORDWRAP | controlstyle,
		xstyle:			0 | controlxstyle,
	},
	c_listbox:			&classData{
		name:			toUTF16("LISTBOX"),
		// we don't use LBS_STANDARD because it sorts (and has WS_BORDER; see above)
		// LBS_NOINTEGRALHEIGHT gives us exactly the size we want
		// LBS_MULTISEL sounds like it does what we want but it actually doesn't; instead, it toggles item selection regardless of modifier state, which doesn't work like anything else (see http://msdn.microsoft.com/en-us/library/windows/desktop/bb775149%28v=vs.85%29.aspx and http://msdn.microsoft.com/en-us/library/windows/desktop/aa511485.aspx)
		style:			_LBS_NOTIFY | _LBS_NOINTEGRALHEIGHT | _WS_VSCROLL | controlstyle,
		xstyle:			_WS_EX_CLIENTEDGE | controlxstyle,
		altStyle:			_LBS_EXTENDEDSEL | _LBS_NOTIFY | _LBS_NOINTEGRALHEIGHT | _WS_VSCROLL | controlstyle,
		appendMsg:		_LB_ADDSTRING,
		insertBeforeMsg:	_LB_INSERTSTRING,
		deleteMsg:		_LB_DELETESTRING,
		selectedIndexMsg:	_LB_GETCURSEL,
		selectedIndexErr:	negConst(_LB_ERR),
		addSpaceErr:		negConst(_LB_ERRSPACE),
		lenMsg:			_LB_GETCOUNT,
	},
	c_progressbar:		&classData{
		name:			toUTF16(x_PROGRESS_CLASS),
		style:			_PBS_SMOOTH | controlstyle,
		xstyle:			0 | controlxstyle,
		doNotLoadFont:	true,
	},
	c_area:			&classData{
		name:			areaWndClass,
		style:			areastyle,
		xstyle:			areaxstyle,
		storeSysData:		true,
		doNotLoadFont:	true,
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

var (
	_blankString = toUTF16("")
	blankString = utf16ToArg(_blankString)
)

func (s *sysData) make(window *sysData) (err error) {
	ret := make(chan uiret)
	defer close(ret)
	ct := classTypes[s.ctype]
	cid := _HMENU(0)
	pwin := uintptr(_NULL)
	if window != nil {		// this is a child control
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
	uitask <- &uimsg{
		call:		_createWindowEx,	
		p:		[]uintptr{
			uintptr(ct.xstyle),
			utf16ToArg(ct.name),
			blankString,		// we set the window text later
			style,
			negConst(_CW_USEDEFAULT),
			negConst(_CW_USEDEFAULT),
			negConst(_CW_USEDEFAULT),
			negConst(_CW_USEDEFAULT),
			pwin,
			uintptr(cid),
			uintptr(hInstance),
			lpParam,
		},
		ret:	ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		if window != nil {
			window.delChild(cid)
		}
		return fmt.Errorf("error actually creating window/control: %v", r.err)
	}
	if !ct.storeSysData {				// regular control; store s.hwnd ourselves
		s.hwnd = _HWND(r.ret)
	} else if s.hwnd != _HWND(r.ret) {	// we store sysData in storeSysData(); sanity check
		panic(fmt.Errorf("hwnd mismatch creating window/control: storeSysData() stored 0x%X but CreateWindowEx() returned 0x%X", s.hwnd, ret))
	}
	if !ct.doNotLoadFont {
		uitask <- &uimsg{
			call:		_sendMessage,
			p:		[]uintptr{
				uintptr(s.hwnd),
				uintptr(_WM_SETFONT),
				uintptr(_WPARAM(controlFont)),
				uintptr(_LPARAM(_TRUE)),
			},
			ret:		ret,
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
func (s *sysData) firstShow() error {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_showWindow,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(nCmdShow),
		},
		ret:		ret,
	}
	<-ret
	uitask <- &uimsg{
		call:		_updateWindow,
		p:		[]uintptr{uintptr(s.hwnd)},
		ret:		ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		return fmt.Errorf("error updating window for the first time: %v", r.err)
	}
	return nil
}

func (s *sysData) show() {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_showWindow,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(_SW_SHOW),
		},
		ret:		ret,
	}
	<-ret
}

func (s *sysData) hide() {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_showWindow,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(_SW_HIDE),
		},
		ret:		ret,
	}
	<-ret
}

func (s *sysData) setText(text string) {
	ptext := toUTF16(text)
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_setWindowText,
		p:		[]uintptr{
			uintptr(s.hwnd),
			utf16ToArg(ptext),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		panic(fmt.Errorf("error setting window/control text: %v", r.err))
	}
}

func (s *sysData) setRect(x int, y int, width int, height int, winheight int) error {
	r1, _, err := _moveWindow.Call(
		uintptr(s.hwnd),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(_TRUE))
	if r1 == 0 {		// failure
		return fmt.Errorf("error setting window/control rect: %v", err)
	}
	return nil
}

func (s *sysData) isChecked() bool {
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
	return r.ret == _BST_CHECKED
}

func (s *sysData) text() (str string) {
	var tc []uint16

	ret := make(chan uiret)
	defer close(ret)
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
	return syscall.UTF16ToString(tc)
}

func (s *sysData) append(what string) {
	pwhat := toUTF16(what)
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].appendMsg),
			uintptr(_WPARAM(0)),
			utf16ToLPARAM(pwhat),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == uintptr(classTypes[s.ctype].addSpaceErr) {
		panic(fmt.Errorf("out of space adding item to combobox/listbox (last error: %v)", r.err))
	} else if r.ret == uintptr(classTypes[s.ctype].selectedIndexErr) {
		panic(fmt.Errorf("failed to add item to combobox/listbox (last error: %v)", r.err))
	}
}

func (s *sysData) insertBefore(what string, index int) {
	pwhat := toUTF16(what)
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].insertBeforeMsg),
			uintptr(_WPARAM(index)),
			utf16ToLPARAM(pwhat),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == uintptr(classTypes[s.ctype].addSpaceErr) {
		panic(fmt.Errorf("out of space adding item to combobox/listbox (last error: %v)", r.err))
	} else if r.ret == uintptr(classTypes[s.ctype].selectedIndexErr) {
		panic(fmt.Errorf("failed to add item to combobox/listbox (last error: %v)", r.err))
	}
}

func (s *sysData) selectedIndex() int {
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
	if r.ret == uintptr(classTypes[s.ctype].selectedIndexErr) {		// no selection or manually entered text (apparently, for the latter)
		return -1
	}
	return int(r.ret)
}

func (s *sysData) selectedIndices() []int {
	if !s.alternate {		// single-selection list box; use single-selection method
		index := s.selectedIndex()
		if index == -1 {
			return nil
		}
		return []int{index}
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
	if r.ret == negConst(_LB_ERR) {
		panic("UI library internal error: LB_ERR from LB_GETSELCOUNT in what we know is a multi-selection listbox")
	}
	if r.ret == 0 {		// nothing selected
		return nil
	}
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
	r = <-ret
	if r.ret == negConst(_LB_ERR) {
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
		if r.ret == negConst(_LB_ERR) {
			panic("UI library internal error: LB_ERR from LB_GETTEXTLEN in what we know is a valid listbox index (came from LB_GETSELITEMS)")
		}
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
		r = <-ret
		if r.ret == negConst(_LB_ERR) {
			panic("UI library internal error: LB_ERR from LB_GETTEXT in what we know is a valid listbox index (came from LB_GETSELITEMS)")
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
		call:		_getClientRect,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(unsafe.Pointer(&rect)),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == 0 {
		return fmt.Errorf("error getting upper-left of window for resize: %v", r.err)
	}
	// 0 because (0,0) is top-left so no winheight
	err := s.setRect(int(rect.left), int(rect.top), width, height, 0)
	if err != nil {
		return fmt.Errorf("error actually resizing window: %v", err)
	}
	return nil
}

func (s *sysData) delete(index int) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].deleteMsg),
			uintptr(_WPARAM(index)),
			uintptr(0),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == uintptr(classTypes[s.ctype].selectedIndexErr) {
		panic(fmt.Errorf("failed to delete item from combobox/listbox (last error: %v)", r.err))
	}
}

func (s *sysData) setIndeterminate() {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_setWindowLong,
		p:		[]uintptr{
			uintptr(s.hwnd),
			negConst(_GWL_STYLE),
			uintptr(classTypes[s.ctype].style | _PBS_MARQUEE),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == 0 {
		panic(fmt.Errorf("error setting progress bar style to enter indeterminate mode: %v", r.err))
	}
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(_PBM_SETMARQUEE),
			uintptr(_WPARAM(_TRUE)),
			uintptr(0),
		},
		ret:		ret,
	}
	<-ret
	s.isMarquee = true
}

func (s *sysData) setProgress(percent int) {
	if percent == -1 {
		s.setIndeterminate()
		return
	}
	ret := make(chan uiret)
	defer close(ret)
	if s.isMarquee {
		// turn off marquee before switching back
		uitask <- &uimsg{
			call:		_sendMessage,
			p:		[]uintptr{
				uintptr(s.hwnd),
				uintptr(_PBM_SETMARQUEE),
				uintptr(_WPARAM(_FALSE)),
				uintptr(0),
			},
			ret:		ret,
		}
		<-ret
		uitask <- &uimsg{
			call:		_setWindowLong,
			p:		[]uintptr{
				uintptr(s.hwnd),
				negConst(_GWL_STYLE),
				uintptr(classTypes[s.ctype].style),
			},
			ret:		ret,
		}
		r := <-ret
		if r.ret == 0 {
			panic(fmt.Errorf("error setting progress bar style to leave indeterminate mode (percent %d): %v", percent, r.err))
		}
		s.isMarquee = false
	}
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(_PBM_SETPOS),
			uintptr(_WPARAM(percent)),
			uintptr(0),
		},
		ret:		ret,
	}
	<-ret
}

func (s *sysData) len() int {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(classTypes[s.ctype].lenMsg),
			uintptr(_WPARAM(0)),
			uintptr(_LPARAM(0)),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == uintptr(classTypes[s.ctype].selectedIndexErr) {
		panic(fmt.Errorf("unexpected error return from sysData.len(); GetLastError() says %v", r.err))
	}
	return int(r.ret)
}

func (s *sysData) setAreaSize(width int, height int) {
	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(msgSetAreaSize),
			uintptr(width),		// WPARAM is UINT_PTR on Windows XP and newer at least, so we're good with this
			uintptr(height),
		},
		ret:		ret,
	}
	<-ret
}

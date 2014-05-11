// 1 march 2014

package ui

import (
	"fmt"
	"unsafe"
	"sync"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
import "C"

type sysData struct {
	cSysData

	id			C.id
	trackingArea	C.id		// for Area
}

type classData struct {
	make		func(parentWindow C.id, alternate bool, s *sysData) C.id
	getinside		func(scrollview C.id) C.id
	show		func(what C.id)
	hide			func(what C.id)
	settextsel		C.SEL
	textsel		C.SEL
	alttextsel		C.SEL
	append		func(id C.id, what string, alternate bool)
	insertBefore	func(id C.id, what string, before int, alternate bool)
	selIndex		func(id C.id) int
	selIndices		func(id C.id) []int
	selTexts		func(id C.id) []string
	delete		func(id C.id, index int)
	len			func(id C.id) int
	selectIndex	func(id C.id, index int, alternate bool)
	selectIndices	func(id C.id, indices []int)
}

var (
	_NSWindow = objc_getClass("NSWindow")
	_NSButton = objc_getClass("NSButton")
	_NSPopUpButton = objc_getClass("NSPopUpButton")
	_NSComboBox = objc_getClass("NSComboBox")
	_NSTextField = objc_getClass("NSTextField")
	_NSSecureTextField = objc_getClass("NSSecureTextField")
	_NSProgressIndicator = objc_getClass("NSProgressIndicator")

	_initWithContentRect = sel_getUid("initWithContentRect:styleMask:backing:defer:")
	_initWithFrame = sel_getUid("initWithFrame:")
	_setDelegate = sel_getUid("setDelegate:")
	_makeKeyAndOrderFront = sel_getUid("makeKeyAndOrderFront:")
	_orderOut = sel_getUid("orderOut:")
	_setHidden = sel_getUid("setHidden:")
	_setTitle = sel_getUid("setTitle:")
	_setStringValue = sel_getUid("setStringValue:")
	_setFrame = sel_getUid("setFrame:")
	_state = sel_getUid("state")
	_title = sel_getUid("title")
	_stringValue = sel_getUid("stringValue")
	_frame = sel_getUid("frame")
	_setContentSize = sel_getUid("setContentSize:")
	_setBezelStyle = sel_getUid("setBezelStyle:")
	_setTarget = sel_getUid("setTarget:")
	_setAction = sel_getUid("setAction:")
	_contentView = sel_getUid("contentView")
	_addSubview = sel_getUid("addSubview:")
	_setButtonType = sel_getUid("setButtonType:")
	_initWithFramePullsDown = sel_getUid("initWithFrame:pullsDown:")
	_setUsesDataSource = sel_getUid("setUsesDataSource:")
	_addItemWithTitle = sel_getUid("addItemWithTitle:")
	_insertItemWithTitleAtIndex = sel_getUid("insertItemWithTitle:atIndex:")
	_removeItemAtIndex = sel_getUid("removeItemAtIndex:")
	_titleOfSelectedItem = sel_getUid("titleOfSelectedItem")
	_indexOfSelectedItem = sel_getUid("indexOfSelectedItem")
	_addItemWithObjectValue = sel_getUid("addItemWithObjectValue:")
	_insertItemWithObjectValueAtIndex = sel_getUid("insertItemWithObjectValue:atIndex:")
	_setEditable = sel_getUid("setEditable:")
	_setBordered = sel_getUid("setBordered:")
	_setDrawsBackground = sel_getUid("setDrawsBackground:")
	_cell = sel_getUid("cell")
	_setLineBreakMode = sel_getUid("setLineBreakMode:")
	_setStyle = sel_getUid("setStyle:")
	_setControlSize = sel_getUid("setControlSize:")
	_setIndeterminate = sel_getUid("setIndeterminate:")
	_setDoubleValue = sel_getUid("setDoubleValue:")
	_numberOfItems = sel_getUid("numberOfItems")
	_selectItemAtIndex = sel_getUid("selectItemAtIndex:")
	_deselectItemAtIndex = sel_getUid("deselectItemAtIndex:")
)

// because the only way to make a new NSControl/NSView is with a frame (it gets overridden later)
func initWithDummyFrame(self C.id) C.id {
	return C.objc_msgSend_rect(self, _initWithFrame,
		C.int64_t(0), C.int64_t(0), C.int64_t(100), C.int64_t(100))
}

func addControl(parentWindow C.id, control C.id) {
	windowView := C.objc_msgSend_noargs(parentWindow, _contentView)
	C.objc_msgSend_id(windowView, _addSubview, control)
}

func controlShow(what C.id) {
	C.objc_msgSend_bool(what, _setHidden, C.BOOL(C.NO))
}

func controlHide(what C.id) {
	C.objc_msgSend_bool(what, _setHidden, C.BOOL(C.YES))
}

const (
	_NSRegularControlSize = 0
)

// By default some controls do not use the correct font.
// These functions set the appropriate control font.
// Which one is used on each control was determined by comparing https://developer.apple.com/library/mac/documentation/UserExperience/Conceptual/AppleHIGuidelines/Characteristics/Characteristics.html#//apple_ref/doc/uid/TP40002721-SW10 to what Interface Builder says for each control.
// (not applicable to ProgressBar, Area)

// Button, Checkbox, Combobox, LineEdit, Label, Listbox
func applyStandardControlFont(id C.id) {
	C.objc_setFont(id, _NSRegularControlSize)
}

var classTypes = [nctypes]*classData{
	c_window:		&classData{
		make:		func(parentWindow C.id, alternate bool, s *sysData) C.id {
			const (
				_NSBorderlessWindowMask = 0
				_NSTitledWindowMask = 1 << 0
				_NSClosableWindowMask = 1 << 1
				_NSMiniaturizableWindowMask = 1 << 2
				_NSResizableWindowMask = 1 << 3
				_NSTexturedBackgroundWindowMask = 1 << 8

				_NSBackingStoreBuffered = 2		// the only backing store method that Apple says we should use (the others are legacy)
			)

			// we have to specify a content rect to start; it will be overridden soon though
			win := C.objc_msgSend_noargs(_NSWindow, _alloc)
			win = C.objc_msgSend_rect_uint_uint_bool(win,
				_initWithContentRect,
				C.int64_t(0), C.int64_t(0), C.int64_t(100), C.int64_t(100),
				C.uintptr_t(_NSTitledWindowMask | _NSClosableWindowMask | _NSMiniaturizableWindowMask | _NSResizableWindowMask),
				C.uintptr_t(_NSBackingStoreBuffered),
				C.BOOL(C.YES))			// defer creation of device until we show the window
			C.objc_msgSend_id(win, _setDelegate, appDelegate)
			// we do not need setAcceptsMouseMovedEvents: here since we are using a tracking rect in Areas for that
			return win
		},
		show:		func(what C.id) {
			C.objc_msgSend_id(what, _makeKeyAndOrderFront, what)
		},
		hide:			func(what C.id) {
			C.objc_msgSend_id(what, _orderOut, what)
		},
		settextsel:		_setTitle,
		textsel:		_title,
	},
	c_button:			&classData{
		make:		func(parentWindow C.id, alternate bool, s *sysData) C.id {
			const (
				_NSRoundedBezelStyle = 1
			)

			button := C.objc_msgSend_noargs(_NSButton, _alloc)
			button = initWithDummyFrame(button)
			C.objc_msgSend_uint(button, _setBezelStyle, C.uintptr_t(_NSRoundedBezelStyle))
			C.objc_msgSend_id(button, _setTarget, appDelegate)
			C.objc_msgSend_sel(button, _setAction, _buttonClicked)
			applyStandardControlFont(button)
			addControl(parentWindow, button)
			return button
		},
		show:		controlShow,
		hide:			controlHide,
		settextsel:		_setTitle,
		textsel:		_title,
	},
	c_checkbox:		&classData{
		make:		func(parentWindow C.id, alternate bool, s *sysData) C.id {
			const (
				_NSSwitchButton = 3
			)

			checkbox := C.objc_msgSend_noargs(_NSButton, _alloc)
			checkbox = initWithDummyFrame(checkbox)
			C.objc_msgSend_uint(checkbox, _setButtonType, C.uintptr_t(_NSSwitchButton))
			applyStandardControlFont(checkbox)
			addControl(parentWindow, checkbox)
			return checkbox
		},
		show:		controlShow,
		hide:			controlHide,
		settextsel:		_setTitle,
		textsel:		_title,
	},
	c_combobox:		&classData{
		make:		func(parentWindow C.id, alternate bool, s *sysData) C.id {
			var combobox C.id

			if alternate {
				combobox = C.objc_msgSend_noargs(_NSComboBox, _alloc)
				combobox = initWithDummyFrame(combobox)
				C.objc_msgSend_bool(combobox, _setUsesDataSource, C.BOOL(C.NO))
			} else {
				combobox = C.objc_msgSend_noargs(_NSPopUpButton, _alloc)
				combobox = C.objc_msgSend_rect_bool(combobox, _initWithFramePullsDown,
					C.int64_t(0), C.int64_t(0), C.int64_t(100), C.int64_t(100),
					C.BOOL(C.NO))
			}
			applyStandardControlFont(combobox)
			addControl(parentWindow, combobox)
			return combobox
		},
		show:		controlShow,
		hide:			controlHide,
		textsel:		_titleOfSelectedItem,
		alttextsel:		_stringValue,
		append:		func(id C.id, what string, alternate bool) {
			str := toNSString(what)
			if alternate {
				C.objc_msgSend_id(id, _addItemWithObjectValue, str)
			} else {
				C.objc_msgSend_id(id, _addItemWithTitle, str)
			}
		},
		insertBefore:	func(id C.id, what string, before int, alternate bool) {
			str := toNSString(what)
			if alternate {
				C.objc_msgSend_id_int(id, _insertItemWithObjectValueAtIndex, str, C.intptr_t(before))
			} else {
				C.objc_msgSend_id_int(id, _insertItemWithTitleAtIndex, str, C.intptr_t(before))
			}
		},
		selIndex:		func(id C.id) int {
			return int(C.objc_msgSend_intret_noargs(id, _indexOfSelectedItem))
		},
		delete:		func(id C.id, index int) {
			C.objc_msgSend_int(id, _removeItemAtIndex, C.intptr_t(index))
		},
		len:			func(id C.id) int {
			return int(C.objc_msgSend_intret_noargs(id, _numberOfItems))
		},
		selectIndex:	func(id C.id, index int, alternate bool) {
			// NSPopUpButton makes this easy
			if !alternate {
				C.objc_msgSend_int(id, _selectItemAtIndex, C.intptr_t(index))
				return
			}
			// NSComboBox doesn't document that we can do [cb selectItemAtIndex:-1], so we have to do this to be safe
			if index == -1 {
				idx := C.objc_msgSend_intret_noargs(id, _indexOfSelectedItem)
				if idx != -1 {
					C.objc_msgSend_int(id, _deselectItemAtIndex, idx)
				}
				return
			}
			C.objc_msgSend_int(id, _selectItemAtIndex, C.intptr_t(index))
		},
	},
	c_lineedit:		&classData{
		make:		func(parentWindow C.id, alternate bool, s *sysData) C.id {
			var lineedit C.id

			if alternate {
				lineedit = C.objc_msgSend_noargs(_NSSecureTextField, _alloc)
			} else {
				lineedit = C.objc_msgSend_noargs(_NSTextField, _alloc)
			}
			lineedit = initWithDummyFrame(lineedit)
			applyStandardControlFont(lineedit)
			addControl(parentWindow, lineedit)
			return lineedit
		},
		show:		controlShow,
		hide:			controlHide,
		settextsel:		_setStringValue,
		textsel:		_stringValue,
		alttextsel:		_stringValue,
	},
	c_label:			&classData{
		make:		func(parentWindow C.id, alternate bool, s *sysData) C.id {
			const (
				_NSLineBreakByWordWrapping = iota
				_NSLineBreakByCharWrapping
				_NSLineBreakByClipping
				_NSLineBreakByTruncatingHead
				_NSLineBreakByTruncatingTail
				_NSLineBreakByTruncatingMiddle
			)

			label := C.objc_msgSend_noargs(_NSTextField, _alloc)
			label = initWithDummyFrame(label)
			C.objc_msgSend_bool(label, _setEditable, C.BOOL(C.NO))
			C.objc_msgSend_bool(label, _setBordered, C.BOOL(C.NO))
			C.objc_msgSend_bool(label, _setDrawsBackground, C.BOOL(C.NO))
			// this disables both word wrap AND ellipsizing in one fell swoop
			// we have to send to the control's cell for this
			C.objc_msgSend_uint(C.objc_msgSend_noargs(label, _cell),
				_setLineBreakMode, _NSLineBreakByClipping)
			// for a multiline label, we either use WordWrapping and send setTruncatesLastVisibleLine: to disable ellipsizing OR use one of those ellipsizing styles
			applyStandardControlFont(label)
			addControl(parentWindow, label)
			return label
		},
		show:		controlShow,
		hide:			controlHide,
		settextsel:		_setStringValue,
		textsel:		_stringValue,
	},
	c_listbox:			&classData{
		make:		makeListbox,
		show:		controlShow,
		hide:			controlHide,
		append:		appendListbox,
		insertBefore:	insertListboxBefore,
		selIndices:	selectedListboxIndices,
		selTexts:		selectedListboxTexts,
		delete:		deleteListbox,
		len:			listboxLen,
		selectIndices:	selectListboxIndices,
	},
	c_progressbar:		&classData{
		make:		func(parentWindow C.id, alternate bool, s *sysData) C.id {
			const (
				_NSProgressIndicatorBarStyle = 0
			)

			pbar := C.objc_msgSend_noargs(_NSProgressIndicator, _alloc)
			pbar = initWithDummyFrame(pbar)
			// NSProgressIndicatorStyle doesn't have an explicit typedef; just use int for now
			C.objc_msgSend_int(pbar, _setStyle, _NSProgressIndicatorBarStyle)
			C.objc_msgSend_uint(pbar, _setControlSize, C.uintptr_t(_NSRegularControlSize))
			C.objc_msgSend_bool(pbar, _setIndeterminate, C.BOOL(C.NO))
			addControl(parentWindow, pbar)
			return pbar
		},
		show:		controlShow,
		hide:			controlHide,
	},
	c_area:			&classData{
		make:		makeArea,
		getinside:		areaInScrollView,
		show:		controlShow,
		hide:			controlHide,
	},
}

// I need to access sysData from appDelegate, but appDelegate doesn't store any data. So, this.
var (
	sysdatas = make(map[C.id]*sysData)
	sysdatalock sync.Mutex
)

func addSysData(key C.id, value *sysData) {
	sysdatalock.Lock()
	sysdatas[key] = value
	sysdatalock.Unlock()
}

func getSysData(key C.id) *sysData {
	sysdatalock.Lock()
	defer sysdatalock.Unlock()

	v, ok := sysdatas[key]
	if !ok {
		panic(fmt.Errorf("internal error: getSysData(%v) unknown", key))
	}
	return v
}

func (s *sysData) make(window *sysData) error {
	var parentWindow C.id

	ct := classTypes[s.ctype]
	if window != nil {
		parentWindow = window.id
	}
	ret := make(chan C.id)
	defer close(ret)
	uitask <- func() {
		ret <- ct.make(parentWindow, s.alternate, s)
	}
	s.id = <-ret
	if ct.getinside != nil {
		uitask <- func() {
			ret <- ct.getinside(s.id)
		}
		addSysData(<-ret, s)
	} else {
		addSysData(s.id, s)
	}
	return nil
}

// used for Windows; nothing special needed elsewhere
func (s *sysData) firstShow() error {
	s.show()
	return nil
}

func (s *sysData) show() {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].show(s.id)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) hide() {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].hide(s.id)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) setText(text string) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		C.objc_msgSend_id(s.id, classTypes[s.ctype].settextsel, toNSString(text))
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) setRect(x int, y int, width int, height int, winheight int) error {
	// winheight - y because (0,0) is the bottom-left corner of the window and not the top-left corner
	// (winheight - y) - height because (x, y) is the bottom-left corner of the control and not the top-left
	C.objc_msgSend_rect(s.id, _setFrame,
		C.int64_t(x), C.int64_t((winheight - y) - height), C.int64_t(width), C.int64_t(height))
	return nil
}

func (s *sysData) isChecked() bool {
	const (
		NSOnState = 1
	)

	ret := make(chan bool)
	defer close(ret)
	uitask <- func() {
		k := C.objc_msgSend_noargs(s.id, _state)
		ret <- uintptr(unsafe.Pointer(k)) == NSOnState
	}
	return <-ret
}

func (s *sysData) text() string {
	sel := classTypes[s.ctype].textsel
	if s.alternate {
		sel = classTypes[s.ctype].alttextsel
	}
	ret := make(chan string)
	defer close(ret)
	uitask <- func() {
		str := C.objc_msgSend_noargs(s.id, sel)
		ret <- fromNSString(str)
	}
	return <-ret
}

func (s *sysData) append(what string) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].append(s.id, what, s.alternate)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) insertBefore(what string, before int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].insertBefore(s.id, what, before, s.alternate)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) selectedIndex() int {
	ret := make(chan int)
	defer close(ret)
	uitask <- func() {
		ret <- classTypes[s.ctype].selIndex(s.id)
	}
	return <-ret
}

func (s *sysData) selectedIndices() []int {
	ret := make(chan []int)
	defer close(ret)
	uitask <- func() {
		ret <- classTypes[s.ctype].selIndices(s.id)
	}
	return <-ret
}

func (s *sysData) selectedTexts() []string {
	ret := make(chan []string)
	defer close(ret)
	uitask <- func() {
		ret <- classTypes[s.ctype].selTexts(s.id)
	}
	return <-ret
}

func (s *sysData) setWindowSize(width int, height int) error {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		// use -[NSWindow setContentSize:], which will resize the window without taking the titlebar as part of the given size and without needing us to consider the window's position (the function takes care of both for us)
		C.objc_msgSend_size(s.id, _setContentSize,
			C.int64_t(width), C.int64_t(height))
		C.objc_msgSend_noargs(s.id, _display)		// TODO needed?
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) delete(index int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].delete(s.id, index)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) setProgress(percent int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		if percent == -1 {
			// At least on Mac OS X 10.8, if the progressbar was already on 0 or 100% when turning on indeterminate mode, the indeterminate animation won't play, leaving just a still progress bar. This is a workaround. Note the selector call order.
			// TODO will the value chosen affect the animation speed?
			C.objc_msgSend_double(s.id, _setDoubleValue, C.double(50))
			C.objc_msgSend_bool(s.id, _setIndeterminate, C.BOOL(C.YES))
		} else {
			C.objc_msgSend_bool(s.id, _setIndeterminate, C.BOOL(C.NO))
			C.objc_msgSend_double(s.id, _setDoubleValue, C.double(percent))
		}
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) len() int {
	ret := make(chan int)
	defer close(ret)
	uitask <- func() {
		ret <- classTypes[s.ctype].len(s.id)
	}
	return <-ret
}

func (s *sysData) setAreaSize(width int, height int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		area := areaInScrollView(s.id)
		C.objc_msgSend_rect(area, _setFrame,
			C.int64_t(0), C.int64_t(0), C.int64_t(width), C.int64_t(height))
		C.objc_msgSend_noargs(area, _display)		// and redraw
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) selectIndex(index int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].selectIndex(s.id, index, s.alternate)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) selectIndices(indices []int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].selectIndices(s.id, indices)
		ret <- struct{}{}
	}
	<-ret
}

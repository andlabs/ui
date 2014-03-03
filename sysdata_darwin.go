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

	id	C.id
}

type classData struct {
	make		func(parentWindow C.id, alternate bool) C.id
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
}

var (
	_NSWindow = objc_getClass("NSWindow")
	_NSButton = objc_getClass("NSButton")
	_NSPopUpButton = objc_getClass("NSPopUpButton")
	_NSComboBox = objc_getClass("NSComboBox")
	_NSTextField = objc_getClass("NSTextField")
	_NSSecureTextField = objc_getClass("NSSecureTextField")

	_initWithContentRect = sel_getUid("initWithContentRect:styleMask:backing:defer:")
	_initWithFrame = sel_getUid("initWithFrame:")
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
	_setFrameDisplay = sel_getUid("setFrame:display:")
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
)

func controlShow(what C.id) {
	C.objc_msgSend_bool(what, _setHidden, C.BOOL(C.NO))
}

func controlHide(what C.id) {
	C.objc_msgSend_bool(what, _setHidden, C.BOOL(C.YES))
}

var classTypes = [nctypes]*classData{
	c_window:		&classData{
		make:		func(parentWindow C.id, alternate bool) C.id {
			const (
				NSBorderlessWindowMask = 0
				NSTitledWindowMask = 1 << 0
				NSClosableWindowMask = 1 << 1
				NSMiniaturizableWindowMask = 1 << 2
				NSResizableWindowMask = 1 << 3
				NSTexturedBackgroundWindowMask = 1 << 8
			)

			// we have to specify a content rect to start; it will be overridden soon though
			win := objc_alloc(_NSWindow)
			win = objc_msgSend_rect_uint_uint_bool(win,
				_initWithContentRect,
				0, 0, 100, 100,
				NSTitledWindowMask | NSClosableWindowMask | NSMiniaturizableWindowMask | NSResizableWindowMask,
				2,					// NSBackingStoreBuffered - the only backing store method that Apple says we should use (the others are legacy)
				C.BOOL(C.YES))			// defer creation of device until we show the window
			objc_setDelegate(win, appDelegate)
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
		make:		func(parentWindow C.id, alternate bool) C.id {
			button := objc_alloc(_NSButton)
			// NSControl requires that we specify a frame; dummy frame for now
			button = objc_msgSend_rect(button, _initWithFrame,
				0, 0, 100, 100)
			objc_msgSend_uint(button, _setBezelStyle, 1)		// NSRoundedBezelStyle
			C.objc_msgSend_id(button, _setTarget, appDelegate)
			C.objc_msgSend_sel(button, _setAction, _buttonClicked)
			windowView := C.objc_msgSend_noargs(parentWindow, _contentView)
			C.objc_msgSend_id(windowView, _addSubview, button)
			return button
		},
		show:		controlShow,
		hide:			controlHide,
		settextsel:		_setTitle,
		textsel:		_title,
	},
	c_checkbox:		&classData{
		make:		func(parentWindow C.id, alternate bool) C.id {
			checkbox := objc_alloc(_NSButton)
			checkbox = objc_msgSend_rect(checkbox, _initWithFrame,
				0, 0, 100, 100)
			objc_msgSend_uint(checkbox, _setButtonType, 3)		// NSSwitchButton
			windowView := C.objc_msgSend_noargs(parentWindow, _contentView)
			C.objc_msgSend_id(windowView, _addSubview, checkbox)
			return checkbox
		},
		show:		controlShow,
		hide:			controlHide,
		settextsel:		_setTitle,
		textsel:		_title,
	},
	c_combobox:		&classData{
		make:		func(parentWindow C.id, alternate bool) C.id {
			var combobox C.id

			if alternate {
				combobox = objc_alloc(_NSComboBox)
				combobox = objc_msgSend_rect(combobox, _initWithFrame,
					0, 0, 100, 100)
				C.objc_msgSend_bool(combobox, _setUsesDataSource, C.BOOL(C.NO))
			} else {
				combobox = objc_alloc(_NSPopUpButton)
				combobox = objc_msgSend_rect_bool(combobox, _initWithFramePullsDown,
					0, 0, 100, 100,
					C.BOOL(C.NO))
			}
			windowView := C.objc_msgSend_noargs(parentWindow, _contentView)
			C.objc_msgSend_id(windowView, _addSubview, combobox)
			return combobox
		},
		show:		controlShow,
		hide:			controlHide,
		// TODO setText
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
	},
	c_lineedit:		&classData{
		make:		func(parentWindow C.id, alternate bool) C.id {
			var lineedit C.id

			if alternate {
				lineedit = objc_alloc(_NSSecureTextField)
			} else {
				lineedit = objc_alloc(_NSTextField)
			}
			lineedit = objc_msgSend_rect(lineedit, _initWithFrame,
				0, 0, 100, 100)
			windowView := C.objc_msgSend_noargs(parentWindow, _contentView)
			C.objc_msgSend_id(windowView, _addSubview, lineedit)
			return lineedit
		},
		show:		controlShow,
		hide:			controlHide,
		settextsel:		_setStringValue,
		textsel:		_stringValue,
		alttextsel:		_stringValue,
	},
	c_label:			&classData{
		make:		func(parentWindow C.id, alternate bool) C.id {
			label := objc_alloc(_NSTextField)
			label = objc_msgSend_rect(label, _initWithFrame,
				0, 0, 100, 100)
			C.objc_msgSend_bool(label, _setEditable, C.BOOL(C.NO))
			C.objc_msgSend_bool(label, _setBordered, C.BOOL(C.NO))
			C.objc_msgSend_bool(label, _setDrawsBackground, C.BOOL(C.NO))
			// TODO others?
			windowView := C.objc_msgSend_noargs(parentWindow, _contentView)
			C.objc_msgSend_id(windowView, _addSubview, label)
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
	},
	c_progressbar:		&classData{
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

func (s *sysData) make(initText string, window *sysData) error {
	var parentWindow C.id

	ct := classTypes[s.ctype]
	if ct.make == nil {
		println(s.ctype, "not implemented")
		return nil
	}
	if window != nil {
		parentWindow = window.id
	}
	ret := make(chan C.id)
	defer close(ret)
	uitask <- func() {
		ret <- ct.make(parentWindow, s.alternate)
	}
	s.id = <-ret
	err := s.setText(initText)
	if err != nil {
		return fmt.Errorf("error setting initial text of new window/control: %v", err)
	}
	addSysData(s.id, s)
	return nil
}

func (s *sysData) show() error {
if classTypes[s.ctype].show == nil { return nil }
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].show(s.id)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) hide() error {
if classTypes[s.ctype].hide == nil { return nil }
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].hide(s.id)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) setText(text string) error {
var zero C.SEL
if classTypes[s.ctype].settextsel == zero { return nil }
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		C.objc_msgSend_id(s.id, classTypes[s.ctype].settextsel, toNSString(text))
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) setRect(x int, y int, width int, height int) error {
if classTypes[s.ctype].make == nil { return nil }
	objc_msgSend_rect(s.id, _setFrame, x, y, width, height)
	return nil
}

func (s *sysData) isChecked() bool {
if classTypes[s.ctype].make == nil { return false }
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
var zero C.SEL
if classTypes[s.ctype].textsel == zero { return "" }
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

func (s *sysData) append(what string) error {
if classTypes[s.ctype].append == nil { return nil }
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].append(s.id, what, s.alternate)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) insertBefore(what string, before int) error {
if classTypes[s.ctype].insertBefore == nil { return nil }
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].insertBefore(s.id, what, before, s.alternate)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) selectedIndex() int {
if classTypes[s.ctype].selIndex == nil { return -1 }
	ret := make(chan int)
	defer close(ret)
	uitask <- func() {
		ret <- classTypes[s.ctype].selIndex(s.id)
	}
	return <-ret
}

func (s *sysData) selectedIndices() []int {
	// TODO
	return nil
}

func (s *sysData) selectedTexts() []string {
	// TODO
	return nil
}

func (s *sysData) setWindowSize(width int, height int) error {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		// we need to get the top left point
		r := C.objc_msgSend_stret_rect_noargs(s.id, _frame)
		objc_msgSend_rect_bool(s.id, _setFrameDisplay,
			int(r.x), int(r.y), width, height,
			C.BOOL(C.YES))		// TODO set to NO to prevent subviews from being redrawn before they are resized?
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) delete(index int) error {
if classTypes[s.ctype].delete == nil { return nil }
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].delete(s.id, index)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) setProgress(percent int) {
	// TODO
}

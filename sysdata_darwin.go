// 1 march 2014

package ui

import (
	"fmt"
	"sync"
)

// #cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit
// #include "objc_darwin.h"
// #include "sysdata_darwin.h"
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
	_startAnimation = sel_getUid("startAnimation:")
	_stopAnimation = sel_getUid("stopAnimation:")
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
	C.addControl(parentWindow, control)
}

func controlShow(what C.id) {
	C.controlShow(what)
}

func controlHide(what C.id) {
	C.controlHide(what)
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
	C.applyStandardControlFont(id)
}

var classTypes = [nctypes]*classData{
	c_window:		&classData{
		make:		func(parentWindow C.id, alternate bool, s *sysData) C.id {
			win := C.makeWindow()
			C.objc_msgSend_id(win, _setDelegate, appDelegate)
			// we do not need setAcceptsMouseMovedEvents: here since we are using a tracking rect in Areas for that
			return win
		},
		show:		func(what C.id) {
			C.windowShow(what)
		},
		hide:			func(what C.id) {
			C.windowHide(what)
		},
		settextsel:		_setTitle,
		textsel:		_title,
	},
	c_button:			&classData{
		make:		func(parentWindow C.id, alternate bool, s *sysData) C.id {
			button := C.makeButton()
			C.buttonSetTargetAction(button, appDelegate)
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
			checkbox := C.makeCheckbox()
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
				combobox = C.makeCombobox(C.YES)
			} else {
				combobox = C.makeCombobox(C.NO)
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
				C.comboboxAppend(id, C.YES, str)
			} else {
				C.comboboxAppend(id, C.NO, str)
			}
		},
		insertBefore:	func(id C.id, what string, before int, alternate bool) {
			str := toNSString(what)
			if alternate {
				C.comboboxInsertBefore(id, C.YES, str, C.intptr_t(before))
			} else {
				C.comboboxInsertBefore(id, C.NO, str, C.intptr_t(before))
			}
		},
		selIndex:		func(id C.id) int {
			return int(C.comboboxSelectedIndex(id))
		},
		delete:		func(id C.id, index int) {
			C.comboboxDelete(id, C.intptr_t(index))
		},
		len:			func(id C.id) int {
			return int(C.comboboxLen(id))
		},
		selectIndex:	func(id C.id, index int, alternate bool) {
			// NSPopUpButton makes this easy
			if alternate {
				C.comboboxSelectIndex(id, C.YES, C.intptr_t(index))
			} else {
				C.comboboxSelectIndex(id, C.NO, C.intptr_t(index))
			}
		},
	},
	c_lineedit:		&classData{
		make:		func(parentWindow C.id, alternate bool, s *sysData) C.id {
			var lineedit C.id

			if alternate {
				lineedit = C.makeLineEdit(C.YES)
			} else {
				lineedit = C.makeLineEdit(C.NO)
			}
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
			label := C.makeLabel()
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
			pbar := C.makeProgressBar()
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
	C.setRect(s.id,
		C.intptr_t(x), C.intptr_t((winheight - y) - height),
		C.intptr_t(width), C.intptr_t(height))
	return nil
}

func (s *sysData) isChecked() bool {
	ret := make(chan bool)
	defer close(ret)
	uitask <- func() {
		ret <- C.isCheckboxChecked(s.id) != C.NO
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
		C.windowSetContentSize(s.id, C.intptr_t(width), C.intptr_t(height))
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
		C.setProgress(s.id, C.intptr_t(percent))
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
		C.setAreaSize(s.id, C.intptr_t(width), C.intptr_t(height))
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

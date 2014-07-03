// 1 march 2014

package ui

import (
	"fmt"
	"sync"
)

// #include "objc_darwin.h"
import "C"

type sysData struct {
	cSysData

	id           C.id
	trackingArea C.id // for Area
}

type classData struct {
	make         func(parentWindow C.id, alternate bool, s *sysData) C.id
	getinside    func(scrollview C.id) C.id
	show         func(what C.id)
	hide         func(what C.id)
	settext      func(what C.id, text C.id)
	text         func(what C.id, alternate bool) C.id
	append       func(id C.id, what string, alternate bool)
	insertBefore func(id C.id, what string, before int, alternate bool)
	selIndex     func(id C.id) int
	selIndices   func(id C.id) []int
	selTexts     func(id C.id) []string
	delete       func(id C.id, index int)
	len          func(id C.id) int
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

// By default some controls do not use the correct font.
// These functions set the appropriate control font.
// Which one is used on each control was determined by comparing https://developer.apple.com/library/mac/documentation/UserExperience/Conceptual/AppleHIGuidelines/Characteristics/Characteristics.html#//apple_ref/doc/uid/TP40002721-SW10 to what Interface Builder says for each control.
// (not applicable to ProgressBar, Area)

// Button, Checkbox, Combobox, LineEdit, Label, Listbox
func applyStandardControlFont(id C.id) {
	C.applyStandardControlFont(id)
}

var classTypes = [nctypes]*classData{
	c_window: &classData{
		make: func(parentWindow C.id, alternate bool, s *sysData) C.id {
			return C.makeWindow(appDelegate)
		},
		show: func(what C.id) {
			C.windowShow(what)
		},
		hide: func(what C.id) {
			C.windowHide(what)
		},
		settext: func(what C.id, text C.id) {
			C.windowSetTitle(what, text)
		},
		text: func(what C.id, alternate bool) C.id {
			return C.windowTitle(what)
		},
	},
	c_button: &classData{
		make: func(parentWindow C.id, alternate bool, s *sysData) C.id {
			button := C.makeButton()
			C.buttonSetTargetAction(button, appDelegate)
			applyStandardControlFont(button)
			addControl(parentWindow, button)
			return button
		},
		show: controlShow,
		hide: controlHide,
		settext: func(what C.id, text C.id) {
			C.buttonSetText(what, text)
		},
		text: func(what C.id, alternate bool) C.id {
			return C.buttonText(what)
		},
	},
	c_checkbox: &classData{
		make: func(parentWindow C.id, alternate bool, s *sysData) C.id {
			checkbox := C.makeCheckbox()
			applyStandardControlFont(checkbox)
			addControl(parentWindow, checkbox)
			return checkbox
		},
		show: controlShow,
		hide: controlHide,
		settext: func(what C.id, text C.id) {
			C.buttonSetText(what, text)
		},
		text: func(what C.id, alternate bool) C.id {
			return C.buttonText(what)
		},
	},
	c_combobox: &classData{
		make: func(parentWindow C.id, alternate bool, s *sysData) C.id {
			combobox := C.makeCombobox(toBOOL(alternate))
			applyStandardControlFont(combobox)
			addControl(parentWindow, combobox)
			return combobox
		},
		show: controlShow,
		hide: controlHide,
		text: func(what C.id, alternate bool) C.id {
			return C.comboboxText(what, toBOOL(alternate))
		},
		append: func(id C.id, what string, alternate bool) {
			C.comboboxAppend(id, toBOOL(alternate), toNSString(what))
		},
		insertBefore: func(id C.id, what string, before int, alternate bool) {
			C.comboboxInsertBefore(id, toBOOL(alternate),
				toNSString(what), C.intptr_t(before))
		},
		selIndex: func(id C.id) int {
			return int(C.comboboxSelectedIndex(id))
		},
		delete: func(id C.id, index int) {
			C.comboboxDelete(id, C.intptr_t(index))
		},
		len: func(id C.id) int {
			return int(C.comboboxLen(id))
		},
	},
	c_lineedit: &classData{
		make: func(parentWindow C.id, alternate bool, s *sysData) C.id {
			lineedit := C.makeLineEdit(toBOOL(alternate))
			applyStandardControlFont(lineedit)
			addControl(parentWindow, lineedit)
			return lineedit
		},
		show: controlShow,
		hide: controlHide,
		settext: func(what C.id, text C.id) {
			C.lineeditSetText(what, text)
		},
		text: func(what C.id, alternate bool) C.id {
			return C.lineeditText(what)
		},
	},
	c_label: &classData{
		make: func(parentWindow C.id, alternate bool, s *sysData) C.id {
			label := C.makeLabel()
			applyStandardControlFont(label)
			addControl(parentWindow, label)
			return label
		},
		show: controlShow,
		hide: controlHide,
		settext: func(what C.id, text C.id) {
			C.lineeditSetText(what, text)
		},
		text: func(what C.id, alternate bool) C.id {
			return C.lineeditText(what)
		},
	},
	c_listbox: &classData{
		make:         makeListbox,
		show:         controlShow,
		hide:         controlHide,
		append:       listboxAppend,
		insertBefore: listboxInsertBefore,
		selIndices:   listboxSelectedIndices,
		selTexts:     listboxSelectedTexts,
		delete:       listboxDelete,
		len:          listboxLen,
	},
	c_progressbar: &classData{
		make: func(parentWindow C.id, alternate bool, s *sysData) C.id {
			pbar := C.makeProgressBar()
			addControl(parentWindow, pbar)
			return pbar
		},
		show: controlShow,
		hide: controlHide,
	},
	c_area: &classData{
		make:      makeArea,
		getinside: areaInScrollView,
		show:      controlShow,
		hide:      controlHide,
	},
}

// I need to access sysData from appDelegate, but appDelegate doesn't store any data. So, this.
var (
	sysdatas    = make(map[C.id]*sysData)
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
	s.id = ct.make(parentWindow, s.alternate, s)
	if ct.getinside != nil {
		addSysData(ct.getinside(s.id), s)
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
	classTypes[s.ctype].show(s.id)
}

func (s *sysData) hide() {
	classTypes[s.ctype].hide(s.id)
}

func (s *sysData) setText(text string) {
	classTypes[s.ctype].settext(s.id, toNSString(text))
}

func (s *sysData) setRect(x int, y int, width int, height int, winheight int) error {
	// winheight - y because (0,0) is the bottom-left corner of the window and not the top-left corner
	// (winheight - y) - height because (x, y) is the bottom-left corner of the control and not the top-left
	C.setRect(s.id,
		C.intptr_t(x), C.intptr_t((winheight-y)-height),
		C.intptr_t(width), C.intptr_t(height))
	return nil
}

func (s *sysData) isChecked() bool {
	return C.isCheckboxChecked(s.id) != C.NO
}

func (s *sysData) text() string {
	str := classTypes[s.ctype].text(s.id, s.alternate)
	return fromNSString(str)
}

func (s *sysData) append(what string) {
	classTypes[s.ctype].append(s.id, what, s.alternate)
}

func (s *sysData) insertBefore(what string, before int) {
	classTypes[s.ctype].insertBefore(s.id, what, before, s.alternate)
}

func (s *sysData) selectedIndex() int {
	return classTypes[s.ctype].selIndex(s.id)
}

func (s *sysData) selectedIndices() []int {
	return classTypes[s.ctype].selIndices(s.id)
}

func (s *sysData) selectedTexts() []string {
	return classTypes[s.ctype].selTexts(s.id)
}

func (s *sysData) setWindowSize(width int, height int) error {
	C.windowSetContentSize(s.id, C.intptr_t(width), C.intptr_t(height))
	return nil
}

func (s *sysData) delete(index int) {
	classTypes[s.ctype].delete(s.id, index)
}

func (s *sysData) setProgress(percent int) {
	C.setProgress(s.id, C.intptr_t(percent))
}

func (s *sysData) len() int {
	return classTypes[s.ctype].len(s.id)
}

func (s *sysData) setAreaSize(width int, height int) {
	C.setAreaSize(s.id, C.intptr_t(width), C.intptr_t(height))
}

func (s *sysData) repaintAll() {
	C.display(s.id)
}

func (s *sysData) center() {
	C.center(s.id)
}

func (s *sysData) setChecked(checked bool) {
	C.setCheckboxChecked(s.id, toBOOL(checked))
}

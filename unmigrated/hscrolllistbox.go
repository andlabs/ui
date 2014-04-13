// WINDOWS (works)

		style:			_LBS_NOTIFY | _LBS_NOINTEGRALHEIGHT | _WS_HSCROLL | _WS_VSCROLL | controlstyle,
		xstyle:			_WS_EX_CLIENTEDGE | controlxstyle,
		altStyle:			_LBS_EXTENDEDSEL | _LBS_NOTIFY | _LBS_NOINTEGRALHEIGHT | _WS_HSCROLL | _WS_VSCROLL | controlstyle,

(call recalcListboxWidth() from sysData.append(), sysData.insertBefore(), and sysData.delete())

// List Boxes do not dynamically handle horizontal scrollbars.
// We have to manually handle this ourselves.
// TODO make this run on the main thread when we switch to uitask taking function literals
// TODO this is inefficient; some caching would be nice
func recalcListboxWidth(hwnd _HWND) {
	var size _SIZE

	ret := make(chan uiret)
	defer close(ret)
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(hwnd),
			uintptr(_LB_GETCOUNT),
			uintptr(0),
			uintptr(0),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == uintptr(_LB_ERR) {		// failure
		panic(fmt.Errorf("error getting number of items for Listbox width calculations: %v", r.err))
	}
	n := int(r.ret)
	uitask <- &uimsg{
		call:		_getWindowDC,
		p:		[]uintptr{uintptr(hwnd)},
		ret:		ret,
	}
	r = <-ret
	if r.ret == 0 {		// failure
		panic(fmt.Errorf("error getting DC for Listbox width calculations: %v", r.err))
	}
	dc := _HANDLE(r.ret)
	uitask <- &uimsg{
		call:		_selectObject,
		p:		[]uintptr{
			uintptr(dc),
			uintptr(controlFont),
		},
		ret:		ret,
	}
	r = <-ret
	if r.ret == 0  {		// failure
		panic(fmt.Errorf("error loading control font into device context for Listbox width calculation: %v", r.err))
	}
	hextent := uintptr(0)
	for i := 0; i < n; i++ {
		uitask <- &uimsg{
			call:		_sendMessage,
			p:		[]uintptr{
				uintptr(hwnd),
				uintptr(_LB_GETTEXTLEN),
				uintptr(_WPARAM(i)),
				uintptr(0),
			},
			ret:		ret,
		}
		r := <-ret
		if r.ret == uintptr(_LB_ERR) {
			panic("UI library internal error: LB_ERR from LB_GETTEXTLEN in what we know is a valid listbox index (came from LB_GETSELITEMS)")
		}
		str := make([]uint16, r.ret)
		uitask <- &uimsg{
			call:		_sendMessage,
			p:		[]uintptr{
				uintptr(hwnd),
				uintptr(_LB_GETTEXT),
				uintptr(_WPARAM(i)),
				uintptr(_LPARAM(unsafe.Pointer(&str[0]))),
			},
			ret:		ret,
		}
		r = <-ret
		if r.ret == uintptr(_LB_ERR) {
			panic("UI library internal error: LB_ERR from LB_GETTEXT in what we know is a valid listbox index (came from LB_GETSELITEMS)")
		}
		// r.ret is still the length of the string; this time without the null terminator
		uitask <- &uimsg{
			call:		_getTextExtentPoint32,
			p:		[]uintptr{
				uintptr(dc),
				uintptr(unsafe.Pointer(&str[0])),
				r.ret,
				uintptr(unsafe.Pointer(&size)),
			},
			ret:		ret,
		}
		r = <-ret
		if r.ret == 0 {		// failure
			panic(fmt.Errorf("error getting width of item %d text for Listbox width calculation: %v", i, r.err))
		}
		if hextent < uintptr(size.cx) {
			hextent = uintptr(size.cx)
		}
	}
	uitask <- &uimsg{
		call:		_releaseDC,
		p:		[]uintptr{
			uintptr(hwnd),
			uintptr(dc),
		},
		ret:		ret,
	}
	r = <-ret
	if r.ret == 0 {		// failure
		panic(fmt.Errorf("error releasing DC for Listbox width calculations: %v", r.err))
	}
	uitask <- &uimsg{
		call:		_sendMessage,
		p:		[]uintptr{
			uintptr(hwnd),
			uintptr(_LB_SETHORIZONTALEXTENT),
			hextent,
			uintptr(0),
		},
		ret:		ret,
	}
	<-ret
}

// DARWIN (does not work)

// NSTableView is actually in a NSScrollView so we have to get it out first
// NSTableView and NSTableColumn both provide sizeToFit methods, but they don't do what we want (NSTableView's sizes to fit the parent; NSTableColumn's sizes to fit the column header)
// We have to get the width manually; see also http://stackoverflow.com/questions/4674163/nstablecolumn-size-to-fit-contents
// We can use the NSTableView sizeToFit to get the height, though.
// TODO this is inefficient!
// TODO move this to listbox_darwin.go
var (
	_dataCellForRow = sel_getUid("dataCellForRow:")
	_cellSize = sel_getUid("cellSize")
	_setMinWidth = sel_getUid("setMinWidth:")
	_setWidth = sel_getUid("setWidth:")
)

func listboxPrefSize(control C.id) (width int, height int) {
	var maxwidth C.int64_t

	listbox := listboxInScrollView(control)
	_, height = controlPrefSize(listbox)
	column := listboxTableColumn(listbox)
	n := C.objc_msgSend_intret_noargs(listbox, _numberOfRows)
	for i := C.intptr_t(0); i < n; i++ {
		cell := C.objc_msgSend_int(column, _dataCellForRow, i)
		csize := C.objc_msgSend_stret_size_noargs(cell, _cellSize)
		if maxwidth < csize.width {
			maxwidth = csize.width
		}
	}
	// and in order for horizontal scrolling to work, we need to set the column width to this
	C.objc_msgSend_cgfloat(column, _setMinWidth, C.double(maxwidth))
	C.objc_msgSend_cgfloat(column, _setWidth, C.double(maxwidth))
	return int(maxwidth), height
}

func (s *sysData) setRect(x int, y int, width int, height int, winheight int) error {
	// winheight - y because (0,0) is the bottom-left corner of the window and not the top-left corner
	// (winheight - y) - height because (x, y) is the bottom-left corner of the control and not the top-left
	C.objc_msgSend_rect(s.id, _setFrame,
		C.int64_t(x), C.int64_t((winheight - y) - height), C.int64_t(width), C.int64_t(height))
	// TODO having this here is a hack; split it into a separate function in listbox_darwin.go
	// the NSTableView:NSTableColumn ratio is what determines horizontal scrolling; see http://stackoverflow.com/questions/7050497/enable-scrolling-for-nstableview
	if s.ctype == c_listbox {
		listbox := listboxInScrollView(s.id)
		C.objc_msgSend_rect(listbox, _setFrame,
			0, 0, C.int64_t(width), C.int64_t(height))
	}
	return nil
}

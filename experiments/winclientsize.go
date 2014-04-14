func (s *sysData) setWindowSize(width int, height int) error {
	var rect _RECT

	ret := make(chan uiret)
	defer close(ret)
	// we need the size to cover only the client rect
	// unfortunately, Windows only provides general client->window rect conversion functions (AdjustWindowRect()/AdjustWindowRectEx()), not any that can pull from an existing window
	// so as long as we don't change the window styles during execution we're good
	// now we tell it to adjust (0,0)..(width,height); this will give us the approrpiate rect to get the proper width/height from
	// then we call SetWindowPos() to set the size without moving the window
	// see also: http://blogs.msdn.com/b/oldnewthing/archive/2003/09/11/54885.aspx
	// TODO do the WM_NCCALCSIZE stuff on that post when adding menus
	rect.Right = int32(width)
	rect.Bottom = int32(height)
	uitask <- &uimsg{
		call:		_adjustWindowRectEx,
		p:		[]uintptr{
			uintptr(unsafe.Pointer(&rect)),
			uintptr(classTypes[c_window].style),
			uintptr(_FALSE),		// TODO change this when adding menus
			uintptr(classTypes[c_window].xstyle),
		},
		ret:		ret,
	}
	r := <-ret
	if r.ret == 0 {		// failure
		panic(fmt.Errorf("error computing window rect from client rect for resize: %v", r.err))
	}
	uitask <- &uimsg{
		call:		_setWindowPos,
		p:		[]uintptr{
			uintptr(s.hwnd),
			uintptr(_NULL),			// not changing the Z-order
			uintptr(0),
			uintptr(0),
			uintptr(rect.Right - rect.Left),
			uintptr(rect.Bottom - rect.Top),
			uintptr(_SWP_NOMOVE | _SWP_NOZORDER),
		},
		ret:		ret,
	}
	r = <-ret
	if r.ret == 0 {		// failure
		return fmt.Errorf("error actually resizing window: %v", r.err)
	}
	return nil
}

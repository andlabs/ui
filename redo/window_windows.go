// 12 july 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
)

type window struct {
	hwnd		uintptr
	shownbefore	bool

	child			Control

	closing		*event
}

const windowclassname = "gouiwindow"
var windowclassptr = syscall.StringToUTF16Ptr(windowclassname)

func makeWindowWindowClass() error {
	var wc s_WNDCLASSW

	wc.lpfnWndProc = syscall.NewCallback(windowWndProc)
	wc.hInstance = hInstance
	wc.hIcon = hDefaultIcon
	wc.hCursor = hArrowCursor
	wc.hbrBackground = c_COLOR_BTNFACE + 1
	wc.lpszClassName = windowclassptr
	res, err := f_RegisterClassW(&wc)
	if res == 0 {
		return fmt.Errorf("error registering Window window class: %v", err)
	}
	return nil
}

func newWindow(title string, width int, height int) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			w := &window{
				// hwnd set in WM_CREATE handler
				closing:	newEvent(),
			}
			hwnd, err := f_CreateWindowExW(
				0,
				windowclassptr,
				syscall.StringToUTF16Ptr(title),
				c_WS_OVERLAPPEDWINDOW,
				c_CW_USEDEFAULT, c_CW_USEDEFAULT,
				uintptr(width), uintptr(height),
				hNULL, hNULL, hInstance, unsafe.Pointer(w))
			if hwnd == hNULL {
				panic(fmt.Errorf("Window creation failed: %v", err))
			} else if hwnd != w.hwnd {
				panic(fmt.Errorf("inconsistency: hwnd returned by CreateWindowEx() (%p) and hwnd stored in window (%p) differ", hwnd, w.hwnd))
			}
			c <- w
		},
		resp:		c,
	}
}

func (w *window) SetControl(control Control) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			if w.child != nil {		// unparent existing control
				w.child.unparent()
			}
			control.unparent()
			control.parent(w)
			w.child = control
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Title() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			c <- getWindowText(w.hwnd)
		},
		resp:		c,
	}
}

func (w *window) SetTitle(title string) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			setWindowText(w.hwnd, title, []t_LRESULT{c_FALSE})
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Show() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			if !w.shownbefore {
				// TODO get rid of need for cast
				f_ShowWindow(w.hwnd, uintptr(nCmdShow))
				updateWindow(w.hwnd, "Window.Show()")
				w.shownbefore = true
			} else {
				f_ShowWindow(w.hwnd, c_SW_SHOW)
			}
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) Hide() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			f_ShowWindow(w.hwnd, c_SW_HIDE)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func doclose(w *window) {
	res, err := f_DestroyWindow(w.hwnd)
	if res == 0 {
		panic(fmt.Errorf("error destroying window: %v", err))
	}
}

func (w *window) Close() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			doclose(w)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (w *window) OnClosing(e func(Doer) bool) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			w.closing.setbool(e)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func windowWndProc(hwnd uintptr, msg t_UINT, wParam t_WPARAM, lParam t_LPARAM) t_LRESULT {
	w := (*window)(unsafe.Pointer(f_GetWindowLongPtrW(hwnd, c_GWLP_USERDATA)))
	if w == nil {
		// the lpParam is available during WM_NCCREATE and WM_CREATE
		if msg == c_WM_NCCREATE {
			storelpParam(hwnd, lParam)
			w := (*window)(unsafe.Pointer(f_GetWindowLongPtrW(hwnd, c_GWLP_USERDATA)))
			w.hwnd = hwnd
		}
		// act as if we're not ready yet, even during WM_NCCREATE (nothing important to the switch statement below happens here anyway)
		return f_DefWindowProcW(hwnd, msg, wParam, lParam)
	}
	switch msg {
	case c_WM_SIZE:
		var r s_RECT

		res, err := f_GetClientRect(w.hwnd, &r)
		if res == 0 {
			panic(fmt.Errorf("error getting client rect for Window in WM_SIZE: %v", err))
		}
		fmt.Printf("new size %d x %d\n", r.right - r.left, r.bottom - r.top)
		return 0
	case c_WM_CLOSE:
		close := w.closing.fire()
		if close {
			doclose(w)
		}
		return 0
	default:
		return f_DefWindowProcW(hwnd, msg, wParam, lParam)
	}
	panic(fmt.Errorf("Window message %d does not return a value (bug in windowWndProc())", msg))
}

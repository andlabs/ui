// 15 july 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
)

type widgetbase struct {
	hwnd	uintptr
}

var emptystr = syscall.StringToUTF16Ptr("")

func newWidget(class *uint16, style uintptr, extstyle uintptr) *widgetbase {
	hwnd, err := f_CreateWindowExW(
		extstyle,
		class, emptystr,
		style | c_WS_CHILD | c_WS_VISIBLE,
		c_CW_USEDEFAULT, c_CW_USEDEFAULT,
//		c_CW_USEDEFAULT, c_CW_USEDEFAULT,
100,100,
		// the following has the consequence of making the control message-only at first
		// this shouldn't cause any problems... hopefully not
		// but see the msgwndproc() for caveat info
		// also don't use low control IDs as they will conflict with dialog boxes (IDCANCEL, etc.)
		msgwin, 100, hInstance, nil)
	if hwnd == hNULL {
		panic(fmt.Errorf("creating control of class %q failed: %v", class, err))
	}
	return &widgetbase{
		hwnd:	hwnd,
	}
}

// these few methods are embedded by all the various Controls since they all will do the same thing

func (w *widgetbase) unparent() {
	res, err := f_SetParent(w.hwnd, msgwin)
	if res == hNULL {		// result type is HWND
		panic(fmt.Errorf("error unparenting control: %v", err))
	}
}

func (w *widgetbase) parent(win *window) {
	res, err := f_SetParent(w.hwnd, win.hwnd)
	if res == hNULL {		// result type is HWND
		panic(fmt.Errorf("error parenting control: %v", err))
	}
}

// don't embed these as exported; let each Control decide if it should

func (w *widgetbase) text() *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			c <- getWindowText(w.hwnd)
		},
		resp:		c,
	}
}

func (w *widgetbase) settext(text string, results ...t_LRESULT) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			setWindowText(w.hwnd, text, append([]t_LRESULT{c_FALSE}, results...))
			c <- struct{}{}
		},
		resp:		c,
	}
}

// all controls that have events receive the events themselves through subclasses
// to do this, all windows (including the message-only window; see http://support.microsoft.com/default.aspx?scid=KB;EN-US;Q104069) forward WM_COMMAND to each control with this function
func forwardCommand(hwnd uintptr, uMsg t_UINT, wParam t_WPARAM, lParam t_LPARAM) t_LRESULT {
	control := uintptr(lParam)
	// don't generate an event if the control (if there is one) is unparented (a child of the message-only window)
	if control != hNULL && f_IsChild(msgwin, control) == 0 {
		return f_SendMessageW(control, msgCOMMAND, wParam, lParam)
	}
	return f_DefWindowProcW(hwnd, uMsg, wParam, lParam)
}

type button struct {
	*widgetbase
	clicked		*event
}

var buttonclass = syscall.StringToUTF16Ptr("BUTTON")

func newButton(text string) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			w := newWidget(buttonclass,
				c_BS_PUSHBUTTON | c_WS_TABSTOP,
				0)
			setWindowText(w.hwnd, text, []t_LRESULT{c_FALSE})
			c <- &button{
				widgetbase:	w,
				clicked:		newEvent(),
			}
		},
		resp:		c,
	}
}

func (b *button) OnClicked(e func(c Doer)) *Request {
	c := make(chan interface{})
	return &Request{
		op:		func() {
			b.clicked.set(e)
			c <- struct{}{}
		},
		resp:		c,
	}
}

func (b *button) Text() *Request {
	return b.text()
}

func (b *button) SetText(text string) *Request {
	return b.settext(text)
}

var buttonsubprocptr = syscall.NewCallback(buttonSubProc)

func buttonSubProc(hwnd uintptr, uMsg t_UINT, wParam t_WPARAM, lParam t_LPARAM, id t_UINT_PTR, data t_DWORD_PTR) t_LRESULT {
	b := (*button)(unsafe.Pointer(uintptr(data)))
	switch uMsg {
	case msgCOMMAND:
		if wParam.HIWORD() == c_BN_CLICKED {
			b.clicked.fire()
			println("button clicked")
			return 0
		}
		// TODO return
	case c_WM_NCDESTROY:
		// TODO remove
		// TODO return
	default:
		// TODO return
	}
	panic(fmt.Errorf("Button message %d does not return a value (bug in buttonSubProc())", uMsg))
}

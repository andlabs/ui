// 11 july 2014

package ui

// wfunc kernel32 GetModuleHandleW *uint16 uintptr
// wfunc kernel32 GetStartupInfoW *s_STARTUPINFOW void
// wfunc user32 LoadIconW uintptr uintptr uintptr
// wfunc user32 LoadCursorW uintptr uintptr uintptr
// wfunc user32 GetMessageW *s_MSG uintptr t_UINT t_UINT t_BOOL

// these two don't technically return void but let's pretend they do since their return values are irrelevant/not indicative of anything useful
// wfunc user32 TranslateMessage *s_MSG void
// wfunc user32 DispatchMessageW *s_MSG void

// wfunc user32 PostMessageW uintptr t_UINT t_WPARAM t_LPARAM uintptr
// wfunc user32 RegisterClassW *s_WNDCLASSW uintptr

// TODO narrow down argument types
// wfunc user32 CreateWindowExW uintptr *uint16 *uint16 uintptr uintptr uintptr uintptr uintptr uintptr uintptr uintptr unsafe.Pointer uintptr

// wfunc user32 DefWindowProcW uintptr t_UINT t_WPARAM t_LPARAM t_LRESULT,noerr

// this one doesn't technically return void but let's pretend it does since its return value is irrelevant/not indicative of anything useful
// wfunc user32 ShowWindow uintptr uintptr void

// wfunc user32 SendMessageW uintptr t_UINT t_WPARAM t_LPARAM t_LRESULT,noerr
// wfunc user32 UpdateWindow uintptr uintptr
// wfunc user32 DestroyWindow uintptr uintptr

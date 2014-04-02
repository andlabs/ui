// 25 february 2014

package ui

import (
	"fmt"
	"syscall"
	"unsafe"
	"io/ioutil"
)

// pretty much every constant here except _WM_USER is from commctrl.h
// TODO for all: filter out constants not available in Windows XP

var (
	// TODO deinitialize at program end?
	comctlManifestCookie uintptr
)

// InitCommonControlsEx constants.
const (
	_ICC_LISTVIEW_CLASSES = 0x00000001
	_ICC_TREEVIEW_CLASSES = 0x00000002
	_ICC_BAR_CLASSES = 0x00000004
	_ICC_TAB_CLASSES = 0x00000008
	_ICC_UPDOWN_CLASS = 0x00000010
	_ICC_PROGRESS_CLASS = 0x00000020
	_ICC_HOTKEY_CLASS = 0x00000040
	_ICC_ANIMATE_CLASS = 0x00000080
	_ICC_WIN95_CLASSES = 0x000000FF
	_ICC_DATE_CLASSES = 0x00000100
	_ICC_USEREX_CLASSES = 0x00000200
	_ICC_COOL_CLASSES = 0x00000400
	_ICC_INTERNET_CLASSES = 0x00000800
	_ICC_PAGESCROLLER_CLASS = 0x00001000
	_ICC_NATIVEFNTCTL_CLASS = 0x00002000
	_ICC_STANDARD_CLASSES = 0x00004000
	_ICC_LINK_CLASS = 0x00008000
)

var (
	_activateActCtx = kernel32.NewProc("ActivateActCtx")
	_createActCtx = kernel32.NewProc("CreateActCtxW")
)

/*
Windows requires a manifest file to enable Common Controls version 6.
The only way to not require an external manifest is to synthesize the manifest ourselves.
We can use the activation context API to load it at runtime.
References:
- http://stackoverflow.com/questions/4308503/how-to-enable-visual-styles-without-a-manifest
- http://support.microsoft.com/kb/830033
*/
func forceCommonControls6() (err error) {
	var (
		// from winbase.h; var because Go won't let me convert this constant
		_INVALID_HANDLE_VALUE = -1
	)

	manifestfile, err := ioutil.TempFile("", "gouicomctl32v6manifest")
	if err != nil {
		return fmt.Errorf("error creating synthesized manifest file: %v", err)
	}
	_, err = manifestfile.Write(manifest)
	if err != nil {
		return fmt.Errorf("error writing synthesized manifest file: %v", err)
	}
	filename := manifestfile.Name()
	// we now have to close the file, otherwise ActivateActCtx() will complain that it's in use by another program
	// if ioutil.TempFile() ever changes so that the file is deleted when it is closed, this will need to change
	manifestfile.Close()

	var actctx struct {
		cbSize				uint32
		dwFlags				uint32
		lpSource				*uint16
		wProcessorArchitecture	uint16
		wLangId				uint16		// originally LANGID
		lpAssemblyDirectory	uintptr		// originally LPCWSTR
		lpResourceName		uintptr		// originally LPCWSTR
		lpApplicationName		uintptr		// originally LPCWSTR
		hModule				_HANDLE		// originally HMODULE
	}

	actctx.cbSize = uint32(unsafe.Sizeof(actctx))
	// TODO set ACTCTX_FLAG_SET_PROCESS_DEFAULT? I can't find a reference to figure out what this means
	actctx.lpSource = syscall.StringToUTF16Ptr(filename)

	r1, _, err := _createActCtx.Call(uintptr(unsafe.Pointer(&actctx)))
	if r1 == uintptr(_INVALID_HANDLE_VALUE) {		// failure
		return fmt.Errorf("error creating activation context for synthesized manifest file: %v", err)
	}
	r1, _, err = _activateActCtx.Call(
		r1,
		uintptr(unsafe.Pointer(&comctlManifestCookie)))
	if r1 == uintptr(_FALSE) {		// failure
		return fmt.Errorf("error activating activation context for synthesized manifest file: %v", err)
	}
	return nil
}

func initCommonControls() (err error) {
	var icc struct {
		dwSize	uint32
		dwICC	uint32
	}

	err = forceCommonControls6()
	if err != nil {
		return fmt.Errorf("error forcing Common Controls version 6 (or newer): %v", err)
	}

	icc.dwSize = uint32(unsafe.Sizeof(icc))
	icc.dwICC = _ICC_PROGRESS_CLASS

	comctl32 = syscall.NewLazyDLL("comctl32.dll")
	r1, _, err := comctl32.NewProc("InitCommonControlsEx").Call(uintptr(unsafe.Pointer(&icc)))
	if r1 == _FALSE {		// failure
		// TODO does it set GetLastError()?
		return fmt.Errorf("error initializing Common Controls (comctl32.dll): %v", err)
	}
	return nil
}

// Common Controls class names.
const (
	_PROGRESS_CLASS = "msctls_progress32"
)

// Shared Common Controls styles.
const (
	_WM_USER = 0x0400
	_CCM_FIRST = 0x2000
	_CCM_SETBKCOLOR = (_CCM_FIRST + 1)
)

// Progress Bar styles.
const (
	_PBS_SMOOTH = 0x01
	_PBS_VERTICAL = 0x04
	_PBS_MARQUEE = 0x08
)

// Progress Bar messages.
const (
	_PBM_SETRANGE = (_WM_USER + 1)
	_PBM_SETPOS = (_WM_USER + 2)
	_PBM_DELTAPOS = (_WM_USER + 3)
	_PBM_SETSTEP = (_WM_USER + 4)
	_PBM_STEPIT = (_WM_USER + 5)
	_PBM_SETRANGE32 = (_WM_USER + 6)
	_PBM_GETRANGE = (_WM_USER + 7)
	_PBM_GETPOS = (_WM_USER + 8)
	_PBM_SETBARCOLOR = (_WM_USER + 9)
	_PBM_SETBKCOLOR = _CCM_SETBKCOLOR
	_PBM_SETMARQUEE = (_WM_USER + 10)
)

var manifest = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
<assemblyIdentity
    version="1.0.0.0"
    processorArchitecture="*"
    name="CompanyName.ProductName.YourApplication"
    type="win32"
/>
<description>Your application description here.</description>
<dependency>
    <dependentAssembly>
        <assemblyIdentity
            type="win32"
            name="Microsoft.Windows.Common-Controls"
            version="6.0.0.0"
            processorArchitecture="*"
            publicKeyToken="6595b64144ccf1df"
            language="*"
        />
    </dependentAssembly>
</dependency>
</assembly>
`)

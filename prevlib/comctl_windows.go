// 25 february 2014

package ui

import (
	"fmt"
	"io/ioutil"
	"syscall"
	"unsafe"
)

// pretty much every constant here except _WM_USER is from commctrl.h, except where noted

var (
	comctlManifestCookie uintptr
)

var (
	_activateActCtx = kernel32.NewProc("ActivateActCtx")
	_createActCtx   = kernel32.NewProc("CreateActCtxW")
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
		cbSize                 uint32
		dwFlags                uint32
		lpSource               *uint16
		wProcessorArchitecture uint16
		wLangId                uint16  // originally LANGID
		lpAssemblyDirectory    uintptr // originally LPCWSTR
		lpResourceName         uintptr // originally LPCWSTR
		lpApplicationName      uintptr // originally LPCWSTR
		hModule                _HANDLE // originally HMODULE
	}

	actctx.cbSize = uint32(unsafe.Sizeof(actctx))
	// make this context the process default, just to be safe
	actctx.dwFlags = _ACTCTX_FLAG_SET_PROCESS_DEFAULT
	actctx.lpSource = toUTF16(filename)

	r1, _, err := _createActCtx.Call(uintptr(unsafe.Pointer(&actctx)))
	// don't negConst() INVALID_HANDLE_VALUE; windowsconstgen was given a pointer by windows.h, and pointers are unsigned, so converting it back to signed doesn't work
	if r1 == _INVALID_HANDLE_VALUE { // failure
		return fmt.Errorf("error creating activation context for synthesized manifest file: %v", err)
	}
	r1, _, err = _activateActCtx.Call(
		r1,
		uintptr(unsafe.Pointer(&comctlManifestCookie)))
	if r1 == uintptr(_FALSE) { // failure
		return fmt.Errorf("error activating activation context for synthesized manifest file: %v", err)
	}
	return nil
}

func initCommonControls() (err error) {
	var icc struct {
		dwSize uint32
		dwICC  uint32
	}

	err = forceCommonControls6()
	if err != nil {
		return fmt.Errorf("error forcing Common Controls version 6 (or newer): %v", err)
	}

	icc.dwSize = uint32(unsafe.Sizeof(icc))
	icc.dwICC = _ICC_PROGRESS_CLASS

	comctl32 = syscall.NewLazyDLL("comctl32.dll")
	r1, _, err := comctl32.NewProc("InitCommonControlsEx").Call(uintptr(unsafe.Pointer(&icc)))
	if r1 == _FALSE { // failure
		return fmt.Errorf("error initializing Common Controls (comctl32.dll); Windows last error: %v", err)
	}
	return nil
}

// Common Controls class names.
const (
	// x (lowercase) prefix to avoid being caught by the constants generator
	x_PROGRESS_CLASS = "msctls_progress32"
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

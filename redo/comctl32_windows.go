// 25 february 2014

package ui

import (
	"fmt"
	"io/ioutil"
	"syscall"
)

// #include "winapi_windows.h"
import "C"

// TODO possibly rewrite the whole file access bits in C

// pretty much every constant here except _WM_USER is from commctrl.h, except where noted

/*
Windows requires a manifest file to enable Common Controls version 6.
The only way to not require an external manifest is to synthesize the manifest ourselves.
We can use the activation context API to load it at runtime.
References:
- http://stackoverflow.com/questions/4308503/how-to-enable-visual-styles-without-a-manifest
- http://support.microsoft.com/kb/830033
The activation context code itself is in comctl32_windows.c.
*/
func initCommonControls() (err error) {
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

	var errmsg *C.char

	errcode := C.initCommonControls(toUTF16(filename), &errmsg)
	if errcode != 0 || errmsg != nil {
		return fmt.Errorf("error actually initializing comctl32.dll: %s: %v", C.GoString(errmsg), syscall.Errno(errcode))
	}
	return nil
}

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

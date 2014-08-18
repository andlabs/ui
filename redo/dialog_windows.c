// 18 august 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

// this should be reasonable
#define NFILENAME 4096

WCHAR *openFile(void)
{
	OPENFILENAMEW ofn;
	DWORD err;
	WCHAR *filenameBuffer;

	// freed on the Go side
	filenameBuffer = (WCHAR *) malloc((NFILENAME + 1) * sizeof (WCHAR));
	if (filenameBuffer == NULL)
		xpanic("memory exhausted in OpenFile()", GetLastError());
	filenameBuffer[0] = L'\0';			// required by GetOpenFileName() to indicate no previous filename
	ZeroMemory(&ofn, sizeof (OPENFILENAMEW));
	ofn.lStructSize = sizeof (OPENFILENAMEW);
	ofn.hwndOwner = NULL;
	ofn.hInstance = hInstance;
	ofn.lpstrFilter = NULL;			// no filters
	ofn.lpstrFile = filenameBuffer;
	ofn.nMaxFile = NFILENAME + 1;	// TODO include + 1?
	ofn.lpstrInitialDir = NULL;			// let system decide
	ofn.lpstrTitle = NULL;			// let system decide
	// TODO OFN_SHAREAWARE?
	// TODO remove OFN_NODEREFERENCELINKS? or does no filters ensure that anyway?
	ofn.Flags = OFN_EXPLORER | OFN_FILEMUSTEXIST | OFN_FORCESHOWHIDDEN | OFN_HIDEREADONLY | OFN_LONGNAMES | OFN_NOCHANGEDIR | OFN_NODEREFERENCELINKS | OFN_NOTESTFILECREATE | OFN_PATHMUSTEXIST;
	if (GetOpenFileNameW(&ofn) == FALSE) {
		// TODO stringify
		err = CommDlgExtendedError();
		if (err == 0)				// user cancelled
			return NULL;
		xpanic("error running open file dialog", GetLastError());
	}
	return filenameBuffer;
}

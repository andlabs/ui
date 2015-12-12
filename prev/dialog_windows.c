// 18 august 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

// this should be reasonable
#define NFILENAME 4096

struct openFileData {
	HWND parent;
	void *f;
	WCHAR *filenameBuffer;
};

static DWORD WINAPI doOpenFile(LPVOID data)
{
	struct openFileData *o = (struct openFileData *) data;
	OPENFILENAMEW ofn;
	DWORD err;

	o->filenameBuffer[0] = L'\0';			// required by GetOpenFileName() to indicate no previous filename
	ZeroMemory(&ofn, sizeof (OPENFILENAMEW));
	ofn.lStructSize = sizeof (OPENFILENAMEW);
	ofn.hwndOwner = o->parent;
	ofn.hInstance = hInstance;
	ofn.lpstrFilter = NULL;			// no filters
	ofn.lpstrFile = o->filenameBuffer;
	ofn.nMaxFile = NFILENAME + 1;	// seems to include null terminator according to docs
	ofn.lpstrInitialDir = NULL;			// let system decide
	ofn.lpstrTitle = NULL;			// let system decide
	// TODO OFN_SHAREAWARE?
	// better question: TODO keep networking?
	ofn.Flags = OFN_EXPLORER | OFN_FILEMUSTEXIST | OFN_FORCESHOWHIDDEN | OFN_HIDEREADONLY | OFN_LONGNAMES | OFN_NOCHANGEDIR | OFN_NODEREFERENCELINKS | OFN_NOTESTFILECREATE | OFN_PATHMUSTEXIST;
	if (GetOpenFileNameW(&ofn) == FALSE) {
		err = CommDlgExtendedError();
		if (err != 0)				// user cancelled
			xpaniccomdlg("error running open file dialog", err);
		free(o->filenameBuffer);		// free now so we can set it to NULL without leaking
		o->filenameBuffer = NULL;
	}
	if (PostMessageW(msgwin, msgOpenFileDone, (WPARAM) (o->filenameBuffer), (LPARAM) (o->f)) == 0)
		xpanic("error posting OpenFile() finished message to message-only window", GetLastError());
	free(o);		// won't free o->f or o->filenameBuffer in above invocation
	return 0;
}

void openFile(HWND hwnd, void *f)
{
	struct openFileData *o;

	// freed by the thread
	o = (struct openFileData *) malloc(sizeof (struct openFileData));
	if (o == NULL)
		xpanic("memory exhausted allocating data structure in OpenFile()", GetLastError());
	o->parent = hwnd;
	o->f = f;
	// freed on the Go side
	o->filenameBuffer = (WCHAR *) malloc((NFILENAME + 1) * sizeof (WCHAR));
	if (o->filenameBuffer == NULL)
		xpanic("memory exhausted allocating filename buffer in OpenFile()", GetLastError());
	if (CreateThread(NULL, 0, doOpenFile, (LPVOID) o, 0, NULL) == NULL)
		xpanic("error creating thread for running OpenFIle()", GetLastError());
}

// 18 august 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

static LRESULT CALLBACK dialogSubProc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam, UINT_PTR id, DWORD_PTR data)
{
	switch (uMsg) {
	case WM_COMMAND:
		// we must re-enable other windows in the right order (see http://blogs.msdn.com/b/oldnewthing/archive/2004/02/27/81155.aspx)
		// see http://stackoverflow.com/questions/25494914/is-there-something-like-cdn-filecancel-analogous-to-cdn-fileok-for-getting-when
		if (HIWORD(wParam) == BN_CLICKED && LOWORD(wParam) == IDCANCEL)
				SendMessageW(msgwin, msgEndModal, 0, 0);
		break;		// let the dialog handle it now
	case WM_NCDESTROY:
		if ((*fv_RemoveWindowSubclass)(hwnd, dialogSubProc, id) == FALSE)
			xpanic("error removing dialog subclass (which was for its own event handler)", GetLastError());
	}
	return (*fv_DefSubclassProc)(hwnd, uMsg, wParam, lParam);
}

// this should be reasonable
#define NFILENAME 4096

static UINT_PTR CALLBACK openSaveFileHook(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	if (uMsg == WM_INITDIALOG) {
		HWND parent;

		parent = GetParent(hwnd);
		if (parent == NULL)
			xpanic("error gettign parent of OpenFile() dialog for event handling", GetLastError());
		if ((*fv_SetWindowSubclass)(parent, dialogSubProc, 0, (DWORD_PTR) NULL) == FALSE)
			xpanic("error subclassing OpenFile() dialog to give it its own event handler", GetLastError());
	} else if (uMsg == WM_NOTIFY) {
		OFNOTIFY *of = (OFNOTIFY *) lParam;

		if (of->hdr.code == CDN_FILEOK)
			SendMessageW(msgwin, msgEndModal, 0, 0);
	}
	return 0;
}

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
	ofn.Flags = OFN_ENABLEHOOK | OFN_EXPLORER | OFN_FILEMUSTEXIST | OFN_FORCESHOWHIDDEN | OFN_HIDEREADONLY | OFN_LONGNAMES | OFN_NOCHANGEDIR | OFN_NODEREFERENCELINKS | OFN_NOTESTFILECREATE | OFN_PATHMUSTEXIST;
	ofn.lpfnHook = openSaveFileHook;
	SendMessageW(msgwin, msgBeginModal, 0, 0);
	if (GetOpenFileNameW(&ofn) == FALSE) {
		err = CommDlgExtendedError();
		if (err == 0)				// user cancelled
			return NULL;
		xpaniccomdlg("error running open file dialog", err);
	}
	return filenameBuffer;
}

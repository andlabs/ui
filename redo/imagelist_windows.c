// 16 august 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

// TODO eliminate duplicate code

HBITMAP unscaledBitmap(void *i, intptr_t dx, intptr_t dy)
{
	BITMAPINFO bi;
	VOID *ppvBits;
	HBITMAP bitmap;

	ZeroMemory(&bi, sizeof (BITMAPINFO));
	bi.bmiHeader.biSize = sizeof (BITMAPINFOHEADER);
	bi.bmiHeader.biWidth = (LONG) dx;
	bi.bmiHeader.biHeight = -((LONG) dy);			// negative height to force top-down drawing;
	bi.bmiHeader.biPlanes = 1;
	bi.bmiHeader.biBitCount = 32;
	bi.bmiHeader.biCompression = BI_RGB;
	bi.bmiHeader.biSizeImage = (DWORD) (dx * dy * 4);
	bitmap = CreateDIBSection(NULL, &bi, DIB_RGB_COLORS, &ppvBits, 0, 0);
	if (bitmap == NULL)
		xpanic("error creating HBITMAP for unscaled ImageList image copy", GetLastError());
	dotoARGB(i, (void *) ppvBits);
	return bitmap;
}

HIMAGELIST newImageList(int width, int height)
{
	HIMAGELIST il;

	// TODO does this strip alpha?
	// sinni800 in irc.freenode.net/#go-nuts suggests our use of *image.RGBA makes this not so much of an issue
	il = (*fv_ImageList_Create)(width, height, ILC_COLOR32, 20, 20);		// should be reasonable
	if (il == NULL)
		xpanic("error creating image list", GetLastError());
	return il;
}

void addImage(HIMAGELIST il, HWND hwnd, HBITMAP bitmap, int origwid, int oright, int width, int height)
{
	BOOL wasScaled = FALSE;
	HDC winDC, scaledDC, origDC;
	HBITMAP scaled;
	HBITMAP prevscaled, prevorig;

	// first we need to scale the bitmap
	if (origwid == width && oright == height) {
		scaled = bitmap;
		goto noscale;
	}
	wasScaled = TRUE;
	winDC = GetDC(hwnd);
	if (winDC == NULL)
		xpanic("error getting DC for window", GetLastError());
	origDC = CreateCompatibleDC(winDC);
	if (origDC == NULL)
		xpanic("error getting DC for original ImageList bitmap", GetLastError());
	prevorig = SelectObject(origDC, bitmap);
	if (prevorig == NULL)
		xpanic("error selecting original ImageList bitmap into DC", GetLastError());
	scaledDC = CreateCompatibleDC(origDC);
	if (scaledDC == NULL)
		xpanic("error getting DC for scaled ImageList bitmap", GetLastError());
	scaled = CreateCompatibleBitmap(origDC, width, height);
	if (scaled == NULL)
		xpanic("error creating scaled ImageList bitmap", GetLastError());
	prevscaled = SelectObject(scaledDC, scaled);
	if (prevscaled == NULL)
		xpanic("error selecting scaled ImageList bitmap into DC", GetLastError());
	if (SetStretchBltMode(scaledDC, COLORONCOLOR) == 0)
		xpanic("error setting scaling mode", GetLastError());
	if (StretchBlt(scaledDC, 0, 0, width, height,
		origDC, 0, 0, origwid, oright,
		SRCCOPY) == 0)
		xpanic("error scaling ImageList bitmap down", GetLastError());
	if (SelectObject(origDC, prevorig) != bitmap)
		xpanic("error selecting previous bitmap into original image's DC", GetLastError());
	if (DeleteDC(origDC) == 0)
		xpanic("error deleting original image's DC", GetLastError());
	if (SelectObject(scaledDC, prevscaled) != scaled)
		xpanic("error selecting previous bitmap into scaled image's DC", GetLastError());
	if (DeleteDC(scaledDC) == 0)
		xpanic("error deleting scaled image's DC", GetLastError());
	if (ReleaseDC(hwnd, winDC) == 0)
		xpanic("error deleting window DC", GetLastError());

noscale:
	if ((*fv_ImageList_Add)(il, scaled, NULL) == -1)
		xpanic("error adding ImageList image to image list", GetLastError());
	if (wasScaled)		// clean up
		if (DeleteObject(scaled) == 0)
			xpanic("error deleting scaled bitmap", GetLastError());
}

void applyImageList(HWND hwnd, UINT uMsg, WPARAM wParam, HIMAGELIST il)
{
	if (SendMessageW(hwnd, uMsg, wParam, (LPARAM) il) == (LRESULT) NULL)
;//TODO		xpanic("error setting image list", GetLastError());
	// TODO free old one here if any
}

static UINT dfcState(int cbstate)
{
	UINT ret;

	ret = DFCS_BUTTONCHECK;
	if ((cbstate & checkboxStateChecked) != 0)
		ret |= DFCS_CHECKED;
	if ((cbstate & checkboxStateHot) != 0)
		ret |= DFCS_HOT;
	if ((cbstate & checkboxStatePushed) != 0)
		ret |= DFCS_PUSHED;
	return ret;
}

static HBITMAP dfcImage(HDC dc, int width, int height, int cbState)
{
	BITMAPINFO bi;
	VOID *ppvBits;
	HBITMAP bitmap;
	RECT r;
	HDC drawDC;
	HBITMAP prevbitmap;

	r.left = 0;
	r.top = 0;
	r.right = width;
	r.bottom = height;
	ZeroMemory(&bi, sizeof (BITMAPINFO));
	bi.bmiHeader.biSize = sizeof (BITMAPINFOHEADER);
	bi.bmiHeader.biWidth = (LONG) width;
	bi.bmiHeader.biHeight = -((LONG) height);			// negative height to force top-down drawing;
	bi.bmiHeader.biPlanes = 1;
	bi.bmiHeader.biBitCount = 32;
	bi.bmiHeader.biCompression = BI_RGB;
	bi.bmiHeader.biSizeImage = (DWORD) (width * height * 4);
	bitmap = CreateDIBSection(NULL, &bi, DIB_RGB_COLORS, &ppvBits, 0, 0);
	if (bitmap == NULL)
		xpanic("error creating HBITMAP for unscaled ImageList image copy", GetLastError());

	drawDC = CreateCompatibleDC(dc);
	if (drawDC == NULL)
		xpanic("error getting DC for checkbox image list bitmap", GetLastError());
	prevbitmap = SelectObject(drawDC, bitmap);
	if (prevbitmap == NULL)
		xpanic("error selecting checkbox image list bitmap into DC", GetLastError());
	if (DrawFrameControl(drawDC, &r, DFC_BUTTON, dfcState(cbState)) == 0)
		xpanic("error drawing checkbox image", GetLastError());
	if (SelectObject(drawDC, prevbitmap) != bitmap)
		xpanic("error selecting previous bitmap into checkbox image's DC", GetLastError());
	if (DeleteDC(drawDC) == 0)
		xpanic("error deleting checkbox image's DC", GetLastError());

	return bitmap;
}

HIMAGELIST checkboxImageList = NULL;

void makeCheckboxImageList_DFC(HWND hwnddc)
{
	int width, height;
	int cbState;
	HDC dc;
	HIMAGELIST il;

	// there's no real metric around
	// let's use SM_CX/YSMICON and hope for the best
	width = GetSystemMetrics(SM_CXSMICON);
	height = GetSystemMetrics(SM_CYSMICON);
	dc = GetDC(hwnddc);
	if (dc == NULL)
		xpanic("error getting DC for making the checkbox image list", GetLastError());
	checkboxImageList = (*fv_ImageList_Create)(width, height, ILC_COLOR32, 20, 20);		// should be reasonable
	if (checkboxImageList == NULL)
		xpanic("error creating checkbox image list", GetLastError());
	for (cbState = 0; cbState < checkboxnStates; cbState++) {
		HBITMAP bitmap;

		bitmap = dfcImage(dc, width, height, cbState);
		if ((*fv_ImageList_Add)(checkboxImageList, bitmap, NULL) == -1)
			xpanic("error adding checkbox image to image list", GetLastError());
		if (DeleteObject(bitmap) == 0)
			xpanic("error deleting checkbox bitmap", GetLastError());
	}
	if (ReleaseDC(hwnddc, dc) == 0)
		xpanic("error deleting checkbox image list DC", GetLastError());
}

static HTHEME theme = NULL;

static void openTheme(HWND hwnd)
{
	if (theme != NULL)
		// TODO save HRESULT
		if (CloseThemeData(theme) != S_OK)
			xpanic("error closing theme", GetLastError());
	// ignore error; if it can't be done, we can fall back to DrawFrameControl()
	theme = OpenThemeData(hwnd, L"button");
}

static int themestates[checkboxnStates] = {
	CBS_UNCHECKEDNORMAL,			// 0
	CBS_CHECKEDNORMAL,				// checked
	CBS_UNCHECKEDHOT,				// hot
	CBS_CHECKEDHOT,					// checked | hot
	CBS_UNCHECKEDPRESSED,			// pushed
	CBS_CHECKEDPRESSED,				// checked | pushed
	CBS_UNCHECKEDPRESSED,			// hot | pushed
	CBS_CHECKEDPRESSED,				// checked | hot | pushed
};

static SIZE getStateSize(HDC dc, int cbState)
{
	SIZE s;

	// TODO use HRESULT
	if (GetThemePartSize(theme, dc, BP_CHECKBOX, themestates[cbState], NULL, TS_DRAW, &s) != S_OK)
		xpanic("error getting theme part size", GetLastError());
	return s;
}

static HBITMAP drawThemeImage(HDC dc, int width, int height, int cbState)
{
	BITMAPINFO bi;
	VOID *ppvBits;
	HBITMAP bitmap;
	RECT r;
	HDC drawDC;
	HBITMAP prevbitmap;

	r.left = 0;
	r.top = 0;
	r.right = width;
	r.bottom = height;
	ZeroMemory(&bi, sizeof (BITMAPINFO));
	bi.bmiHeader.biSize = sizeof (BITMAPINFOHEADER);
	bi.bmiHeader.biWidth = (LONG) width;
	bi.bmiHeader.biHeight = -((LONG) height);			// negative height to force top-down drawing;
	bi.bmiHeader.biPlanes = 1;
	bi.bmiHeader.biBitCount = 32;
	bi.bmiHeader.biCompression = BI_RGB;
	bi.bmiHeader.biSizeImage = (DWORD) (width * height * 4);
	bitmap = CreateDIBSection(NULL, &bi, DIB_RGB_COLORS, &ppvBits, 0, 0);
	if (bitmap == NULL)
		xpanic("error creating HBITMAP for unscaled ImageList image copy", GetLastError());

	drawDC = CreateCompatibleDC(dc);
	if (drawDC == NULL)
		xpanic("error getting DC for checkbox image list bitmap", GetLastError());
	prevbitmap = SelectObject(drawDC, bitmap);
	if (prevbitmap == NULL)
		xpanic("error selecting checkbox image list bitmap into DC", GetLastError());
	// TODO get HRESULT
	if (DrawThemeBackground(theme, drawDC, BP_CHECKBOX, themestates[cbState], &r, NULL) != S_OK)
		xpanic("error drawing checkbox image", GetLastError());
	if (SelectObject(drawDC, prevbitmap) != bitmap)
		xpanic("error selecting previous bitmap into checkbox image's DC", GetLastError());
	if (DeleteDC(drawDC) == 0)
		xpanic("error deleting checkbox image's DC", GetLastError());

	return bitmap;
}


void makeCheckboxImageList_theme(HWND hwnddc)
{
	int width, height;
	int cbState;
	SIZE size;
	HDC dc;

	// first, make sure that all things have the same size
	dc = GetDC(hwnddc);
	if (dc == NULL)
		xpanic("error getting DC for making the checkbox image list", GetLastError());
	size = getStateSize(dc, 0);
	for (cbState = 1; cbState < checkboxnStates; cbState++) {
		SIZE against;

		against = getStateSize(dc, cbState);
		if (size.cx != against.cx || size.cy != against.cy)
			xpanic("size mismatch in checkbox states", GetLastError());
	}
	width = (int) size.cx;
	height = (int) size.cy;

	// NOW draw
	checkboxImageList = (*fv_ImageList_Create)(width, height, ILC_COLOR32, 20, 20);		// should be reasonable
	if (checkboxImageList == NULL)
		xpanic("error creating checkbox image list", GetLastError());
	for (cbState = 0; cbState < checkboxnStates; cbState++) {
		HBITMAP bitmap;

		bitmap = drawThemeImage(dc, width, height, cbState);
		if ((*fv_ImageList_Add)(checkboxImageList, bitmap, NULL) == -1)
			xpanic("error adding checkbox image to image list", GetLastError());
		if (DeleteObject(bitmap) == 0)
			xpanic("error deleting checkbox bitmap", GetLastError());
	}
	if (ReleaseDC(hwnddc, dc) == 0)
		xpanic("error deleting checkbox image list DC", GetLastError());
}

void makeCheckboxImageList(HWND hwnddc)
{
	if (theme == NULL)		// try to open the theme
		openTheme(hwnddc);
	if (theme != NULL) {		// use the theme
		makeCheckboxImageList_theme(hwnddc);
		return;
	}
	// couldn't open; fall back
	makeCheckboxImageList_DFC(hwnddc);
}

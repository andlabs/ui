// 16 august 2014

#define UNICODE
#define _UNICODE
#define STRICT
#define STRICT_TYPED_ITEMIDS
#define CINTERFACE
// get Windows version right; right now Windows XP
#define WINVER 0x0501
#define _WIN32_WINNT 0x0501
#define _WIN32_WINDOWS 0x0501		/* according to Microsoft's winperf.h */
#define _WIN32_IE 0x0600			/* according to Microsoft's sdkddkver.h */
#define NTDDI_VERSION 0x05010000	/* according to Microsoft's sdkddkver.h */
#include <windows.h>
#include <commctrl.h>
#include <stdint.h>
#include <uxtheme.h>
#include <string.h>
#include <wchar.h>
#include <windowsx.h>
#include <vsstyle.h>
#include <vssym32.h>
#include <oleacc.h>

enum {
        checkboxStateChecked = 1 << 0,
        checkboxStateHot = 1 << 1,
        checkboxStatePushed = 1 << 2,
        checkboxnStates = 1 << 3,
};

#define xpanic(...) abort()
#define xpanichresult(...) abort()

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

static void dfcImage(HDC dc, RECT *r, int cbState, HTHEME theme)
{
	if (DrawFrameControl(dc, r, DFC_BUTTON, dfcState(cbState)) == 0)
		xpanic("error drawing checkbox image", GetLastError());
}

static void dfcSize(HDC dc, int *width, int *height, HTHEME theme)
{
	// there's no real metric around
	// let's use SM_CX/YSMICON and hope for the best
	*width = GetSystemMetrics(SM_CXSMICON);
	*height = GetSystemMetrics(SM_CYSMICON);
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

static SIZE getStateSize(HDC dc, int cbState, HTHEME theme)
{
	SIZE s;
	HRESULT res;

	res = GetThemePartSize(theme, dc, BP_CHECKBOX, themestates[cbState], NULL, TS_DRAW, &s);
	if (res != S_OK)
		xpanichresult("error getting theme part size", res);
	return s;
}

static void themeImage(HDC dc, RECT *r, int cbState, HTHEME theme)
{
	HRESULT res;

	res = DrawThemeBackground(theme, dc, BP_CHECKBOX, themestates[cbState], r, NULL);
	if (res != S_OK)
		xpanichresult("error drawing checkbox image", res);
}

static void themeSize(HDC dc, int *width, int *height, HTHEME theme)
{
	SIZE size;
	int cbState;

	size = getStateSize(dc, 0, theme);
	for (cbState = 1; cbState < checkboxnStates; cbState++) {
		SIZE against;

		against = getStateSize(dc, cbState, theme);
		if (size.cx != against.cx || size.cy != against.cy)
			xpanic("size mismatch in checkbox states", GetLastError());
	}
	*width = (int) size.cx;
	*height = (int) size.cy;
}

static HBITMAP makeCheckboxImageListEntry(HDC dc, int width, int height, int cbState, void (*drawfunc)(HDC, RECT *, int, HTHEME), HTHEME theme)
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
	(*drawfunc)(drawDC, &r, cbState, theme);
	if (SelectObject(drawDC, prevbitmap) != bitmap)
		xpanic("error selecting previous bitmap into checkbox image's DC", GetLastError());
	if (DeleteDC(drawDC) == 0)
		xpanic("error deleting checkbox image's DC", GetLastError());

	return bitmap;
}

static HIMAGELIST newCheckboxImageList(HWND hwnddc, void (*sizefunc)(HDC, int *, int *, HTHEME), void (*drawfunc)(HDC, RECT *, int, HTHEME), HTHEME theme, int *width, int *height)
{
	int cbState;
	HDC dc;
	HIMAGELIST il;

	dc = GetDC(hwnddc);
	if (dc == NULL)
		xpanic("error getting DC for making the checkbox image list", GetLastError());
	(*sizefunc)(dc, width, height, theme);
	il = ImageList_Create(*width, *height, ILC_COLOR32, 20, 20);		// should be reasonable
	if (il == NULL)
		xpanic("error creating checkbox image list", GetLastError());
	for (cbState = 0; cbState < checkboxnStates; cbState++) {
		HBITMAP bitmap;

		bitmap = makeCheckboxImageListEntry(dc, *width, *height, cbState, drawfunc, theme);
		if (ImageList_Add(il, bitmap, NULL) == -1)
			xpanic("error adding checkbox image to image list", GetLastError());
		if (DeleteObject(bitmap) == 0)
			xpanic("error deleting checkbox bitmap", GetLastError());
	}
	if (ReleaseDC(hwnddc, dc) == 0)
		xpanic("error deleting checkbox image list DC", GetLastError());
	return il;
}

HIMAGELIST makeCheckboxImageList(HWND hwnddc, HTHEME *theme, int *width, int *height)
{
	if (*theme != NULL) {
		HRESULT res;

		res = CloseThemeData(*theme);
		if (res != S_OK)
			xpanichresult("error closing theme", res);
		*theme = NULL;
	}
	// ignore error; if it can't be done, we can fall back to DrawFrameControl()
	if (*theme == NULL)		// try to open the theme
		*theme = OpenThemeData(hwnddc, L"button");
	if (*theme != NULL)		// use the theme
		return newCheckboxImageList(hwnddc, themeSize, themeImage, *theme, width, height);
	// couldn't open; fall back
	return newCheckboxImageList(hwnddc, dfcSize, dfcImage, *theme, width, height);
}

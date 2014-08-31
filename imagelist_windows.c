// 16 august 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

// TODO top left pixel of checkbox state 0 not drawn?

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
	// image lists use non-premultiplied RGBA - see http://stackoverflow.com/a/25578789/3408572
	// the TRUE here does the conversion
	dotoARGB(i, (void *) ppvBits, TRUE);
	return bitmap;
}

HIMAGELIST newImageList(int width, int height)
{
	HIMAGELIST il;

	// this handles alpha properly; see https://web.archive.org/web/20100512144953/http://msdn.microsoft.com/en-us/library/ms997646.aspx#xptheming_topic13 and http://stackoverflow.com/a/2640897/3408572
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

void applyImageList(HWND hwnd, UINT uMsg, WPARAM wParam, HIMAGELIST il, HIMAGELIST old)
{
	if (SendMessageW(hwnd, uMsg, wParam, (LPARAM) il) != (LRESULT) old)
		xpanic("error setting image list", GetLastError());
	if (old != NULL && (*fv_ImageList_Destroy)(old) == 0)
		xpanic("error freeing old checkbox image list", GetLastError());

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

static HIMAGELIST newCheckboxImageList(HWND hwnddc, void (*sizefunc)(HDC, int *, int *, HTHEME), void (*drawfunc)(HDC, RECT *, int, HTHEME), HTHEME theme)
{
	int width, height;
	int cbState;
	HDC dc;
	HIMAGELIST il;

	dc = GetDC(hwnddc);
	if (dc == NULL)
		xpanic("error getting DC for making the checkbox image list", GetLastError());
	(*sizefunc)(dc, &width, &height, theme);
	il = (*fv_ImageList_Create)(width, height, ILC_COLOR32, 20, 20);		// should be reasonable
	if (il == NULL)
		xpanic("error creating checkbox image list", GetLastError());
	for (cbState = 0; cbState < checkboxnStates; cbState++) {
		HBITMAP bitmap;

		bitmap = makeCheckboxImageListEntry(dc, width, height, cbState, drawfunc, theme);
		if ((*fv_ImageList_Add)(il, bitmap, NULL) == -1)
			xpanic("error adding checkbox image to image list", GetLastError());
		if (DeleteObject(bitmap) == 0)
			xpanic("error deleting checkbox bitmap", GetLastError());
	}
	if (ReleaseDC(hwnddc, dc) == 0)
		xpanic("error deleting checkbox image list DC", GetLastError());
	return il;
}

HIMAGELIST makeCheckboxImageList(HWND hwnddc, HTHEME *theme)
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
		return newCheckboxImageList(hwnddc, themeSize, themeImage, *theme);
	// couldn't open; fall back
	return newCheckboxImageList(hwnddc, dfcSize, dfcImage, *theme);
}

// 16 august 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

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
	if (winDC == NULL)
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
	if (DeleteDC(winDC) == 0)
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

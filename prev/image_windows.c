// 16 august 2014

#include "winapi_windows.h"
#include "_cgo_export.h"

// TODO rename to images_windows.c?

HBITMAP toBitmap(void *i, intptr_t dx, intptr_t dy)
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
		xpanic("error creating HBITMAP in toBitmap()", GetLastError());
	// image lists use non-premultiplied RGBA - see http://stackoverflow.com/a/25578789/3408572
	// the TRUE here does the conversion
	dotoARGB(i, (void *) ppvBits, TRUE);
	return bitmap;
}

void freeBitmap(uintptr_t bitmap)
{
	if (DeleteObject((HBITMAP) bitmap) == 0)
		xpanic("error deleting bitmap in freeBitmap()", GetLastError());
}

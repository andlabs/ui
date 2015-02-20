// 19 october 2014
#include "../wininclude_windows.h"

// #qo LIBS: user32 kernel32 gdi32 comctl32 uxtheme ole32 oleaut32 oleacc uuid msimg32

#include "main.h"

HWND tablehwnd = NULL;
BOOL msgfont = FALSE;

HBITMAP mkbitmap(void);

BOOL mainwinCreate(HWND hwnd, LPCREATESTRUCT lpcs)
{
	intptr_t c;
	intptr_t row, col;

	tablehwnd = CreateWindowExW(0,
		tableWindowClass, L"Main Window",
		WS_CHILD | WS_VISIBLE | WS_HSCROLL | WS_VSCROLL,
		CW_USEDEFAULT, CW_USEDEFAULT,
		400, 400,
		hwnd, NULL, lpcs->hInstance, NULL);
	if (tablehwnd == NULL)
		panic("(test program) error creating Table");
	SendMessageW(tablehwnd, tableAddColumn, tableColumnText, (LPARAM) L"Column");
	SendMessageW(tablehwnd, tableAddColumn, tableColumnImage, (LPARAM) L"Column 2");
	SendMessageW(tablehwnd, tableAddColumn, tableColumnCheckbox, (LPARAM) L"Column 3");
	if (msgfont) {
		NONCLIENTMETRICSW ncm;
		HFONT font;

		ZeroMemory(&ncm, sizeof (NONCLIENTMETRICSW));
		ncm.cbSize = sizeof (NONCLIENTMETRICSW);
		if (SystemParametersInfoW(SPI_GETNONCLIENTMETRICS, sizeof (NONCLIENTMETRICSW), &ncm, sizeof (NONCLIENTMETRICSW)) == 0)
			panic("(test program) error getting non-client metrics");
		font = CreateFontIndirectW(&ncm.lfMessageFont);
		if (font == NULL)
			panic("(test program) error creating lfMessageFont HFONT");
		SendMessageW(tablehwnd, WM_SETFONT, (WPARAM) font, TRUE);
	}
	c = 100;
	SendMessageW(tablehwnd, tableSetRowCount, 0, (LPARAM) (&c));
	row = 2;
	col = 1;
	SendMessageW(tablehwnd, tableSetSelection, (WPARAM) (&row), (LPARAM) (&col));
	SetFocus(tablehwnd);
	return TRUE;
}

void mainwinDestroy(HWND hwnd)
{
        DestroyWindow(tablehwnd);
        PostQuitMessage(0);
}

void mainwinResize(HWND hwnd, UINT state, int cx, int cy)
{
        if (tablehwnd != NULL)
                MoveWindow(tablehwnd, 0, 0, cx, cy, TRUE);
}

BOOL checkboxstates[100];

LRESULT CALLBACK mainwndproc(HWND hwnd, UINT uMsg, WPARAM wParam, LPARAM lParam)
{
	NMHDR *nmhdr = (NMHDR *) lParam;
	tableNM *nm = (tableNM *) lParam;
	WCHAR *text;
	int n;

	if (uMsg == WM_CREATE)
		ZeroMemory(checkboxstates, 100 * sizeof (BOOL));
	switch (uMsg) {
	HANDLE_MSG(hwnd, WM_CREATE, mainwinCreate);
	HANDLE_MSG(hwnd, WM_SIZE, mainwinResize);
	HANDLE_MSG(hwnd, WM_DESTROY, mainwinDestroy);
	case WM_NOTIFY:
		if (nmhdr->hwndFrom != tablehwnd)
			break;
		switch (nmhdr->code) {
		case tableNotificationGetCellData:
			switch (nm->columnType) {
			case tableColumnText:
				n = _scwprintf(L"mainwin (%d,%d)", nm->row, nm->column);
				text = (WCHAR *) malloc((n + 1) * sizeof (WCHAR));
				if (text == NULL)
					panic("(table program) error allocating string");
				_swprintf(text, L"mainwin (%d,%d)", nm->row, nm->column);
				return (LRESULT) text;
			case tableColumnImage:
				return (LRESULT) mkbitmap();
			case tableColumnCheckbox:
				return (LRESULT) (checkboxstates[nm->row]);
			}
			panic("(test program) unreachable");
		case tableNotificationFinishedWithCellData:
			switch (nm->columnType) {
			case tableColumnText:
				free((void *) (nm->data));
				break;
			case tableColumnImage:
				if (DeleteObject((HBITMAP) (nm->data)) == 0)
					panic("(test program) error deleting cell image");
				break;
			}
			return 0;
		case tableNotificationCellCheckboxToggled:
			checkboxstates[nm->row] = !checkboxstates[nm->row];
			return 0;
		}
		break;
	}
	return DefWindowProcW(hwnd, uMsg, wParam, lParam);
}

int main(int argc, char *argv[])
{
	HWND mainwin;
	MSG msg;
	INITCOMMONCONTROLSEX icc;
	WNDCLASSW wc;

	if (argc != 1)
		msgfont = TRUE;
	ZeroMemory(&icc, sizeof (INITCOMMONCONTROLSEX));
	icc.dwSize = sizeof (INITCOMMONCONTROLSEX);
	icc.dwICC = ICC_LISTVIEW_CLASSES;
	if (InitCommonControlsEx(&icc) == 0)
		panic("(test program) error initializing comctl32.dll");
	initTable(NULL, _TrackMouseEvent);
	ZeroMemory(&wc, sizeof (WNDCLASSW));
	wc.lpszClassName = L"mainwin";
	wc.lpfnWndProc = mainwndproc;
	wc.hIcon = LoadIcon(NULL, IDI_APPLICATION);
	wc.hCursor = LoadCursor(NULL, IDC_ARROW);
	wc.hbrBackground = (HBRUSH) (COLOR_BTNFACE + 1);
	wc.hInstance = GetModuleHandle(NULL);
	if (RegisterClassW(&wc) == 0)
		panic("(test program) error registering main window class");
	mainwin = CreateWindowExW(0,
		L"mainwin", L"Main Window",
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT,
		400, 400,
		NULL, NULL, GetModuleHandle(NULL), NULL);
	if (mainwin == NULL)
		panic("(test program) error creating main window");
	ShowWindow(mainwin, SW_SHOWDEFAULT);
	if (UpdateWindow(mainwin) == 0)
		panic("(test program) error updating window");
	while (GetMessageW(&msg, NULL, 0, 0) > 0) {
		TranslateMessage(&msg);
		DispatchMessageW(&msg);
	}
	return 0;
}

// from tango-icon-theme-0.8.90/16x16/status/audio-volume-high.png (public domain)
COLORREF iconpix[] = {
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1003060A, 0x0, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x34101210, 0x5020202, 0x0, 0x0, 0x0, 0x0, 0x9E203E66, 0x78182F4D, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x32111110, 0xF051534F, 0x8030303, 0x0, 0x0, 0x0, 0x0, 0x1C050B12, 0xE92F5C96, 0x2B08111B, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x32111110, 0xF1575855, 0xFE545652, 0x8030303, 0x0, 0x0, 0x67152942, 0x440E1B2C, 0x0, 0x6214263F, 0xCB295083, 0x4000102,
	0xF3525450, 0xFB545653, 0xF0555653, 0x8830302E, 0xF1575855, 0xFFBEBFBC, 0xFC565854, 0x8030303, 0x0, 0x0, 0x420C192A, 0xE12E5991, 0xF03060A, 0x0, 0xC6284E7F, 0x52102034,
	0xFB555653, 0xFFE5E5E3, 0xFFD1D2CF, 0xFE585A56, 0xFFC7C8C5, 0xFFF7F7F6, 0xFC5A5C58, 0x8030303, 0x68142943, 0x1C060B12, 0x0, 0x9A1E3D63, 0x81193353, 0x0, 0x831A3454, 0x891C3658,
	0xFB545653, 0xFFFAFAF9, 0xFFFAFAFA, 0xFF585A56, 0xFFF8F8F8, 0xFFFFFFFF, 0xFC5B5D59, 0x8030303, 0x7B19314F, 0xAB23436E, 0x0, 0x4A0F1D30, 0xBA264978, 0x0, 0x4F101F33, 0xBD264B7A,
	0xFB545653, 0xFFCFD0CD, 0xFFD3D4D1, 0xFF585A56, 0xFFD2D2D0, 0xFFD7D7D5, 0xFC565854, 0x8030303, 0x1C060B12, 0xEF305F9A, 0x1000001, 0x19050A10, 0xEF305F9A, 0x1000001, 0x1B050A11, 0xF0315F9A,
	0xFB545653, 0xFFC2C3C0, 0xFFC6C7C4, 0xFF585A56, 0xFFCDCECB, 0xFFD1D2CF, 0xFC565854, 0x8030303, 0x19050A10, 0xF231609C, 0x1000001, 0x1603090E, 0xF1315F9B, 0x1000001, 0x1804090F, 0xF231609C,
	0xFB545653, 0xFFB5B7B3, 0xFFB9BBB7, 0xFF575955, 0xFFC8C8C6, 0xFFCCCCCA, 0xFC565854, 0x8030303, 0x71172D49, 0xB2244772, 0x0, 0x470D1C2E, 0xBE264B7A, 0x0, 0x4C0F1E31, 0xC0274C7B,
	0xFB545653, 0xFFA9ABA7, 0xFFA0A29E, 0xFF575955, 0xFFA9ABA8, 0xFFC6C7C4, 0xFC565854, 0x8030303, 0x75172E4B, 0x23070E17, 0x0, 0x921D3A5E, 0x861A3556, 0x0, 0x801A3252, 0x8C1C375A,
	0xFA535551, 0xFC555753, 0xF7545652, 0x98353534, 0xF3565854, 0xFFA8A8A6, 0xFC565854, 0x8030303, 0x0, 0x0, 0x390C1725, 0xE62F5B94, 0x1404080D, 0x0, 0xC0274C7B, 0x57112238,
	0x7020202, 0x8020202, 0x3010100, 0x0, 0x3C141414, 0xF3565854, 0xFE545652, 0x8030303, 0x0, 0x0, 0x70162C48, 0x4E0F1F32, 0x0, 0x57112238, 0xD32B5388, 0x6010203,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x3C141414, 0xF5515350, 0x8030303, 0x0, 0x0, 0x0, 0x0, 0x1604080E, 0xE72F5B95, 0x330A1420, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3F151515, 0x5020202, 0x0, 0x0, 0x0, 0x0, 0x991F3C62, 0x831A3454, 0x0, 0x0,
	0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1B050A11, 0x4000102, 0x0, 0x0,
};

HBITMAP mkbitmap(void)
{
	BITMAPINFO bi;
	void *ppvBits;
	HBITMAP b;

	ZeroMemory(&bi, sizeof (BITMAPINFO));
	bi.bmiHeader.biSize = sizeof (BITMAPINFOHEADER);
	bi.bmiHeader.biWidth = 16;
	bi.bmiHeader.biHeight = -16;			// negative height to force top-down drawing
	bi.bmiHeader.biPlanes = 1;
	bi.bmiHeader.biBitCount = 32;
	bi.bmiHeader.biCompression = BI_RGB;
	bi.bmiHeader.biSizeImage = 16 * 16 * 4;
	b = CreateDIBSection(NULL, &bi, DIB_RGB_COLORS, &ppvBits, 0, 0);
	if (b == 0)
		panic("test bitmap creation failed");
	memcpy(ppvBits, iconpix, bi.bmiHeader.biSizeImage);
	return b;
}

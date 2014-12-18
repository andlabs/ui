// 16 august 2014

// TODO instead of caching checkbox images, draw them on the fly, because they could be transparent

// TODO actually make this
#define panichresult(a, b) panic(a)

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
		panic("error drawing Table checkbox image with DrawFrameControl()");
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
		panichresult("error getting theme part size for Table checkboxes", res);
	return s;
}

static void themeImage(HDC dc, RECT *r, int cbState, HTHEME theme)
{
	HRESULT res;

	res = DrawThemeBackground(theme, dc, BP_CHECKBOX, themestates[cbState], r, NULL);
	if (res != S_OK)
		panichresult("error drawing Table checkbox image from theme", res);
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
			// TODO make this use a no-information (or two ints) panic()
			panic("size mismatch in Table checkbox states");
	}
	*width = (int) size.cx;
	*height = (int) size.cy;
}

static void makeCheckboxImage(struct table *t, HDC dc, int cbState, void (*drawfunc)(HDC, RECT *, int, HTHEME))
{
	BITMAPINFO bi;
	VOID *ppvBits;
	HBITMAP bitmap;
	RECT r;
	HDC drawDC;
	HBITMAP prevbitmap;

	r.left = 0;
	r.top = 0;
	r.right = t->checkboxWidth;
	r.bottom = t->checkboxHeight;
	ZeroMemory(&bi, sizeof (BITMAPINFO));
	bi.bmiHeader.biSize = sizeof (BITMAPINFOHEADER);
	bi.bmiHeader.biWidth = (LONG) (t->checkboxWidth);
	bi.bmiHeader.biHeight = -((LONG) (t->checkboxHeight));			// negative height to force top-down drawing;
	bi.bmiHeader.biPlanes = 1;
	bi.bmiHeader.biBitCount = 32;
	bi.bmiHeader.biCompression = BI_RGB;
	bi.bmiHeader.biSizeImage = (DWORD) (t->checkboxWidth * t->checkboxHeight * 4);
	bitmap = CreateDIBSection(NULL, &bi, DIB_RGB_COLORS, &ppvBits, 0, 0);
	if (bitmap == NULL)
		panic("error creating HBITMAP for Table checkbox image");

	drawDC = CreateCompatibleDC(dc);
	if (drawDC == NULL)
		panic("error getting DC for drawing Table checkbox image");
	prevbitmap = SelectObject(drawDC, bitmap);
	if (prevbitmap == NULL)
		panic("error selecting Table checkbox image list bitmap into DC");
	(*drawfunc)(drawDC, &r, cbState, t->theme);
	if (SelectObject(drawDC, prevbitmap) != bitmap)
		panic("error selecting previous bitmap into Table checkbox image's DC");
	if (DeleteDC(drawDC) == 0)
		panic("error deleting Table checkbox image's DC");

	t->checkboxImages[cbState] = bitmap;
}

static void getCheckboxImages(struct table *t, void (*sizefunc)(HDC, int *, int *, HTHEME), void (*drawfunc)(HDC, RECT *, int, HTHEME))
{
	int cbState;
	HDC dc;

	dc = GetDC(t->hwnd);
	if (dc == NULL)
		panic("error getting DC for making Table checkbox images");
	(*sizefunc)(dc, &(t->checkboxWidth), &(t->checkboxHeight), t->theme);
	for (cbState = 0; cbState < checkboxnStates; cbState++)
		makeCheckboxImage(t, dc, cbState, drawfunc);
	if (ReleaseDC(t->hwnd, dc) == 0)
		panic("error deleting Table DC for making checkbox images");
}

static void makeCheckboxImages(struct table *t)
{
	if (t->theme != NULL) {
		HRESULT res;

		res = CloseThemeData(t->theme);
		if (res != S_OK)
			panichresult("error closing theme", res);
		t->theme = NULL;
	}
	// ignore error; if it can't be done, we can fall back to DrawFrameControl()
	if (t->theme == NULL)		// try to open the theme
		t->theme = OpenThemeData(t->hwnd, L"button");
	if (t->theme != NULL) {		// use the theme
		getCheckboxImages(t, themeSize, themeImage);
		return;
	}
	// couldn't open; fall back
	getCheckboxImages(t, dfcSize, dfcImage);
}

static void freeCheckboxImages(struct table *t)
{
	int cbState;

	for (cbState = 0; cbState < checkboxnStates; cbState++)
		if (DeleteObject(t->checkboxImages[cbState]) == 0)
			panic("error freeing Table checkbox image");
}

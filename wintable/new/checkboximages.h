// 16 august 2014

enum {
        checkboxStateChecked = 1 << 0,
        checkboxStateHot = 1 << 1,
        checkboxStatePushed = 1 << 2,
        checkboxnStates = 1 << 3,
};

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

static void drawFrameControlCheckbox(HDC dc, RECT *r, int cbState)
{
	if (DrawFrameControl(dc, r, DFC_BUTTON, dfcState(cbState)) == 0)
		panic("error drawing Table checkbox image with DrawFrameControl()");
}

static void getFrameControlCheckboxSize(HDC dc, int *width, int *height, HTHEME theme)
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

static void drawThemeCheckbox(HDC dc, RECT *r, int cbState, HTHEME theme)
{
	HRESULT res;

	res = DrawThemeBackground(theme, dc, BP_CHECKBOX, themestates[cbState], r, NULL);
	if (res != S_OK)
		panichresult("error drawing Table checkbox image from theme", res);
}

static void getThemeCheckboxSize(HDC dc, int *width, int *height, HTHEME theme)
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

static void drawCheckbox(struct table *t, HDC dc, int x, int y, int cbState)
{
	RECT r;

	r.left = x;
	r.top = y;
	r.right = r.bottom + t->checkboxWidth;
	r.bottom = r.top + t->checkboxHeight;
	if (t->theme != NULL) {
		drawThemeCheckbox(dc, &r, cbState, t->theme);
		return;
	}
	drawFrameControlCheckbox(dc, &r, cbState);
}

static void freeCheckboxThemeData(struct table *t)
{
	if (t->theme != NULL) {
		HRESULT res;

		res = CloseThemeData(t->theme);
		if (res != S_OK)
			panichresult("error closing Table checkbox theme", res);
		t->theme = NULL;
	}
}

static void loadCheckboxThemeData(struct table *t)
{
	HDC dc;

	freeCheckboxThemeData(t);
	dc = GetDC(t->hwnd);
	if (dc == NULL)
		panic("error getting Table DC for loading checkbox theme data");
	// ignore error; if it can't be done, we can fall back to DrawFrameControl()
	if (t->theme == NULL)		// try to open the theme
		t->theme = OpenThemeData(t->hwnd, L"button");
	if (t->theme != NULL)		// use the theme
		getThemeCheckboxSize(dc, &(t->checkboxWidth), &(t->checkboxHeight), t->theme);
	else						// couldn't open; fall back
		getFrameControlCheckboxSize(dc, &(t->checkboxWidth), &(t->checkboxHeight), t->theme);
	if (ReleaseDC(t->hwnd, dc) == 0)
		panic("error releasing Table DC for loading checkbox theme data");
}

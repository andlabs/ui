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

static void getFrameControlCheckboxSize(HDC dc, int *width, int *height)
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

static void drawCheckbox(struct table *t, HDC dc, RECT *r, int cbState)
{
	if (t->theme != NULL) {
		drawThemeCheckbox(dc, r, cbState, t->theme);
		return;
	}
	drawFrameControlCheckbox(dc, r, cbState);
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
		getFrameControlCheckboxSize(dc, &(t->checkboxWidth), &(t->checkboxHeight));
	if (ReleaseDC(t->hwnd, dc) == 0)
		panic("error releasing Table DC for loading checkbox theme data");
}

static void redrawCheckboxRect(struct table *t, LPARAM lParam)
{
	struct rowcol rc;
	RECT r;

	rc = lParamToRowColumn(t, lParam);
	if (rc.row == -1 && rc.column == -1)
		return;
	if (t->columnTypes[rc.column] != tableColumnCheckbox)
		return;
	if (!rowColumnToClientRect(t, rc, &r))
		return;
	// TODO only the checkbox rect?
	if (InvalidateRect(t->hwnd, &r, TRUE) == 0)
		panic("error redrawing Table checkbox rect for mouse events");
}

HANDLER(checkboxMouseMoveHandler)
{
	// don't actually check to see if the mouse is in the checkbox rect
	// if there's scrolling without mouse movement, that will change
	// instead, drawCell() will handle it
	if (!t->checkboxMouseOverLast) {
		t->checkboxMouseOverLast = TRUE;
		retrack(t);
	} else
		redrawCheckboxRect(t, t->checkboxMouseOverLastPoint);
	t->checkboxMouseOverLastPoint = lParam;
	redrawCheckboxRect(t, t->checkboxMouseOverLastPoint);
	*lResult = 0;
	return TRUE;
}

HANDLER(checkboxMouseLeaveHandler)
{
	if (t->checkboxMouseOverLast)
		redrawCheckboxRect(t, t->checkboxMouseOverLastPoint);
	// TODO remember what I wanted to do here in the case of a held mouse button
	t->checkboxMouseOverLast = FALSE;
	*lResult = 0;
	return TRUE;
}

HANDLER(checkboxMouseDownHandler)
{
	struct rowcol rc;
	RECT r;
	POINT pt;

	rc = lParamToRowColumn(t, lParam);
	if (rc.row == -1 || rc.column == -1)
		return FALSE;
	if (t->columnTypes[rc.column] != tableColumnCheckbox)
		return FALSE;
	if (!rowColumnToClientRect(t, rc, &r))
		return FALSE;
	toCheckboxRect(t, &r, 0);
	pt.x = GET_X_LPARAM(lParam);
	pt.y = GET_Y_LPARAM(lParam);
	if (PtInRect(&r, pt) == 0)
		return FALSE;
	t->checkboxMouseDown = TRUE;
	t->checkboxMouseDownRow = rc.row;
	t->checkboxMouseDownColumn = rc.column;
	// TODO redraw the whole cell?
	if (InvalidateRect(t->hwnd, &r, TRUE) == 0)
		panic("error redrawing Table checkbox after mouse down");
	*lResult = 0;
	return TRUE;
}

HANDLER(checkboxMouseUpHandler)
{
	struct rowcol rc;
	RECT r;
	POINT pt;

	if (!t->checkboxMouseDown)
		return FALSE;
	// the logic behind goto wrongUp is that the mouse must be released on the same checkbox
	rc = lParamToRowColumn(t, lParam);
	if (rc.row == -1 || rc.column == -1)
		goto wrongUp;
	if (rc.row != t->checkboxMouseDownRow || rc.column != t->checkboxMouseDownColumn)
		goto wrongUp;
	if (t->columnTypes[rc.column] != tableColumnCheckbox)
		goto wrongUp;
	if (!rowColumnToClientRect(t, rc, &r))
		goto wrongUp;
	toCheckboxRect(t, &r, 0);
	pt.x = GET_X_LPARAM(lParam);
	pt.y = GET_Y_LPARAM(lParam);
	if (PtInRect(&r, pt) == 0)
		goto wrongUp;
	notify(t, tableNotificationCellCheckboxToggled, rc.row, rc.column, 0);
	t->checkboxMouseDown = FALSE;
	// TODO redraw the whole cell?
	if (InvalidateRect(t->hwnd, &r, TRUE) == 0)
		panic("error redrawing Table checkbox after mouse up");
	// TODO really only the row? no way to specify column too?
	NotifyWinEvent(EVENT_OBJECT_STATECHANGE, t->hwnd, OBJID_CLIENT, rc.row);
	*lResult = 0;
	return TRUE;
wrongUp:
	if (t->checkboxMouseDown) {
		rc.row = t->checkboxMouseDownRow;
		rc.column = t->checkboxMouseDownColumn;
		if (rowColumnToClientRect(t, rc, &r))
			// TODO only the checkbox rect?
			if (InvalidateRect(t->hwnd, &r, TRUE) == 0)
				panic("error redrawing Table checkbox rect for aborted mouse up event");
	}
	// if we landed on another checkbox, be sure to draw that one too
	if (t->checkboxMouseOverLast)
		redrawCheckboxRect(t, t->checkboxMouseOverLastPoint);
	t->checkboxMouseDown = FALSE;
	return FALSE;		// TODO really?
}

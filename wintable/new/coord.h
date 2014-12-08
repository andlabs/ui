// 4 december 2014

typedef struct rowcol rowcol;

struct rowcol {
	intptr_t row;
	intptr_t column;
};

static rowcol clientCoordToRowColumn(struct table *t, POINT pt)
{
	// TODO
}

// same as client coordinates, but stored in a lParam (like the various mouse messages provide)
static rowcol lParamToRowColumn(struct table *t, LPARAM lParam)
{
	POINT pt;

	pt.x = GET_X_LPARAM(lParam);
	pt.y = GET_Y_LPARAM(lParam);
	return clientCoordToRowColumn(t, pt);
}

// returns TRUE if the row is visible and thus has client coordinates; FALSE otherwise
static BOOL rowColumnToClientCoord(struct table *t, rowcol rc, struct POINT *pt)
{
	// TODO
}

// TODO idealCoordToRowColumn/rowColumnToIdealCoord?

// TODO find a better place for this
static LONG rowHeight(struct table *t, HDC dc, BOOL select)
{
	BOOL release;
	HFONT prevfont, newfont;
	TEXTMETRICW tm;

	release = FALSE;
	if (dc == NULL) {
		dc = GetDC(t->hwnd);
		if (dc == NULL)
			panic("error getting Table DC for rowHeight()");
		release = TRUE;
	}
	if (select)
		prevfont = selectFont(t, dc, &newfont);
	if (GetTextMetricsW(dc, &tm) == 0)
		panic("error getting text metrics for rowHeight()");
	if (select)
		deselectFont(dc, prevfont, newfont);
	if (release)
		if (ReleaseDC(t->hwnd, dc) == 0)
			panic("error releasing Table DC for rowHeight()");
	return tm.tmHeight;
}

#define rowht(t) rowHeight(t, NULL, TRUE)

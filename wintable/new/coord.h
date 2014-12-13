// 4 december 2014

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

struct rowcol {
	intptr_t row;
	intptr_t column;
};

static struct rowcol clientCoordToRowColumn(struct table *t, POINT pt)
{
	RECT r;
	struct rowcol rc;
	intptr_t i;
	RECT colrect;

	// initial values for the PtInRect() check
	rc.row = -1;
	rc.column = -1;

	if (GetClientRect(t->hwnd, &r) == 0)
		panic("error getting Table client rect in clientCoordToRowColumn()");
	r.top += t->headerHeight;
	if (PtInRect(&r, pt) == 0)
		return rc;

	// the row is easy
	pt.y -= t->headerHeight;
	rc.row = (pt.y / rowht(t)) + t->vscrollpos;

	// the column... not so much
	// we scroll p.x, then subtract column widths until we cross the left edge of the control
	pt.x += t->hscrollpos;
	rc.column = 0;
	for (i = 0; i < t->nColumns; i++) {
		// TODO error check
		SendMessage(t->header, HDM_GETITEMRECT, (WPARAM) i, (LPARAM) (&colrect));
		pt.x -= colrect.right - colrect.left;
		// use <, not <=, here:
		// assume r.left and t->hscrollpos == 0;
		// given the first column is 100 wide,
		// pt.x == 0 (first pixel of col 0) -> p.x - 100 == -100 < 0 -> break
		// pt.x == 99 (last pixel of col 0) -> p.x - 100 == -1 < 0 -> break
		// pt.x == 100 (first pixel of col 1) -> p.x - 100 == 0 >= 0 -> next column
		if (pt.x < r.left)
			break;
		rc.column++;
	}
	// TODO what happens if the break was never taken?

	return rc;
}

// same as client coordinates, but stored in a lParam (like the various mouse messages provide)
static struct rowcol lParamToRowColumn(struct table *t, LPARAM lParam)
{
	POINT pt;

	pt.x = GET_X_LPARAM(lParam);
	pt.y = GET_Y_LPARAM(lParam);
	return clientCoordToRowColumn(t, pt);
}

// returns TRUE if the row is visible and thus has client coordinates; FALSE otherwise
static BOOL rowColumnToClientCoord(struct table *t, struct rowcol rc, POINT *pt)
{
	// TODO
}

// TODO idealCoordToRowColumn/rowColumnToIdealCoord?

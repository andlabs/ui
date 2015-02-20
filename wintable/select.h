// 13 december 2014

// damn winsock
static void doselect(struct table *t, intptr_t row, intptr_t column)
{
	RECT r, client;
	intptr_t oldrow;
	LONG width, height;
	struct rowcol rc;
	BOOL dovscroll;
	intptr_t i;
	intptr_t xpos;
	LONG clientWidth;

	// check existing selection to see if it's valid
	if (t->selectedRow == -1 && t->selectedColumn != -1)
		panic("sanity check failure: old Table selection invalid (row == -1, column != -1)");
	if (t->selectedRow != -1 && t->selectedColumn == -1)
		panic("sanity check failure: old Table selection invalid (row != -1, column == -1)");
	if (t->selectedRow >= t->count)
		panic("sanity check failure: old Table selection invalid (row out of range)");
	if (t->selectedColumn >= t->nColumns)
		panic("sanity check failure: old Table selection invalid (column out of range)");

	oldrow = t->selectedRow;
	t->selectedRow = row;
	t->selectedColumn = column;

	// check new selection to see if it's valid
	if (t->selectedRow == -1 && t->selectedColumn != -1)
		panic("sanity check failure: new Table selection invalid (row == -1, column != -1)");
	if (t->selectedRow != -1 && t->selectedColumn == -1)
		panic("sanity check failure: new Table selection invalid (row != -1, column == -1)");
	if (t->selectedRow >= t->count)
		panic("sanity check failure: new Table selection invalid (row out of range)");
	if (t->selectedColumn >= t->nColumns)
		panic("sanity check failure: new Table selection invalid (column out of range)");

	// do this even if we don't scroll before; noScroll depends on it
	if (GetClientRect(t->hwnd, &client) == 0)
		panic("error getting Table client rect in doselect()");
	client.top += t->headerHeight;
	height = rowht(t);

	// only scroll if we selected something
	if (t->selectedRow == -1 || t->selectedColumn == -1)
		goto noScroll;

	// first vertically scroll to the new row to make it fully visible (or as visible as possible)
	if (t->selectedRow < t->vscrollpos)
		vscrollto(t, t->selectedRow);
	else {
		rc.row = t->selectedRow;
		rc.column = t->selectedColumn;
		// first assume entirely outside the client area
		dovscroll = TRUE;
		if (rowColumnToClientRect(t, rc, &r))
			// partially outside the client area?
			if (r.bottom <= client.bottom)		// <= here since we are comparing bottoms (which are the first pixels outside the rectangle)
				dovscroll = FALSE;
		if (dovscroll)
			vscrollto(t, t->selectedRow - t->vpagesize + 1);		// + 1 because apparently just t->selectedRow - t->vpagesize results in no scrolling (t->selectedRow - t->vpagesize == t->vscrollpos)...
	}

	// now see if the cell we want is to the left of offscreen, in which case scroll to its x-position
	xpos = 0;
	for (i = 0; i < t->selectedColumn; i++)
		xpos += columnWidth(t, i);
	if (xpos < t->hscrollpos)
		hscrollto(t, xpos);
	else {
		// if the full cell is not visible, scroll to the right just enough to make it fully visible (or as visible as possible)
		width = columnWidth(t, t->selectedColumn);
		clientWidth = client.right - client.left;
		if (xpos + width > t->hscrollpos + clientWidth)			// > because both sides deal with the first pixel outside
			// if the column is too wide, then just make it occupy the whole visible area (left-aligned)
			if (width > clientWidth)			// TODO >= ?
				hscrollto(t, xpos);
			else
				// TODO don't use t->hpagesize here? depends if other code uses it
				hscrollto(t, (xpos + width) - t->hpagesize);
	}

noScroll:
	// now redraw the old and new /rows/
	// we do this after scrolling so the rectangles to be invalidated make sense
	r.left = client.left;
	r.right = client.right;
	if (oldrow != -1 && oldrow >= t->vscrollpos) {
		r.top = client.top + ((oldrow - t->vscrollpos) * height);
		r.bottom = r.top + height;
		if (InvalidateRect(t->hwnd, &r, TRUE) == 0)
			panic("error queueing previously selected row for redraw in doselect()");
	}
	// t->selectedRow must be visible by this point; we scrolled to it
	if (t->selectedRow != -1 && t->selectedRow != oldrow) {
		r.top = client.top + ((t->selectedRow - t->vscrollpos) * height);
		r.bottom = r.top + height;
		if (InvalidateRect(t->hwnd, &r, TRUE) == 0)
			panic("error queueing newly selected row for redraw in doselect()");
	}

	// TODO what about deselect/defocus?
	// TODO notify on the old row too?
	NotifyWinEvent(EVENT_OBJECT_SELECTION, t->hwnd, OBJID_CLIENT, t->selectedRow);
	// TODO send EVENT_OBJECT_STATECHANGED too?
	// TODO send EVENT_OBJECT_FOCUS

	// TODO before or after NotifyWinEvent()? (see what other things I'm doing)
	notify(t, tableNotificationSelectionChanged, t->selectedRow, t->selectedColumn, 0);
}

// TODO make this needless
HANDLER(checkboxMouseDownHandler);

// TODO which WM_xBUTTONDOWNs?
HANDLER(mouseDownSelectHandler)
{
	struct rowcol rc;

	// TODO separate this from here
	// TODO other mouse buttons?
	// don't check SetFocus()'s error (http://stackoverflow.com/questions/24073695/winapi-can-setfocus-return-null-without-an-error-because-thats-what-im-see)
	// TODO NotifyWinEvent() here?
	SetFocus(t->hwnd);
	rc = lParamToRowColumn(t, lParam);
	// don't check if lParamToRowColumn() returned row -1 or column -1; we want deselection behavior
	doselect(t, rc.row, rc.column);
	// TODO separate this from here
	checkboxMouseDownHandler(t, uMsg, wParam, lParam, lResult);
	*lResult = 0;
	return TRUE;
}

/*
the routine below is intended to simulate the comctl32.dll listview keyboard navigation rules, at least as far as vertical navigation is concerned.
horizontal scrolling is different because unlike the comctl32 listview, we say that a single column in each row has the keyboard focus, so left and right navigate between columns here, instead of scrolling left/right by pixels.
	TODO provide an override for scrolling by pixels?
	TODO any other keyboard shortcuts?
		TODO browser keys
		TODO media navigation keys
			TODO XBUTTON1/2?
		TODO clear keys?

keyboard selection behaviors of the windows 7 listview:
with 100 items (0-99), the window currently shows items 30 through 47 as well as having item 48 partially visible
-  item 30:
	- page up -> item 13
	- page down -> item 47
- item 31:
	- page up -> item 30
	- page down -> item 47
- item 42:
	- page up -> item 30
	- page down -> item 47
- item 46:
	- page up -> item 30
	- page down -> item 47
- item 47:
	- page up: -> item 30
	- page down: -> item 64

when nothing is selected:
- down selects item 0 regardless of scroll
- up selects nothing regardless of scroll
- page down selects the last fully visible item depending on scroll
	- so with the above configuration:
		- item 0 -> item 17
		- item 30 -> item 47
		- item 80 -> item 97
- page up selects item 0 regardless of scroll
- home selects item 0 regardless of scroll
- end selects the last item regardless of scroll

for left and right we will simulate up and down, respectively (so right selects row 0 column 0); remember that you can't have either row or column be -1 but not both

TODO what happens if page up and page down are pressed with an item selected and the scroll in a different position?
*/

HANDLER(keyDownSelectHandler)
{
	intptr_t row;
	intptr_t column;

	if (t->count == 0 || t->nColumns == 0)		// no items to select
		return FALSE;
	row = t->selectedRow;
	column = t->selectedColumn;
	switch (wParam) {
	case VK_UP:
		if (row == -1)
			return FALSE;
		row--;
		if (row < 0)
			row = 0;
		break;
	case VK_DOWN:
		if (row == -1) {
			row = 0;
			column = 0;
		} else {
			row++;
			if (row >= t->count)
				row = t->count - 1;
		}
		break;
	case VK_LEFT:
		if (column == -1)
			return FALSE;
		column--;
		if (column < 0)
			column = 0;
		break;
	case VK_RIGHT:
		if (column == -1) {
			row = 0;
			column = 0;
		} else {
			column++;
			if (column >= t->nColumns)
				column = t->nColumns - 1;
		}
		break;
	case VK_HOME:
		row = 0;
		if (column == -1)
			column = 0;
		break;
	case VK_END:
		row = t->count - 1;
		if (column == -1)
			column = 0;
		break;
	case VK_PRIOR:
		if (row == -1) {
			row = 0;
			column = 0;
		} else {
			row = t->vscrollpos;
			if (row == t->selectedRow)
				// TODO investigate why the - 1 is needed here and below
				// TODO if this is a misunderstanding of how t->vpagesize works, figure out what happens if there is no partially visible row, and what is supposed to happen
				row -= t->vpagesize - 1;
			if (row < 0)
				row = 0;
		}
		break;
	case VK_NEXT:
		if (row == -1) {
			row = t->vscrollpos + t->vpagesize - 1;
			// TODO ensusre this is the case with the real list view
			if (row >= t->count)
				row = t->count - 1;
			column = 0;
		} else {
			row = t->vscrollpos + t->vpagesize - 1;
			if (row == t->selectedRow)
				row += t->vpagesize - 1;
			if (row >= t->count)
				row = t->count - 1;
		}
		break;
	default:
		return FALSE;
	}
	doselect(t, row, column);
	*lResult = 0;
	return TRUE;
}

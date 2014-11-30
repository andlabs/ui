// 30 november 2014

static void finishSelect(struct table *t, intptr_t prev)
{
	if (t->selected < 0)
		t->selected = 0;
	if (t->selected >= t->count)
		t->selected = t->count - 1;

	// always redraw the old and new rows to avoid artifacts when scrolling, even if they are the same (since the focused column may have changed)
	redrawRow(t, prev);
	if (prev != t->selected)
		redrawRow(t, t->selected);

	// if we need to scroll, the scrolling will force a redraw, so we don't have to worry about doing so ourselves
	if (t->selected < t->firstVisible)
		vscrollto(t, t->selected);
	// note that this is not lastVisible(t) because the last visible row may only be partially visible and we want selections to make them fully visible
	else if (t->selected >= (t->firstVisible + t->pagesize))
		vscrollto(t, t->selected - t->pagesize + 1);
}

// TODO isolate functionality so other keyboard event handlers can run
static void keySelect(struct table *t, WPARAM wParam, LPARAM lParam)
{
	intptr_t prev;

	// TODO figure out correct behavior with nothing selected
	if (t->count == 0)		// don't try to do anything if there's nothing to do
		return;
	prev = t->selected;
	switch (wParam) {
	case VK_UP:
		t->selected--;
		break;
	case VK_DOWN:
		t->selected++;
		break;
	case VK_PRIOR:
		t->selected -= t->pagesize;
		break;
	case VK_NEXT:
		t->selected += t->pagesize;
		break;
	case VK_HOME:
		t->selected = 0;
		break;
	case VK_END:
		t->selected = t->count - 1;
		break;
	case VK_LEFT:
		t->focusedColumn--;
		if (t->focusedColumn < 0)
			if (t->nColumns == 0)		// peg at -1
				t->focusedColumn = -1;
			else
				t->focusedColumn = 0;
		break;
	case VK_RIGHT:
		t->focusedColumn++;
		if (t->focusedColumn >= t->nColumns)
			if (t->nColumns == 0)		// peg at -1
				t->focusedColumn = -1;
			else
				t->focusedColumn = t->nColumns - 1;
		break;
	// TODO keyboard shortcuts for going to the first/last column?
	default:
		// don't touch anything
		return;
	}
	finishSelect(t, prev);
}

// TODO rename
static void selectItem(struct table *t, WPARAM wParam, LPARAM lParam)
{
	intptr_t prev;

	prev = t->selected;
	lParamToRowColumn(t, lParam, &(t->selected), &(t->focusedColumn));
	finishSelect(t, prev);
}

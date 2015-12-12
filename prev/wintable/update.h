// 8 january 2015

// Whenever a number of things in the Table changes, the update() function needs to be called to update any metrics and scrolling positions.
// The control font changing is the big one, as that comes with a flag that decides whether or not to redraw everything. We'll need to respect that here.

// For my personal convenience, each invocation of update() and updateAll() will be suffixed with a DONE comment once I have reasoned that the chosen function is correct and that no additional redrawing is necessary.

// TODO actually use redraw here
static void update(struct table *t, BOOL redraw)
{
	RECT client;
	intptr_t i;
	intptr_t height;

	// before we do anything we need the client rect
	if (GetClientRect(t->hwnd, &client) == 0)
		panic("error getting Table client rect in update()");

	// the first step is to figure out how wide the whole table is
	// TODO count dividers?
	t->width = 0;
	for (i = 0; i < t->nColumns; i++)
		t->width += columnWidth(t, i);

	// now we need to figure out how much of the width of the table can be seen at once
	t->hpagesize = client.right - client.left;
	// this part is critical: if we resize the columns to less than the client area width, then the following hscrollby() will make t->hscrollpos negative, which does very bad things
	// we do this regardless of which of the two has changed, just to be safe
	if (t->hpagesize > t->width)
		t->hpagesize = t->width;

	// now we do a dummy horizontal scroll to apply the new width and horizontal page size
	// this will also reposition and resize the header (the latter in case the font changed), which will be important for the next step
	hscrollby(t, 0);

	// now that we have the new height of the header, we can fix up vertical scrolling
	// so let's take the header height away from the client area
	client.top += t->headerHeight;
	// and update our page size appropriately
	height = client.bottom - client.top;
	t->vpagesize = height / rowht(t);
	// and do a dummy vertical scroll to apply that
	vscrollby(t, 0);
}

// this is the same as update(), but also redraws /everything/
// as such, redraw is TRUE
static void updateAll(struct table *t)
{
	update(t, TRUE);
	if (InvalidateRect(t->hwnd, NULL, TRUE) == 0)
		panic("error queueing all of Table for redraw in updateAll()");
}

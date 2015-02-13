// 4 december 2014

typedef BOOL (*handlerfunc)(struct table *, UINT, WPARAM, LPARAM, LRESULT *);
#define HANDLER(name) static BOOL name(struct table *t, UINT uMsg, WPARAM wParam, LPARAM lParam, LRESULT *lResult)

static BOOL runHandlers(const handlerfunc list[], struct table *t, UINT uMsg, WPARAM wParam, LPARAM lParam, LRESULT *lResult)
{
	const handlerfunc *p;

	for (p = list; *p != NULL; p++)
		if ((*(*p))(t, uMsg, wParam, lParam, lResult))
			return TRUE;
	return FALSE;
}

// memory allocation stuff
// each of these functions do an implicit ZeroMemory()
// these also make tableRealloc(NULL, ...)/tableFree(NULL) act like realloc(NULL, ...)/free(NULL) (that is, same as tableAlloc(...)/malloc(...) and a no-op, respectively)
// we /would/ use LocalAlloc() here because
// - whether the malloc() family supports the last-error code is undefined
// - the HeapAlloc() family DOES NOT support the last-error code; you're supposed to use Windows exceptions, and I can't find a clean way to do this with MinGW-w64 that doesn't rely on inline assembly or external libraries (unless they added __try/__except blocks)
// - there's no VirtualReAlloc() to complement VirtualAlloc() and I'm not sure if we can even get the original allocated size back out reliably to write it ourselves (http://blogs.msdn.com/b/oldnewthing/archive/2012/03/16/10283988.aspx)
// alas, LocalAlloc() doesn't want to work on real Windows 7 after a few times, throwing up ERROR_NOT_ENOUGH_MEMORY after three (3) ints or so :|
// we'll use malloc() until then
// needless to say, TODO

static void *tableAlloc(size_t size, const char *panicMessage)
{
//	HLOCAL out;
	void *out;

//	out = LocalAlloc(LMEM_FIXED | LMEM_ZEROINIT, size);
	out = malloc(size);
	if (out == NULL)
		panic(panicMessage);
	ZeroMemory(out, size);
	return (void *) out;
}

static void *tableRealloc(void *p, size_t size, const char *panicMessage)
{
//	HLOCAL out;
	void *out;

	if (p == NULL)
		return tableAlloc(size, panicMessage);
//	out = LocalReAlloc((HLOCAL) p, size, LMEM_ZEROINIT);
	out = realloc(p, size);
	if (out == NULL)
		panic(panicMessage);
	// TODO zero the extra memory
	return (void *) out;
}

static void tableFree(void *p, const char *panicMessage)
{
	if (p == NULL)
		return;
//	if (LocalFree((HLOCAL) p) != NULL)
//		panic(panicMessage);
	free(p);
}

// font selection

static HFONT selectFont(struct table *t, HDC dc, HFONT *newfont)
{
	HFONT prevfont;

	// copy it in casse we get a WM_SETFONT before this call's respective deselectFont() call
	*newfont = t->font;
	if (*newfont == NULL) {
		// get it on demand in the (unlikely) event it changes while this Table is alive
		*newfont = GetStockObject(SYSTEM_FONT);
		if (*newfont == NULL)
			panic("error getting default font for selecting into Table DC");
	}
	prevfont = (HFONT) SelectObject(dc, *newfont);
	if (prevfont == NULL)
		panic("error selecting Table font into Table DC");
	return prevfont;
}

static void deselectFont(HDC dc, HFONT prevfont, HFONT newfont)
{
	if (SelectObject(dc, prevfont) != newfont)
		panic("error deselecting Table font from Table DC");
	// doin't delete newfont here, even if it is the system font (see http://msdn.microsoft.com/en-us/library/windows/desktop/dd144925%28v=vs.85%29.aspx)
}

// and back to other functions

static LONG columnWidth(struct table *t, intptr_t n)
{
	RECT r;

	if (SendMessageW(t->header, HDM_GETITEMRECT, (WPARAM) n, (LPARAM) (&r)) == 0)
		panic("error getting Table column width");
	return r.right - r.left;
}

/* TODO:
http://blogs.msdn.com/b/oldnewthing/archive/2003/10/13/55279.aspx
http://blogs.msdn.com/b/oldnewthing/archive/2003/10/14/55286.aspx
we'll need to make sure that initial edge case works properly
(TODO get the linked article in the latter)
also implement retrack() as so, in the WM_MOUSEMOVE handler
*/
static void retrack(struct table *t)
{
	TRACKMOUSEEVENT tm;

	ZeroMemory(&tm, sizeof (TRACKMOUSEEVENT));
	tm.cbSize = sizeof (TRACKMOUSEEVENT);
	tm.dwFlags = TME_LEAVE;		// TODO TME_NONCLIENT as well?
	tm.hwndTrack = t->hwnd;
	if ((*tableTrackMouseEvent)(&tm) == 0)
		panic("error retracking Table mouse events");
}

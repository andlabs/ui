// 4 december 2014

typedef BOOL (*handlerfunc)(struct table *, UINT, WPARAM, LPARAM, LRESULT *);
#define HANDLER(name) static BOOL name(struct table *t, UINT uMsg, WPARAM wParam, LPARAM lParam, LRESULT *lResult)

static BOOL runHandlers(const handlerfunc list[], struct table *t, UINT uMsg, WPARAM wParam, LPARAM lParam, LRESULT *lResult)
{
	handlerfunc *p;

	for (p = list; *p != NULL; p++)
		if ((*(*p))(t, uMsg, wParam, lParam, lResult))
			return TRUE;
	return FALSE;
}

// memory allocation stuff
// each of these functions do an implicit ZeroMemory()
// these also make tableRealloc(NULL, ...) act like realloc(NULL, ...) (that is, same as tableAlloc(...)/malloc(...))
// we're using LocalAlloc() because:
// - whether the malloc() family supports the last-error code is undefined
// - the HeapAlloc() family DOES NOT support the last-error code; you're supposed to use Windows exceptions, and I can't find a clean way to do this with MinGW-w64 that doesn't rely on inline assembly or external libraries (unless they added __try/__except blocks)
// - there's no VirtualReAlloc() to complement VirtualAlloc() and I'm not sure if we can even get the original allocated size back out reliably to write it ourselves (http://blogs.msdn.com/b/oldnewthing/archive/2012/03/16/10283988.aspx)
// needless to say, TODO

static void *tableAlloc(size_t size, const char *panicMessage)
{
	HLOCAL out;

	out = LocalAlloc(LMEM_FIXED | LMEM_ZEROINIT, size);
	if (out == NULL)
		panic(panicMessage);
	return (void *) out;
}

static void *tableRealloc(void *p, size_t size, const char *panicMessage)
{
	HLOCAL out;

	if (p == NULL)
		return tableAlloc(size, panicMessage);
	out = LocalReAlloc((HLOCAL) p, size, LMEM_ZEROINIT);
	if (out == NULL)
		panic(panicMessage);
	return (void *) out;
}

static void tableFree(void *p, const char *panicMessage)
{
	if (LocalFree((HLOCAL) p) != NULL)
		panic(panicMessage);
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

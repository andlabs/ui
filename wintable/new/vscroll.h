// 9 december 2014

// forward declaration needed here
static void repositionHeader(struct table *);

static struct scrollParams vscrollParams(struct table *t)
{
	struct scrollParams p;

	ZeroMemory(&p, sizeof (struct scrollParams));
	p.pos = &(t->vscrollpos);
	p.pagesize = t->vpagesize;
	p.length = t->count;
	p.scale = rowht(t);
	p.post = NULL;
	return p;
}

static void vscrollto(struct table *t, intptr_t pos)
{
	struct scrollParams p;

	p = vscrollParams(t);
	scrollto(t, SB_VERT, &p, pos);
}

static void vscrollby(struct table *t, intptr_t delta)
{
	struct scrollParams p;

	p = vscrollParams(t);
	scrollby(t, SB_VERT, &p, delta);
}

static void vscroll(struct table *t, WPARAM wParam, LPARAM lParam)
{
	struct scrollParams p;

	p = vscrollParams(t);
	scroll(t, SB_VERT, &p, wParam, lParam);
}

// TODO WM_MOUSEWHEEL
HANDLER(vscrollHandler)
{
	if (uMsg != WM_VSCROLL)
		return FALSE;
	vscroll(t, wParam, lParam);
	*lResult = 0;
	return TRUE;
}

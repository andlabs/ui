// 9 december 2014

// forward declaration needed here
static void repositionHeader(struct table *);

static struct scrollParams hscrollParams(struct table *t)
{
	struct scrollParams p;

	ZeroMemory(&p, sizeof (struct scrollParams));
	p.pos = &(t->hscrollpos);
	p.pagesize = t->hpagesize;
	p.length = t->width;
	p.scale = 1;
	p.post = repositionHeader;
	p.wheelCarry = &(t->hwheelCarry);
	return p;
}

static void hscrollto(struct table *t, intptr_t pos)
{
	struct scrollParams p;

	p = hscrollParams(t);
	scrollto(t, SB_HORZ, &p, pos);
}

static void hscrollby(struct table *t, intptr_t delta)
{
	struct scrollParams p;

	p = hscrollParams(t);
	scrollby(t, SB_HORZ, &p, delta);
}

static void hscroll(struct table *t, WPARAM wParam, LPARAM lParam)
{
	struct scrollParams p;

	p = hscrollParams(t);
	scroll(t, SB_HORZ, &p, wParam, lParam);
}

// TODO find out if we can indicriminately check for WM_WHEELHSCROLL
HANDLER(hscrollHandler)
{
	if (uMsg != WM_HSCROLL)
		return FALSE;
	hscroll(t, wParam, lParam);
	*lResult = 0;
	return TRUE;
}

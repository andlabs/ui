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

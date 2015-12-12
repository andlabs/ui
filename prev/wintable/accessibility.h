// 24 december 2014

// implement MSAA-conformant accessibility
// we need to use MSAA because UI Automation is too new for us
// unfortunately, MSAA's documentation is... very poor. very ambiguous, very inconsistent (just run through this file's commit history and watch the TODO progression to see)
// resources:
// http://msdn.microsoft.com/en-us/library/ms971338.aspx
// http://msdn.microsoft.com/en-us/library/windows/desktop/cc307844.aspx
// http://msdn.microsoft.com/en-us/library/windows/desktop/cc307847.aspx
// http://blogs.msdn.com/b/saraford/archive/2004/08/20/which-controls-support-which-msaa-properties-and-how-these-controls-implement-msaa-properties.aspx
// http://msdn.microsoft.com/en-us/library/ms971325
// http://msdn.microsoft.com/en-us/library/windows/desktop/dd318017%28v=vs.85%29.aspx
// http://msdn.microsoft.com/en-us/library/windows/desktop/dd373624%28v=vs.85%29.aspx

// notes:
// - TODO figure out what to do about header
// - a row extends as far right as the right edge of the last cell in the row; anything to the right of that is treated as table space (just like with mouse selection)
// 	- this has the added effect that hit-testing can only ever return either the table or a cell, never a row
// - cells have no children; checkbox cells are themselves the accessible object
// 	- TODO if we ever add combobox columns, this will need to change somehow
// - only Table and Cell can have focus; only Row can have selection
// 	- TODO allow selecting a cell?

// TODOs:
// - make sure E_POINTER and RPC_E_DISCONNECTED are correct returns for IAccessible
// - return last error on newTableAcc() in all accessible functions
// - figure out what should be names and what should be values
// - figure out what to do about that header row
// - http://acccheck.codeplex.com/

// uncomment this to debug table linked list management
//#define TABLE_DEBUG_LINKEDLIST

// TODO get rid of this
typedef struct tableAccWhat tableAccWhat;

struct tableAccWhat {
	LONG role;
	intptr_t row;
	intptr_t column;
};

struct tableAcc {
	const IAccessibleVtbl *vtbl;
	ULONG refcount;
	struct table *t;
	IAccessible *std;
	tableAccWhat what;

	// the list of currently active accessibility objects is a doubly linked list
	struct tableAcc *prev;
	struct tableAcc *next;
};

#ifdef TABLE_DEBUG_LINKEDLIST
void list(struct table *t)
{
	struct tableAcc *ta;

	printf("\n");
	if (t->firstAcc == NULL) {
		printf("\tempty\n");
		return;
	}
	printf("\t-> ");
	for (ta = t->firstAcc; ta != NULL; ta = ta->next)
		printf("%p ", ta);
	printf("\n\t<- ");
	for (ta = t->firstAcc; ta->next != NULL; ta = ta->next)
		;
	for (; ta != NULL; ta = ta->prev)
		printf("%p ", ta);
	printf("\n");
}
#endif

// called after each allocation
static struct tableAcc *newTableAcc(struct table *t, LONG role, intptr_t row, intptr_t column);

// common validation for accessibility functions that take varChild
// also normalizes what as if varChild == CHILDID_SELF
static HRESULT normalizeWhat(struct tableAcc *ta, VARIANT varChild, tableAccWhat *what)
{
	LONG cid;

	if (varChild.vt != VT_I4)
		return E_INVALIDARG;
	cid = varChild.lVal;
	if (cid == CHILDID_SELF)
		return S_OK;
	cid--;
	if (cid < 0)
		return E_INVALIDARG;
	switch (what->role) {
	case ROLE_SYSTEM_TABLE:
		// TODO +1?
		if (cid >= ta->t->count)
			return E_INVALIDARG;
		what->role = ROLE_SYSTEM_ROW;
		what->row = (intptr_t) cid;
		what->column = -1;
		break;
	case ROLE_SYSTEM_ROW:
	case ROLE_SYSTEM_CELL:
		// TODO
		// TODO also figure out what to do if the current row/cell become invalid (rows being removed, etc.)
		break;
	}
	return S_OK;
}

#define TA ((struct tableAcc *) this)

static HRESULT STDMETHODCALLTYPE tableAccQueryInterface(IAccessible *this, REFIID riid, void **ppvObject)
{
	if (ppvObject == NULL)
		return E_POINTER;
	if (IsEqualIID(riid, &IID_IUnknown) ||
		IsEqualIID(riid, &IID_IDispatch) ||
		IsEqualIID(riid, &IID_IAccessible)) {
		IAccessible_AddRef(this);
		*ppvObject = (void *) this;
		return S_OK;
	}
	*ppvObject = NULL;
	return E_NOINTERFACE;
}

// TODO use InterlockedIncrement()/InterlockedDecrement() for these?

static ULONG STDMETHODCALLTYPE tableAccAddRef(IAccessible *this)
{
	TA->refcount++;
	// TODO correct?
	return TA->refcount;
}

static ULONG STDMETHODCALLTYPE tableAccRelease(IAccessible *this)
{
	TA->refcount--;
	if (TA->refcount == 0) {
		struct tableAcc *prev, *next;

#ifdef TABLE_DEBUG_LINKEDLIST
if (TA->t != NULL) { printf("before delete:"); list(TA->t); }
#endif
		if (TA->t != NULL && TA->t->firstAcc == TA)
			TA->t->firstAcc = TA->next;
		prev = TA->prev;
		next = TA->next;
		if (prev != NULL)
			prev->next = next;
		if (next != NULL)
			next->prev = prev;
#ifdef TABLE_DEBUG_LINKEDLIST
if (TA->t != NULL) { printf("after delete:"); list(TA->t); }
#endif
		if (TA->std != NULL)
			IAccessible_Release(TA->std);
		tableFree(TA, "error freeing Table accessibility object");
		return 0;
	}
	return TA->refcount;
}

// IDispatch
// TODO make sure relegating these to the standard accessibility object is harmless

static HRESULT STDMETHODCALLTYPE tableAccGetTypeInfoCount(IAccessible *this, UINT *pctinfo)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_GetTypeInfoCount(TA->std, pctinfo);
}

static HRESULT STDMETHODCALLTYPE tableAccGetTypeInfo(IAccessible *this, UINT iTInfo, LCID lcid, ITypeInfo **ppTInfo)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_GetTypeInfo(TA->std, iTInfo, lcid, ppTInfo);
}

static HRESULT STDMETHODCALLTYPE tableAccGetIDsOfNames(IAccessible *this, REFIID riid, LPOLESTR *rgszNames, UINT cNames, LCID lcid, DISPID *rgDispId)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_GetIDsOfNames(TA->std, riid, rgszNames, cNames, lcid, rgDispId);
}

static HRESULT STDMETHODCALLTYPE tableAccInvoke(IAccessible *this, DISPID dispIdMember, REFIID riid, LCID lcid, WORD wFlags, DISPPARAMS *pDispParams, VARIANT *pVarResult, EXCEPINFO *pExcepInfo, UINT *puArgErr)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_Invoke(TA->std, dispIdMember, riid, lcid, wFlags, pDispParams, pVarResult, pExcepInfo, puArgErr);
}

// IAccessible

static HRESULT STDMETHODCALLTYPE tableAccget_accParent(IAccessible *this, IDispatch **ppdispParent)
{
	if (ppdispParent == NULL)
		return E_POINTER;
	// TODO set ppdispParent to zero?
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	// TODO check if row/column is still valid
	switch (TA->what.role) {
	case ROLE_SYSTEM_TABLE:
		// defer to standard accessible object
		// TODO [EDGE CASE/POOR DOCUMENTATION?] https://msdn.microsoft.com/en-us/library/ms971325 says "Returns the IDispatch interface of the Table object."; isn't that just returning self?
		return IAccessible_get_accParent(TA->std, ppdispParent);
	case ROLE_SYSTEM_ROW:
		*ppdispParent = (IDispatch *) newTableAcc(TA->t, ROLE_SYSTEM_TABLE, -1, -1);
		return S_OK;
	case ROLE_SYSTEM_CELL:
		*ppdispParent = (IDispatch *) newTableAcc(TA->t, ROLE_SYSTEM_ROW, TA->what.row, -1);
		return S_OK;
	}
	// TODO actually do this right
	// TODO un-GetLastError() this
	panic("impossible blah blah blah TODO write this");
	return E_FAIL;
}

static HRESULT STDMETHODCALLTYPE tableAccget_accChildCount(IAccessible *this, long *pcountChildren)
{
	if (pcountChildren == NULL)
		return E_POINTER;
	// TODO set pcountChildren to zero?
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	switch (TA->what.role) {
	case ROLE_SYSTEM_TABLE:
		// TODO header row
		*pcountChildren = (long) (TA->t->count);
		return S_OK;
	case ROLE_SYSTEM_ROW:
		*pcountChildren = (long) (TA->t->nColumns);
		return S_OK;
	case ROLE_SYSTEM_CELL:
		*pcountChildren = 0;
		return S_OK;
	}
	// TODO actually do this right
	// TODO un-GetLastError() this
	panic("impossible blah blah blah TODO write this");
	return E_FAIL;
}

// TODO [EDGE CASE/NOT DOCUMENTED/CHECK SAMPLE] what happens if CHILDID_SELF is passed?
// TODO [EDGE CASE/NOT DOCUMENTED/CHECK SAMPLE] what SHOULD happen if an out of bounds ID is passed?
static HRESULT STDMETHODCALLTYPE tableAccget_accChild(IAccessible *this, VARIANT varChild, IDispatch **ppdispChild)
{
	LONG cid;

	if (ppdispChild == NULL)
		return E_POINTER;
	*ppdispChild = NULL;
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	if (varChild.vt != VT_I4)
		return E_INVALIDARG;
	cid = varChild.lVal;
	if (cid < 0)
		// TODO really?
		return E_INVALIDARG;
	if (cid == CHILDID_SELF)
		return E_FAIL;		// TODO
	cid--;
	switch (TA->what.role) {
	case ROLE_SYSTEM_TABLE:
		// TODO table header
		if (TA->t->count == 0)
			return S_FALSE;
		if (cid > TA->t->count - 1)
			// TODO really?
			return E_INVALIDARG;
		*ppdispChild = (IDispatch *) newTableAcc(TA->t, ROLE_SYSTEM_ROW, cid, -1);
		return S_OK;
	case ROLE_SYSTEM_ROW:
		// TODO verify that row is still valid
		if (TA->t->nColumns == 0)
			return S_FALSE;
		if (cid > TA->t->nColumns - 1)
			// TODO really?
			return E_INVALIDARG;
		*ppdispChild = (IDispatch *) newTableAcc(TA->t, ROLE_SYSTEM_CELL, TA->what.row, cid);
	case ROLE_SYSTEM_CELL:
		// TODO verify that row/column are still valid?
		return S_FALSE;
	}
	// TODO actually do this right
	// TODO un-GetLastError() this
	panic("impossible blah blah blah TODO write this");
	return E_FAIL;
}

static HRESULT STDMETHODCALLTYPE tableAccget_accName(IAccessible *this, VARIANT varChild, BSTR *pszName)
{
	HRESULT hr;
	tableAccWhat what;

	if (pszName == NULL)
		return E_POINTER;
	// TODO double-check that this must be set to zero
	*pszName = NULL;
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	switch (what.role) {
	case ROLE_SYSTEM_TABLE:
		// defer to standard accessible object
		return IAccessible_get_accName(TA->std, varChild, pszName);
	case ROLE_SYSTEM_ROW:
		// TODO
		return S_FALSE;
	case ROLE_SYSTEM_CELL:
		// TODO
		return S_FALSE;
	}
	// TODO actually do this right
	// TODO un-GetLastError() this
	panic("impossible blah blah blah TODO write this");
	return E_FAIL;
}

// this differs quite much from what is described at https://msdn.microsoft.com/en-us/library/ms971325
static HRESULT STDMETHODCALLTYPE tableAccget_accValue(IAccessible *this, VARIANT varChild, BSTR *pszValue)
{
	HRESULT hr;
	tableAccWhat what;
	WCHAR *text;

	if (pszValue == NULL)
		return E_POINTER;
	// TODO set pszValue to zero?
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	switch (what.role) {
	case ROLE_SYSTEM_TABLE:
		// TODO really?
		return IAccessible_get_accValue(TA->std, varChild, pszValue);
	case ROLE_SYSTEM_ROW:
		// TODO
		return DISP_E_MEMBERNOTFOUND;
	case ROLE_SYSTEM_CELL:
		switch (TA->t->columnTypes[what.column]) {
		case tableColumnText:
			text = getCellText(TA->t, what.row, what.column);
			// TODO check for error
			*pszValue = SysAllocString(text);
			returnCellData(TA->t, what.row, what.column, text);
			return S_OK;
		case tableColumnImage:
			// TODO
			return DISP_E_MEMBERNOTFOUND;
		case tableColumnCheckbox:
			// TODO!!!!!!
			return DISP_E_MEMBERNOTFOUND;
		}
	}
	// TODO actually do this right
	// TODO un-GetLastError() this
	panic("impossible blah blah blah TODO write this");
	return E_FAIL;
}

static HRESULT STDMETHODCALLTYPE tableAccget_accDescription(IAccessible *this, VARIANT varChild, BSTR *pszDescription)
{
	HRESULT hr;
	tableAccWhat what;

	if (pszDescription == NULL)
		return E_POINTER;
	*pszDescription = NULL;
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	// don't support descriptions anyway; do return the above errors just to be safe
	return DISP_E_MEMBERNOTFOUND;
}

static HRESULT STDMETHODCALLTYPE tableAccget_accRole(IAccessible *this, VARIANT varChild, VARIANT *pvarRole)
{
	HRESULT hr;
	tableAccWhat what;

	if (pvarRole == NULL)
		return E_POINTER;
	pvarRole->vt = VT_EMPTY;
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	pvarRole->vt = VT_I4;
	pvarRole->lVal = what.role;
	return S_OK;
}

// TODO reason about STATE_SYSTEM_INVISIBLE and STATE_SYSTEM_OFFSCREEN
static HRESULT STDMETHODCALLTYPE tableAccget_accState(IAccessible *this, VARIANT varChild, VARIANT *pvarState)
{
	HRESULT hr;
	tableAccWhat what;
	LONG state;

	if (pvarState == NULL)
		return E_POINTER;
	pvarState->vt = VT_EMPTY;
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;

	state = 0;
	switch (what.role) {
	case ROLE_SYSTEM_TABLE:
		hr = IAccessible_get_accState(TA->std, varChild, pvarState);
		if (hr != S_OK)
			return hr;
		// TODO make sure pvarState->vt == VT_I4 (what to return otherwise?)
		state |= pvarState->lVal;
		break;
	case ROLE_SYSTEM_ROW:
		state |= STATE_SYSTEM_SELECTABLE;
		if (TA->t->selectedRow == what.row)
			state |= STATE_SYSTEM_SELECTED;
		break;
	case ROLE_SYSTEM_CELL:
		if (TA->t->columnTypes[what.column] == tableColumnCheckbox) {
			// TODO is there no STATE_SYSTEM_CHECKABLE?
			if (isCheckboxChecked(TA->t, what.row, what.column))
				state |= STATE_SYSTEM_CHECKED;
		}
		state |= STATE_SYSTEM_FOCUSABLE;
		if (TA->t->selectedRow == what.row && TA->t->selectedColumn == what.column)
			state |= STATE_SYSTEM_FOCUSED;
		if (TA->t->columnTypes[what.column] != tableColumnCheckbox)
			state |= STATE_SYSTEM_READONLY;
		break;
	}
	pvarState->vt = VT_I4;
	pvarState->lVal = state;
	return S_OK;
}

static HRESULT STDMETHODCALLTYPE tableAccget_accHelp(IAccessible *this, VARIANT varChild, BSTR *pszHelp)
{
	HRESULT hr;
	tableAccWhat what;

	if (pszHelp == NULL)
		return E_POINTER;
	*pszHelp = NULL;
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	// don't support help anyway; do return the above errors just to be safe
	return DISP_E_MEMBERNOTFOUND;
}

// TODO Inspect.exe seems to ignore the DISP_E_MEMBERNOTFOUND and just tells us the help topic is the empty string; make sure this works right
static HRESULT STDMETHODCALLTYPE tableAccget_accHelpTopic(IAccessible *this, BSTR *pszHelpFile, VARIANT varChild, long *pidTopic)
{
	HRESULT hr;
	tableAccWhat what;

	if (pszHelpFile == NULL || pidTopic == NULL)
		return E_POINTER;
	// TODO set pszHelpFile and pidTopic to zero?
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	// don't support Windows Help (the super-old .hlp help files) topics anyway; do return the above errors just to be safe
	// TODO [EDGE CASE??] or should we defer back to the standard accessible object? get_accHelp() was explicitly documented as not being supported by the standard/common controls, but this one isn't...
	return DISP_E_MEMBERNOTFOUND;
}

static HRESULT STDMETHODCALLTYPE tableAccget_accKeyboardShortcut(IAccessible *this, VARIANT varChild, BSTR *pszKeyboardShortcut)
{
	HRESULT hr;
	tableAccWhat what;

	if (pszKeyboardShortcut == NULL)
		return E_POINTER;
	// TODO set pszKeyboardShortcut to zero?
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	// defer to the standard accessible object for the table itself in case a program assigns an access key somehow (adjacent label?); MSDN says to, anyway
	if (what.role == ROLE_SYSTEM_TABLE)
		return IAccessible_get_accKeyboardShortcut(TA->std, varChild, pszKeyboardShortcut);
	if (what.role == ROLE_SYSTEM_CELL)
		;	// TODO implement this for checkbox cells?
	return DISP_E_MEMBERNOTFOUND;
}

// TODO TEST THIS
// TODO [EDGE CASE??] no parents?
static HRESULT STDMETHODCALLTYPE tableAccget_accFocus(IAccessible *this, VARIANT *pvarChild)
{
	HRESULT hr;

	if (pvarChild == NULL)
		return E_POINTER;
	// TODO set pvarChild to empty?
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	// TODO verify that TA is still pointing to a valid row/column

	// first see if the control has the focus
	// this is why a standard accessible object is needed on all accessible objects
	hr = IAccessible_get_accFocus(TA->std, pvarChild);
	// check the pvarChild type instead of hr
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dd318479%28v=vs.85%29.aspx does this
	// TODO [EDGE CASE] figure out why
	if (pvarChild->vt != VT_I4)
		return hr;

	switch (TA->what.role) {
	case ROLE_SYSTEM_TABLE:
		if (TA->t->selectedRow != -1 && TA->t->selectedColumn != -1)
			goto selectedCell;
		goto self;
	case ROLE_SYSTEM_ROW:
		if (TA->t->selectedRow != TA->what.row)
			goto nothing;
		goto selectedCell;
	case ROLE_SYSTEM_CELL:
		if (TA->t->selectedRow != TA->what.row)
			goto nothing;
		if (TA->t->selectedColumn != TA->what.column)
			goto nothing;
		goto self;
	}
	// TODO actually do this right
	// TODO un-GetLastError() this
	panic("impossible blah blah blah TODO write this");
	return E_FAIL;

nothing:
	pvarChild->vt = VT_EMPTY;
	// TODO really this one?
	return S_FALSE;

self:
	pvarChild->vt = VT_I4;
	pvarChild->lVal = CHILDID_SELF;
	return S_OK;

selectedCell:
	pvarChild->vt = VT_I4;
	pvarChild->pdispVal = (IDispatch *) newTableAcc(TA->t, ROLE_SYSTEM_CELL, TA->what.row, TA->what.column);
	return S_OK;
}

// note: https://msdn.microsoft.com/en-us/library/ms971325 is geared toward cell-based selection
// we have row-based selection, so only Tables implement this method, and they return a row
static HRESULT STDMETHODCALLTYPE tableAccget_accSelection(IAccessible *this, VARIANT *pvarChildren)
{
	if (pvarChildren == NULL)
		return E_POINTER;
	// TOOD set pvarChildren to VT_EMPTY?
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	if (TA->what.role != ROLE_SYSTEM_TABLE)
		// TODO [EDGE CASE] implement this for row anyway? how?
		return DISP_E_MEMBERNOTFOUND;
	if (TA->t->selectedRow == -1) {
		pvarChildren->vt = VT_EMPTY;
		return S_OK;
	}
	pvarChildren->vt = VT_DISPATCH;
	pvarChildren->pdispVal = (IDispatch *) newTableAcc(TA->t, ROLE_SYSTEM_ROW, TA->t->selectedRow, -1);
	return S_OK;
}

static HRESULT STDMETHODCALLTYPE tableAccget_accDefaultAction(IAccessible *this, VARIANT varChild, BSTR *pszDefaultAction)
{
	HRESULT hr;
	tableAccWhat what;

	if (pszDefaultAction == NULL)
		return E_POINTER;
	*pszDefaultAction = NULL;
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	if (what.role == ROLE_SYSTEM_CELL)
		;	// TODO implement this for checkbox cells?
	return DISP_E_MEMBERNOTFOUND;
}

// TODO should this method result in an event?
// TODO [EDGE CASE] how do we deselect? in the table or in the row? wouldn't this go against multiple selection?
// TODO require cell rows to be selected before focusing?
static HRESULT STDMETHODCALLTYPE tableAccaccSelect(IAccessible *this, long flagsSelect, VARIANT varChild)
{
	HRESULT hr;
	tableAccWhat what;

	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	if (what.role == ROLE_SYSTEM_TABLE)		// defer to the standard accessible object
		return IAccessible_accSelect(TA->std, flagsSelect, varChild);
	// reject flags that are only applicable to multiple selection
	if ((flagsSelect & (SELFLAG_EXTENDSELECTION | SELFLAG_ADDSELECTION | SELFLAG_REMOVESELECTION)) != 0)
		return E_INVALIDARG;
	// and do nothing if a no-op
	if (flagsSelect == SELFLAG_NONE)
		return S_FALSE;
	// TODO cast ~ expressions to the correct type
	switch (what.role) {
	case ROLE_SYSTEM_ROW:
		// reject any other flag
		if ((flagsSelect & (~SELFLAG_TAKESELECTION)) != 0)
			return E_INVALIDARG;
		if ((flagsSelect & SELFLAG_TAKESELECTION) != 0) {
			if (TA->t->nColumns == 0)		// can't select
				return S_FALSE;
			// if no column selected, select first (congruent to behavior of certain keyboard events)
			// TODO handle GetLastError()
			if (TA->t->selectedColumn == -1)
				doselect(TA->t, what.row, TA->t->selectedColumn);
			else
				doselect(TA->t, what.row, 0);
			return S_OK;
		}
		return S_FALSE;
	case ROLE_SYSTEM_CELL:
		// reject any other flag
		if ((flagsSelect & (~SELFLAG_TAKEFOCUS)) != 0)
			return E_INVALIDARG;
		if ((flagsSelect & SELFLAG_TAKEFOCUS) != 0) {
			doselect(TA->t, what.row, what.column);
			return S_OK;
		}
		return S_FALSE;
	}
	// TODO actually do this right
	// TODO un-GetLastError() this
	panic("impossible blah blah blah TODO write this");
	return E_FAIL;
}

static HRESULT STDMETHODCALLTYPE tableAccaccLocation(IAccessible *this, long *pxLeft, long *pyTop, long *pcxWidth, long *pcyHeight, VARIANT varChild)
{
	HRESULT hr;
	tableAccWhat what;
	RECT r;
	POINT pt;
	struct rowcol rc;

	if (pxLeft == NULL || pyTop == NULL || pcxWidth == NULL || pcyHeight == NULL)
		return E_POINTER;
	// TODO set the out parameters to zero?
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	switch (what.role) {
	case ROLE_SYSTEM_TABLE:
		return IAccessible_accLocation(TA->std, pxLeft, pyTop, pcxWidth, pcyHeight, varChild);
	case ROLE_SYSTEM_ROW:
		// TODO actually write this
		return E_FAIL;
	case ROLE_SYSTEM_CELL:
		rc.row = what.row;
		rc.column = what.column;
		if (!rowColumnToClientRect(TA->t, rc, &r)) {
			// TODO [EDGE CASE] what do we do here?
			// TODO we have to return something indicating that the object is off-screen
		}
		// TODO [EDGE CASE] intersect with client rect?
		break;
	}
	pt.x = r.left;
	pt.y = r.top;
	if (ClientToScreen(TA->t->hwnd, &pt) == 0)
		return HRESULT_FROM_WIN32(GetLastError());
	*pxLeft = pt.x;
	*pyTop = pt.y;
	*pcxWidth = r.right - r.left;
	*pcyHeight = r.bottom - r.top;
	return S_OK;
}

static HRESULT STDMETHODCALLTYPE tableAccaccNavigate(IAccessible *this, long navDir, VARIANT varStart, VARIANT *pvarEndUpAt)
{
	HRESULT hr;
	tableAccWhat what;
	intptr_t row = -1;
	intptr_t column = -1;

	if (pvarEndUpAt == NULL)
		return E_POINTER;
	// TODO set pvarEndUpAt to an invalid value?
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varStart, &what);
	if (hr != S_OK)
		return hr;
	switch (what.role) {
	case ROLE_SYSTEM_TABLE:
		switch (navDir) {
		case NAVDIR_FIRSTCHILD:
			// TODO header row
			if (TA->t->count == 0)
				goto nowhere;
			row = 0;
			goto specificRow;
		case NAVDIR_LASTCHILD:
			// TODO header row
			if (TA->t->count == 0)
				goto nowhere;
			row = TA->t->count - 1;
			goto specificRow;
		}
		// otherwise, defer to the standard accessible object
		return IAccessible_accNavigate(TA->std, navDir, varStart, pvarEndUpAt);
	case ROLE_SYSTEM_ROW:
		row = what.row;
		switch (navDir) {
		case NAVDIR_UP:
		case NAVDIR_PREVIOUS:
			if (row == 0)		// can't go up
				goto nowhere;
			row--;
			// row should still be valid because normalizeWhat() returns an error if the current row is no longer valid, and if that's valid, the row above it should also be valid
			goto specificRow;
		case NAVDIR_DOWN:
		case NAVDIR_NEXT:
			if (row == TA->t->count - 1)		// can't go down
				goto nowhere;
			row++;
			// row should still be valid by the above conjecture
			goto specificRow;
		case NAVDIR_LEFT:
		case NAVDIR_RIGHT:
			goto nowhere;
// TODO this doesn't actually exist yet https://msdn.microsoft.com/en-us/library/ms971325 talks about it
//		case NAVDIR_PARENT:
//			goto tableItself;
		case NAVDIR_FIRSTCHILD:
			if (TA->t->nColumns == 0)
				goto nowhere;
			column = 0;
			goto specificCell;
		case NAVDIR_LASTCHILD:
			if (TA->t->nColumns == 0)
				goto nowhere;
			column = TA->t->nColumns - 1;
			goto specificCell;
		}
		// TODO differentiate between unsupported navigation directions and invalid navigation directions
		goto nowhere;
	case ROLE_SYSTEM_CELL:
		row = what.row;
		column = what.column;
		switch (navDir) {
		case NAVDIR_UP:
			if (row == 0)		// can't go up
				goto nowhere;
			row--;
			goto specificCell;
		case NAVDIR_DOWN:
			if (row == TA->t->count - 1)		// can't go down
				goto nowhere;
			row++;
			goto specificCell;
		case NAVDIR_LEFT:
		case NAVDIR_PREVIOUS:
			if (column == 0)		// can't go left
				goto nowhere;
			column--;
			goto specificCell;
		case NAVDIR_RIGHT:
		case NAVDIR_NEXT:
			if (column == TA->t->nColumns - 1)		// can't go right
				goto nowhere;
			column++;
			goto specificCell;
// TODO this doesn't actually exist yet https://msdn.microsoft.com/en-us/library/ms971325 talks about it
//		case NAVDIR_PARENT:
//			goto specificRow;
		case NAVDIR_FIRSTCHILD:
		case NAVDIR_LASTCHILD:
			goto nowhere;
		}
		// TODO differentiate between unsupported navigation directions and invalid navigation directions
		goto nowhere;
	}
	// TODO actually do this right
	// TODO un-GetLastError() this
	panic("impossible blah blah blah TODO write this");
	return E_FAIL;

nowhere:
	pvarEndUpAt->vt = VT_EMPTY;
	return S_FALSE;

tableItself:
	pvarEndUpAt->vt = VT_DISPATCH;
	pvarEndUpAt->pdispVal = (IDispatch *) newTableAcc(TA->t, ROLE_SYSTEM_TABLE, -1, -1);
	return S_OK;

specificRow:
	pvarEndUpAt->vt = VT_DISPATCH;
	pvarEndUpAt->pdispVal = (IDispatch *) newTableAcc(TA->t, ROLE_SYSTEM_ROW, row, -1);
	return S_OK;

specificCell:
	pvarEndUpAt->vt = VT_DISPATCH;
	pvarEndUpAt->pdispVal = (IDispatch *) newTableAcc(TA->t, ROLE_SYSTEM_CELL, row, column);
	return S_OK;
}

// TODO [EDGE CASE??] should this ever return parents?
static HRESULT STDMETHODCALLTYPE tableAccaccHitTest(IAccessible *this, long xLeft, long yTop, VARIANT *pvarChild)
{
	POINT pt;
	struct rowcol rc;
	RECT r;

	if (pvarChild == NULL)
		return E_POINTER;
	// TODO set pvarChild to an invalid value?
	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;

	pt.x = xLeft;
	pt.y = yTop;
	if (ScreenToClient(TA->t->hwnd, &pt) == 0)
		return HRESULT_FROM_WIN32(GetLastError());

	// first see if the point is even IN the table
	if (GetClientRect(TA->t->hwnd, &r) == 0)
		return HRESULT_FROM_WIN32(GetLastError());
	r.top += TA->t->headerHeight;
	if (PtInRect(&r, pt) == 0)
		goto outside;

	// now see if we're in a cell or in the table
	// TODO also handle GetLastError() here
	rc = clientCoordToRowColumn(TA->t, pt);
	switch (TA->what.role) {
	case ROLE_SYSTEM_TABLE:
		// either the table or the cell
		if (rc.row == -1 || rc.column == -1)
			goto self;
		goto specificCell;
	case ROLE_SYSTEM_ROW:
		// a specific cell, but only if in the same row
		// TODO actually do we really need these spurious rc.column ==/!= -1 checks?
		if (rc.row == TA->what.row) {
			if (rc.column == -1)
				// TODO de-GetLastError() this
				panic("impossible situation TODO write this");
			goto specificCell;
		}
		goto outside;
	case ROLE_SYSTEM_CELL:
		if (rc.row == TA->what.row && rc.column == TA->what.column)
			goto self;
		goto outside;
	}
	// TODO actually do this right
	// TODO un-GetLastError() this
	panic("impossible blah blah blah TODO write this");
	return E_FAIL;

outside:
	pvarChild->vt = VT_EMPTY;
	return S_FALSE;

self:
	pvarChild->vt = VT_I4;
	pvarChild->lVal = CHILDID_SELF;
	return S_OK;

specificCell:
	pvarChild->vt = VT_DISPATCH;
	// TODO GetLastError() here too
	pvarChild->pdispVal = (IDispatch *) newTableAcc(TA->t, ROLE_SYSTEM_CELL, rc.row, rc.column);
	return S_OK;
}

static HRESULT STDMETHODCALLTYPE tableAccaccDoDefaultAction(IAccessible *this, VARIANT varChild)
{
	HRESULT hr;
	tableAccWhat what;

	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	if (what.role == ROLE_SYSTEM_CELL)
		;	// TODO implement this for checkbox cells?
	return DISP_E_MEMBERNOTFOUND;
}

// inconsistencies, inconsistencies
// https://msdn.microsoft.com/en-us/library/windows/desktop/dd318491%28v=vs.85%29.aspx says to just return E_NOTIMPL and not even bother with an implementation; in fact it doesn't even *have* the documentation anymore
// http://blogs.msdn.com/b/saraford/archive/2004/08/20/which-controls-support-which-msaa-properties-and-how-these-controls-implement-msaa-properties.aspx says never to return E_NOTIMPL from an IAccessible method (but it also discounts RPC_E_DISCONNECTED (not explicitly), so I'm guessing this is a much older document)
// let's just do what our put_accValue() does and do full validation, then just return DISP_E_MEMBERNOTFOUND
// I really hope UI Automation isn't so ambiguous and inconsistent... too bad I'm still choosing to support Windows XP while its market share (compared to *every other OS ever*) is still as large as it is
static HRESULT STDMETHODCALLTYPE tableAccput_accName(IAccessible *this, VARIANT varChild, BSTR szName)
{
	HRESULT hr;
	tableAccWhat what;

	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	// don't support setting values anyway; do return the above errors just to be safe
	return DISP_E_MEMBERNOTFOUND;
}

static HRESULT STDMETHODCALLTYPE tableAccput_accValue(IAccessible *this, VARIANT varChild, BSTR szValue)
{
	HRESULT hr;
	tableAccWhat what;

	if (TA->t == NULL || TA->std == NULL)
		return RPC_E_DISCONNECTED;
	what = TA->what;
	hr = normalizeWhat(TA, varChild, &what);
	if (hr != S_OK)
		return hr;
	// don't support setting values anyway; do return the above errors just to be safe
	// TODO defer ROW_SYSTEM_TABLE to the standard accessible object?
	// TODO implement for checkboxes?
	return DISP_E_MEMBERNOTFOUND;
}

static const IAccessibleVtbl tableAccVtbl = {
	.QueryInterface = tableAccQueryInterface,
	.AddRef = tableAccAddRef,
	.Release = tableAccRelease,
	.GetTypeInfoCount = tableAccGetTypeInfoCount,
	.GetTypeInfo = tableAccGetTypeInfo,
	.GetIDsOfNames = tableAccGetIDsOfNames,
	.Invoke = tableAccInvoke,
	.get_accParent = tableAccget_accParent,
	.get_accChildCount = tableAccget_accChildCount,
	.get_accChild = tableAccget_accChild,
	.get_accName = tableAccget_accName,
	.get_accValue = tableAccget_accValue,
	.get_accDescription = tableAccget_accDescription,
	.get_accRole = tableAccget_accRole,
	.get_accState = tableAccget_accState,
	.get_accHelp = tableAccget_accHelp,
	.get_accHelpTopic = tableAccget_accHelpTopic,
	.get_accKeyboardShortcut = tableAccget_accKeyboardShortcut,
	.get_accFocus = tableAccget_accFocus,
	.get_accSelection = tableAccget_accSelection,
	.get_accDefaultAction = tableAccget_accDefaultAction,
	.accSelect = tableAccaccSelect,
	.accLocation = tableAccaccLocation,
	.accNavigate = tableAccaccNavigate,
	.accHitTest = tableAccaccHitTest,
	.accDoDefaultAction = tableAccaccDoDefaultAction,
	.put_accName = tableAccput_accName,
	.put_accValue = tableAccput_accValue,
};

static struct tableAcc *newTableAcc(struct table *t, LONG role, intptr_t row, intptr_t column)
{
	struct tableAcc *ta;
	HRESULT hr;
	IAccessible *std;

	ta = (struct tableAcc *) tableAlloc(sizeof (struct tableAcc), "error creating Table accessibility object");
	ta->vtbl = &tableAccVtbl;
	ta->refcount = 1;
	ta->t = t;
	// TODO adjust last argument
	hr = CreateStdAccessibleObject(t->hwnd, OBJID_CLIENT, &IID_IAccessible, (void **) (&std));
	if (hr != S_OK)
		// TODO panichresult
		panic("error creating standard accessible object for Table");
	ta->std = std;
	ta->what.role = role;
	ta->what.row = row;
	ta->what.column = column;

#ifdef TABLE_DEBUG_LINKEDLIST
printf("before add:"); list(t);
#endif
	ta->next = t->firstAcc;
	if (t->firstAcc != NULL)
		t->firstAcc->prev = ta;
	t->firstAcc = ta;
#ifdef TABLE_DEBUG_LINKEDLIST
printf("after add:"); list(t);
#endif

	return ta;
}

static void invalidateTableAccs(struct table *t)
{
	struct tableAcc *ta;

	for (ta = t->firstAcc; ta != NULL; ta = ta->next) {
		ta->t = NULL;
		IAccessible_Release(ta->std);
		ta->std = NULL;
	}
	t->firstAcc = NULL;
	NotifyWinEvent(EVENT_OBJECT_DESTROY, t->hwnd, OBJID_CLIENT, CHILDID_SELF);
}

HANDLER(accessibilityHandler)
{
	struct tableAcc *ta;

	if (uMsg != WM_GETOBJECT)
		return FALSE;
	// OBJID_CLIENT evaluates to an expression of type LONG
	// the documentation for WM_GETOBJECT says to cast "it" to a DWORD before comparing
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dd373624%28v=vs.85%29.aspx casts them both to DWORDs; let's do that
	// its two siblings only cast lParam, resulting in an erroneous DWORD to LONG comparison
	// The Old New Thing book does not cast anything
	// Microsoft's MSAA sample casts lParam to LONG instead!
	// (As you can probably tell, the biggest problem with MSAA is that its documentation is ambiguous and/or self-contradictory...)
	if (((DWORD) lParam) != ((DWORD) OBJID_CLIENT))
		return FALSE;
	ta = newTableAcc(t, ROLE_SYSTEM_TABLE, -1, -1);
	*lResult = LresultFromObject(&IID_IAccessible, wParam, (LPUNKNOWN) (ta));
	// TODO check *lResult
	// TODO adjust pointer
	IAccessible_Release((IAccessible *) ta);
	return TRUE;
}

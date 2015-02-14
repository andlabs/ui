// 24 december 2014

// TODOs:
// - make sure E_POINTER and RPC_E_DISCONNECTED are correct returns for IAccessible

// uncomment this to debug table linked list management
#define TABLE_DEBUG_LINKEDLIST

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
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accParent(TA->std, ppdispParent);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accChildCount(IAccessible *this, long *pcountChildren)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
//TODO
if (pcountChildren == NULL)
return E_POINTER;
*pcountChildren = 0;
return S_OK;
	return IAccessible_get_accChildCount(TA->std, pcountChildren);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accChild(IAccessible *this, VARIANT varChild, IDispatch **ppdispChild)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accChild(TA->std, varChild, ppdispChild);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accName(IAccessible *this, VARIANT varChild, BSTR *pszName)
{
printf("get_accName() t %p std %p\n", TA->t, TA->std);
	if (TA->t == NULL || TA->std == NULL) {
printf("returning RPC_E_DISCONNECTED\n");
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
printf("running main function\n");
//TODO
if (pszName == NULL)
return E_POINTER;
*pszName = SysAllocString(L"accessible table");
return S_OK;
	return IAccessible_get_accName(TA->std, varChild, pszName);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accValue(IAccessible *this, VARIANT varChild, BSTR *pszValue)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accValue(TA->std, varChild, pszValue);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accDescription(IAccessible *this, VARIANT varChild, BSTR *pszDescription)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accDescription(TA->std, varChild, pszDescription);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accRole(IAccessible *this, VARIANT varChild, VARIANT *pvarRole)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
//TODO
if (pvarRole == NULL)
return E_POINTER;
pvarRole->vt = VT_I4;
pvarRole->lVal = TA->what.role;
return S_OK;
	return IAccessible_get_accRole(TA->std, varChild, pvarRole);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accState(IAccessible *this, VARIANT varChild, VARIANT *pvarState)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accState(TA->std, varChild, pvarState);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accHelp(IAccessible *this, VARIANT varChild, BSTR *pszHelp)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accHelp(TA->std, varChild, pszHelp);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accHelpTopic(IAccessible *this, BSTR *pszHelpFile, VARIANT varChild, long *pidTopic)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accHelpTopic(TA->std, pszHelpFile, varChild, pidTopic);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accKeyboardShortcut(IAccessible *this, VARIANT varChild, BSTR *pszKeyboardShortcut)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accKeyboardShortcut(TA->std, varChild, pszKeyboardShortcut);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accFocus(IAccessible *this, VARIANT *pvarChild)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accFocus(TA->std, pvarChild);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accSelection(IAccessible *this, VARIANT *pvarChildren)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accSelection(TA->std, pvarChildren);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accDefaultAction(IAccessible *this, VARIANT varChild, BSTR *pszDefaultAction)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_get_accDefaultAction(TA->std, varChild, pszDefaultAction);
}

static HRESULT STDMETHODCALLTYPE tableAccaccSelect(IAccessible *this, long flagsSelect, VARIANT varChild)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_accSelect(TA->std, flagsSelect, varChild);
}

static HRESULT STDMETHODCALLTYPE tableAccaccLocation(IAccessible *this, long *pxLeft, long *pyTop, long *pcxWidth, long *pcyHeight, VARIANT varChild)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_accLocation(TA->std, pxLeft, pyTop, pcxWidth, pcyHeight, varChild);
}

static HRESULT STDMETHODCALLTYPE tableAccaccNavigate(IAccessible *this, long navDir, VARIANT varStart, VARIANT *pvarEndUpAt)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_accNavigate(TA->std, navDir, varStart, pvarEndUpAt);
}

static HRESULT STDMETHODCALLTYPE tableAccaccHitTest(IAccessible *this, long xLeft, long yTop, VARIANT *pvarChild)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_accHitTest(TA->std, xLeft, yTop, pvarChild);
}

static HRESULT STDMETHODCALLTYPE tableAccaccDoDefaultAction(IAccessible *this, VARIANT varChild)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_accDoDefaultAction(TA->std, varChild);
}

static HRESULT STDMETHODCALLTYPE tableAccput_accName(IAccessible *this, VARIANT varChild, BSTR szName)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_put_accName(TA->std, varChild, szName);
}

static HRESULT STDMETHODCALLTYPE tableAccput_accValue(IAccessible *this, VARIANT varChild, BSTR szValue)
{
	if (TA->t == NULL || TA->std == NULL) {
		// TODO set values on error
		return RPC_E_DISCONNECTED;
	}
	return IAccessible_put_accValue(TA->std, varChild, szValue);
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
printf("new ta %p\n", ta);
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
printf("creating ta\n");
	ta = newTableAcc(t, ROLE_SYSTEM_TABLE, -1, -1);
printf("ta %p\n", ta);
	*lResult = LresultFromObject(&IID_IAccessible, wParam, (LPUNKNOWN) (ta));
printf("lResult %I32d\n", *lResult);
	// TODO check *lResult
	// TODO adjust pointer
	IAccessible_Release((IAccessible *) ta);
	return TRUE;
}

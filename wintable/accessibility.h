// 24 december 2014

struct tableAcc {
	const IAccessibleVtbl *vtbl;
	ULONG refcount;
	struct table *t;
	IAccessible *std;

	LONG role;
	intptr_t row;
	intptr_t column;
};

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
	return IAccessible_GetTypeInfoCount(TA->std, pctinfo);
}

static HRESULT STDMETHODCALLTYPE tableAccGetTypeInfo(IAccessible *this, UINT iTInfo, LCID lcid, ITypeInfo **ppTInfo)
{
	return IAccessible_GetTypeInfo(TA->std, iTInfo, lcid, ppTInfo);
}

static HRESULT STDMETHODCALLTYPE tableAccGetIDsOfNames(IAccessible *this, REFIID riid, LPOLESTR *rgszNames, UINT cNames, LCID lcid, DISPID *rgDispId)
{
	return IAccessible_GetIDsOfNames(TA->std, riid, rgszNames, cNames, lcid, rgDispId);
}

static HRESULT STDMETHODCALLTYPE tableAccInvoke(IAccessible *this, DISPID dispIdMember, REFIID riid, LCID lcid, WORD wFlags, DISPPARAMS *pDispParams, VARIANT *pVarResult, EXCEPINFO *pExcepInfo, UINT *puArgErr)
{
	return IAccessible_Invoke(TA->std, dispIdMember, riid, lcid, wFlags, pDispParams, pVarResult, pExcepInfo, puArgErr);
}

// IAccessible

static HRESULT STDMETHODCALLTYPE tableAccget_accParent(IAccessible *this, IDispatch **ppdispParent)
{
	return IAccessible_get_accParent(TA->std, ppdispParent);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accChildCount(IAccessible *this, long *pcountChildren)
{
	return IAccessible_get_accChildCount(TA->std, pcountChildren);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accChild(IAccessible *this, VARIANT varChild, IDispatch **ppdispChild)
{
	return IAccessible_get_accChild(TA->std, varChild, ppdispChild);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accName(IAccessible *this, VARIANT varChild, BSTR *pszName)
{
	return IAccessible_get_accName(TA->std, varChild, pszName);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accValue(IAccessible *this, VARIANT varChild, BSTR *pszValue)
{
	return IAccessible_get_accValue(TA->std, varChild, pszValue);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accDescription(IAccessible *this, VARIANT varChild, BSTR *pszDescription)
{
	return IAccessible_get_accDescription(TA->std, varChild, pszDescription);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accRole(IAccessible *this, VARIANT varChild, VARIANT *pvarRole)
{
	return IAccessible_get_accRole(TA->std, varChild, pvarRole);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accState(IAccessible *this, VARIANT varChild, VARIANT *pvarState)
{
	return IAccessible_get_accState(TA->std, varChild, pvarState);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accHelp(IAccessible *this, VARIANT varChild, BSTR *pszHelp)
{
	return IAccessible_get_accHelp(TA->std, varChild, pszHelp);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accHelpTopic(IAccessible *this, BSTR *pszHelpFile, VARIANT varChild, long *pidTopic)
{
	return IAccessible_get_accHelpTopic(TA->std, pszHelpFile, varChild, pidTopic);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accKeyboardShortcut(IAccessible *this, VARIANT varChild, BSTR *pszKeyboardShortcut)
{
	return IAccessible_get_accKeyboardShortcut(TA->std, varChild, pszKeyboardShortcut);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accFocus(IAccessible *this, VARIANT *pvarChild)
{
	return IAccessible_get_accFocus(TA->std, pvarChild);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accSelection(IAccessible *this, VARIANT *pvarChildren)
{
	return IAccessible_get_accSelection(TA->std, pvarChildren);
}

static HRESULT STDMETHODCALLTYPE tableAccget_accDefaultAction(IAccessible *this, VARIANT varChild, BSTR *pszDefaultAction)
{
	return IAccessible_get_accDefaultAction(TA->std, varChild, pszDefaultAction);
}

static HRESULT STDMETHODCALLTYPE tableAccaccSelect(IAccessible *this, long flagsSelect, VARIANT varChild)
{
	return IAccessible_accSelect(TA->std, flagsSelect, varChild);
}

static HRESULT STDMETHODCALLTYPE tableAccaccLocation(IAccessible *this, long *pxLeft, long *pyTop, long *pcxWidth, long *pcyHeight, VARIANT varChild)
{
	return IAccessible_accLocation(TA->std, pxLeft, pyTop, pcxWidth, pcyHeight, varChild);
}

static HRESULT STDMETHODCALLTYPE tableAccaccNavigate(IAccessible *this, long navDir, VARIANT varStart, VARIANT *pvarEndUpAt)
{
	return IAccessible_accNavigate(TA->std, navDir, varStart, pvarEndUpAt);
}

static HRESULT STDMETHODCALLTYPE tableAccaccHitTest(IAccessible *this, long xLeft, long yTop, VARIANT *pvarChild)
{
	return IAccessible_accHitTest(TA->std, xLeft, yTop, pvarChild);
}

static HRESULT STDMETHODCALLTYPE tableAccaccDoDefaultAction(IAccessible *this, VARIANT varChild)
{
	return IAccessible_accDoDefaultAction(TA->std, varChild);
}

static HRESULT STDMETHODCALLTYPE tableAccput_accName(IAccessible *this, VARIANT varChild, BSTR szName)
{
	return IAccessible_put_accName(TA->std, varChild, szName);
}

static HRESULT STDMETHODCALLTYPE tableAccput_accValue(IAccessible *this, VARIANT varChild, BSTR szValue)
{
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

static struct tableAcc *newTableAcc(struct table *t)
{
	struct tableAcc *ta;
	HRESULT hr;
	IAccessible *std;

	ta = (struct tableAcc *) tableAlloc(sizeof (struct tableAcc), "error creating Table accessibility object");
	ta->vtbl = &tableAccVtbl;
	// TODO
	IAccessible_AddRef((IAccessible *) ta);
	ta->t = t;
	hr = CreateStdAccessibleObject(t->hwnd, OBJID_CLIENT, &IID_IAccessible, (void *) (&std));
	if (hr != S_OK)
		// TODO panichresult
		panic("error creating standard accessible object for Table");
	ta->std = std;
	ta->role = ROLE_SYSTEM_TABLE;
	ta->row = -1;
	ta->column = -1;
	return ta;
}

static void freeTableAcc(struct tableAcc *ta)
{
	ta->t = NULL;
	// TODO
	IAccessible_Release((IAccessible *) ta);
}

HANDLER(accessibilityHandler)
{
	if (uMsg != WM_GETOBJECT)
		return FALSE;
	// OBJID_CLIENT evaluates to an expression of type LONG
	// the documentation for WM_GETOBJECT says to cast "it" to a DWORD before comparing
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dd373624%28v=vs.85%29.aspx casts them both to DWORDs; let's do that
	// its two siblings only cast lParam, resulting in an erroneous DWORD to LONG comparison
	// The Old New Thing book does not cast anything
	// Microsoft's MSAA sample casts lParam to LONG instead!
	if (((DWORD) lParam) != ((DWORD) OBJID_CLIENT))
		return FALSE;
	*lResult = LresultFromObject(&IID_IAccessible, wParam, (LPUNKNOWN) (t->ta));
	// TODO check *lResult
	return TRUE;
}

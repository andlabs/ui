// 24 december 2014

struct tableAcc {
	IAccessibleVtbl *vtbl;
	ULONG refcount;
	struct table *t;
};

static HRESULT STDMETHODCALLTYPE tableAccQueryInterface(IUnknown *this, REFIID riid, void **ppvObject)
{
	if (ppvObject == NULL)
		return E_POINTER;
	if (IsEqualIID(riid, IID_IUnknown) ||
0)//		IsEqualIID(riid, IID_IDispatch)
{//		IsEqualIID(riid, IID_IAccessible) {
		*ppvObject = (void *) this;
		return S_OK;
	}
	*ppvObject = NULL;
	return E_NOINTERFACE;
}

#define TA ((struct tableAcc *) this)

// TODO use InterlockedIncrement()/InterlockedDecrement() for these?

static ULONG STDMETHODCALLTYPE tableAccAddRef(IUnknown *this)
{
	TA->refcount++;
	return TA->refcount;
}

static ULONG STDMETHODCALLTYPE tableAccRelease(IUnknown *this)
{
	TA->refcount--;
	if (TA->refcount == 0) {
		tableFree(TA, "error freeing Table accessibility object");
		return 0;
	}
	return TA->refcount;
}

static const IAccessibleVtbl tableAccVtbl = {
	.QueryInterface = tableAccQueryInterface,
	.AddRef = tableAccAddRef,
	.Release = tableAccRelease,
};

static struct tableAcc *newTableAcc(struct table *t)
{
	struct tableAcc *ta;

	ta = (struct tableAcc *) tableAlloc(sizeof (struct tableAcc), "error creating Table accessibility object");
	ta->vtbl = &tableAccVtbl;
	ta->vtbl->AddRef(vtbl);
	ta->t = t;
	return ta;
}

static void freeTableAcc(struct tableAcc *ta)
{
	ta->t = NULL;
	ta->vtbl->Release(ta);
}

HANDLER(accessibilityHandler)
{
	if (uMsg != WM_GETOBJECT)
		return FALSE;
	if (wParam != OBJID_CLIENT)
		return FALSE;
	*lResult = LresultFromObject(IID_IUnknown, wParam, t->ta);
	// TODO check *lResult
	return TRUE;
}

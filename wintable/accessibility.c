// 9 november 2014
#define UNICODE
#define _UNICODE
#define STRICT
#define STRICT_TYPED_ITEMIDS
#define CINTERFACE
// get Windows version right; right now Windows XP
#define WINVER 0x0501
#define _WIN32_WINNT 0x0501
#define _WIN32_WINDOWS 0x0501		/* according to Microsoft's winperf.h */
#define _WIN32_IE 0x0600			/* according to Microsoft's sdkddkver.h */
#define NTDDI_VERSION 0x05010000	/* according to Microsoft's sdkddkver.h */
#include <windows.h>
#include <commctrl.h>
#include <stdint.h>
#include <uxtheme.h>
#include <string.h>
#include <wchar.h>
#include <windowsx.h>
#include <vsstyle.h>
#include <vssym32.h>
#include <oleacc.h>

struct tableAccessible {
	IAccessibleVtbl vtbl;
	volatile ULONG refcount;		// TODO ensure this is aligned
	struct table *t;
};

static IAccessibleVtbl aaccessible = {
	// IUnknkown
	.QueryInterface = taQueryInterface,
	.AddRef = taAddRef,
	.Release = taRelease,
	// IDispatch
	.GetTypeInfoCount = taGetTypeInfoCount,
	.GetTypeInfo = taGetTypeInfo,
	.GetIDsOfNames = taGetIDsOfNames,
	.Invoke = taInvoke,
	// IAccessible
	...
};

HRESULT STDMETHODCALLTYPE taQueryInterface(IUnknown *this, REFIID riid, void **ppvObject)
{
	if (ppvObject == NULL)
		return E_POINTER;
	// we're required to return the same pointer for IUnknown
	// since this is a straight singly-derived interface inheritance, we can exploit the structure layout and just return the same pointer for everything
	// at least I hope... (TODO)
	if (IsEqualIID(riid, IID_IUnknown) ||
		IsEqualIID(riid, IID_IDispatch) ||
		IsEqualIID(riid, IID_IAccessible)) {
		this->AddRef(this);
		*ppvObject = (void *) this;
		return S_OK;
	}
	// we're not making a special class for this
	*ppvObject = NULL;
	return E_NOINTERFACE;
}

ULONG STDMETHODCALLTYPE taAddRef(IUnknown *this)
{
	// TODO is the signed conversion safe?
	return (ULONG) InterlockedIncrement((volatile LONG *) (&(((tableAccessible *) this)->refcount)));
}

ULONG STDMETHODCALLTYPE taRelease(IUnknown *this)
{
	ULONG rc;

	rc = (ULONG) InterlockedDecrement((volatile LONG *) (&(((tableAccessible *) this)->refcount)));
	// TODO pull the refcount back out?
	if (rc == 0)
		free((tableAccessible *) this);
	// TODO pull the refcount back out?
	return rc;
}

// here's the IDispatch member functions
// we actually /don't/ need to define any of these!
// see also http://msdn.microsoft.com/en-us/library/windows/desktop/cc307844.aspx

HRESULT STDMETHODCALLTYPE taGetTypeInfoCount(IDispatch *this, UINT *pctinfo)
{
	if (pctinfo == NULL)
		return E_INVALIDARG;
	// TODO really set this to zero?
	*pctinfo = 0;
	return E_NOTIMPL;
}

HRESULT STDMETHODCALLTYPE taGetTypeInfo(IDispatch *this, UINT iTInfo, LCID lcid, ITypeInfo **ppTInfo)
{
	if (pctinfo == NULL)
		return E_INVALIDARG;
	*ppTInfo = NULL;
	// let's do this just to be safe
	if (iTInfo == 0)
		return DISP_E_BADINDEX;
	return E_NOTIMPL;
}

HRESULT STDMETHODCALLTYPE taGetIDsOfNames(IDispatch *this, REFIID riid, LPOLESTR *rgszNames, UINT cNames, LCID lcid, DISPID *rgDispId)
{
	// rgDispId is an array of LONGs; setting it to NULL is useless
	// TODO should we clear the array?
	return E_NOTIMPL;
}

HRESULT STDMETHODCALLTYPE taInvoke(IDispatch *this, DISPID dispIdMember, REFIID riid, LCID lcid, WORD wFlags, DISPPARAMS *pDispParams, VARIANT *pVarResult, EXCEPINFO *pExcepInfo, UINT *puArgErr)
{
	// TODO set anything to NULL or 0?
	return E_NOTIMPL;
}

// ok that's it for IDispatch; now for IAccessible!

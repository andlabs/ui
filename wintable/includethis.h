// 17 february 2015

#ifndef __GO_UI_WINTABLE_INCLUDETHIS_H__
#define __GO_UI_WINTABLE_INCLUDETHIS_H__

#define tableWindowClass L"gouitable"

// start at WM_USER + 20 just in case for whatever reason we ever get the various dialog manager messages (see also http://blogs.msdn.com/b/oldnewthing/archive/2003/10/21/55384.aspx)
// each of these returns nothing unless otherwise indicated
enum {
	// wParam - one of the type constants
	// lParam - column name as a Unicode string
	tableAddColumn = WM_USER + 20,
	// wParam - 0
	// lParam - pointer to intptr_t containing new count
	tableSetRowCount,
	// wParam - pointer to intptr_t where selected row will be stored
	// lParam - pointer to intptr_t where selected column will be stored
	// both will be -1 for no selection
	// if either is NULL, that value is not written
	tableGetSelection,
	// wParam - pointer to intptr_t containing selected row
	// lParam - pointer to intptr_t containing selected column
	// if lParam is NULL, do not change selected column (selects column 0 if nothing previously selected; TODO explicitly document this?)
	// TODO allow wParam to be NULL too; should both being NULL select nothing or keep the current selection?
	// this WILL result in a selection changed notification (TODO work into the package ui Table)
	tableSetSelection,
};

enum {
	tableColumnText,
	tableColumnImage,
	tableColumnCheckbox,
	nTableColumnTypes,
};

// notification codes
// note that these are positive; see http://blogs.msdn.com/b/oldnewthing/archive/2009/08/21/9877791.aspx
// each of these is of type tableNM
// all fields except data will always be set
enum {
	// data parameter is always 0
	// for tableColumnText return should be WCHAR *
	// for tableColumnImage return should be HBITMAP
	// for tableColumnCheckbox return is nonzero for checked, zero for unchecked
	tableNotificationGetCellData,
	// data parameter is pointer, same as tableNotificationGetCellData
	// not sent for checkboxes
	// no return
	tableNotificationFinishedWithCellData,
	// data is zero
	// no return
	tableNotificationCellCheckboxToggled,
	// sent even on deselection (in that case, row == -1 and column == -1)
	// data is zero
	// no return
	tableNotificationSelectionChanged,
};

typedef struct tableNM tableNM;

struct tableNM {
	NMHDR nmhdr;
	intptr_t row;
	intptr_t column;
	int columnType;
	uintptr_t data;
};

// TODO have hInstance passed in
extern void initTable(void (*panicfunc)(const char *msg, DWORD lastError), BOOL (*WINAPI tme)(LPTRACKMOUSEEVENT));

#endif

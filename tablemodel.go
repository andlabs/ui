// 24 august 2018

package ui

// #include <stdlib.h>
// #include "ui.h"
// #include "util.h"
// extern int doTableModelNumColumns(uiTableModelHandler *, uiTableModel *);
// extern uiTableValueType doTableModelColumnType(uiTableModelHandler *, uiTableModel *, int);
// extern int doTableModelNumRows(uiTableModelHandler *, uiTableModel *);
// extern uiTableValue *doTableModelCellValue(uiTableModelHandler *mh, uiTableModel *m, int row, int column);
// extern void doTableModelSetCellValue(uiTableModelHandler *, uiTableModel *, int, int, const uiTableValue *);
// static inline uiTableModelHandler *allocTableModelHandler(void)
// {
// 	uiTableModelHandler *mh;
// 
// 	mh = (uiTableModelHandler *) pkguiAlloc(sizeof (uiTableModelHandler));
// 	mh->NumColumns = doTableModelNumColumns;
// 	mh->ColumnType = doTableModelColumnType;
// 	mh->NumRows = doTableModelNumRows;
// 	mh->CellValue = doTableModelCellValue;
// 	mh->SetCellValue = doTableModelSetCellValue;
// 	return mh;
// }
// static inline void freeTableModelHandler(uiTableModelHandler *mh)
// {
// 	free(mh);
// }
import "C"

// TableValue is a type that represents a piece of data that can come
// out of a TableModel.
type TableValue interface {
	toLibui() *C.uiTableValue
}

// TableString is a TableValue that stores a string. TableString is
// used for displaying text in a Table.
type TableString string

func (s TableString) toLibui() *C.uiTableValue {
	cs := C.CString(string(s))
	defer freestr(cs)
	return C.uiNewTableValueString(cs)
}

// TableImage is a TableValue that represents an Image. Ownership
// of the Image is not copied; you must keep it alive alongside the
// TableImage.
type TableImage struct {
	I	*Image
}

func (i TableImage) toLibui() *C.uiTableValue {
	return C.uiNewTableValueImage(i.I.i)
}

// TableInt is a TableValue that stores integers. These are used for
// progressbars. Due to current limitations of libui, they also
// represent checkbox states, via TableFalse and TableTrue.
type TableInt int

// TableFalse and TableTrue are the Boolean constants for TableInt.
const (
	TableFalse TableInt = 0
	TableTrue TableInt = 1
)

func (i TableInt) toLibui() *C.uiTableValue {
	return C.uiNewTableValueInt(C.int(i))
}

// TableColor is a TableValue that represents a color.
type TableColor struct {
	R	float64
	G	float64
	B	float64
	A	float64
}

func (c TableColor) toLibui() *C.uiTableValue {
	return C.uiNewTableValueColor(C.double(c.R), C.double(c.G), C.double(c.B), C.double(c.A))
}

func tableValueFromLibui(value *C.uiTableValue) TableValue {
	if value == nil {
		return nil
	}
	switch C.uiTableValueGetType(value) {
	case C.uiTableValueTypeString:
		cs := C.uiTableValueString(value)
		return TableString(C.GoString(cs))
	case C.uiTableValueTypeImage:
		panic("TODO")
	case C.uiTableValueTypeInt:
		return TableInt(C.uiTableValueInt(value))
	case C.uiTableValueTypeColor:
		panic("TODO")
	}
	panic("unreachable")
}

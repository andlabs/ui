// 24 august 2018

package ui

// #include "pkgui.h"
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

// no need to lock these; only the GUI thread can access them
var modelhandlers = make(map[*C.uiTableModel]TableModelHandler)
var models = make(map[*C.uiTableModel]*TableModel)

// TableModel is an object that provides the data for a Table.
// This data is returned via methods you provide in the
// TableModelHandler interface.
//
// TableModel represents data using a table, but this table does
// not map directly to Table itself. Instead, you can have data
// columns which provide instructions for how to render a given
// Table's column — for instance, one model column can be used
// to give certain rows of a Table a different background color.
// Row numbers DO match with uiTable row numbers.
//
// Once created, the number and data types of columns of a
// TableModel cannot change.
//
// Row and column numbers start at 0. A TableModel can be
// associated with more than one Table at a time.
type TableModel struct {
	m	*C.uiTableModel
}

// TableModelHandler defines the methods that TableModel
// calls when it needs data.
type TableModelHandler interface {
	// ColumnTypes returns a slice of value types of the data
	// stored in the model columns of the TableModel.
	// Each entry in the slice should ideally be a zero value for
	// the TableValue type of the column in question; the number
	// of elements in the slice determines the number of model
	// columns in the TableModel. The returned slice must remain
	// constant through the lifetime of the TableModel. This
	// method is not guaranteed to be called depending on the
	// system.
	ColumnTypes(m *TableModel) []TableValue

	// NumRows returns the number or rows in the TableModel.
	// This value must be non-negative.
	NumRows(m *TableModel) int

	// CellValue returns a TableValue corresponding to the model
	// cell at (row, column). The type of the returned TableValue
	// must match column's value type. Under some circumstances,
	// nil may be returned; refer to the various methods that add
	// columns to Table for details.
	CellValue(m *TableModel, row, column int) TableValue

	// SetCellValue changes the model cell value at (row, column)
	// in the TableModel. Within this function, either do nothing
	// to keep the current cell value or save the new cell value as
	// appropriate. After SetCellValue is called, the Table will
	// itself reload the table cell. Under certain conditions, the
	// TableValue passed in can be nil; refer to the various
	// methods that add columns to Table for details.
	SetCellValue(m *TableModel, row, column int, value TableValue)
}

//export pkguiDoTableModelNumColumns
func pkguiDoTableModelNumColumns(umh *C.uiTableModelHandler, um *C.uiTableModel) C.int {
	mh := modelhandlers[um]
	return C.int(len(mh.ColumnTypes(models[um])))
}

//export pkguiDoTableModelColumnType
func pkguiDoTableModelColumnType(umh *C.uiTableModelHandler, um *C.uiTableModel, n C.int) C.uiTableValueType {
	mh := modelhandlers[um]
	c := mh.ColumnTypes(models[um])
	switch c[n].(type) {
	case TableString:
		return C.uiTableValueTypeString
	case TableImage:
		return C.uiTableValueTypeImage
	case TableInt:
		return C.uiTableValueTypeInt
	case TableColor:
		return C.uiTableValueTypeColor
	}
	panic("unreachable")
}

//export pkguiDoTableModelNumRows
func pkguiDoTableModelNumRows(umh *C.uiTableModelHandler, um *C.uiTableModel) C.int {
	mh := modelhandlers[um]
	return C.int(mh.NumRows(models[um]))
}

//export pkguiDoTableModelCellValue
func pkguiDoTableModelCellValue(umh *C.uiTableModelHandler, um *C.uiTableModel, row, column C.int) *C.uiTableValue {
	mh := modelhandlers[um]
	v := mh.CellValue(models[um], int(row), int(column))
	if v == nil {
		return nil
	}
	return v.toLibui()
}

//export pkguiDoTableModelSetCellValue
func pkguiDoTableModelSetCellValue(umh *C.uiTableModelHandler, um *C.uiTableModel, row, column C.int, value *C.uiTableValue) {
	mh := modelhandlers[um]
	v := tableValueFromLibui(value)
	mh.SetCellValue(models[um], int(row), int(column), v)
}

// NewTableModel creates a new TableModel.
func NewTableModel(handler TableModelHandler) *TableModel {
	m := &TableModel{
		m:	C.uiNewTableModel(&C.pkguiTableModelHandler),
	}
	modelhandlers[m.m] = handler
	models[m.m] = m
	return m
}

// Free frees m. It is an error to Free any models associated with a
// Table.
func (m *TableModel) Free() {
	delete(models, m.m)
	delete(modelhandlers, m.m)
	C.uiFreeTableModel(m.m)
}

// RowInserted tells any Tables associated with m that a new row
// has been added to m at index index. You call this method when
// the number of rows in your model has changed; after calling it,
// NumRows should returm the new row count.
func (m *TableModel) RowInserted(index int) {
	C.uiTableModelRowInserted(m.m, C.int(index))
}

// RowChanged tells any Tables associated with m that the data in
// the row at index has changed. You do not need to call this in
// your SetCellValue handlers, but you do need to call this if your
// data changes at some other point.
func (m *TableModel) RowChanged(index int) {
	C.uiTableModelRowChanged(m.m, C.int(index))
}

// RowDeleted tells any Tables associated with m that the row at
// index index has been deleted. You call this function when the
// number of rows in your model has changed; after calling it,
// NumRows should returm the new row count.
func (m *TableModel) RowDeleted(index int) {
	C.uiTableModelRowDeleted(m.m, C.int(index))
}

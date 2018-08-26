// 26 august 2018

package ui

import (
	"unsafe"
)

// #include "pkgui.h"
import "C"

// TableModelColumnNeverEditable and
// TableModelColumnAlwaysEditable are the value of an editable
// model column parameter to one of the Table create column
// functions; if used, that jparticular Table colum is not editable
// by the user and always editable by the user, respectively.
const (
	TableModelColumnNeverEditable = -1
	TableModelColumnAlwaysEditable = -2
)

// TableTextColumnOptionalParams are the optional parameters
// that control the appearance of the text column of a Table.
type TableTextColumnOptionalParams struct {
	// ColorModelColumn is the model column containing the
	// text color of this Table column's text, or -1 to use the
	// default color.
	//
	// If CellValue for this column for any cell returns nil, that
	// cell will also use the default text color.
	ColorModelColumn		int
}

func (p *TableTextColumnOptionalParams) toLibui() *C.uiTableTextColumnOptionalParams {
	if p == nil {
		return nil
	}
	cp := C.pkguiAllocTableTextColumnOptionalParams()
	cp.ColorModelColumn = C.int(p.ColorModelColumn)
	return cp
}

// TableParams defines the parameters passed to NewTable.
type TableParams struct {
	// Model is the TableModel to use for this uiTable.
	// This parameter cannot be nil.
	Model		*TableModel

	// RowBackgroundColorModelColumn is a model column
	// number that defines the background color used for the
	// entire row in the Table, or -1 to use the default color for
	// all rows.
	//
	// If CellValue for this column for any row returns NULL, that
	// row will also use the default background color.
	RowBackgroundColorModelColumn		int
}

func (p *TableParams) toLibui() *C.uiTableParams {
	cp := C.pkguiAllocTableParams()
	cp.Model = p.Model.m
	cp.RowBackgroundColorModelColumn = C.int(p.RowBackgroundColorModelColumn)
	return cp
}

// Table is a Control that shows tabular data, allowing users to
// manipulate rows of such data at a time.
type Table struct {
	ControlBase
	t	*C.uiTable
}

// NewTable creates a new Table with the specified parameters.
func NewTable(p *TableParams) *Table {
	t := new(Table)

	cp := p.toLibui()
	t.t = C.uiNewTable(cp)
	C.pkguiFreeTableParams(cp)

	t.ControlBase = NewControlBase(t, uintptr(unsafe.Pointer(t.t)))
	return t
}

// AppendTextColumn appends a text column to t. name is
// displayed in the table header. textModelColumn is where the text
// comes from. If a row is editable according to
// textEditableModelColumn, SetCellValue is called with
// textModelColumn as the column.
func (t *Table) AppendTextColumn(name string, textModelColumn int, textEditableModelColumn int, textParams *TableTextColumnOptionalParams) {
	cname := C.CString(name)
	defer freestr(cname)
	cp := textParams.toLibui()
	defer C.pkguiFreeTableTextColumnOptionalParams(cp)
	C.uiTableAppendTextColumn(t.t, cname, C.int(textModelColumn), C.int(textEditableModelColumn), cp)
}

// AppendImageColumn appends an image column to t.
// Images are drawn at icon size, appropriate to the pixel density
// of the screen showing the Table.
func (t *Table) AppendImageColumn(name string, imageModelColumn int) {
	cname := C.CString(name)
	defer freestr(cname)
	C.uiTableAppendImageColumn(t.t, cname, C.int(imageModelColumn))
}

// AppendImageTextColumn appends a column to t that
// shows both an image and text.
func (t *Table) AppendImageTextColumn(name string, imageModelColumn int, textModelColumn int, textEditableModelColumn int, textParams *TableTextColumnOptionalParams) {
	cname := C.CString(name)
	defer freestr(cname)
	cp := textParams.toLibui()
	defer C.pkguiFreeTableTextColumnOptionalParams(cp)
	C.uiTableAppendImageTextColumn(t.t, cname, C.int(imageModelColumn), C.int(textModelColumn), C.int(textEditableModelColumn), cp)
}

// AppendCheckboxColumn appends a column to t that
// contains a checkbox that the user can interact with (assuming the
// checkbox is editable). SetCellValue will be called with
// checkboxModelColumn as the column in this case.
func (t *Table) AppendCheckboxColumn(name string, checkboxModelColumn int, checkboxEditableModelColumn int) {
	cname := C.CString(name)
	defer freestr(cname)
	C.uiTableAppendCheckboxColumn(t.t, cname, C.int(checkboxModelColumn), C.int(checkboxEditableModelColumn))
}

// AppendCheckboxTextColumn appends a column to t
// that contains both a checkbox and text.
func (t *Table) AppendCheckboxTextColumn(name string, checkboxModelColumn int, checkboxEditableModelColumn int, textModelColumn int, textEditableModelColumn int, textParams *TableTextColumnOptionalParams) {
	cname := C.CString(name)
	defer freestr(cname)
	cp := textParams.toLibui()
	defer C.pkguiFreeTableTextColumnOptionalParams(cp)
	C.uiTableAppendCheckboxTextColumn(t.t, cname, C.int(checkboxModelColumn), C.int(checkboxEditableModelColumn), C.int(textModelColumn), C.int(textEditableModelColumn), cp)
}

// AppendProgressBarColumn appends a column to t
// that displays a progress bar. These columns work like
// ProgressBar: a cell value of 0..100 displays that percentage, and
// a cell value of -1 displays an indeterminate progress bar.
func (t *Table) AppendProgressBarColumn(name string, progressModelColumn int) {
	cname := C.CString(name)
	defer freestr(cname)
	C.uiTableAppendProgressBarColumn(t.t, cname, C.int(progressModelColumn))
}

// AppendButtonColumn appends a column to t
// that shows a button that the user can click on. When the user
// does click on the button, SetCellValue is called with a nil
// value and buttonModelColumn as the column.
// CellValue on buttonModelColumn should return the text to show
// in the button.
func (t *Table) AppendButtonColumn(name string, buttonModelColumn int, buttonClickableModelColumn int) {
	cname := C.CString(name)
	defer freestr(cname)
	C.uiTableAppendButtonColumn(t.t, cname, C.int(buttonModelColumn), C.int(buttonClickableModelColumn))
}

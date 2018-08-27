// 26 august 2018

// +build OMIT

// TODO possible bug in libui: the checkboxes on macOS retain their values when they shouldn't

package main

import (
	"fmt"

	"github.com/andlabs/ui"
)

type modelHandler struct {
	row9Text		string
	yellowRow	int
	checkStates	[15]int
}

func newModelHandler() *modelHandler {
	m := new(modelHandler)
	m.row9Text = "You can edit this one"
	m.yellowRow = -1
	return m
}

func (mh *modelHandler) ColumnTypes(m *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""),		// column 0 text
		ui.TableString(""),		// column 1 text
		ui.TableString(""),		// column 2 text
		ui.TableColor{},			// row background color
		ui.TableColor{},			// column 1 text color
		ui.TableImage{},		// column 1 image
		ui.TableString(""),		// column 4 button text
		ui.TableInt(0),			// column 3 checkbox state
		ui.TableInt(0),			// column 5 progress
	}
}

func (mh *modelHandler) NumRows(m *ui.TableModel) int {
	return 15
}

func (mh *modelHandler) CellValue(m *ui.TableModel, row, column int) ui.TableValue {
	if column == 3 {
		if row == mh.yellowRow {
			return ui.TableColor{1, 1, 0, 1}
		}
		if row == 3 {
			return ui.TableColor{1, 0, 0, 1}
		}
		if row == 11 {
			return ui.TableColor{0, 0.5, 1, 0.5}
		}
		return nil
	}
	if column == 4 {
		if (row % 2) == 1 {
			return ui.TableColor{0.5, 0, 0.75, 1}
		}
		return nil
	}
/*	if column == 5 {
		if row < 8 {
			return ui.TableImage{img[0]}
		}
		return ui.TableImage{img[1]}
	}
*/	if column == 7 {
		return ui.TableInt(mh.checkStates[row])
	}
	if column == 8 {
		if row == 0 {
			return ui.TableInt(0)
		}
		if row == 13 {
			return ui.TableInt(100)
		}
		if row == 14 {
			return ui.TableInt(-1)
		}
		return ui.TableInt(50)
	}
	switch column {
	case 0:
		return ui.TableString(fmt.Sprintf("Row %d", row))
	case 2:
		if row == 9 {
			return ui.TableString(mh.row9Text)
		}
		return ui.TableString("Editing this won't change anything")
	case 1:
		return ui.TableString("Colors!")
	case 6:
		return ui.TableString("Make Yellow")
	}
	panic("unreachable")
}

func (mh *modelHandler) SetCellValue(m *ui.TableModel, row, column int, value ui.TableValue) {
}

func setupUI() {
	mainwin := ui.NewWindow("libui Control Gallery", 640, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	mh := newModelHandler()
	model := ui.NewTableModel(mh)

	table := ui.NewTable(&ui.TableParams{
		Model:	model,
		RowBackgroundColorModelColumn:	3,
	})
	mainwin.SetChild(table)
	mainwin.SetMargined(true)

	table.AppendTextColumn("Column 1",
		0, ui.TableModelColumnNeverEditable, nil)

/*	table.AppendImageTextColumn("Column 2",
		5,
		1, ui.TableModelColumnNeverEditable, &ui.TableTextColumnOptionalParams{
			ColorModelColumn:		4,
		});
*/	table.AppendTextColumn("Editable",
		2, ui.TableModelColumnAlwaysEditable, nil)

	table.AppendCheckboxColumn("Checkboxes",
		7, ui.TableModelColumnAlwaysEditable)
	table.AppendButtonColumn("Buttons",
		6, ui.TableModelColumnAlwaysEditable)

	table.AppendProgressBarColumn("Progress Bar",
		8)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}

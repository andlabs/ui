// 12 august 2018

// +build OMIT

package main

import (
	"github.com/andlabs/ui"
)

func makeBasicControlsPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)

	hbox.Append(ui.NewButton("Button"), false)
	hbox.Append(ui.NewCheckbox("Checkbox"), false)

	vbox.Append(ui.NewLabel("This is a label. Right now, labels can only span one line."), false)

	vbox.Append(ui.NewHorizontalSeparator(), false)

	group := ui.NewGroup("Entries")
	group.SetMargined(true)
	vbox.Append(group, true)

group.SetChild(ui.NewNonWrappingMultilineEntry())
/*
	entryForm = uiNewForm();
	uiFormSetPadded(entryForm, 1);
	uiGroupSetChild(group, uiControl(entryForm));

	uiFormAppend(entryForm,
		"Entry",
		uiControl(uiNewEntry()),
		0);
	uiFormAppend(entryForm,
		"Password Entry",
		uiControl(uiNewPasswordEntry()),
		0);
	uiFormAppend(entryForm,
		"Search Entry",
		uiControl(uiNewSearchEntry()),
		0);
	uiFormAppend(entryForm,
		"Multiline Entry",
		uiControl(uiNewMultilineEntry()),
		1);
	uiFormAppend(entryForm,
		"Multiline Entry No Wrap",
		uiControl(uiNewNonWrappingMultilineEntry()),
		1);
*/

	return vbox
}

func makeNumbersPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	group := ui.NewGroup("Numbers")
	group.SetMargined(true)
	hbox.Append(group, true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	group.SetChild(vbox)

	spinbox := ui.NewSpinbox(0, 100)
	slider := ui.NewSlider(0, 100)
	pbar := ui.NewProgressBar()
	spinbox.OnChanged(func(*ui.Spinbox) {
		slider.SetValue(spinbox.Value())
		pbar.SetValue(spinbox.Value())
	})
	slider.OnChanged(func(*ui.Slider) {
		spinbox.SetValue(slider.Value())
		pbar.SetValue(slider.Value())
	})
	vbox.Append(spinbox, false)
	vbox.Append(slider, false)
	vbox.Append(pbar, false)

	ip := ui.NewProgressBar()
	ip.SetValue(-1)
	vbox.Append(ip, false)

	group = ui.NewGroup("Lists")
	group.SetMargined(true)
	hbox.Append(group, true)

	vbox = ui.NewVerticalBox()
	vbox.SetPadded(true)
	group.SetChild(vbox)

	cbox := ui.NewCombobox()
	cbox.Append("Combobox Item 1")
	cbox.Append("Combobox Item 2")
	cbox.Append("Combobox Item 3")
	vbox.Append(cbox, false)

	ecbox := ui.NewEditableCombobox()
	ecbox.Append("Editable Item 1")
	ecbox.Append("Editable Item 2")
	ecbox.Append("Editable Item 3")
	vbox.Append(ecbox, false)

	rb := ui.NewRadioButtons()
	rb.Append("Radio Button 1")
	rb.Append("Radio Button 2")
	rb.Append("Radio Button 3")
	vbox.Append(rb, false)

	return hbox
}

func makeDataChoosersPage() ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	vbox.Append(ui.NewDatePicker(), false)
	vbox.Append(ui.NewTimePicker(), false)
	vbox.Append(ui.NewDateTimePicker(), false)
/*
	uiBoxAppend(vbox,
		uiControl(uiNewFontButton()),
		0);
	uiBoxAppend(vbox,
		uiControl(uiNewColorButton()),
		0);

	uiBoxAppend(hbox,
		uiControl(uiNewVerticalSeparator()),
		0);

	vbox = uiNewVerticalBox();
	uiBoxSetPadded(vbox, 1);
	uiBoxAppend(hbox, uiControl(vbox), 1);

	grid = uiNewGrid();
	uiGridSetPadded(grid, 1);
	uiBoxAppend(vbox, uiControl(grid), 0);

	button = uiNewButton("Open File");
	entry = uiNewEntry();
	uiEntrySetReadOnly(entry, 1);
	uiButtonOnClicked(button, onOpenFileClicked, entry);
	uiGridAppend(grid, uiControl(button),
		0, 0, 1, 1,
		0, uiAlignFill, 0, uiAlignFill);
	uiGridAppend(grid, uiControl(entry),
		1, 0, 1, 1,
		1, uiAlignFill, 0, uiAlignFill);

	button = uiNewButton("Save File");
	entry = uiNewEntry();
	uiEntrySetReadOnly(entry, 1);
	uiButtonOnClicked(button, onSaveFileClicked, entry);
	uiGridAppend(grid, uiControl(button),
		0, 1, 1, 1,
		0, uiAlignFill, 0, uiAlignFill);
	uiGridAppend(grid, uiControl(entry),
		1, 1, 1, 1,
		1, uiAlignFill, 0, uiAlignFill);

	msggrid = uiNewGrid();
	uiGridSetPadded(msggrid, 1);
	uiGridAppend(grid, uiControl(msggrid),
		0, 2, 2, 1,
		0, uiAlignCenter, 0, uiAlignStart);

	button = uiNewButton("Message Box");
	uiButtonOnClicked(button, onMsgBoxClicked, NULL);
	uiGridAppend(msggrid, uiControl(button),
		0, 0, 1, 1,
		0, uiAlignFill, 0, uiAlignFill);
	button = uiNewButton("Error Box");
	uiButtonOnClicked(button, onMsgBoxErrorClicked, NULL);
	uiGridAppend(msggrid, uiControl(button),
		1, 0, 1, 1,
		0, uiAlignFill, 0, uiAlignFill);
*/

	return hbox
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

	tab := ui.NewTab()
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	tab.Append("Basic Controls", makeBasicControlsPage())
	tab.SetMargined(0, true)

	tab.Append("Numbers and Lists", makeNumbersPage());
	tab.SetMargined(1, true)

	tab.Append("Data Choosers", makeDataChoosersPage());
	tab.SetMargined(2, true)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}

// 19 august 2018

// +build OMIT

package main

// TODO probably a bug in libui: changing the font away from skia leads to a crash

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

var (
	fontButton *ui.FontButton
	alignment *ui.Combobox

	attrstr *ui.AttributedString
)

func appendWithAttributes(what string, attrs ...ui.Attribute) {
	start := len(attrstr.String())
	end := start + len(what)
	attrstr.AppendUnattributed(what)
	for _, a := range attrs {
		attrstr.SetAttribute(a, start, end)
	}
}

func makeAttributedString() {
	attrstr = ui.NewAttributedString(
		"Drawing strings with package ui is done with the ui.AttributedString and ui.DrawTextLayout objects.\n" +
		"ui.AttributedString lets you have a variety of attributes: ")

	appendWithAttributes("font family", ui.TextFamily("Courier New"))
	attrstr.AppendUnattributed(", ")

	appendWithAttributes("font size", ui.TextSize(18))
	attrstr.AppendUnattributed(", ")

	appendWithAttributes("font weight", ui.TextWeightBold)
	attrstr.AppendUnattributed(", ")

	appendWithAttributes("font italicness", ui.TextItalicItalic)
	attrstr.AppendUnattributed(", ")

	appendWithAttributes("font stretch", ui.TextStretchCondensed)
	attrstr.AppendUnattributed(", ")

	appendWithAttributes("text color", ui.TextColor{0.75, 0.25, 0.5, 0.75})
	attrstr.AppendUnattributed(", ")

	appendWithAttributes("text background color", ui.TextBackground{0.5, 0.5, 0.25, 0.5})
	attrstr.AppendUnattributed(", ")

	appendWithAttributes("underline style", ui.UnderlineSingle)
	attrstr.AppendUnattributed(", ")

	attrstr.AppendUnattributed("and ")
	appendWithAttributes("underline color",
		ui.UnderlineDouble,
		ui.UnderlineColorCustom{1.0, 0.0, 0.5, 1.0})
	attrstr.AppendUnattributed(". ")

	attrstr.AppendUnattributed("Furthermore, there are attributes allowing for ")
	appendWithAttributes("special underlines for indicating spelling errors",
		ui.UnderlineSuggestion,
		ui.UnderlineColorSpelling)
	attrstr.AppendUnattributed(" (and other types of errors) ")

	attrstr.AppendUnattributed("and control over OpenType features such as ligatures (for instance, ")
	appendWithAttributes("afford", ui.OpenTypeFeatures{
		ui.ToOpenTypeTag('l', 'i', 'g', 'a'):		0,
	})
	attrstr.AppendUnattributed(" vs. ")
	appendWithAttributes("afford", ui.OpenTypeFeatures{
		ui.ToOpenTypeTag('l', 'i', 'g', 'a'):		1,
	})
	attrstr.AppendUnattributed(").\n")

	attrstr.AppendUnattributed("Use the controls opposite to the text to control properties of the text.")
}

type areaHandler struct{}

func (areaHandler) Draw(a *ui.Area, p *ui.AreaDrawParams) {
	tl := ui.DrawNewTextLayout(&ui.DrawTextLayoutParams{
		String:		attrstr,
		DefaultFont:	fontButton.Font(),
		Width:		p.AreaWidth,
		Align:		ui.DrawTextAlign(alignment.Selected()),
	})
	defer tl.Free()
	p.Context.Text(tl, 0, 0)
}

func (areaHandler) MouseEvent(a *ui.Area, me *ui.AreaMouseEvent) {
	// do nothing
}

func (areaHandler) MouseCrossed(a *ui.Area, left bool) {
	// do nothing
}

func (areaHandler) DragBroken(a *ui.Area) {
	// do nothing
}

func (areaHandler) KeyEvent(a *ui.Area, ke *ui.AreaKeyEvent) (handled bool) {
	// reject all keys
	return false
}

func setupUI() {
	makeAttributedString()

	mainwin := ui.NewWindow("libui Text-Drawing Example", 640, 480, true)
	mainwin.SetMargined(true)
	mainwin.OnClosing(func(*ui.Window) bool {
		mainwin.Destroy()
		ui.Quit()
		return false
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	mainwin.SetChild(hbox)

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	hbox.Append(vbox, false)

	area := ui.NewArea(areaHandler{})

	fontButton = ui.NewFontButton()
	fontButton.OnChanged(func(*ui.FontButton) {
		area.QueueRedrawAll()
	})
	vbox.Append(fontButton, false)

	form := ui.NewForm()
	form.SetPadded(true)
	// TODO on OS X if this is set to 1 then the window can't resize; does the form not have the concept of stretchy trailing space?
	vbox.Append(form, false)

	alignment = ui.NewCombobox()
	// note that the items match with the values of the uiDrawTextAlign values
	alignment.Append("Left")
	alignment.Append("Center")
	alignment.Append("Right")
	alignment.SetSelected(0)		// start with left alignment
	alignment.OnSelected(func(*ui.Combobox) {
		area.QueueRedrawAll()
	})
	form.Append("Alignment", alignment, false)

	hbox.Append(area, true)

	mainwin.Show()
}

func main() {
	ui.Main(setupUI)
}

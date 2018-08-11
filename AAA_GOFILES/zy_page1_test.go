// 12 december 2015

package ui

var page1 *Box

func makePage1(w *Window) {
	var xbutton *Button

	page1 = newVerticalBox()

	entry := NewEntry()
	page1.Append(entry, false)

	spaced := NewCheckbox("Spaced")
	spaced.OnToggled(func(*Checkbox) {
		setSpaced(spaced.Checked())
	})
	label := NewLabel("Label")

	hbox := newHorizontalBox()
	getButton := NewButton("Get Window Text")
	getButton.OnClicked(func(*Button) {
		entry.SetText(w.Title())
	})
	setButton := NewButton("Set Window Text")
	setButton.OnClicked(func(*Button) {
		w.SetTitle(entry.Text())
	})
	hbox.Append(getButton, true)
	hbox.Append(setButton, true)
	page1.Append(hbox, false)

	hbox = newHorizontalBox()
	getButton = NewButton("Get Button Text")
	xbutton = getButton
	getButton.OnClicked(func(*Button) {
		entry.SetText(xbutton.Text())
	})
	setButton = NewButton("Set Button Text")
	setButton.OnClicked(func(*Button) {
		xbutton.SetText(entry.Text())
	})
	hbox.Append(getButton, true)
	hbox.Append(setButton, true)
	page1.Append(hbox, false)

	hbox = newHorizontalBox()
	getButton = NewButton("Get Checkbox Text")
	getButton.OnClicked(func(*Button) {
		entry.SetText(spaced.Text())
	})
	setButton = NewButton("Set Checkbox Text")
	setButton.OnClicked(func(*Button) {
		spaced.SetText(entry.Text())
	})
	hbox.Append(getButton, true)
	hbox.Append(setButton, true)
	page1.Append(hbox, false)

	hbox = newHorizontalBox()
	getButton = NewButton("Get Label Text")
	getButton.OnClicked(func(*Button) {
		entry.SetText(label.Text())
	})
	setButton = NewButton("Set Label Text")
	setButton.OnClicked(func(*Button) {
		label.SetText(entry.Text())
	})
	hbox.Append(getButton, true)
	hbox.Append(setButton, true)
	page1.Append(hbox, false)

	hbox = newHorizontalBox()
	getButton = NewButton("Get Group Text")
	getButton.OnClicked(func(*Button) {
		entry.SetText(page2group.Title())
	})
	setButton = NewButton("Set Group Text")
	setButton.OnClicked(func(*Button) {
		page2group.SetTitle(entry.Text())
	})
	hbox.Append(getButton, true)
	hbox.Append(setButton, true)
	page1.Append(hbox, false)

	hbox = newHorizontalBox()
	hbox.Append(spaced, true)
	getButton = NewButton("On")
	getButton.OnClicked(func(*Button) {
		spaced.SetChecked(true)
	})
	hbox.Append(getButton, false)
	getButton = NewButton("Off")
	getButton.OnClicked(func(*Button) {
		spaced.SetChecked(false)
	})
	hbox.Append(getButton, false)
	getButton = NewButton("Show")
	getButton.OnClicked(func(*Button) {
		// TODO
	})
	hbox.Append(getButton, false)
	page1.Append(hbox, false)

	testBox := newHorizontalBox()
	ybutton := NewButton("Button")
	testBox.Append(ybutton, true)
	getButton = NewButton("Show")
	getButton.OnClicked(func(*Button) {
		ybutton.Show()
	})
	testBox.Append(getButton, false)
	getButton = NewButton("Hide")
	getButton.OnClicked(func(*Button) {
		ybutton.Hide()
	})
	testBox.Append(getButton, false)
	getButton = NewButton("Enable")
	getButton.OnClicked(func(*Button) {
		ybutton.Enable()
	})
	testBox.Append(getButton, false)
	getButton = NewButton("Disable")
	getButton.OnClicked(func(*Button) {
		ybutton.Disable()
	})
	testBox.Append(getButton, false)
	page1.Append(testBox, false)

	hbox = newHorizontalBox()
	getButton = NewButton("Show")
	getButton.OnClicked(func(*Button) {
		testBox.Show()
	})
	hbox.Append(getButton, false)
	getButton = NewButton("Hide")
	getButton.OnClicked(func(*Button) {
		testBox.Hide()
	})
	hbox.Append(getButton, false)
	getButton = NewButton("Enable")
	getButton.OnClicked(func(*Button) {
		testBox.Enable()
	})
	hbox.Append(getButton, false)
	getButton = NewButton("Disable")
	getButton.OnClicked(func(*Button) {
		testBox.Disable()
	})
	hbox.Append(getButton, false)
	page1.Append(hbox, false)

	page1.Append(label, false)
}

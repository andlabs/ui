// 12 december 2015

package ui

var page2group *Group

var (
	movingLabel   *Label
	movingBoxes   [2]*Box
	movingCurrent int
)

func moveLabel(*Button) {
	from := movingCurrent
	to := 0
	if from == 0 {
		to = 1
	}
	movingBoxes[from].Delete(0)
	movingBoxes[to].Append(movingLabel, false)
	movingCurrent = to
}

var moveBack bool

const (
	moveOutText  = "Move Page 1 Out"
	moveBackText = "Move Page 1 Back"
)

func movePage1(b *Button) {
	if moveBack {
		mainbox.Delete(1)
		mainTab.InsertAt("Page 1", 0, page1)
		b.SetText(moveOutText)
		moveBack = false
		return
	}
	mainTab.Delete(0)
	mainbox.Append(page1, true)
	b.SetText(moveBackText)
	moveBack = true
}

func makePage2() *Box {
	page2 := newVerticalBox()

	group := newGroup("Moving Label")
	page2group = group
	page2.Append(group, false)
	vbox := newVerticalBox()
	group.SetChild(vbox)

	hbox := newHorizontalBox()
	button := NewButton("Move the Label!")
	button.OnClicked(moveLabel)
	hbox.Append(button, true)
	hbox.Append(NewLabel(""), true)
	vbox.Append(hbox, false)

	hbox = newHorizontalBox()
	movingBoxes[0] = newVerticalBox()
	hbox.Append(movingBoxes[0], true)
	movingBoxes[1] = newVerticalBox()
	hbox.Append(movingBoxes[1], true)
	vbox.Append(hbox, false)

	movingCurrent = 0
	movingLabel = NewLabel("This label moves!")
	movingBoxes[movingCurrent].Append(movingLabel, false)

	hbox = newHorizontalBox()
	button = NewButton(moveOutText)
	button.OnClicked(movePage1)
	hbox.Append(button, false)
	page2.Append(hbox, false)
	moveBack = false

	hbox = newHorizontalBox()
	hbox.Append(NewLabel("Label Alignment Test"), false)
	button = NewButton("Open Menued Window")
	button.OnClicked(func(*Button) {
		w := NewWindow("Another Window", 100, 100, true)
		b := NewVerticalBox()
		b.Append(NewEntry(), false)
		b.Append(NewButton("Button"), false)
		b.SetPadded(true)
		w.SetChild(b)
		w.SetMargined(true)
		w.Show()
	})
	hbox.Append(button, false)
	button = NewButton("Open Menuless Window")
	button.OnClicked(func(*Button) {
		w := NewWindow("Another Window", 100, 100, true)
		//TODO		w.SetChild(makePage6())
		w.SetMargined(true)
		w.Show()
	})
	hbox.Append(button, false)
	button = NewButton("Disabled Menued")
	button.OnClicked(func(*Button) {
		w := NewWindow("Another Window", 100, 100, true)
		w.Disable()
		w.Show()
	})
	hbox.Append(button, false)
	button = NewButton("Disabled Menuless")
	button.OnClicked(func(*Button) {
		w := NewWindow("Another Window", 100, 100, false)
		w.Disable()
		w.Show()
	})
	hbox.Append(button, false)
	page2.Append(hbox, false)

	nestedBox := newHorizontalBox()
	innerhbox := newHorizontalBox()
	innerhbox.Append(NewButton("These"), false)
	button = NewButton("buttons")
	button.Disable()
	innerhbox.Append(button, false)
	nestedBox.Append(innerhbox, false)
	innerhbox = newHorizontalBox()
	innerhbox.Append(NewButton("are"), false)
	innerhbox2 := newHorizontalBox()
	button = NewButton("in")
	button.Disable()
	innerhbox2.Append(button, false)
	innerhbox.Append(innerhbox2, false)
	nestedBox.Append(innerhbox, false)
	innerhbox = newHorizontalBox()
	innerhbox2 = newHorizontalBox()
	innerhbox2.Append(NewButton("nested"), false)
	innerhbox3 := newHorizontalBox()
	button = NewButton("boxes")
	button.Disable()
	innerhbox3.Append(button, false)
	innerhbox2.Append(innerhbox3, false)
	innerhbox.Append(innerhbox2, false)
	nestedBox.Append(innerhbox, false)
	page2.Append(nestedBox, false)

	hbox = newHorizontalBox()
	button = NewButton("Enable Nested Box")
	button.OnClicked(func(*Button) {
		nestedBox.Enable()
	})
	hbox.Append(button, false)
	button = NewButton("Disable Nested Box")
	button.OnClicked(func(*Button) {
		nestedBox.Disable()
	})
	hbox.Append(button, false)
	page2.Append(hbox, false)

	disabledTab := newTab()
	disabledTab.Append("Disabled", NewButton("Button"))
	disabledTab.Append("Tab", NewLabel("Label"))
	disabledTab.Disable()
	page2.Append(disabledTab, true)

	entry := NewEntry()
	readonly := NewEntry()
	entry.OnChanged(func(*Entry) {
		readonly.SetText(entry.Text())
	})
	readonly.SetText("If you can see this, uiEntryReadOnly() isn't working properly.")
	readonly.SetReadOnly(true)
	if readonly.ReadOnly() {
		readonly.SetText("")
	}
	page2.Append(entry, false)
	page2.Append(readonly, false)

	hbox = newHorizontalBox()
	button = NewButton("Show Button 2")
	button2 := NewButton("Button 2")
	button.OnClicked(func(*Button) {
		button2.Show()
	})
	button2.Hide()
	hbox.Append(button, true)
	hbox.Append(button2, false)
	page2.Append(hbox, false)

	return page2
}

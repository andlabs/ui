package main

import (
	"fmt"
	"time"

	"github.com/andlabs/ui"
)

// Example showing how to update the UI using the QueueMain function
// especially if the update is coming from another goroutine
//
// see QueueMain in 'main.go' for detailed description

var countLabel *ui.Label
var count int

func setupUI() {
	mainWindow := ui.NewWindow("libui Updating UI", 640, 480, true)
	mainWindow.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainWindow.Destroy()
		return true
	})

	vbContainer := ui.NewVerticalBox()
	vbContainer.SetPadded(true)

	inputGroup := ui.NewGroup("Input")
	inputGroup.SetMargined(true)

	vbInput := ui.NewVerticalBox()
	vbInput.SetPadded(true)

	inputForm := ui.NewForm()
	inputForm.SetPadded(true)

	message := ui.NewEntry()
	message.SetText("Hello World")
	inputForm.Append("What message do you want to show?", message, false)

	showMessageButton := ui.NewButton("Show message")
	clearMessageButton := ui.NewButton("Clear message")

	vbInput.Append(inputForm, false)
	vbInput.Append(showMessageButton, false)
	vbInput.Append(clearMessageButton, false)

	inputGroup.SetChild(vbInput)

	messageGroup := ui.NewGroup("Message")
	messageGroup.SetMargined(true)

	vbMessage := ui.NewVerticalBox()
	vbMessage.SetPadded(true)

	messageLabel := ui.NewLabel("")

	vbMessage.Append(messageLabel, false)

	messageGroup.SetChild(vbMessage)

	countGroup := ui.NewGroup("Counter")
	countGroup.SetMargined(true)

	vbCounter := ui.NewVerticalBox()
	vbCounter.SetPadded(true)

	countLabel = ui.NewLabel(fmt.Sprintf("%d", count))

	vbCounter.Append(countLabel, false)
	countGroup.SetChild(vbCounter)

	vbContainer.Append(inputGroup, false)
	vbContainer.Append(messageGroup, false)
	vbContainer.Append(countGroup, false)

	mainWindow.SetChild(vbContainer)

	showMessageButton.OnClicked(func(*ui.Button) {
		// Update the UI directly as it is called from the main thread
		messageLabel.SetText(message.Text())
	})

	clearMessageButton.OnClicked(func(*ui.Button) {
		// Update the UI directly as it is called from the main thread
		messageLabel.SetText("")
	})

	mainWindow.Show()

	// Counting and updating the UI from another goroutine
	go counter()
}

func counter() {
	for {
		time.Sleep(1 * time.Second)
		count++

		// Update the UI using the QueueMain function
		ui.QueueMain(func() {
			countLabel.SetText(fmt.Sprintf("%d", count))
		})
	}
}

func main() {
	count = 0

	ui.Main(setupUI)
}

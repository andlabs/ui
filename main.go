// 11 february 2014
package main

func main() {
	w := NewWindow("Main Window", 320, 240)
	w.Closing = make(chan struct{})
	b := NewButton("Click Me")
	w.SetControl(b)
	err := w.Open()
	if err != nil {
		panic(err)
	}
mainloop:
	for {
		select {
		case <-w.Closing:
			break mainloop
		case <-b.Clicked:
			println("clicked")
		}
	}
	w.Close()
}


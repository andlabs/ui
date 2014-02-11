// 11 february 2014
package main

func main() {
	w := NewWindow("Main Window")
	w.Closing = make(chan struct{})
	w.Open()
	<-w.Closing
	w.Close()
}


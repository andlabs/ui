// 11 december 2015

package ui

import "time"
import "testing"

func TestIt(t *testing.T) {
	err := Main(func() {
		w := NewWindow("Hello", 320, 240, false)
		stop := make(chan struct{})
		w.OnClosing(func(w *Window) bool {
			stop <- struct{}{}
			Quit()
			return true
		})
		p := NewProgressBar()
		w.SetChild(p)
		go func() {
			value := 0
			ticker := time.NewTicker(time.Second / 2)
			for {
				select {
				case <-ticker.C:
					QueueMain(func() {
						value++
						if value > 100 {
							value = 0
						}
						p.SetValue(value)
					})
				case <-stop:
					return
				}
			}
		}()
		w.Show()
	})
	if err != nil {
		t.Fatal(err)
	}
}

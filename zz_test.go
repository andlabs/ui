// 11 december 2015

package ui

import "testing"

func TestIt(t *testing.T) {
	err := Main(func() {
		t.Log("we're here")
		Quit()
	})
	if err != nil {
		t.Fatal(err)
	}
}

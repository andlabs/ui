// +build !windows,!darwin,!plan9

// 16 february 2014
//package ui
package main

import (
	"fmt"
)

type sysData struct {
	cSysData

	widget		*gtkWidget
	container		*gtkWidget	// for moving
}

type classData struct {
	make	func() *gtkWidget
	setText	func(widget *gtkWidget, text string)
	// ...
	signals	map[string]func(*sysData) func() bool
}

var classTypes = [nctypes]*classData{
	c_window:	&classData{
		make:	gtk_window_new,
		setText:	gtk_window_set_title,
		signals:	map[string]func(*sysData) func() bool{
			"delete-event":		func(w *sysData) func() bool {
				return func() bool {
					if w.event != nil {
						w.event <- struct{}{}
					}
					return true		// do not close the window
				}
			},
		},
	},
	c_button:		&classData{
//		make:	gtk_button_new,
	},
	c_checkbox:	&classData{
	},
	c_combobox:	&classData{
	},
	c_lineedit:	&classData{
	},
	c_label:		&classData{
	},
	c_listbox:		&classData{
	},
}

func (s *sysData) make(initText string, window *sysData) error {
	ct := classTypes[s.ctype]
	if ct.make == nil {		// not yet implemented
		println(s.ctype, "not implemented")
		return nil
	}
	ret := make(chan *gtkWidget)
	defer close(ret)
	uitask <- func() {
		ret <- ct.make()
	}
	s.widget = <-ret
println(s.widget)
	if window == nil {
		uitask <- func() {
			fixed := gtk_fixed_new()
			gtk_container_add(s.widget, fixed)
			for signal, generator := range ct.signals {
				g_signal_connect(s.widget, signal, generator(s))
			}
			ret <- fixed
		}
		s.container = <-ret
	} else {
		s.container = window.container
	}
	err := s.setText(initText)
	if err != nil {
		return fmt.Errorf("error setting initial text of new window/control: %v", err)
	}
	return nil
}

func (s *sysData) show() error {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		gtk_widget_show(s.widget)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) hide() error {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		gtk_widget_hide(s.widget)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) setText(text string) error {
if classTypes[s.ctype] == nil || classTypes[s.ctype].setText == nil { return nil }
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].setText(s.widget, text)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) setRect(x int, y int, width int, height int) error {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		gtk_fixed_put(s.container, s.widget, x, y)
		gtk_widget_set_size_request(s.widget, width, height)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) isChecked() bool {
	// TODO
	return false
}

func (s *sysData) text() string {
	// TODO
	return ""
}

func (s *sysData) append(what string) error {
	// TODO
	return nil
}

func (s *sysData) insertBefore(what string, before int) error {
	// TODO
	return nil
}

func (s *sysData) selectedIndex() int {
	// TODO
	return -1
}

func (s *sysData) selectedIndices() []int {
	// TODO
	return nil
}

func (s *sysData) selectedTexts() []string {
	// TODO
	return nil
}

func (s *sysData) setWindowSize(width int, height int) error {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		gtk_window_resize(s.widget, width, height)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) delete(index int) error {
	// TODO
	return nil
}

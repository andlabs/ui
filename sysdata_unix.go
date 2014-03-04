// +build !windows,!darwin,!plan9

// 16 february 2014

//
package ui

import (
	"fmt"
)

type sysData struct {
	cSysData

	widget    *gtkWidget
	container *gtkWidget // for moving
}

type classData struct {
	make     func() *gtkWidget
	makeAlt  func() *gtkWidget
	setText  func(widget *gtkWidget, text string)
	text     func(widget *gtkWidget) string
	append   func(widget *gtkWidget, text string)
	insert   func(widget *gtkWidget, index int, text string)
	selected func(widget *gtkWidget) int
	selMulti func(widget *gtkWidget) []int
	smtexts  func(widget *gtkWidget) []string
	delete   func(widget *gtkWidget, index int)
	// ...
	signals map[string]func(*sysData) func() bool
}

var classTypes = [nctypes]*classData{
	c_window: &classData{
		make:    gtk_window_new,
		setText: gtk_window_set_title,
		text:    gtk_window_get_title,
		signals: map[string]func(*sysData) func() bool{
			"delete-event": func(s *sysData) func() bool {
				return func() bool {
					s.signal()
					return true // do not close the window
				}
			},
			"configure-event": func(s *sysData) func() bool {
				return func() bool {
					if s.container != nil && s.resize != nil { // wait for init
						width, height := gtk_window_get_size(s.widget)
						// top-left is (0,0) so no need for winheight
						err := s.resize(0, 0, width, height, 0)
						if err != nil {
							panic("child resize failed: " + err.Error())
						}
					}
					// returning false indicates that we continue processing events related to configure-event; if we choose not to, then after some controls have been added, the layout fails completely and everything stays in the starting position/size
					// TODO make sure this is the case
					return false
				}
			},
		},
	},
	c_button: &classData{
		make:    gtk_button_new,
		setText: gtk_button_set_label,
		text:    gtk_button_get_label,
		signals: map[string]func(*sysData) func() bool{
			"clicked": func(s *sysData) func() bool {
				return func() bool {
					s.signal()
					return true // do not close the window
				}
			},
		},
	},
	c_checkbox: &classData{
		make:    gtk_check_button_new,
		setText: gtk_button_set_label,
	},
	c_combobox: &classData{
		make:    gtk_combo_box_text_new,
		makeAlt: gtk_combo_box_text_new_with_entry,
		// TODO setText
		text:     gtk_combo_box_text_get_active_text,
		append:   gtk_combo_box_text_append_text,
		insert:   gtk_combo_box_text_insert_text,
		selected: gtk_combo_box_get_active,
		delete:   gtk_combo_box_text_remove,
	},
	c_lineedit: &classData{
		make:    gtk_entry_new,
		makeAlt: gtkPasswordEntryNew,
		setText: gtk_entry_set_text,
		text:    gtk_entry_get_text,
	},
	c_label: &classData{
		make:    gtk_label_new,
		setText: gtk_label_set_text,
		text:    gtk_label_get_text,
	},
	c_listbox: &classData{
		make:    gListboxNewSingle,
		makeAlt: gListboxNewMulti,
		// TODO setText
		text:     gListboxText,
		append:   gListboxAppend,
		insert:   gListboxInsert,
		selected: gListboxSelected,
		selMulti: gListboxSelectedMulti,
		smtexts:  gListboxSelMultiTexts,
		delete:   gListboxDelete,
	},
	c_progressbar: &classData{
		make: gtk_progress_bar_new,
	},
}

func (s *sysData) make(initText string, window *sysData) error {
	ct := classTypes[s.ctype]
	ret := make(chan *gtkWidget)
	defer close(ret)
	uitask <- func() {
		if s.alternate {
			ret <- ct.makeAlt()
			return
		}
		ret <- ct.make()
	}
	s.widget = <-ret
	if window == nil {
		uitask <- func() {
			fixed := gtk_fixed_new()
			gtk_container_add(s.widget, fixed)
			// TODO return the container before assigning the signals?
			for signal, generator := range ct.signals {
				g_signal_connect(s.widget, signal, generator(s))
			}
			ret <- fixed
		}
		s.container = <-ret
	} else {
		s.container = window.container
		uitask <- func() {
			gtk_container_add(s.container, s.widget)
			for signal, generator := range ct.signals {
				g_signal_connect(s.widget, signal, generator(s))
			}
			ret <- nil
		}
		<-ret
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
	if classTypes[s.ctype].setText == nil { // does not have concept of text
		return nil
	}
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].setText(s.widget, text)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) setRect(x int, y int, width int, height int, winheight int) error {
	gtk_fixed_move(s.container, s.widget, x, y)
	gtk_widget_set_size_request(s.widget, width, height)
	return nil
}

func (s *sysData) isChecked() bool {
	ret := make(chan bool)
	defer close(ret)
	uitask <- func() {
		ret <- gtk_toggle_button_get_active(s.widget)
	}
	return <-ret
}

func (s *sysData) text() string {
	ret := make(chan string)
	defer close(ret)
	uitask <- func() {
		ret <- classTypes[s.ctype].text(s.widget)
	}
	return <-ret
}

func (s *sysData) append(what string) error {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].append(s.widget, what)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) insertBefore(what string, before int) error {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].insert(s.widget, before, what)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) selectedIndex() int {
	ret := make(chan int)
	defer close(ret)
	uitask <- func() {
		ret <- classTypes[s.ctype].selected(s.widget)
	}
	return <-ret
}

func (s *sysData) selectedIndices() []int {
	ret := make(chan []int)
	defer close(ret)
	uitask <- func() {
		ret <- classTypes[s.ctype].selMulti(s.widget)
	}
	return <-ret
}

func (s *sysData) selectedTexts() []string {
	ret := make(chan []string)
	defer close(ret)
	uitask <- func() {
		ret <- classTypes[s.ctype].smtexts(s.widget)
	}
	return <-ret
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
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].delete(s.widget, index)
		ret <- struct{}{}
	}
	<-ret
	return nil
}

func (s *sysData) setProgress(percent int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		gtk_progress_bar_set_fraction(s.widget, percent)
		ret <- struct{}{}
	}
	<-ret
}

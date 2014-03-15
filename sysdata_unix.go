// +build !windows,!darwin,!plan9

// 16 february 2014

package ui

import (
	"time"
)

type sysData struct {
	cSysData

	widget		*gtkWidget
	container		*gtkWidget	// for moving
	pulse		chan bool		// for sysData.progressPulse()
}

type classData struct {
	make	func() *gtkWidget
	makeAlt	func() *gtkWidget
	setText	func(widget *gtkWidget, text string)
	text		func(widget *gtkWidget) string
	append	func(widget *gtkWidget, text string)
	insert	func(widget *gtkWidget, index int, text string)
	selected	func(widget *gtkWidget) int
	selMulti	func(widget *gtkWidget) []int
	smtexts	func(widget *gtkWidget) []string
	delete	func(widget *gtkWidget, index int)
	len		func(widget *gtkWidget) int
	// ...
	signals	callbackMap
}

var classTypes = [nctypes]*classData{
	c_window:		&classData{
		make:		gtk_window_new,
		setText:		gtk_window_set_title,
		text:			gtk_window_get_title,
		signals:		callbackMap{
			"delete-event":		window_delete_event_callback,
			"configure-event":	window_configure_event_callback,
		},
	},
	c_button:			&classData{
		make:		gtk_button_new,
		setText:		gtk_button_set_label,
		text:			gtk_button_get_label,
		signals:		callbackMap{
			"clicked":		button_clicked_callback,
		},
	},
	c_checkbox:		&classData{
		make:		gtk_check_button_new,
		setText:		gtk_button_set_label,
		text:			gtk_button_get_label,
	},
	c_combobox:		&classData{
		make:		gtk_combo_box_text_new,
		makeAlt:		gtk_combo_box_text_new_with_entry,
		// TODO setText
		text:			gtk_combo_box_text_get_active_text,
		append:		gtk_combo_box_text_append_text,
		insert:		gtk_combo_box_text_insert_text,
		selected:		gtk_combo_box_get_active,
		delete:		gtk_combo_box_text_remove,
		len:			gtkComboBoxLen,
	},
	c_lineedit:		&classData{
		make:		gtk_entry_new,
		makeAlt:		gtkPasswordEntryNew,
		setText:		gtk_entry_set_text,
		text:			gtk_entry_get_text,
	},
	c_label:			&classData{
		make:		gtk_label_new,
		setText:		gtk_label_set_text,
		text:			gtk_label_get_text,
	},
	c_listbox:			&classData{
		make:		gListboxNewSingle,
		makeAlt:		gListboxNewMulti,
		// TODO setText
		text:			gListboxText,
		append:		gListboxAppend,
		insert:		gListboxInsert,
		selMulti:		gListboxSelectedMulti,
		smtexts:		gListboxSelMultiTexts,
		delete:		gListboxDelete,
		len:			gListboxLen,
	},
	c_progressbar:		&classData{
		make:		gtk_progress_bar_new,
	},
	c_area:			&classData{
		make:		gtkAreaNew,
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
			for signame, sigfunc := range ct.signals {
				g_signal_connect(s.widget, signame, sigfunc, s)
			}
			ret <- fixed
		}
		s.container = <-ret
	} else {
		s.container = window.container
		uitask <- func() {
			gtk_container_add(s.container, s.widget)
			for signame, sigfunc := range ct.signals {
				g_signal_connect(s.widget, signame, sigfunc, s)
			}
			ret <- nil
		}
		<-ret
	}
	s.setText(initText)
	return nil
}

// used for Windows; nothing special needed elsewhere
func (s *sysData) firstShow() error {
	s.show()
	return nil
}

func (s *sysData) show() {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		gtk_widget_show(s.widget)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) hide() {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		gtk_widget_hide(s.widget)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) setText(text string) {
	if classTypes[s.ctype].setText == nil {		// does not have concept of text
		return
	}
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].setText(s.widget, text)
		ret <- struct{}{}
	}
	<-ret
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

func (s *sysData) append(what string) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].append(s.widget, what)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) insertBefore(what string, before int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].insert(s.widget, before, what)
		ret <- struct{}{}
	}
	<-ret
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

func (s *sysData) delete(index int) {
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		classTypes[s.ctype].delete(s.widget, index)
		ret <- struct{}{}
	}
	<-ret
}

// With GTK+, we must manually pulse the indeterminate progressbar ourselves. This goroutine does that.
func (s *sysData) progressPulse() {
	pulse := func() {
		ret := make(chan struct{})
		defer close(ret)
		uitask <- func() {
			gtk_progress_bar_pulse(s.widget)
			ret <- struct{}{}
		}
		<-ret
	}

	var ticker *time.Ticker
	var tickchan <-chan time.Time

	// the default on Windows
	const pulseRate = 30 * time.Millisecond

	for {
		select {
		case start := <-s.pulse:
			if start {
				ticker = time.NewTicker(pulseRate)
				tickchan = ticker.C
				pulse()			// start the pulse animation now, not 30ms later
			} else {
				if ticker != nil {
					ticker.Stop()
				}
				ticker = nil
				tickchan = nil
				s.pulse <- true		// notify sysData.setProgress()
			}
		case <-tickchan:
			pulse()
		}
	}
}

func (s *sysData) setProgress(percent int) {
	if s.pulse == nil {
		s.pulse = make(chan bool)
		go s.progressPulse()
	}
	if percent == -1 {
		s.pulse <- true
		return
	}
	s.pulse <- false
	<-s.pulse			// wait for sysData.progressPulse() to register that
	ret := make(chan struct{})
	defer close(ret)
	uitask <- func() {
		gtk_progress_bar_set_fraction(s.widget, percent)
		ret <- struct{}{}
	}
	<-ret
}

func (s *sysData) len() int {
	ret := make(chan int)
	defer close(ret)
	uitask <- func() {
		ret <- classTypes[s.ctype].len(s.widget)
	}
	return <-ret
}

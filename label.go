// 14 february 2014

package ui

// A Label is a static line of text used to mark other controls.
// Label text is drawn on a single line; text that does not fit is truncated.
// A Label can appear in one of two places: bound to a control or standalone.
// This determines the vertical alignment of the label.
type Label struct {
	created  bool
	sysData  *sysData
	initText string
	standalone	bool
}

// NewLabel creates a new Label with the specified text.
// The label is set to be bound to a control, so its vertical position depends on its vertical cell size in an implementation-defined manner.
func NewLabel(text string) *Label {
	return &Label{
		sysData:  mksysdata(c_label),
		initText: text,
	}
}

// NewStandaloneLabel creates a new Label with the specified text.
// The label is set to be standalone, so its vertical position will always be at the top of the vertical space assigned to it.
func NewStandaloneLabel(text string) *Label {
	return &Label{
		sysData:  mksysdata(c_label),
		initText: text,
		standalone:	true,
	}
}

// SetText sets the Label's text.
func (l *Label) SetText(text string) {
	if l.created {
		l.sysData.setText(text)
		return
	}
	l.initText = text
}

// Text returns the Label's text.
func (l *Label) Text() string {
	if l.created {
		return l.sysData.text()
	}
	return l.initText
}

func (l *Label) make(window *sysData) error {
	l.sysData.alternate = l.standalone
	err := l.sysData.make(window)
	if err != nil {
		return err
	}
	l.sysData.setText(l.initText)
	l.created = true
	return nil
}

func (l *Label) allocate(x int, y int, width int, height int, d *sysSizeData) []*allocation {
	return []*allocation{&allocation{
		x:       x,
		y:       y,
		width:   width,
		height:  height,
		this:		l,
	}}
}

func (l *Label) preferredSize(d *sysSizeData) (width int, height int) {
	return l.sysData.preferredSize(d)
}

func (l *Label) commitResize(a *allocation, d *sysSizeData) {
	l.sysData.commitResize(a, d)
}

func (l *Label) getAuxResizeInfo(d *sysSizeData) {
	l.sysData.getAuxResizeInfo(d)
}

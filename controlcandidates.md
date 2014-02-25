WINDOWS
- DateTime Picker
- ListView for Tables
- either Property Sheets or Tabs for Tabs
- either Rebar or Toolbar for Toolbars
- Status Bar
- Tooltip (should be a property of each control)
- Trackbar for Sliders
	- cannot automatically snap to custom step; need to do it manually
- Tree View
- Up-Down Control for Spinners
- maybe:
	- swap ComboBox for ComboBoxEx (probably only if requested enough)
	- IP Address control (iff GTK+ and Cocoa have it; maybe not necessary if we allow arbitrary target addresses?)
	- ListView for its Icon View?
	- something similar to Task Dialog might be useful to have as a convenience template later

GTK+
- GtkNotebook for Tabs
- GtkScale for Sliders
	- cannot automatically snap to INTEGERS (let alone to custom steps); need to do it manually
	- natural size is 0x0 for some reason
- GtkSpinButton for Spinners
- GtkStatusBar
- GtkToolbar
- maybe:
	- GtkFontButton would be nice but unless ComboBoxEx provides it Windows doesn't
		- same for GtkColorButton
	- GtkIconView
	- GtkSeparator (I think Windows makes this a mode of Static controls?)
- notes to self:
	- groupbox is GtkFrame
	- GtkTreeView can do tree views and Tables

COCOA
- NSOutlineView for tree views
- NSSlider for Sliders
- NSStatusBar
- NSStepper for Spinners
- NSTabView for Tabs
- NSTableView for Tables
- NSToolbar
- maybe:
	- NSBrowser seems nice...???
	- NSCollectionView for Icon View?
	- NSColorWell is the color button
	- NSOpenGLView for OpenGL; need to see how much OpenGL-specific stuff I need to expose
	- NSRuleEditor/NSPredicateEditor look nice too but
- notes to self:
	- groupbox is NSBox
	- non-editable combobox is NSPopUpButton
	- editable combobox is NSCombobox
	- don't look at NSForm; though it arranges in the ideal form layout, it only allows single-line text entry fields as controls
	- NSSecureTextField does password entries
	- NSProgressIndicator for ProgressBar
- TODO:
	- what does NSPathControl look like?

# Slider Capabilities
Capability | Windows | GTK+ | Cocoa
-|-|-|-
Data Type | int | float | TODO
Can Simulate ints? | yes | TODO | TODO
Mouse Step Snap | 1, fixed | something; likely 0.1 but not sure | TODO
Keyboard Step Snap | configurable | configurable | TODO
Current Value Display | tooltip during drag | label, always visible | TODO
Ticks | configurable display, configurable interval | TODO | TODO
Can Catch Mouse Events to Snap? | I think this is how to do it | TODO | TODO
Preferred Size | given in UI guidelines | natural: 0x0; minimum: TODO | TODO

# Spinner Capabilities
Capability | Windows | GTK+ | Cocoa
-|-|-|-
Data Type | int | float | TODO
Can Simulate ints? | yes | yes | TODO
Mouse Step Snap | 1, fixed | configurable | TODO
Keyboard Step Snap | 1, fixed | configurable (uses same value as mouse) | TODO
Can Catch Events To Snap? | TODO | no need | TODO
Preferred Size | TODO | TODO | TODO

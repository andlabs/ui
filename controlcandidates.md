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
- TODO
	- commcntl.h has stuff on a font control that isn't documented?
		- actually not a control, but localization support: http://msdn.microsoft.com/en-us/library/windows/desktop/bb775454%28v=vs.85%29.aspx
- notes to self:
	- OpenGL: http://msdn.microsoft.com/en-us/library/windows/desktop/dd374379%28v=vs.85%29.aspx

GTK+
- GtkCalendar for date selection (TODO doesn't handle times)
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
	- OpenGL is done outside GTK+: https://projects.gnome.org/gtkglext/
		- only an issue if I want to provide OpenGL by default...
		- http://stackoverflow.com/questions/3815806/gtk-and-opengl-bindings suggest GtkGLArea is better but that seems to be a Mono thing? also indicates Clutter (with its Cogl) is not an option because it locks you out of using the OpenGL API directly
			- er no, the Mono thing is just the homepage... but it doesn't say if this targets GTK+ 2 or GTK+ 3, hm. (also it appears to not have been updated since Precise; in Ubuntu it's libgtkgl)
			- and gtkglext doesn't support GTK+ 3 officially anyway
			- and cairo doesn't seem to support OpenGL explicitly so it looks like I will need to communicate with glx directly: http://stackoverflow.com/questions/17628241/how-can-i-use-gtk3-and-opengl-together
				- except replace glx with EGL/GLES2 because of Wayland: http://wayland.freedesktop.org/faq.html#heading_toc_j_0 (assuming EGL/GLES2 can work on X11)

COCOA
- NSDatePicker for date/time selection
- NSOutlineView for tree views
- NSSlider for Sliders
- NSStatusBar
- NSStepper for Spinners
	- TODO does this require me to manually pair it with a single-line text entry field?
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
	- don't look at NSForm; though it arranges in the ideal form layout, it only allows single-line text entry fields as controls
- TODO:
	- what does NSPathControl look like?

# Slider Capabilities
Capability | Windows | GTK+ | Cocoa
----- | ----- | ----- | -----
Data Type | int | float | float
Can Simulate ints? | yes | TODO | TODO
Mouse Step Snap | 1, fixed | something; likely 0.1 but not sure | yes (`setAllowsTickMarkValuesOnly:`); caveat: must specify an exact number of ticks (see below)
Keyboard Step Snap | configurable | configurable | TODO (same as mouse?)
Current Value Display | tooltip during drag | label, always visible | TODO
Tooltips? | TODO | TODO | TODO
Ticks | configurable display, configurable interval | TODO | configurable display; configurable COUNT (not interval!)
Can Catch Mouse Events to Snap? | I think this is how to do it | TODO | TODO
Preferred Size | given in UI guidelines | natural: 0x0; minimum: TODO | TODO

# Spinner Capabilities
Capability | Windows | GTK+ | Cocoa
----- | ----- | ----- | -----
Data Type | int | float | flaot
Can Simulate ints? | yes | yes | TODO
Mouse Step Snap | 1, fixed | configurable | configurable
Keyboard Step Snap | 1, fixed | configurable (uses same value as mouse) | TODO (same as mouse?)
Can Catch Events To Snap? | TODO | no need | TODO
Preferred Size | TODO | TODO | TODO


# Dialog box hijack
## Open/Save Dialogs
  | Windows | GTK+ | Cocoa
----- | ----- | ----- | -----
Directories | xxx | open and save | xxx
Network vs. local only (URI vs. filename) | xxx | yes (default local only; if local only, changing to, say, smb://127.0.0.1/ will pop up an error box; if nonlocal allowed, filename can be null) | xxx
Multiple selection | yes | yes | xxx
Hidden files | xxx | hidden by default; can be switched on in code (but is a no-op?) and also by the user | xxx
Overwrite confirmation | xxx | available; must be explicitly enabled | xxx
New Folder button | xxx | optional (I think enabled by default? should do it explicitly to be safe, anyway) | xxx
Preview widget | xxx | user-defined | xxx
Extra widget | xxx | user-defined | xxx
File filters | Specified by "patterns" (consisting of filename characters and * but not space; I assume the only safe ones are *.ext and *.*); multiple patterns separated by semicolons; can have custom labels | Specified by MIME type (handles subtypes; has wildcards) or pattern ("shell-style glob", so I assume over whole basename) or by custom function; can have multiple of the above; can have custom labels; also has a shortcut to add all gdk-pixbuf-supported formats | xxx
File filter list format | `"Label\0Filter-list\0Label\0Filter-list\0...Label\0FIlter-list\0\0"`; filter for all files is canonically `"All Files\0*.*\0\0"` in the docs (specifically this due to handling of shortcut links); also provides a way for users to write in their own filters | Add or remove individual GtkFileFIlter objects; can select one specified in the list to show by default; default behavior is all files; if selected one when none has been specified, filter selection disabled; filter for all files specified in docs under gtk_file_filter_new() (except doesn't set a name) | xxx
Default file name | settable | settable | xxx
Initial directory | complex rules that have changed over time; we can pass an absolute filename (the previous filename or a default filename) and have its path used (if we specify just a path it will either be used as the filename or the program will crash); or we can give it a directory; or Windows will remember for us for some time, or... | pass previous filename or URI to show; overrides default file name; intended only for saving files (so I don't know if it's possible to remember current directory for opening??????); effect of passing containing directory undocumented(???? in my tests the given folder itself is selected) | xxx
Confirmation and cancel buttons | xxx | GTK_STOCK_OPEN, GTK_STOCK_SAVE, GTK_STOCK_SAVE_AS / GTK_STOCK_CANCEL | xxx
Returned filename rules | xxxx | memory provided by GTK+ itself (so no need to worry about size limits); can return a single filename or URI or a GSList of filenames or URIs | xxx
Window title | optional; defaults to either Open or Save As | required(?) | xxx

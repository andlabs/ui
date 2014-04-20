general list:
- Window.SizeToFit() or WIndow.OptimalSize() (use: `Window.SetOptimalSize())`) for sizing a window to the control's interest
	- with the current code, will be a bit of a kludge, because preferredSize() assumes it's running on the main thread without locks
- Control.Show()/Control.Hide()
- Groupbox
- character-limited entry fields (not for passwords), numeric entry fields, multiline entry fields
	- possible rename of LineEdit?
		- especially for password fields - NewPasswordEntry()?
- padding and spacing in Stack
- allow Combobox to have initial settings
- Combobox and Listbox insertions and deletions should allow bulk (...string/...int)
- Combobox/Listbox.DeleteAll
- Combobox/Listbox.Select (with Listbox.Select allowing bulk)
	- Checkbox.Check or Checkbox.SetChecked
- Listbox.SelectAll
- Listbox/Combobox.Index(n)
	- Index(n) is the name used by reflect.Value; use a different one?
- figure out where to auto-place windows in Cocoa (also window coordinates are still not flipped properly so (0,0) on screen is the bottom-left)
	- also provide a method to center windows; Cocoa provides one for us but
		- GTK+ too: gtk_window_set_position(window, GTK_WIN_POS_CENTER) (via http://stackoverflow.com/questions/16832581/how-to-make-a-gtkwindow-background-transparent-on-linux)
- make Combobox and Listbox satisfy sort.Interface?
- should a noneditable Combobox be allowed to return to unselected mode by the user?
- provide a way for MouseEvent/KeyEvent to signal that the keypress caused the Area to gain/lose focus
	- provide an event for leaving focus so a focus rectangle can be drawn
- when adding menus:
	- provide automated About, Preferneces, and Quit that place these in the correct location
		- Quit should pulse AppQuit
- will probably want to bring back Event() as NewEvent() should that facility be necesary for menus, etc.

issues of policy:
- LineEdit heights on Windows seem too big; either that or LineEdit, Button, and Label text is not vertically centered properly
	- are Checkboxes and Comboboxes too small?
- consolidate scroll view code in GTK+? it's not a lot, but it's rather divergent...

problem points:
- because the main event loop is not called if initialization fails, it is presently impossible for MsgBoxError() to work if UI initialization fails; this basically means we cannot allow initializiation to fail on Mac OS X if we want to be able to report UI init failures to the user with one (which would be desirable, maybe (would violate Windows HIG?))
- make sure GTK+ documentation version point differences (x in 4.3.x) don't matter
	- I found a GTK+ version number meaning page somewhere; have to find it again (TODO)

twists of fate:
- listboxes spanning the vertical height of the window don't always align with the bottom border of the edit control attached to the bottom of the window...
	- this specifically only happens when the window has an odd height; I don't think this can be fixed unless we explicitly ignore the extra pixel everywhere
- need a way to get ideal size for all controls on Windows, not just push buttons (Microsoft...)
- Cocoa controls have padding around them; this padding is **opaque** so we can't just use the control's cell rect and some shuffling around

style changes:
- make specific wording in documentation consistent (make/create, etc.)
	- document minor details like wha thappens on specific events so that they are guaranteed to work the same on all platforms (are there any left?)
		- what happens when the user clicks and drags on a listbox
			- I think this is a platform behavior...
	- should field descriptions in method comments include the receiver name? (for instance e.Held vs. Held) - see what Go's own documentation does
	- need to figure out exactly how to indicate that a struct{}{} is sent on an event channel (I use about six or so different wordings so far...)
	- "package ui", "the package"
		- also "library" both in docs and comments and code, etc.
- make passing of parameters and type conversions of parameters to uitask on Windows consistent: explicit _WPARAM(xxx)/_LPARAM(xxx)/uintptr(xxx), for example
	- do this for type signatures in exported functions: (err error) or just error?
	- do this for the names of GTK+ helper functions (gtkXXX or gXXX)
	- areaView -> goArea (or change the class name to be like the delegate name?) in area_darwin.go?

far off:
- localization
- strip unused constants from the Windows files
- combine more Windows files; rename some?
- tab stops
	- http://blogs.msdn.com/b/oldnewthing/archive/2003/10/21/55384.aspx
- rename Stack to Box?
- maybe change multiple selection lists to checkbox lists?
	- windows HIG refernece: http://msdn.microsoft.com/en-us/library/windows/desktop/aa511485.aspx - conflicting, confusing info
	- gtk+ HIG reference: https://developer.gnome.org/hig-book/3.4/controls-lists.html.en
	- mac HIG reference: ???
- go over the old new thing's scrollbar series to make sure I'm doing everything right with scrollbars in Windows Areas
- change the MsgBox() calls to encourage good alert dialog design??????? maybe? TODO
- make gcc (Unix)/clang (Mac OS X) pedantic about warnings/errors; also -Werror
	- problem: cgo-generated files trip -Werror up; I can't seem to turn off unused argument warnings with the -Wall/-Wextra/-pedantic options

big things:
- make sure every sysData function only performs a single invocation to uitask; see http://blogs.msdn.com/b/oldnewthing/archive/2005/10/10/479124.aspx#479182
	- windows: this requires major restructuring
	- gtk, mac: this just requires checking
- steamroll ALL errors, especially on windows
	- gtk: no way to catch errors
	- cocoa: discouraged
- make fully lock free
	- prerequisite is the above two
	- locks are used because of initial state; we can override by creating controls at construct time
		- cocoa, gtk: no real issues
		- windows: now required to specify no parent window at create time and set both the parent window AND the child window ID later
			- http://msdn.microsoft.com/en-us/library/windows/desktop/ms633541%28v=vs.85%29.aspx
				- don't worry about UI state messages yet; this is before opening the UI anyway (these might be needed when we add tab stops)
			- http://msdn.microsoft.com/en-us/library/windows/desktop/ms644898%28v=vs.85%29.aspx GWLP_ID
- preferred sizes in general are awkward: on Windows, no text-based size calculation is performed, so we have weird things like Labels always having the same width (so if you place a Label in a Stack by itself and forget to make it stretchy, it'll be truncated on Windows (but not on GTK+ or OS X?!))

big dumb things:
- for our two custom window classes on Windows, we should allocate extra space in the window class's info structure and then use SetWindowLongPtrW() during WM_CREATE to store the sysData and not have to make a new window class each time; this might also fix the s != nil && s.hwnd != 0 special cases in the Area WndProc if done right
	- references: https://github.com/glfw/glfw/blob/master/src/win32_window.c#L182, http://www.catch22.net/tuts/custom-controls
	- this is a bit flakier as SetWindowLongPtr() can fail, and it can also succeed in such a way that the last error is unreliable
		- and doesn't exist on 32-bit windows; we will need special code for detecting 32-bit/64-bit (see http://bugs.winehq.org/show_bug.cgi?id=30556 and GerbilSoft in irc.badnik.net/#retro)
		- also CreateWindow() and CreateWindowEx() docs differ in indicating which messages are sent but ultimately send the same set; WM_GETMINMAXINFO is sent first so that throws a wrench in the whole point, AND we'll need a way to properly differentiate between custom classes and controls...
			- tl;dr what started as a somewhat quick change was really way too much effort for only potential/theoretical gain; approach if someone actually DOES hit Go's syscall.NewCallback() limit
			- theoretically I could just have a s != nil check in WM_GETMINMAXINFO only, but not too hot on the idea :/
				- raymond chen does it here: http://blogs.msdn.com/b/oldnewthing/archive/2005/04/22/410773.aspx (check the implementation of Window::s_WndProc())
				- ...and suggests we do it here http://blogs.msdn.com/b/oldnewthing/archive/2014/02/03/10496248.aspx (**NOTE THE DATE**) - the comments on this one provide some potential ideas, including IIntrospect's comment about HCBT_CREATEWND; later Raymond says we should not worry about SetWindowLongPtr() failing
				- and raymond suggests GWL_USERDATA here: http://blogs.msdn.com/b/oldnewthing/archive/2005/03/03/384285.aspx
- listboxes should have horizontal scrollbars on all platforms; this is way too hard on OS X and doesn't work; my code is in experiments/
	- also moved the Windows code there for the sake of efficiency
	- GTK+ works just fine though
- window sizes need to not include the window decoration; while Mac OS X and GTK+ both obey this, I've only had issues with Windows; check the experiments/ folder
	- also will need to be documented in window.go

specifics:

WINDOWS
<br>TODO re-evaluate this list for Common Controls 6-only stuff (controls and features) I passed up on the first itme
- DateTime Picker
- ListView for Tables
- either Property Sheets or Tabs for Tabs
- either Rebar or Toolbar for Toolbars
- Status Bar
- Tooltip (should be a property of each control)
- Trackbar for Sliders
	- cannot automatically snap to custom step; need to do it manually
- Tree View
	- http://blogs.msdn.com/b/oldnewthing/archive/2011/11/25/10241394.aspx
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
	- groupbox is a mode of the BUTTON class (????)
	- OpenGL: http://msdn.microsoft.com/en-us/library/windows/desktop/dd374379%28v=vs.85%29.aspx
	- don't use ES_NUMBER for number-only text boxes, as "it is still possible to paste non-digits into the edit control." (though a commenter on MSDN says that's wrong?)
		- might want to just have spinners and not numeric text boxes???

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

ALL THREE; LIST BY FEATURE (TODO merge the above into this)
- LineEdit prompts
	- Windows: http://msdn.microsoft.com/en-us/library/windows/desktop/bb761639%28v=vs.85%29.aspx
		being able to use comctl32.dll version 6 now that win2k support is gone is a good thing, isn't it?
	- GTK+: TODO
	- Cocoa: TODO

OTHER NOTES
- on webpage controls: http://blogs.msdn.com/b/oldnewthing/archive/2005/02/11/371042.aspx

THOUGHTS
- privilege escalation requests
	- windows vista and up use a button with a shield icon and osme other API that I don't know before performing the privileged action
	- typically graphical unix programs nowadays use a graphical sudo prompt provided by the desktop environment (gksu or kdesudo or what not) either before the program even runs or after clicking a button to start the process with elevated privileges
		- however, GTK+ also provides https://developer.gnome.org/gtk3/3.4/GtkLockButton.html, which is similar to the Cocoa one below:
	- Cocoa uses a button with a lock on it that will disable the entire form when unlocked, then ask for permission when clicked

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
Directories | no (separate facility provided by the shell) | open and save | open only
Network vs. local only (URI vs. filename) | Network button enabled by default; can be switched off (**TODO** how are network filenames returned?) | yes (default local only; if local only, changing to, say, smb://127.0.0.1/ will pop up an error box; if nonlocal allowed, filename can be null) | xxx
Multiple selection | yes | yes | open only
Hidden files | user-specified; can be switched on in code (but is a no-op?) | hidden by default; can be switched on in code (but is a no-op?) and also by the user | xxx
Overwrite confirmation | available; must be explicitly enabled | available; must be explicitly enabled | xxx
New Folder button | xxx | optional (I think enabled by default? should do it explicitly to be safe, anyway) | optional
Preview widget | xxx | yes; optional, custom | xxx
Extra custom widget | xxx | yes; optional | yes; optional
File filters | Specified by "patterns" (consisting of filename characters and * but not space; I assume the only safe ones are *.ext and *.*); multiple patterns separated by semicolons; can have custom labels | Specified by MIME type (handles subtypes; has wildcards) or pattern ("shell-style glob", so I assume over whole basename) or by custom function; can have multiple of the above; can have custom labels; also has a shortcut to add all gdk-pixbuf-supported formats | Specified by "UTI"s or by individual filename extensions (format not documented but appears to be just the extension without embellishments); cannot have labels; 1:1 filter:extension mapping.
File filter list format | `"Label\0Filter-list\0Label\0Filter-list\0...Label\0FIlter-list\0\0"`; filter for all files is canonically `"All Files\0*.*\0\0"` in the docs (specifically this due to handling of shortcut links); also provides a way for users to write in their own filters | Add or remove individual GtkFileFIlter objects; can select one specified in the list to show by default; default behavior is all files; if selected one when none has been specified, filter selection disabled; filter for all files specified in docs under gtk_file_filter_new() (except doesn't set a name) | NSArray of filter strings, or nil for All Files. There is no provision to have an "all files" option: you either specify a set of filters or you don't. (See filename extension auto-append below.). All filters are applied at once; there is no way to select. We might need to introduce an accessory panel (extra widget) to fake the filtering rules of other platforms...
Default file name | settable | settable | settable (as the filename label)
Initial directory | complex rules that have changed over time; we can pass an absolute filename (the previous filename or a default filename) and have its path used (if we specify just a path it will either be used as the filename or the program will crash); or we can give it a directory; or Windows will remember for us for some time, or... | pass previous filename or URI to show; overrides default file name; intended only for saving files (so I don't know if it's possible to remember current directory for opening??????); effect of passing containing directory undocumented(???? in my tests the given folder itself is selected) | has some rules; there is a way to specify a custom one; seems to have the undocumented effect that it selects the file if a file is named
Confirmation and cancel buttons | xxx | GTK_STOCK_OPEN, GTK_STOCK_SAVE, GTK_STOCK_SAVE_AS / GTK_STOCK_CANCEL | cancel button predefined; confirmation button can be changed (setPrompt:) but **TODO** the docs imply prompt is actually a global property?
Returned filename rules | xxxx | memory provided by GTK+ itself (so no need to worry about size limits); can return a single filename or URI or a GSList of filenames or URIs | xxx
Window title | optional; defaults to either Open or Save As | required(?) | optional for save (defaults to Save); unknown (**TODO**) for Open
Prompt to create new files | available; must be explicitly enabled; seems to only apply to Open File dialogs (**TODO**) | xxx | xxx
Adds file to a Recent Documents list | available; must be explicitly disabled | xxx | xxx
Allows nonexistent files to be created in Open dialogs | yes; can be switched off | xxx | xxx
"Open as read-only"/"Save as read-only" checkbox | provided; default; can be switched off | xxx | xxx
Navigating changes the current working directory of the program | yes; can be switched off for Save dialogs only (????) | xxx | xxx
Link following | For .lnk files, enabled by default iff a filter is specified; the All Files filter above is listed as being necessary to follow links; can be shut off with a flag in all cases | xxx | For Finder aliases, setResolvesAliases:
Help button | Available; old-style dialog boxes need a parent window (Explorer-style ones don't; they just need a hook function) | xxx | xxx
Extension auto-appending | Optional; three-character maximum; doesn't seem to be available on a per-filter basis | xxx | **NOT OPTIONAL.** The only way to avoid this is to not specify any filters. You can't even circumvent this with a delegate. If the user specifies another extension, they are asked to choose one if setAllowsOtherFileTypes: is set. (**TODO** could we use nameFieldStringValue to circumvent?)
Other labels | xxx | xxx | field before filename entry can be changed; also can provide an additional optional message
Multiple selection allows empty selection | xxx | xxx | xxx

TODO
* Windows: OFN_NOTESTFILECREATE might be necessary
* Windows: OFN_NOVALIDATE - see what happens without a hook
* Windows: OFN_SHAREAWARE - this is a weird one but it's network related
* Windows: templates seem to be how to provide extra parameters, but their usage isn't documented on the OPENFILENAME struct help page; check the rest of MSDN
* Mac OS X: turn on both setExtensionHidden: and setCanSelectHiddenExtension: to show the extnesion in the dialog
* Mac OS X: turn on setTreatsFilePackagesAsDirectories: since file packages (bundles) are an OS X-specific concept

## The Scrollbar Series
This actually turns out to be one of the very first things that Raymond ever blogged about, so if you just go to the last page of posts on The Old New Thing, it'll be there. But for my own convenience:
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/07/23/54576.aspx
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/07/25/54582.aspx
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/07/29/54591.aspx
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/07/30/54600.aspx
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/07/31/54601.aspx
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/08/05/54602.aspx
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/08/05/54610.aspx
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/08/07/54615.aspx
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/08/07/54617.aspx
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/08/11/54624.aspx
[23:59] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/08/11/54629.aspx
[00:00] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/08/13/54639.aspx
[00:00] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/08/15/54647.aspx
[00:00] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/08/18/54668.aspx
[00:00] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/09/09/54826.aspx
[00:00] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/09/13/54917.aspx
[00:00] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/09/13/54917.aspx
[00:00] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/09/15/54925.aspx
[00:00] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/09/17/54944.aspx
[00:00] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/09/17/54945.aspx
[00:00] <andlabs> http://blogs.msdn.com/b/oldnewthing/archive/2003/09/17/54946.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/10/16/55344.aspx
- not really part of, so to speak, but still http://blogs.msdn.com/b/oldnewthing/archive/2004/05/10/129068.aspx


TODO SORT
(TODO the original time formats article)
- http://blogs.msdn.com/b/oldnewthing/archive/2003/08/26/54708.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/09/05/54805.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/10/01/55108.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/10/13/55279.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/10/14/55286.aspx
	- TODO get article on spurious mouse moves
- http://blogs.msdn.com/b/oldnewthing/archive/2003/10/21/55384.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/10/24/55413.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/10/27/55461.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/10/29/55479.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/11/07/55619.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/12/29/46371.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2003/12/30/46594.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2004/01/21/61101.aspx (on maximum pathname length)
- http://blogs.msdn.com/b/oldnewthing/archive/2004/02/27/81155.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2004/03/29/101121.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2004/04/19/115912.aspx
- http://blogs.msdn.com/b/oldnewthing/archive/2004/07/01/170856.aspx

interesting things
- low-tech models http://blogs.msdn.com/b/oldnewthing/archive/2003/10/07/55202.aspx
- JIT http://blogs.msdn.com/b/oldnewthing/archive/2003/11/04/55560.aspx
- comments on frameworks (!) http://blogs.msdn.com/b/oldnewthing/archive/2003/11/10/55635.aspx
- dialog box entry autoselection http://blogs.msdn.com/b/oldnewthing/archive/2003/11/14/55678.aspx
- more on dialog codes http://blogs.msdn.com/b/oldnewthing/archive/2003/11/26/55872.aspx
- custom window class styles http://blogs.msdn.com/b/oldnewthing/archive/2003/12/03/55927.aspx
- plug and play device changes http://blogs.msdn.com/b/oldnewthing/archive/2003/12/05/55936.aspx
- more JIT-related http://blogs.msdn.com/b/oldnewthing/archive/2003/12/08/55954.aspx
- suspending threads http://blogs.msdn.com/b/oldnewthing/archive/2003/12/09/55988.aspx
- UI design http://blogs.msdn.com/b/oldnewthing/archive/2003/12/09/55995.aspx
- passing data to child processes (is this useful for detecting multiple instances?) http://blogs.msdn.com/b/oldnewthing/archive/2003/12/11/56043.aspx
- on the need for cbSize fields http://blogs.msdn.com/b/oldnewthing/archive/2003/12/12/56061.aspx
- notes on critical sections http://blogs.msdn.com/b/oldnewthing/archive/2003/12/18/44425.aspx http://blogs.msdn.com/b/oldnewthing/archive/2003/12/22/45152.aspx
- on COM interfaces and control panels (also uses activation contexts for SxS manifests) http://blogs.msdn.com/b/oldnewthing/archive/2003/12/26/45979.aspx
- return address attacks http://blogs.msdn.com/b/oldnewthing/archive/2004/01/01/47042.aspx
- on GetProcAddress() http://blogs.msdn.com/b/oldnewthing/archive/2004/01/12/57833.aspx
- DllMain() http://blogs.msdn.com/b/oldnewthing/archive/2004/01/27/63401.aspx http://blogs.msdn.com/b/oldnewthing/archive/2004/01/28/63880.aspx
- string resources http://blogs.msdn.com/b/oldnewthing/archive/2004/01/30/65013.aspx
- text/unicode macros http://blogs.msdn.com/b/oldnewthing/archive/2004/02/12/71851.aspx
- version numbers http://blogs.msdn.com/b/oldnewthing/archive/2004/02/13/72476.aspx
- RECT endpoint-exclusivity (everything does this, not just Windows, so it's useful to point people to) http://blogs.msdn.com/b/oldnewthing/archive/2004/02/18/75652.aspx
- on fixing a problem everywhere (see RECT above) http://blogs.msdn.com/b/oldnewthing/archive/2004/02/19/76368.aspx
- proper use of COM http://blogs.msdn.com/b/oldnewthing/archive/2004/02/20/77120.aspx
- invalid thread and process IDs http://blogs.msdn.com/b/oldnewthing/archive/2004/02/23/78395.aspx
- on GetDesktopWindow() http://blogs.msdn.com/b/oldnewthing/archive/2004/02/24/79212.aspx
- links in notification bubbles http://blogs.msdn.com/b/oldnewthing/archive/2004/02/25/79842.aspx
- dialog boxes are initially hidden http://blogs.msdn.com/b/oldnewthing/archive/2004/03/11/87941.aspx
- Win16 resource format changes (probably irrelevant since NE vs. PE) http://blogs.msdn.com/b/oldnewthing/archive/2004/03/19/92648.aspx
- COM is fickle because people are fickle http://blogs.msdn.com/b/oldnewthing/archive/2004/03/26/96777.aspx
- on reference counting http://blogs.msdn.com/b/oldnewthing/archive/2004/04/06/108395.aspx
- taskbar program names http://blogs.msdn.com/b/oldnewthing/archive/2004/04/08/109775.aspx
- random number generator seeds http://blogs.msdn.com/b/oldnewthing/archive/2004/04/12/111596.aspx
- Unicode collation http://blogs.msdn.com/b/oldnewthing/archive/2004/04/13/112234.aspx
- on Unicode digits http://blogs.msdn.com/b/oldnewthing/archive/2004/04/16/114651.aspx
- accessibility http://blogs.msdn.com/b/oldnewthing/archive/2004/04/23/118893.aspx
- HWND_BROADCAST is not for multiple instances http://blogs.msdn.com/b/oldnewthing/archive/2004/05/05/126427.aspx
- COM threading http://blogs.msdn.com/b/oldnewthing/archive/2004/05/11/129760.aspx
- disabling vs hiding http://blogs.msdn.com/b/oldnewthing/archive/2004/05/14/131881.aspx
- when to put ... after a command name http://blogs.msdn.com/b/oldnewthing/archive/2004/05/17/133181.aspx
- COM destructors http://blogs.msdn.com/b/oldnewthing/archive/2004/05/20/135841.aspx http://blogs.msdn.com/b/oldnewthing/archive/2004/05/21/136701.aspx
- high performance multithreading http://blogs.msdn.com/b/oldnewthing/archive/2004/05/28/143769.aspx
- getting icons for known file types http://blogs.msdn.com/b/oldnewthing/archive/2004/06/01/145428.aspx
- COM threading models http://blogs.msdn.com/b/oldnewthing/archive/2004/06/02/146671.aspx
	- why am I posting so much on COM? in case I need to use shell functions in the future that aren't flat functions...
- thread messages and message reentrancy (and some DDE for good measure) http://blogs.msdn.com/b/oldnewthing/archive/2004/06/08/150929.aspx
- windows language IDs http://blogs.msdn.com/b/oldnewthing/archive/2004/06/09/151689.aspx
- combobox/listbox speed improvements http://blogs.msdn.com/b/oldnewthing/archive/2004/06/10/152612.aspx
- dialog templates
	- http://blogs.msdn.com/b/oldnewthing/archive/2004/06/17/158175.aspx
	- http://blogs.msdn.com/b/oldnewthing/archive/2004/06/18/159248.aspx
	- http://blogs.msdn.com/b/oldnewthing/archive/2004/06/21/161375.aspx
	- http://blogs.msdn.com/b/oldnewthing/archive/2004/06/22/162360.aspx
	- http://blogs.msdn.com/b/oldnewthing/archive/2004/06/23/163596.aspx
	- http://blogs.msdn.com/b/oldnewthing/archive/2004/06/24/164737.aspx
- COM and DLL unloading http://blogs.msdn.com/b/oldnewthing/archive/2004/06/28/167800.aspx
- thread safety vs reentrancy http://blogs.msdn.com/b/oldnewthing/archive/2004/06/29/168719.aspx
- shell and COM allocation http://blogs.msdn.com/b/oldnewthing/archive/2004/07/05/173226.aspx
- minimized window coordinates (should I move this one to the higher list?) http://blogs.msdn.com/b/oldnewthing/archive/2004/07/07/175285.aspx

articles on time
- http://blogs.msdn.com/b/oldnewthing/archive/2004/02/26/80492.aspx

up to
http://blogs.msdn.com/b/oldnewthing/default.aspx?PageIndex=374
last one: Why does icon text get a solid background if drop shadows are disabled?

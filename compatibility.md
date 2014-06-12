# Useful things in newer versions

## Windows
### Windows Vista

### Windows 7

### Windows 8

### Windows 8.1

## GTK+
TODO what ships with Ubuntu Quantal (12.10)?

### GTK+ 3.6
ships with: Ubuntu Raring (13.04)

- GtkEntry and GtkTextView have input purposes and input hints for external input methods but do not change input themselves
	- according to Company, we connect to insert-text for that
- GtkLevelBar
- GtkMenuButton
- **GtkSearchEntry**

### GTK+ 3.8
ships with: Ubuntu Saucy (13.10)

Not many interesting new things to us here, unless you count widget-internal tickers and single-click instead of double-click to select list items (a la KDE)... and oh yeah, also widget opacity.

### GTK+ 3.10
ships with: Ubuntu Trusty (14.04 LTS)

- tab character stops in GtkEntry
- GtkHeaderBar
	- intended for titlebar overrides; GtkInfoBar is what I keep thinking GtkHeaderBar is
- **GtkListBox**
- GtkRevealer for smooth animations of disclosure triangles
- GtkSearchBar for custom search popups
- **GtkStack and GtkStackSwitcher**
- titlebar overrides (seems to be the hot new thing)

### GTK+ 3.12
not yet in Ubuntu Utopic (14.10)

- GtkActionBar (basically like the bottom-of-the-window toolbars in Mac programs)
- gtk_get_locale_direction(), for internationalization
- more control over GtkHeaderBar
- **GtkPopover**
	- GtkPopovers on GtkMenuButtons
- GtkStack signaling
- **gtk_tree_path_new_from_indicesv()** (for when we add Table if we have trees too)

## Cocoa
### Mac OS X 10.7+

- Foundation ([full details](https://developer.apple.com/library/mac/releasenotes/Foundation/RN-FoundationOlderNotes/#//apple_ref/doc/uid/TP40008080-TRANSLATED_CHAPTER_965-TRANSLATED_DEST_0))
	- can use [ClassName class] == nil to test for feature existence at class level
	- NSUserDefaults is now thread-safe and respects user IDs
	- NSProcessInfo is now thread-safe
	- NSProcessInfo has automatic termination; might want to watch for this
	- NSLinguisticTrigger
	- NSCalendar changes (important for date/time pickers)
- AppKit ([full details](https://developer.apple.com/library/mac/releasenotes/AppKit/RN-AppKitOlderNotes/#X10_7Notes)
	- Cocoa decides when to draw focus rings
	- lots of resolution independence stuff
	- window edge resize/live resize
	- restorable state
	- toolbar versioning
	- better scroll wheel management
	- multi-image dragging (TODO didn't GTK+ add this too)
	- built-in full screen (TODO what was fullscreen like before?)
	- NSTextFinder for a standard search for text bar
	- NSPasteboard: no longer need to invoke -[NSPasteboard types]
	- NSPasteboards use UTIs for type information
	- NOTE FROM LEOPARD: full-width cells, source lists
	- view-based NSTableViews and NSOutlineViews
	- URLs returned may not be file:// anymore! TODO note in futureplans.md
	- **NSPopover**
	- vertical text support improved
	- **placeholder text in NSTextField now shown when focused**
	- NSTextView now has an optional formatting bar

### Mac OS X 10.8+

### Mac OS X 10.9+

# ui: platform-native GUI library for Go

## Feature requests wanted! (Really; IDK what to add next!)

This is a library that aims to provide simple GUI software development in Go. It runs on/requires:

- Windows: cgo, MinGW-w64 (see note below), Windows XP and newer
	- **Note**: Notice I specifically said [MinGW-w64](http://mingw-w64.sourceforge.net/) here. This is important: regular MinGW is missing various recent header files which package ui uses, and thus won't build with it. Make sure your MinGW is that version instead. If you're running on Windows and are not sure what to download, get the mingw-builds distribution.
- Mac OS X: cgo, Mac OS X 10.7 and newer
- other Unixes: cgo, GTK+ 3.4 and newer

Go 1.3 is required. Note that vanilla 1.3 has a bug in Mac OS X cgo; the next release will fix it.

(this README needs some work)

Be sure to have at least each outermost Window escaping to the heap until a good resolution to Go issue 8310 comes out.

prevlib.tar contains the previous version of the library as it stood when I restarted; don't bother using it.

# Documentation

The in-code documentation needs improvement. I have written a [tutorial](https://github.com/andlabs/ui/wiki/Getting-Started) in the Wiki.

# Updates

**4 November 2014**<br>Added two new controls, Spinbox (which allows numeric entry with up and down buttons) and ProgressBar (which measures progress). Both aren't fully fleshed out, but are good enough for general use now.

**28 October 2014**<br>Mac OS X resizing issues should be (mostly?) fixed now. Textbox still doesn't work right...

**24 October 2014**<br>Textbox, a multi-line version of TextField, has been added. (Note that it may not work properly on Mac OS X; this is being investigated.) In addition, excess space around controls on Mac OS X should be settled now.

**18 October 2014**<br>The container system was rewritten entirely. You can now set a margin on Windows and Groups and spacing between controls ("padding") on Stacks, Grids, and SimpleGrids. Margins on Tabs will come soon. The work needed to change this will make future additions (like Popover and Spinbox) easier/more sensible. (The Mac OS X code is still glitchy; mind the dust.)

As part of the change, standalone Labels have been removed. All Labels now behave like standalone labels. A new layout container, Form, will be introduced in the near future to allow proper layout of widgets with labels.

**3 September 2014**<br>The new GtkGrid-style Grid is now implemented! See its documentation for more details. Also, debugging spew has been removed.

**31 August 2014**<br>Grid is now renamed SimpleGrid in preparation for implementing a more [GtkGrid](https://developer.gnome.org/gtk3/unstable/GtkGrid.html)-like Grid. Mind the change.

# Screenshots
The example widget gallery on GTK+ in the Adwaita theme (3.13/master):

![widget gallery example](https://raw.githubusercontent.com/andlabs/ui/master/examples/widgetgallery/widgetgallery.png)

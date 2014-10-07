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

**3 September 2014**<br>The new GtkGrid-style Grid is now implemented! See its documentation for more details. Also, debugging spew has been removed.

**31 August 2014**<br>Grid is now renamed SimpleGrid in preparation for implementing a more [GtkGrid](https://developer.gnome.org/gtk3/unstable/GtkGrid.html)-like Grid. Mind the change.

# Screenshots
The example widget gallery on GTK+ in the Adwaita theme (3.13/master):

![widget gallery example](https://raw.githubusercontent.com/andlabs/ui/master/examples/widgetgallery/widgetgallery.png)

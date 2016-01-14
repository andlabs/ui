# ui: platform-native GUI library for Go

# This package is still very much WIP.

As of December 2015 the previous package ui API that has been around since this repo was started is no longer being supported. It is being replaced with a much more stable API built around my libui; see below.

If you still want to use the old package ui, you can get the package under the `pre-libui` tag. Keep in mind that it's not stable, buggy, and **no longer supported**. If you do continue, make sure that instances of `ui.Window` escape to the heap to avoid some of the issues.

If you want to play around with this new package ui, you'll need to install libui manually. Clone that repo and `make` (with GNU make) libui, then:

- On Windows, merely copy out\libui.dll to the root of this repo.
	- Go 1.5 is adequate.
- On OS X, copy out/libui.A.dylib to the root of this repo as libui.A.dylib and symlink it to libui.dylib
	- You must also be running Go 1.6 Beta 2 or newer due to more Go bugs.
- On other Unixes, copy out/libui.so.0 to the root of this repo as libui.so.0 and symlink it to libui.so
	- Go 1.5 is adequate.

and then copy ui.h to the top of this repo as well. (You may symlink any files instead of copying if so choose.)

Stable releases of package ui will have all these files built in; these steps are only necessary for master builds.

# New README

This is a library that aims to provide simple GUI software development in Go.

It is based on my [libui](https://github.com/andlabs/libui), a simple cross-platform library that does the same thing, but written in C. **You must include this library in your binary distributions.**

It runs on/requires:

- Windows: cgo, Windows Vista and newer
- Mac OS X: cgo, Mac OS X 10.7 and newer
- other Unixes: cgo, GTK+ 3.4 and newer

It also requires Go 1.6 or newer (due to various bugs in cgo that were fixed only starting with 1.6).

(this README needs some work)

# Installation

# Documentation

The in-code documentation needs improvement. I have written a [tutorial](https://github.com/andlabs/ui/wiki/Getting-Started) in the Wiki.

# Updates

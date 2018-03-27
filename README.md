# ui: platform-native GUI library for Go

# Update 17 February 2018
I fixed the Enter+Escape crashing bug on Windows, and applied the resultant Alpha 3.5 binary release to package ui. However, build issues prevented a linux/386 binary from being made, so the API updates won't come yet. The next Alpha release, which will use semver and thus be called v0.4.0, should hopefully have no such issues. Sorry!

# Update 5 June 2016: You can FINALLY `go get` this package!

`go get` should work out of the box for the following configurations:

* darwin/amd64
* linux/386
* linux/amd64
* windows/386
* windows/amd64

Everything is now fully static — no DLLs or shared objects anymore!

Note that these might not fully work right now, as the libui Alpha 3.1 API isn't fully implemented yet, and there might be residual binding problems. Hopefully none which require an Alpha 3.2...

# New README

This is a library that aims to provide simple GUI software development in Go.

It is based on my [libui](https://github.com/andlabs/libui), a simple cross-platform library that does the same thing, but written in C. **You must include this library in your binary distributions.**

It runs on/requires:

- Windows: cgo, Windows Vista SP2 with Platform Update and newer
- Mac OS X: cgo, Mac OS X 10.8 and newer
- other Unixes: cgo, GTK+ 3.10 and newer
	- Debian, Ubuntu, etc.: `sudo apt-get install libgtk-3-dev`
	- Red Hat/Fedora, etc.: `sudo dnf install gtk3-devel`
	- TODO point out this is fine for most people but refer to distro docs if more control is needed, including cross-compilation instructions
	- TODO clean this part up and put it in the appropriate place (maybe libui itself)

It also requires Go 1.6 or newer (due to various bugs in cgo that were fixed only starting with 1.6).

(this README needs some work)

# Installation

# Documentation

The in-code documentation needs improvement. I have written a [tutorial](https://github.com/andlabs/ui/wiki/Getting-Started) in the Wiki.

# Updates

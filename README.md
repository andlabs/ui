Woah, lots of attention! Thanks!

## Updates

- **11 June 2014**<br>**I have decided to remove Mac OS X 10.6 support** because it's only causing problems for building (and everyone else says I should anyway, including Mac developers!). This does break my original goal, but I'm going to have to break things sooner or later. Please let me know if any of you actually use this package on 10.6. (I personally don't like it when programs require 10.7 (or iOS 7, for that matter), but what are you gonna do?)

## README

This is a placeholder README; the previous file (olddocs/oldREADME.md) was rather long and confusing. I'll be rewriting it properly soon.

Until then, here's the important things you need to know:
- **this package is very much incomplete; see `stable.md` for a list of what is guaranteed to not change at the API level — for everything newer, you have been warned!**
- this package requires Go 1.3, which is presently available as a RC build (source builds from go tip will work too)
	- I don't think the Windows side uses any Go 1.3 features, but just to be safe I'm going to say express caution
	- Unix builds need 1.3 to fix some type-checker bugs in cgo
	- Mac OS X builds need 1.3 because Go 1.3 adds Objective-C support to cgo
- the Windows build does not need cgo unless you want to regenerate the `zconstants_windows_*.go` files; the other targets **do**
- my plan is to target all versions of OSs that Go itself supports. I will, however, make concessions where appropriate. This means:
	- Windows: Windows XP or newer
	- Unix: this is trickier; I decided to settle on GTK+ 3.4 or newer as Ubuntu 12.04 LTS ships with it
	- Mac OS X: Mac OS X 10.7 or newer (Go supports 10.6 but this is a pain to compile Cocoa programs for due to flaws in the later header files)
- for the Windows build, you won't need to provide a comctl32.dll version 6 manifest, as the package produces its own
	- comctl32.dll version 6 *is* required for proper functioning!

[andlabs/wakeup](https://github.com/andlabs/wakeup) is a repository that provides a sample application.

If you are feeling adventurous, running `./test.sh` (which accepts `go build` options) from within the package directory will build a test program which I use to make sure everything works. (I'm not sure how to do automated tests for a package like this, so `go test` will say no tests found for now; sorry.) If you are cross-compiling to Windows, you will need to have a very specific Go setup which allows multiple cross-compilation setups in a single installation; this requires [a CL which won't be in Go 1.3 but may appear in Go 1.4 if accepted](https://codereview.appspot.com/93580043) and both windows/386 and windows/amd64 set up for cgo. (This is because `./test.sh` on Windows targets invariably regenerates the `zconstants_windows_*.go` files; there is no option to turn it off lest I become complacent and use it myself.)

Finally, please send documentation suggestions! I'm taking the documentation of this package very seriously because I don't want to make **anything** ambiguous. (Trust me, ambiguity in API documentation was a pain when writing this...)

Thanks!

(Note: I temporarily disabled Travis.ci; if I can figure out how to do good cross-compiles with it, then I can put it back.)

## Screenshots

You asked for them; here they are.

Image | Description
----- | -----
<img src="http://andlabs.lostsig.com/screenshots/20140608/uiwin7.png" width="400px"> | The test program on Windows 7
<img src="http://andlabs.lostsig.com/screenshots/20140608/uimac.png" width="400px"> | The test program on Mac OS X 10.8
<img src="http://andlabs.lostsig.com/screenshots/20140608/uikde.png" width="400px"> | The test program on Ubuntu 14.04 with KDE and the oxygen-gtk theme

## Big Note

Well, it's come to this. Race conditions in th eevent handler and GTK+ not supposed to be use like this make me have no choice but to make the API completely single-threaded.

Most of the documentation is updated properly. The package documentation and dialog box documentation is not; the latter is something I'm not sure how to solve within the designs of all three supported environments (which are all wildly different with regards to dialog boxes and code modality). The test program and wakeup have been updated accordingly. **The Window MsgBox() methods will not work properly on Mac OS X until further notice.**

I'm not sure if I'm going to continue developing this package or write a whole new API from scratch; like skelterjohn's go.uik, this would have manual widget drawing and a single look and feel across all platforms, because native single-threaded APIs cannot be made natively multithreaded. (The best hope is Windows, where each thread can have GUI state, followed by Mac OS X, where you can decide to push things onto the main thread and wait for them to run.)

The biggest problem with a concurrent event loop, which I don't know how to solve, is something I call atomicity of event handlers. Let's say you have an event handler;

```go
func button1Clicked() {
	button2.Disable()
	timer.Stop()
}
```

As things stand, it is entirely possible for the scheduler to interrupt the event handler at any point. If it's interrupted before calling `button2.Disable()`, the user may be able to click button 2 before it gets disabled, leading to an event you don't want. And if it's interrupted before the call to `timer.Stop()`, the timer might have a chance to fire, giving you a spurious timer interrupt.

I currently have no idea how to avoid this problem, but it's the only problem I can think of when designing a concurrent user interface.

It is now July. I've spent five months developing and tuning eveyrthing and its sudden popularity last month caught me well off-guard. I thought I had everything ready for the music editing software that I wanted to make, but these problems have resulted in me creating an API that no longer satisfies its goal of working the Go way.

If anyone knows a possible solution to the above event problem, ***PLEASE*** let me know.

I'm not sure where I'm going to take things from here. Until then, thanks for all the support. I appreciate it. The project was fun, and I learned a lot of things I've always wanted to learn.

I know one thing, though: I don't want to abandon desktop programming in Go.

- andlabs

## Old README

**Note to ALL users: [please read and comment; the design of the package is fatally flawed but I want to know what people think of the fix](http://andlabs.lostsig.com/blog/2014/06/27/61/my-ui-package-lessons-learned-about-threading-and-plans-for-the-future).**

**Note to Mac users: there is [a bug in Go 1.3 stable](https://code.google.com/p/go/issues/detail?id=8238) that causes cgo to crash trying to build this package. Please follow the linked bug report for detials.**

Woah, lots of attention! Thanks!

## Old Updates

- **26 June 2014**
	- Controls in Windows can now be spaced apart more naturally. Call `w.SetSpaced(true)` to opt in. **Whether this will remain opt-in or whether the name will change is still unknown at this point.**
	- There's a new function `Layout()` which provides high-level layout creation. The function was written by [boppreh](https://github.com/boppreh) and details can be found [here](https://github.com/andlabs/ui/pull/19). **Whether this function will stay in the main package or be moved to a subpackage is still unknown.**
	- There is now `Checkbox.SetChecked()` to set the check state of a Checkbox programmatically.

- **25 June 2014**<br>Labels by default now align themselves relative to the control they are next to. There is a new function `NewStandaloneLabel()` which returns a label whose text is aligned to the top-left corner of the alloted space regardless.

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

// 13 december 2015

package ui

// #cgo CFLAGS: -mmacosx-version-min=10.8 -DMACOSX_DEPLOYMENT_TARGET=10.8
// #cgo LDFLAGS: ${SRCDIR}/libui_darwin_amd64.a -framework Foundation -framework AppKit -lpthread -mmacosx-version-min=10.8
// #include <CoreFoundation/CoreFoundation.h>
// #include <pthread.h>
// extern void _CFRunLoopSetCurrent(CFRunLoopRef);
// extern pthread_t _CFMainPThread;
import "C"

// OS X cares very deeply if we don't run on the very first thread the OS creates
// why? who knows. it's stupid and completely indefensible. let's use undocumented APIs to get around it.
// apple uses them too: http://www.opensource.apple.com/source/kext_tools/kext_tools-19.2/kextd_main.c?txt
// apple HAS SUGGESTED them too: http://lists.apple.com/archives/darwin-development/2002/Sep/msg00250.html
// gstreamer uses them too: http://cgit.freedesktop.org/gstreamer/gst-plugins-good/tree/sys/osxvideo/osxvideosink.m
func ensureMainThread() {
	// TODO set to nil like the apple code?
	C._CFRunLoopSetCurrent(C.CFRunLoopGetMain())
	// TODO is this part necessary?
	C._CFMainPThread = C.pthread_self()
}

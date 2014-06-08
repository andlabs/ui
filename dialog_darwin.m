// 15 may 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#import <AppKit/NSAlert.h>

// see delegateuitask_darwin.m
// in this case, NSWindow.h includes NSApplication.h

#ifdef MAC_OS_X_VERSION_10_7
#undef MAC_OS_X_VERSION_MIN_REQUIRED
#undef MAC_OS_X_VERSION_MAX_ALLOWED
#define MAC_OS_X_VERSION_MIN_REQUIRED MAC_OS_X_VERSION_10_7
#define MAC_OS_X_VERSION_MAX_ALLOWED MAC_OS_X_VERSION_10_7
#endif
#import <AppKit/NSApplication.h>
#undef MAC_OS_X_VERSION_MIN_REQUIRED
#undef MAC_OS_X_VERSION_MAX_ALLOWED
#define MAC_OS_X_VERSION_MIN_REQUIRED MAC_OS_X_VERSION_10_6
#define MAC_OS_X_VERSION_MAX_ALLOWED MAC_OS_X_VERSION_10_6

#import <AppKit/NSWindow.h>

#define to(T, x) ((T *) (x))
#define toNSWindow(x) to(NSWindow, (x))

static void alert(id parent, NSString *primary, NSString *secondary, NSAlertStyle style, void *chan)
{
	NSAlert *box;

	box = [NSAlert new];
	[box setMessageText:primary];
	if (secondary != nil)
		[box setInformativeText:secondary];
	[box setAlertStyle:style];
	// TODO is there a named constant? will also need to be changed when we add different dialog types
	[box addButtonWithTitle:@"OK"];
	if (parent == nil)
		dialog_send(chan, (intptr_t) [box runModal]);
	else
		[box beginSheetModalForWindow:toNSWindow(parent)
			modalDelegate:[NSApp delegate]
			didEndSelector:@selector(alertDidEnd:returnCode:contextInfo:)
			contextInfo:chan];
}

void msgBox(id parent, id primary, id secondary, void *chan)
{
	alert(parent, (NSString *) primary, (NSString *) secondary, NSInformationalAlertStyle, chan);
}

void msgBoxError(id parent, id primary, id secondary, void *chan)
{
	alert(parent, (NSString *) primary, (NSString *) secondary, NSCriticalAlertStyle, chan);
}

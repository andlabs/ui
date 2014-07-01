// 15 may 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#import <AppKit/NSAlert.h>
#import <AppKit/NSWindow.h>
#import <AppKit/NSApplication.h>

#define to(T, x) ((T *) (x))
#define toNSWindow(x) to(NSWindow, (x))

static intptr_t alert(id parent, NSString *primary, NSString *secondary, NSAlertStyle style)
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
		return (intptr_t) [box runModal];
	else {
		NSInteger ret;

		[box beginSheetModalForWindow:toNSWindow(parent)
			modalDelegate:[NSApp delegate]
			didEndSelector:@selector(alertDidEnd:returnCode:contextInfo:)
			contextInfo:&ret];
		// TODO
		return (intptr_t) ret;
	}
}

void msgBox(id parent, id primary, id secondary)
{
	alert(parent, (NSString *) primary, (NSString *) secondary, NSInformationalAlertStyle);
}

void msgBoxError(id parent, id primary, id secondary)
{
	alert(parent, (NSString *) primary, (NSString *) secondary, NSCriticalAlertStyle);
}

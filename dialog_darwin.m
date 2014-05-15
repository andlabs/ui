// 15 may 2014

#include "objc_darwin.h"
#include "dialog_darwin.h"
#include <AppKit/NSAlert.h>

static void alert(NSString *primary, NSString *secondary, NSAlertStyle style)
{
	NSAlert *box;

	box = [NSAlert new];
	[box setMessageText:primary];
	if (secondary != nil)
		[box setInformativeText:secondary];
	[box setAlertStyle:style];
	// TODO is there a named constant? will also need to be changed when we add different dialog types
	[box addButtonWithTitle:@"OK"];
	[box runModal];
}

void msgBox(id primary, id secondary)
{
	alert((NSString *) primary, (NSString *) secondary, NSInformationalAlertStyle);
}

void msgBoxError(id primary, id secondary)
{
	alert((NSString *) primary, (NSString *) secondary, NSCriticalAlertStyle);
}

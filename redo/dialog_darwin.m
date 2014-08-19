// 19 august 2014

#import "objc_darwin.h"
#import <Cocoa/Cocoa.h>

char *openFile(void)
{
	NSOpenPanel *op;
	NSInteger ret;

	op = [NSOpenPanel openPanel];
	[op setCanChooseFiles:YES];
	[op setCanChooseDirectories:NO];
	[op setResolvesAliases:NO];
	[op setAllowsMultipleSelection:NO];
	[op setShowsHiddenFiles:YES];
	[op setCanSelectHiddenExtension:NO];
	[op setExtensionHidden:NO];
	[op setAllowsOtherFileTypes:YES];
	[op setTreatsFilePackagesAsDirectories:YES];
	// disable custom events
	// TODO doesn't seem to work
	dispatch_suspend(dispatch_get_main_queue());
	ret = [op runModal];
	dispatch_resume(dispatch_get_main_queue());
	if (ret != NSFileHandlingPanelOKButton)
		return NULL;
	// string freed on the Go side
	return strdup([[[op URL] path] UTF8String]);
}

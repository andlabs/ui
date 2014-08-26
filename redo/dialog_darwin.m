// 19 august 2014

#import "objc_darwin.h"
#import <Cocoa/Cocoa.h>

#define toNSWindow(x) ((NSWindow *) (x))

void openFile(id parent, void *data)
{
	NSOpenPanel *op;

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
	[op beginSheetModalForWindow:toNSWindow(parent) completionHandler:^(NSInteger ret){
		if (ret != NSFileHandlingPanelOKButton) {
			finishOpenFile(NULL, data);
			return;
		}
		// string freed on the Go side
		finishOpenFile(strdup([[[op URL] path] UTF8String]), data);
	}];
}

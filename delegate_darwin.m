// 13 may 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#include <Foundation/NSObject.h>
#include <Foundation/NSValue.h>
#include <Foundation/NSNotification.h>
#include <AppKit/NSApplication.h>
#include <AppKit/NSWindow.h>

@interface appDelegate : NSObject
@end

@implementation appDelegate

- (void)uitask:(NSValue *)fp
{
	appDelegate_uitask([fp pointerValue]);
}

- (BOOL)windowShouldClose:(id)win
{
	appDelegate_windowShouldClose(win);
	return NO;		// don't close
}

- (void)windowDidResize:(NSNotification *)n
{
	appDelegate_windowDidResize([n object]);
}

- (void)buttonClicked:(id)button
{
	appDelegate_buttonClicked(button);
}

- (NSApplicationTerminateReply)applicationShouldTerminate:(NSApplication *)app
{
	appDelegate_applicationShouldTerminate();
	return NSTerminateCancel;
}

@end

id makeAppDelegate(void)
{
	return [appDelegate new];
}

id windowGetContentView(id window)
{
	return [((NSWindow *) window) contentView];
}

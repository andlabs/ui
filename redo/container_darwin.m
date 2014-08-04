// 4 august 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#include <Cocoa/Cocoa.h>

// calling -[className] on the content views of NSWindow, NSTabItem, and NSBox all return NSView, so I'm assuming I just need to override these
// fortunately, in the case of NSTabView, this -[setFrame:] is called when resizing and when changing tabs, so we can indeed use this directly there
@interface goContainerView : NSView {
@public
	void *gocontainer;
}
@end

@implementation goContainerView

- (void)setFrame:(NSRect)r
{
	[super setFrame:r];
	if (self->gocontainer != NULL)
		containerResized(self->gocontainer, (intptr_t) r.size.width, (intptr_t) r.size.height);
}

@end

id newContainerView(void *gocontainer)
{
	goContainerView *c;

	c = [[goContainerView alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	c->gocontainer = gocontainer;
	return (id) c;
}

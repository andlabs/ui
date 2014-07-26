// 25 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSTabView(x) ((NSTabView *) (x))

@interface goTabView : NSTabView {
@public
	void *gotab;
}
@end

@implementation goTabView

- (void)setFrame:(NSRect)r
{
	NSRect content;

	[super setFrame:r];
	content = [self contentRect];
	tabResized(self->gotab, (intptr_t) content.size.width, (intptr_t) content.size.height);
}

@end

id newTab(void *gotab)
{
	goTabView *t;

	t = [[goTabView alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	setStandardControlFont((id) t);		// safe; same selector provided by NSTabView
	t->gotab = gotab;
	return (id) t;
}

id tabAppend(id t, char *name)
{
	NSTabViewItem *i;

	i = [[NSTabViewItem alloc] initWithIdentifier:nil];
	[i setLabel:[NSString stringWithUTF8String:name]];
	[toNSTabView(t) addTabViewItem:i];
	return (id) [i view];
}

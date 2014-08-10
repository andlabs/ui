// 25 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSTabView(x) ((NSTabView *) (x))
#define toNSView(x) ((NSView *) (x))

id newTab(void)
{
	NSTabView *t;

	t = [[NSTabView alloc] initWithFrame:NSZeroRect];
	setStandardControlFont((id) t);		// safe; same selector provided by NSTabView
	return (id) t;
}

void tabAppend(id t, char *name, id view)
{
	NSTabViewItem *i;

	i = [[NSTabViewItem alloc] initWithIdentifier:nil];
	[i setLabel:[NSString stringWithUTF8String:name]];
	[i setView:toNSView(view)];
	[toNSTabView(t) addTabViewItem:i];
}

struct xsize tabPreferredSize(id control)
{
	NSTabView *tv;
	NSSize s;
	struct xsize t;

	tv = toNSTabView(control);
	s = [tv minimumSize];
	t.width = (intptr_t) s.width;
	t.height = (intptr_t) s.height;
	return t;
}

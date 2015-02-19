// 29 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSTableView(x) ((NSTableView *) (x))

// NSTableColumn provides no provision to store an integer data
// it does provide an identifier tag, but that's a NSString, and I'd rather not risk the conversion overhead
@interface goTableColumn : NSTableColumn {
@public
	intptr_t gocolnum;
}
@end

@implementation goTableColumn
@end

@interface goTableDataSource : NSObject <NSTableViewDataSource, NSTableViewDelegate> {
@public
	void *gotable;
}
@end

@implementation goTableDataSource

- (NSInteger)numberOfRowsInTableView:(NSTableView *)view
{
	return (NSInteger) goTableDataSource_getRowCount(self->gotable);
}

- (id)tableView:(NSTableView *)view objectValueForTableColumn:(NSTableColumn *)col row:(NSInteger)row
{
	void *ret;
	NSString *s;
	intptr_t colnum;
	char *str;
	int type = colTypeText;

	colnum = ((goTableColumn *) col)->gocolnum;
	ret = goTableDataSource_getValue(self->gotable, (intptr_t) row, colnum, &type);
	switch (type) {
	case colTypeImage:
		// TODO free the returned image when done somehow
		return (id) ret;
	case colTypeCheckbox:
		if (ret == NULL)
			return nil;
		return (id) [NSNumber numberWithInt:1];
	}
	str = (char *) ret;
	s = [NSString stringWithUTF8String:str];
	free(str);		// allocated with C.CString() on the Go side
	return (id) s;
}

- (void)tableView:(NSTableView *)view setObjectValue:(id)value forTableColumn:(NSTableColumn *)col row:(NSInteger)row
{
	intptr_t colnum;
	NSNumber *number = (NSNumber *) value;	// thanks to mikeash in irc.freenode.net/#macdev

	colnum = ((goTableColumn *) col)->gocolnum;
	goTableDataSource_toggled(self->gotable, (intptr_t) row, colnum, [number boolValue]);
}

- (void)tableViewSelectionDidChange:(NSNotification *)note
{
	tableSelectionChanged(self->gotable);
}

@end

id newTable(void)
{
	NSTableView *t;

	t = [[NSTableView alloc] initWithFrame:NSZeroRect];
	[t setAllowsColumnReordering:NO];
	[t setAllowsColumnResizing:YES];
	[t setAllowsMultipleSelection:NO];
	[t setAllowsEmptySelection:YES];
	[t setAllowsColumnSelection:NO];
	return (id) t;
}

void tableAppendColumn(id t, intptr_t colnum, char *name, int type, BOOL editable)
{
	goTableColumn *c;
	NSImageCell *ic;
	NSButtonCell *bc;
	NSLineBreakMode lbm = NSLineBreakByTruncatingTail;		// default for most types

	c = [[goTableColumn alloc] initWithIdentifier:nil];
	c->gocolnum = colnum;
	switch (type) {
	case colTypeImage:
		ic = [[NSImageCell alloc] initImageCell:nil];
		// this is the behavior we want, which differs from the Interface Builder default of proportionally down
		[ic setImageScaling:NSImageScaleProportionallyUpOrDown];
		// these two, however, ARE Interface Builder defaults
		[ic setImageFrameStyle:NSImageFrameNone];
		[ic setImageAlignment:NSImageAlignCenter];
		[c setDataCell:ic];
		break;
	case colTypeCheckbox:
		bc = [[NSButtonCell alloc] init];
		[bc setBezelStyle:NSRegularSquareBezelStyle];		// not explicitly stated as such in Interface Builder; extracted with a test program
		[bc setButtonType:NSSwitchButton];
		[bc setBordered:NO];
		[bc setTransparent:NO];
		[bc setAllowsMixedState:NO];
		[bc setTitle:@""];						// no label
		lbm = NSLineBreakByWordWrapping;		// Interface Builder sets this mode for this type
		[c setDataCell:bc];
		break;
	}
	// otherwise just use the current cell
	[c setEditable:editable];
	[[c headerCell] setStringValue:[NSString stringWithUTF8String:name]];
	setSmallControlFont((id) [c headerCell]);
	setStandardControlFont((id) [c dataCell]);
	// the following are according to Interface Builder
	// for the header cell, a stub program had to be written because Interface Builder doesn't support editing header cells directly
	[[c headerCell] setScrollable:NO];
	[[c headerCell] setWraps:NO];
	[[c headerCell] setLineBreakMode:NSLineBreakByTruncatingTail];
	[[c headerCell] setUsesSingleLineMode:NO];
	[[c headerCell] setTruncatesLastVisibleLine:NO];
	[[c dataCell] setScrollable:NO];
	[[c dataCell] setWraps:NO];
	[[c dataCell] setLineBreakMode:lbm];
	[[c dataCell] setUsesSingleLineMode:NO];
	[[c dataCell] setTruncatesLastVisibleLine:NO];
	[toNSTableView(t) addTableColumn:c];
}

void tableUpdate(id t)
{
	[toNSTableView(t) reloadData];
}

// also sets the delegate
void tableMakeDataSource(id table, void *gotable)
{
	goTableDataSource *model;

	model = [goTableDataSource new];
	model->gotable = gotable;
	[toNSTableView(table) setDataSource:model];
	[toNSTableView(table) setDelegate:model];
}

// -[NSTableView sizeToFit] does not actually size to fit
// -[NSTableColumn sizeToFit] is just for the header
// -[NSTableColumn sizeToFit] can work for guessing but overrides user settings
// -[[NSTableColumn headerCell] cellSize] does NOT (despite being documented as returning the minimum needed size)
// Let's write our own to prefer:
// - width of the sum of all columns's current widths
// - height of 5 rows (arbitrary count; seems reasonable) + header view height
// Hopefully this is reasonable.
struct xsize tablePreferredSize(id control)
{
	NSTableView *t;
	NSArray *columns;
	struct xsize s;
	NSUInteger i, n;
	NSTableColumn *c;

	t = toNSTableView(control);
	columns = [t tableColumns];
	n = [columns count];
	s.width = 0;
	for (i = 0; i < n; i++) {
		CGFloat width;

		c = (NSTableColumn *) [columns objectAtIndex:i];
		s.width += (intptr_t) [c width];
	}
	s.height = 5 * (intptr_t) [t rowHeight];
	s.height += (intptr_t) [[t headerView] frame].size.height;
	return s;
}

intptr_t tableSelected(id table)
{
	return (intptr_t) [toNSTableView(table) selectedRow];
}

void tableSelect(id table, intptr_t row)
{
	[toNSTableView(table) deselectAll:table];
	if (row != -1)
		[toNSTableView(table) selectRowIndexes:[NSIndexSet indexSetWithIndex:((NSUInteger) row)] byExtendingSelection:NO];
}

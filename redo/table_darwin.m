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

@interface goTableDataSource : NSObject <NSTableViewDataSource> {
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
	char *str;
	NSString *s;
	intptr_t colnum;

	colnum = ((goTableColumn *) col)->gocolnum;
	str = goTableDataSource_getValue(self->gotable, (intptr_t) row, colnum);
	s = [NSString stringWithUTF8String:str];
	free(str);		// allocated with C.CString() on the Go side
	return s;
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

void tableAppendColumn(id t, intptr_t colnum, char *name)
{
	goTableColumn *c;

	c = [[goTableColumn alloc] initWithIdentifier:nil];
	c->gocolnum = colnum;
	[c setEditable:NO];
	[[c headerCell] setStringValue:[NSString stringWithUTF8String:name]];
	setSmallControlFont((id) [c headerCell]);
	setStandardControlFont((id) [c dataCell]);
	// TODO text layout options
	[toNSTableView(t) addTableColumn:c];
}

void tableUpdate(id t)
{
	[toNSTableView(t) reloadData];
}

void tableMakeDataSource(id table, void *gotable)
{
	goTableDataSource *model;

	model = [goTableDataSource new];
	model->gotable = gotable;
	[toNSTableView(table) setDataSource:model];
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

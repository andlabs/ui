// 29 july 2014

#import "objc_darwin.h"
#import "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define toNSTableView(x) ((NSTableView *) (x))
#define toNSView(x) ((NSView *) (x))

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

	// TODO there has to be a better way to get the column index
	str = goTableDataSource_getValue(self->gotable, (intptr_t) row, (intptr_t) [[view tableColumns] indexOfObject:col]);
	s = [NSString stringWithUTF8String:str];
	free(str);		// allocated with C.CString() on the Go side
	return s;
}

@end

id newTable(void)
{
	NSTableView *t;

	// TODO makerect
	t = [[NSTableView alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	[t setAllowsColumnReordering:NO];
	[t setAllowsColumnResizing:YES];
	// TODO make this an option on all platforms
	[t setAllowsMultipleSelection:NO];
	[t setAllowsEmptySelection:YES];
	[t setAllowsColumnSelection:NO];
	// TODO check against interface builder
	return (id) t;
}

// TODO scroll view

void tableAppendColumn(id t, char *name)
{
	NSTableColumn *c;

	c = [[NSTableColumn alloc] initWithIdentifier:nil];
	[c setEditable:NO];
	[[c headerCell] setStringValue:[NSString stringWithUTF8String:name]];
	// TODO other options
	[toNSTableView(t) addTableColumn:c];
}

void tableUpdate(id t)
{
	[toNSTableView(t) reloadData];
}

// TODO SPLIT
id newScrollView(id content)
{
	NSScrollView *sv;

	sv = [[NSScrollView alloc] initWithFrame:NSMakeRect(0, 0, 100, 100)];
	[sv setDocumentView:toNSView(content)];
	return (id) sv;
}

void tableMakeDataSource(id table, void *gotable)
{
	goTableDataSource *model;

	model = [goTableDataSource new];
	model->gotable = gotable;
	[toNSTableView(table) setDataSource:model];
}

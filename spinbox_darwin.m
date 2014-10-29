// 29 october 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#import <Cocoa/Cocoa.h>

@interface goSpinbox : NSObject {
@public
	void *gospinbox;
}
@property NSInteger integerValue;
@end

@implementation goSpinbox
@synthesize integerValue;
@end

id newSpinboxStepper(void)
{
	NSStepper *s;

	s = [[NSStepper alloc] initWithFrame:NSZeroRect];
	[s setMinValue:0];
	[s setMaxValue:100];
	[s setIncrement:1];
	[s setAutorepeat:YES];		// hold mouse button to step repeatedly
	return (id) s;
}

id spinboxSetup(id textfield, id stepper, void *gospinbox)
{
	goSpinbox *g;

	g = [goSpinbox new];
	g->gospinbox = gospinbox;
	// TODO set any options?
	[textfield bind:@"integerValue" toObject:g withKeyPath:@"integerValue" options:nil];
	[stepper bind:@"integerValue" toObject:g withKeyPath:@"integerValue" options:nil];
	return (id) g;
}

// 29 october 2014

#include "objc_darwin.h"
#include "_cgo_export.h"
#import <Cocoa/Cocoa.h>

#define togoSpinbox(x) ((goSpinbox *) (x))

@interface goSpinbox : NSObject {
@public
	void *gospinbox;
	NSTextField *textfield;
	NSNumberFormatter *formatter;
	NSStepper *stepper;
}
@property NSInteger integerValue;
@property NSInteger minimum;
@property NSInteger maximum;
@end

@implementation goSpinbox
@synthesize integerValue;
@synthesize minimum;
@synthesize maximum;
@end

id newSpinbox(void *gospinbox)
{
	goSpinbox *s;

	s = [goSpinbox new];
	s->gospinbox = gospinbox;
	s->textfield = (NSTextField *) newTextField();
	s->formatter = [NSNumberFormatter new];
	[s->formatter setAllowsFloats:NO];
	[s->textfield setFormatter:s->formatter];
	s->stepper = [[NSStepper alloc] initWithFrame:NSZeroRect];
	[s->stepper setAutorepeat:YES];		// hold mouse button to step repeatedly

	[s setMinimum:0];
	[s setMaximum:100];

	[s->textfield bind:@"integerValue" toObject:s withKeyPath:@"integerValue" options:nil];
	[s->stepper bind:@"integerValue" toObject:s withKeyPath:@"integerValue" options:nil];
//	[s->formatter bind:@"minimum" toObject:s withKeyPath:@"minimum" options:nil];
	[s->stepper bind:@"minValue" toObject:s withKeyPath:@"minimum" options:nil];
//	[s->formatter bind:@"maximum" toObject:s withkeyPath:@"maximum" options:nil];
	[s->stepper bind:@"maxValue" toObject:s withKeyPath:@"maximum" options:nil];

	return (id) s;
}

id spinboxTextField(id spinbox)
{
	return (id) (togoSpinbox(spinbox)->textfield);
}

id spinboxStepper(id spinbox)
{
	return (id) (togoSpinbox(spinbox)->stepper);
}

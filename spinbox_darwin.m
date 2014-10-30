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

	NSInteger value;
	NSInteger minimum;
	NSInteger maximum;
}
@end

@implementation goSpinbox

- (id)init
{
	self = [super init];
	if (self == nil)
		return nil;

	self->textfield = (NSTextField *) newTextField();
	self->formatter = [NSNumberFormatter new];
	[self->formatter setAllowsFloats:NO];
	[self->textfield setFormatter:self->formatter];
	self->stepper = [[NSStepper alloc] initWithFrame:NSZeroRect];
	[self->stepper setIncrement:1];
	[self->stepper setAutorepeat:YES];		// hold mouse button to step repeatedly

	// TODO how SHOULD the formatter treat invald input?

	[self setMinimum:0];
	[self setMaximum:100];
	[self setValue:0];

	[self->textfield setDelegate:self];
	[self->stepper setTarget:self];
	[self->stepper setAction:@selector(stepperClicked:)];

	return self;
}

- (void)setValue:(NSInteger)value
{
	self->value = value;
	[self->textfield setIntegerValue:value];
	[self->stepper setIntegerValue:value];
}

- (void)setMinimum:(NSInteger)min
{
	self->minimum = min;
	[self->formatter setMinimum:[NSNumber numberWithInteger:min]];
	[self->stepper setMinValue:((double) min)];
}

- (void)setMaximum:(NSInteger)max
{
	self->maximum = max;
	[self->formatter setMaximum:[NSNumber numberWithInteger:max]];
	[self->stepper setMaxValue:((double) max)];
}

- (IBAction)stepperClicked:(id)sender
{
	[self setValue:[self->stepper integerValue]];
}

- (void)controlTextDidChange:(NSNotification *)note
{
	[self setValue:[self->textfield integerValue]];
}

@end

id newSpinbox(void *gospinbox)
{
	goSpinbox *s;

	s = [goSpinbox new];
	s->gospinbox = gospinbox;
	return s;
}

id spinboxTextField(id spinbox)
{
	return (id) (togoSpinbox(spinbox)->textfield);
}

id spinboxStepper(id spinbox)
{
	return (id) (togoSpinbox(spinbox)->stepper);
}

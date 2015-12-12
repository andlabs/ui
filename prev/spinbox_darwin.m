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

- (id)initWithMinimum:(NSInteger)minimum maximum:(NSInteger)maximum
{
	self = [super init];
	if (self == nil)
		return nil;

	self->textfield = (NSTextField *) newTextField();

	self->formatter = [NSNumberFormatter new];
	[self->formatter setFormatterBehavior:NSNumberFormatterBehavior10_4];
	[self->formatter setLocalizesFormat:NO];
	[self->formatter setUsesGroupingSeparator:NO];
	[self->formatter setHasThousandSeparators:NO];
	[self->formatter setAllowsFloats:NO];
	// TODO partial string validation?
	[self->textfield setFormatter:self->formatter];

	self->stepper = [[NSStepper alloc] initWithFrame:NSZeroRect];
	[self->stepper setIncrement:1];
	[self->stepper setValueWraps:NO];
	[self->stepper setAutorepeat:YES];		// hold mouse button to step repeatedly

	// TODO how SHOULD the formatter treat invald input?

	[self setMinimum:minimum];
	[self setMaximum:maximum];
	[self setValue:self->minimum];

	[self->textfield setDelegate:self];
	[self->stepper setTarget:self];
	[self->stepper setAction:@selector(stepperClicked:)];

	return self;
}

- (void)setValue:(NSInteger)value
{
	self->value = value;
	if (self->value < self->minimum)
		self->value = self->minimum;
	if (self->value > self->maximum)
		self->value = self->maximum;
	[self->textfield setIntegerValue:self->value];
	[self->stepper setIntegerValue:self->value];
}

- (void)setMinimum:(NSInteger)min
{
	self->minimum = min;
	[self->formatter setMinimum:[NSNumber numberWithInteger:self->minimum]];
	[self->stepper setMinValue:((double) (self->minimum))];
}

- (void)setMaximum:(NSInteger)max
{
	self->maximum = max;
	[self->formatter setMaximum:[NSNumber numberWithInteger:self->maximum]];
	[self->stepper setMaxValue:((double) (self->maximum))];
}

- (IBAction)stepperClicked:(id)sender
{
	[self setValue:[self->stepper integerValue]];
	spinboxChanged(self->gospinbox);
}

- (void)controlTextDidChange:(NSNotification *)note
{
	[self setValue:[self->textfield integerValue]];
	spinboxChanged(self->gospinbox);
}

@end

id newSpinbox(void *gospinbox, intmax_t minimum, intmax_t maximum)
{
	goSpinbox *s;

	s = [[goSpinbox new] initWithMinimum:((NSInteger) minimum) maximum:((NSInteger) maximum)];
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

intmax_t spinboxValue(id spinbox)
{
	return (intmax_t) (togoSpinbox(spinbox)->value);
}

void spinboxSetValue(id spinbox, intmax_t value)
{
	[togoSpinbox(spinbox) setValue:((NSInteger) value)];
}

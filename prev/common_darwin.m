// 11 august 2014

#include "objc_darwin.h"
#include <Cocoa/Cocoa.h>

void disableAutocorrect(id onwhat)
{
	NSTextView *tv;

	tv = (NSTextView *) onwhat;
	[tv setEnabledTextCheckingTypes:0];
	[tv setAutomaticDashSubstitutionEnabled:NO];
	// don't worry about automatic data detection; it won't change stringValue (thanks pretty_function in irc.freenode.net/#macdev)
	[tv setAutomaticSpellingCorrectionEnabled:NO];
	[tv setAutomaticTextReplacementEnabled:NO];
	[tv setAutomaticQuoteSubstitutionEnabled:NO];
	[tv setAutomaticLinkDetectionEnabled:NO];
	[tv setSmartInsertDeleteEnabled:NO];
}

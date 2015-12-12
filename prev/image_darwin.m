// 16 august 2014

#import "objc_darwin.h"
#import <Cocoa/Cocoa.h>

#define toNSInteger(x) ((NSInteger) (x))

id toTableImage(void *pixels, intptr_t width, intptr_t height, intptr_t stride)
{
	NSBitmapImageRep *bitmap;
	NSImage *image;

	// we can't just hand it pixels; we need to make a copy
	bitmap = [[NSBitmapImageRep alloc]
		initWithBitmapDataPlanes:NULL
		pixelsWide:toNSInteger(width)
		pixelsHigh:toNSInteger(height)
		bitsPerSample:8
		samplesPerPixel:4
		hasAlpha:YES
		isPlanar:NO
		colorSpaceName:NSDeviceRGBColorSpace
		bitmapFormat:0
		bytesPerRow:toNSInteger(stride)
		bitsPerPixel:32];
	memcpy((void *) [bitmap bitmapData], pixels, [bitmap bytesPerPlane]);
	image = [[NSImage alloc] initWithSize:NSMakeSize((CGFloat) width, (CGFloat) height)];
	[image addRepresentation:bitmap];
	// TODO release bitmap?
	return (id) image;
}

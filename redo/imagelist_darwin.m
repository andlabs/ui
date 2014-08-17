// 16 august 2014

#import "objc_darwin.h"
#import <Cocoa/Cocoa.h>

#define toNSInteger(x) ((NSInteger) (x))

// TODO top two pixels of 16x16 images are green?

id toImageListImage(void *pixels, intptr_t width, intptr_t height, intptr_t stride)
{
	unsigned char *planes[1];			// NSBitmapImageRep wants an array of planes; we have one plane
	NSBitmapImageRep *bitmap;
	NSImage *image;

	planes[0] = (unsigned char *) pixels;
	bitmap = [[NSBitmapImageRep alloc]
		initWithBitmapDataPlanes:planes
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
	image = [[NSImage alloc] initWithSize:NSMakeSize((CGFloat) width, (CGFloat) height)];
	[image addRepresentation:bitmap];
	return (id) image;
}

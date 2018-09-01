// 21 august 2018

package ui

import (
	"image"
)

// #include "pkgui.h"
import "C"

// Image stores an image for display on screen.
// 
// Images are built from one or more representations, each with the
// same aspect ratio but a different pixel size. Package ui
// automatically selects the most appropriate representation for
// drawing the image when it comes time to draw the image; what
// this means depends on the pixel density of the target context.
// Therefore, one can use Image to draw higher-detailed images on
// higher-density displays. The typical use cases are either:
// 
// 	- have just a single representation, at which point all screens
// 	  use the same image, and thus uiImage acts like a simple
// 	  bitmap image, or
// 	- have two images, one at normal resolution and one at 2x
// 	  resolution; this matches the current expectations of some
// 	  desktop systems at the time of writing (mid-2018)
//
// Image allocates OS resources; you must explicitly free an Image
// when you are finished with it.
type Image struct {
	i	*C.uiImage
}

// NewImage creates a new Image with the given width and
// height. This width and height should be the size in points of the
// image in the device-independent case; typically this is the 1x size.
func NewImage(width, height float64) *Image {
	return &Image{
		i:	C.uiNewImage(C.double(width), C.double(height)),
	}
}

// Free frees the Image.
func (i *Image) Free() {
	C.uiFreeImage(i.i)
}

// Append adds the given image as a representation of the Image.
func (i *Image) Append(img *image.RGBA) {
	cpix := C.CBytes(img.Pix)
	defer C.free(cpix)
	C.uiImageAppend(i.i, cpix,
		C.int(img.Rect.Dx()),
		C.int(img.Rect.Dy()),
		C.int(img.Stride))
}

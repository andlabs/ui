// 16 august 2014

package ui

import (
	"image"
)

// ImageList is a list of images that can be used in the rows of a Table or Tree.
// ImageList maintains a copy of each image added.
// Images in an ImageList will be automatically scaled to the needed size.
type ImageList interface {
	// Append inserts an image into the ImageList.
	Append(i *image.RGBA)

	// Len returns the number of images in the ImageList.
	Len() ImageIndex

	imageListApply
}

// NewImageList creates a new ImageList.
// The ImageList is initially empty.
func NewImageList() ImageList {
	return newImageList()
}

// ImageIndex is a special type used to denote an entry in a Table or Tree's ImageList.
type ImageIndex int

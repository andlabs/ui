// 16 august 2014

package ui

import (
	"image"
	"unsafe"
)

// #include "objc_darwin.h"
import "C"

type imagelist struct {
	list []C.id
}

func newImageList() ImageList {
	return new(imagelist)
}

func (i *imagelist) Append(img *image.RGBA) {
	id := C.toImageListImage(
		unsafe.Pointer(pixelData(img)), C.intptr_t(img.Rect.Dx()), C.intptr_t(img.Rect.Dy()), C.intptr_t(img.Stride))
	i.list = append(i.list, id)
}

func (i *imagelist) Len() ImageIndex {
	return ImageIndex(len(i.list))
}

type imageListApply interface {
	apply(*[]C.id)
}

func (i *imagelist) apply(out *[]C.id) {
	*out = make([]C.id, len(i.list))
	copy(*out, i.list)
}

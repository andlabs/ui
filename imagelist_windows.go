// 16 august 2014

package ui

import (
	"image"
	"unsafe"
)

// #include "winapi_windows.h"
import "C"

type imagelist struct {
	list   []C.HBITMAP
	width  []int
	height []int
}

func newImageList() ImageList {
	return new(imagelist)
}

func (i *imagelist) Append(img *image.RGBA) {
	i.list = append(i.list, C.unscaledBitmap(unsafe.Pointer(img), C.intptr_t(img.Rect.Dx()), C.intptr_t(img.Rect.Dy())))
	i.width = append(i.width, img.Rect.Dx())
	i.height = append(i.height, img.Rect.Dy())
}

func (i *imagelist) Len() ImageIndex {
	return ImageIndex(len(i.list))
}

type imageListApply interface {
	apply(C.HWND, C.UINT)
}

func (i *imagelist) apply(hwnd C.HWND, uMsg C.UINT) {
	width := C.GetSystemMetrics(C.SM_CXSMICON)
	height := C.GetSystemMetrics(C.SM_CYSMICON)
	il := C.newImageList(width, height)
	for index := range i.list {
		C.addImage(il, hwnd, i.list[index], C.int(i.width[index]), C.int(i.height[index]), width, height)
	}
	C.SendMessageW(hwnd, uMsg, 0, C.LPARAM(uintptr(unsafe.Pointer(il))))
}

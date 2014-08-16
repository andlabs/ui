# Images in Tables and Trees

In yet another "blame Windows" moment:

```go
type ImageList interface {
	Add(image *image.Image)
	Remove(index int)
	Len() int
}
func NewImageList() ImageList

type Table/Tree interface {
	// ...
	LoadImageList(ImageList)
}

type ImageIndex int
```

For Table, a field of type ImageIndex represents an index into the ImageList.

For Tree, there is a field ImageIndex which contains the index of the entry.

Note the name of the methods on both being LoadImageList(): the image list is copied, and any future changes will **not** be reflected. (This is to accomodate check boxes in Table, which must be done manually.)

Icons scale automatically to the best possible size.

On Windows this is GetSystemMetrics(SM_CX/YSMICON)

On GTK+ this will be determined by trial and error
	[11:14] <LiamW> andlabs: probably
	[11:14] <LiamW> GTK_ICON_SIZE_SMALL_TOOLBAR
	[11:15] <baedert> the procedure is to just try all of them and see what looks best :p

On Mac OS X this is rowHeight

TODO which side does the scaling, C or Go?

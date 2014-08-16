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

TODO appropriate size of images in an ImageList

# Tree

Unlike Table, Tree can only store a set of a single data type. (Blame Windows.)

```go
type TreeData struct {
	Checked		bool
	Image		ImageIndex
	Text			string
	Children		[]TreeData		// TODO does this need to be *[]TreeData?
}
```

(the facilities for Images has yet to be designed)

Tree itself will operate similarly to Table:

```go
type Tree struct {
	Control
	sync.Locker		// with Unlock() refreshing the view
	Data() *[]TreeData
	SetHasCheckboxes(bool)
	SetHasImages(bool)
}
```

By default, a Tree only shows the Text field.

A Tree path is just an `[]int` with each element set to the consecutive index in Children. For example:

```go
i := []int{3, 4, 5}
value := tree.Data()[i[0]].Children[i[1]].Children[i[2]].Text
```

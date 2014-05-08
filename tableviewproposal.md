# andlabs/ui table view/tree view proposal
<tt>(text that looks like this is optional)</tt>

```go
type TableView struct {
    Selected chan struct{}
    Data     interface{}
    // contains hidden or unexported fields
}
```
A TableView is a Control that displays rows of structured data organized by columns.

Data is a slice of objects of a structure containing multiple fields. Each field represents a column in the table. Column names are, by default, the exported field name; the struct tag ui:"Custom name" can be used to specify a custom field name. For example:
```go
type Person struct {
    Name        string
    Address     Address
    PhoneNumber PhoneNumber `ui:"Phone Number"`
}
```
Data is displayed using the fmt package's %v rule. The structure must satisfy sync.Locker.

<tt>If one of the members is of type slice of the structure type, then any element of the main slice with a Children whose length is nonzero represents child nodes. For example:
```go
type File struct {
    Filename string
    Size     int64
    Type     FileType
    Contents []File
}
```
In this case, File.Contents specifies children of the parent File.</tt>

```go
func NewTableView(initData interface{}) *TableView
```
Creates a new TableView with the specified initial data. This also determines the data type of the TableView; after this, all accesses to the data are made through the Data field of TableView. NewTableView() panics if initData is nil or not a slice of structures. The slice may be empty. (TODO slice of pointers to structures?) <tt>NewTableView() also panics if the structure has more than one possible children field.</tt>

```go
// if trees are not supported
func (t *TableView) Append(items ...interface{})
func (t *TableView) InsertBefore(index int, items ...interface{})
func (t *TableView) Delete(indices ...int)

// if trees are supported
func (t *TableView) Append(path []int, items ...interface{})
func (t *TableView) InsertBefore(path []int, items ...interface{})
func (t *TableView) Delete(path []int, indices ...int)
```
Standard methods to manipulate data in the TableView. These methods hold the write lock upon entry and release it upon exit. They panic if any index is invalid. <tt>path specifies which node of the tree to append to. If path has length zero, the operation is performed on the top level; if path has length one, the operation is performed on the children of the first entry in the list; and so on and so forth. Each element of path is the index relative to the first item at the level (so []int{4, 2, 1} specifies the fifth entry's third child's second child's children).</tt>

```go
func (t *TableView) Lock()
func (t *TableView) Unlock()
func (t *TableView) RLock()
func (t *TableView) RUnlock()
```
For more complex manipulations, TableView acts as a sync.RWMutex. Any goroutine holding the read lock may access t.Data, but cannot change it. Any goroutine holding the regular lock may modify t.Data. Before t.Unlock() returns, it automatically refreshes the TaleView with the new contents of Data.

```go
// if trees are not supported
func (t *TableView) Selection() []int
func (t *TableView) Select(indices ...int)

// if trees are supported
func (t *TableView) Selection() [][]int
func (t *TableView) Select(indices ...[]int)

// or should these be SelectedIndices() and SelectIndices() for consistency?
```
Methods that act on TableView row selection. These methods hold the read lock on entry and release it on exit. <tt>Each entry in the returned slice consists of a path followed by the selected index of the child. A slice of length 1 indicates that a top-level entry has been selected. The slices shall not be of length zero; passing one in will panic. (TODO this means that multiple children node will have a copy of path each; that should be fixed...)</tt>


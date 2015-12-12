# Sidebar Control

```go
type Sidebar interface {
	Control

	AppendCategory(text string)
	DeleteCategory(index int)

	AppendItem(category int, name string)
	DeleteItem(category int, index int)

	Selection() (category int, index int)		// or Selected()?
	Select(category int, index int)

	OnSelected(func())
}
```

Simple two-level sidebars.

Could have images on each item in the future.

## Mac OS X
Source List NSTableView (need to see how this will work)

## GTK+
GTK_STYLE_CLASS_SIDEBAR (available in 3.4); see how GtkPlacesSidebar implements this
	- other programs that do: Rhythmbox

## Windows
????

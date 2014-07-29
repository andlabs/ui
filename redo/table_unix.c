/* 29 july 2014 */

#include "gtk_unix.h"
#include "_cgo_export.h"

void tableAppendColumn(GtkTreeView *table, gchar *name)
{
	GtkTreeViewColumn *col;
	GtkCellRenderer *renderer;

	renderer = gtk_cell_renderer_text_new();
	col = gtk_tree_view_column_new_with_attributes(name, renderer,
		/* TODO */
		NULL);
	gtk_tree_view_append_column(table, col);
}

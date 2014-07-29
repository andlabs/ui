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

static void goTableModel_initGtkTreeModel(GtkTreeModelIface *);

G_DEFINE_TYPE_WITH_CODE(goTableModel, goTableModel, G_TYPE_OBJECT,
	G_IMPLEMENT_INTERFACE(GTK_TYPE_TREE_MODEL, goTableModel_initGtkTreeModel))

static void goTableModel_init(goTableModel *t)
{
	/* do nothing */
}

static void goTableModel_dispose(GObject *obj)
{
	G_OBJECT_CLASS(goTableModel_parent_class)->dispose(obj);
}

/* and now for the interface function definitions */

static void goTableModel_finalize(GObject *obj)
{
	G_OBJECT_CLASS(goTableModel_parent_class)->finalize(obj);
}

static GtkTreeModelFlags goTableModel_get_flags(GtkTreeModel *model)
{
	return GTK_TREE_MODEL_LIST_ONLY;
}

static void goTableModel_initGtkTreeModel(GtkTreeModelIface *interface)
{
	GtkTreeModelIface *chain;

	chain = (GtkTreeModelIface *) g_type_interface_peek_parent(interface);
#define DEF(x) interface->x = goTableModel_ ## x;
#define CHAIN(x) interface->x = chain->x;
	/* signals */
	CHAIN(row_changed)
	CHAIN(row_inserted)
	CHAIN(row_has_child_toggled)
	CHAIN(row_deleted)
	CHAIN(rows_reordered)
	/* vtable */
	DEF(get_flags)
	CHAIN(get_n_columns)
	CHAIN(get_column_type)
	CHAIN(get_iter)
	CHAIN(get_path)
	CHAIN(get_value)
	CHAIN(iter_next)
	CHAIN(iter_previous)
	CHAIN(iter_children)
	CHAIN(iter_has_child)
	CHAIN(iter_n_children)
	CHAIN(iter_nth_child)
	CHAIN(iter_parent)
	CHAIN(ref_node)
	CHAIN(unref_node)
}

static GParamSpec *goTableModelProperties[2];

static void goTableModel_set_property(GObject *obj, guint id, const GValue *value, GParamSpec *pspec)
{
	goTableModel *t = (goTableModel *) obj;

	if (id == 1) {
		t->gotable = (void *) g_value_get_pointer(value);
		return;
	}
	G_OBJECT_WARN_INVALID_PROPERTY_ID(t, id, pspec);
}

static void goTableModel_get_property(GObject *obj, guint id, GValue *value, GParamSpec *pspec)
{
	G_OBJECT_WARN_INVALID_PROPERTY_ID((goTableModel *) obj, id, pspec);
}

static void goTableModel_class_init(goTableModelClass *class)
{
	G_OBJECT_CLASS(class)->dispose = goTableModel_dispose;
	G_OBJECT_CLASS(class)->finalize = goTableModel_finalize;
	G_OBJECT_CLASS(class)->set_property = goTableModel_set_property;
	G_OBJECT_CLASS(class)->get_property = goTableModel_get_property;

	goTableModelProperties[1] = g_param_spec_pointer(
		"gotable",
		"go-table",
		"Go-side *table value",
		G_PARAM_WRITABLE | G_PARAM_CONSTRUCT_ONLY | G_PARAM_STATIC_STRINGS);
	g_object_class_install_properties(G_OBJECT_CLASS(class), 2, goTableModelProperties);
}

goTableModel *newTableModel(void *gotable)
{
	return (goTableModel *) g_object_new(goTableModel_get_type(), "gotable", (gpointer) gotable, NULL);
}

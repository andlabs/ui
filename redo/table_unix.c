// +build !windows,!darwin

/* 29 july 2014 */

#include "gtk_unix.h"
#include "_cgo_export.h"

void tableAppendColumn(GtkTreeView *table, gint index, gchar *name)
{
	GtkTreeViewColumn *col;
	GtkCellRenderer *renderer;

	renderer = gtk_cell_renderer_text_new();
	col = gtk_tree_view_column_new_with_attributes(name, renderer,
		"text", index,
		NULL);
	/* allow columns to be resized */
	gtk_tree_view_column_set_resizable(col, TRUE);
	gtk_tree_view_append_column(table, col);
}

/*
how our GtkTreeIters are stored:
	stamp: either GOOD_STAMP or BAD_STAMP
	user_data: row index
TODO verify range of gint
*/
#define GOOD_STAMP 0x1234
#define BAD_STAMP 0x5678

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

static void goTableModel_finalize(GObject *obj)
{
	G_OBJECT_CLASS(goTableModel_parent_class)->finalize(obj);
}

/* and now for the interface function definitions */

static GtkTreeModelFlags goTableModel_get_flags(GtkTreeModel *model)
{
	return GTK_TREE_MODEL_LIST_ONLY;
}

/* get_n_columns in Go */

static GType goTableModel_get_column_type(GtkTreeModel *model, gint column)
{
	return G_TYPE_STRING;
}

static gboolean goTableModel_get_iter(GtkTreeModel *model, GtkTreeIter *iter, GtkTreePath *path)
{
	goTableModel *t = (goTableModel *) model;
	gint index;

	if (gtk_tree_path_get_depth(path) != 1)
		goto bad;
	index = gtk_tree_path_get_indices(path)[0];
	if (index < 0)
		goto bad;
	if (index >= goTableModel_getRowCount(t->gotable))
		goto bad;
	iter->stamp = GOOD_STAMP;
	iter->user_data = (gpointer) index;
	return TRUE;
bad:
	iter->stamp = BAD_STAMP;
	return FALSE;
}

static GtkTreePath *goTableModel_get_path(GtkTreeModel *model, GtkTreeIter *iter)
{
	if (iter->stamp != GOOD_STAMP)
		return NULL;		/* TODO is this right? */
	return gtk_tree_path_new_from_indices((gint) iter->user_data, -1);
}

static void goTableModel_get_value(GtkTreeModel *model, GtkTreeIter *iter, gint column, GValue *value)
{
	goTableModel *t = (goTableModel *) model;
	gchar *str;

	/* TODO what if iter is invalid? */
	/* we (actually cgo) allocated str with malloc(), not g_malloc(), so let's free it explicitly and give the GValue a copy to be safe */
	str = goTableModel_do_get_value(t->gotable, (gint) iter->user_data, column);
	/* value is uninitialized */
	g_value_init(value, G_TYPE_STRING);
	g_value_set_string(value, str);
	free(str);
}

static gboolean goTableModel_iter_next(GtkTreeModel *model, GtkTreeIter *iter)
{
	goTableModel *t = (goTableModel *) model;
	gint index;

	if (iter->stamp != GOOD_STAMP)
		return FALSE;		/* TODO correct? */
	index = (gint) iter->user_data;
	index++;
	iter->user_data = (gpointer) index;
	if (index >= goTableModel_getRowCount(t->gotable)) {
		iter->stamp = BAD_STAMP;
		return FALSE;
	}
	return TRUE;
}

static gboolean goTableModel_iter_previous(GtkTreeModel *model, GtkTreeIter *iter)
{
	goTableModel *t = (goTableModel *) model;
	gint index;

	if (iter->stamp != GOOD_STAMP)
		return FALSE;		/* TODO correct? */
	index = (gint) iter->user_data;
	index--;
	iter->user_data = (gpointer) index;
	if (index < 0) {
		iter->stamp = BAD_STAMP;
		return FALSE;
	}
	return TRUE;
}

static gboolean goTableModel_iter_children(GtkTreeModel *model, GtkTreeIter *child, GtkTreeIter *parent)
{
	goTableModel *t = (goTableModel *) model;

	if (parent == NULL && goTableModel_getRowCount(t->gotable) > 0) {
		child->stamp = GOOD_STAMP;
		child->user_data = 0;
		return TRUE;
	}
	child->stamp = BAD_STAMP;
	return FALSE;
}

static gboolean goTableModel_iter_has_child(GtkTreeModel *model, GtkTreeIter *iter)
{
	return FALSE;
}

static gint goTableModel_iter_n_children(GtkTreeModel *model, GtkTreeIter *iter)
{
	goTableModel *t = (goTableModel *) model;

	if (iter == NULL)
		return goTableModel_getRowCount(t->gotable);
	return 0;
}

static gboolean goTableModel_iter_nth_child(GtkTreeModel *model, GtkTreeIter *child, GtkTreeIter *parent, gint n)
{
	goTableModel *t = (goTableModel *) model;

	if (parent == NULL && n >= 0 && n < goTableModel_getRowCount(t->gotable)) {
		child->stamp = GOOD_STAMP;
		child->user_data = (gpointer) n;
		return TRUE;
	}
	child->stamp = BAD_STAMP;
	return FALSE;
}

static gboolean goTableModel_iter_parent(GtkTreeModel *model, GtkTreeIter *parent, GtkTreeIter *child)
{
	parent->stamp = BAD_STAMP;
	return FALSE;
}

/* end of interface definitions */

static void goTableModel_initGtkTreeModel(GtkTreeModelIface *interface)
{
	/* don't chain; we have nothing to chain to */
#define DEF(x) interface->x = goTableModel_ ## x;
	DEF(get_flags)
	DEF(get_n_columns)
	DEF(get_column_type)
	DEF(get_iter)
	DEF(get_path)
	DEF(get_value)
	DEF(iter_next)
	DEF(iter_previous)
	DEF(iter_children)
	DEF(iter_has_child)
	DEF(iter_n_children)
	DEF(iter_nth_child)
	DEF(iter_parent)
	/* no need for ref_node and unref_node */
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

/* somewhat naive, but the only alternatives seem to be unloading/reloading the model (or the view!), which is bleh */
void tableUpdate(goTableModel *t, gint old, gint new)
{
	gint i;
	gint nUpdate;
	GtkTreePath *path;
	GtkTreeIter iter;

	iter.stamp = GOOD_STAMP;
	/* first, append extra items */
	if (old < new) {
		for (i = old; i < new; i++) {
			path = gtk_tree_path_new_from_indices(i, -1);
			iter.user_data = (gpointer) i;
			g_signal_emit_by_name(t, "row-inserted", path, &iter);
		}
		nUpdate = old;
	} else
		nUpdate = new;
	/* next, update existing items */
	for (i = 0; i < nUpdate; i++) {
		path = gtk_tree_path_new_from_indices(i, -1);
		iter.user_data = (gpointer) i;
		g_signal_emit_by_name(t, "row-updated", path, &iter);
	}
	/* finally, remove deleted items */
	if (old > new)
		for (i = new; i < old; i++) {
			/* note that we repeatedly remove the row at index new, as that changes with each removal; NOT i */
			path = gtk_tree_path_new_from_indices(new, -1);
			/* row-deleted has no iter */
			g_signal_emit_by_name(t, "row-deleted", path);
		}
}

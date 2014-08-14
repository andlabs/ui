/* 13 august 2014 */

#include "gtk_unix.h"
#include "_cgo_export.h"

#define GOCONTAINER_TYPE (goContainer_get_type())
#define GOCONTAINER(obj) (G_TYPE_CHECK_INSTANCE_CAST((obj), GOCONTAINER_TYPE, goContainer))
#define IS_GOCONTAINER(obj) (G_TYPE_CHECK_INSTANCE_TYPE((obj), GOCONTAINER_TYPE))
#define GOCONTAINER_CLASS(class) (G_TYPE_CHECK_CLASS_CAST((class), GOCONTAINER_TYPE, goContainerClass))
#define IS_GOCONTAINER_CLASS(class) (G_TYPE_CHECK_CLASS_TYPE((class), GOCONTAINER_TYPE))
#define GOCONTAINER_GET_CLASS(obj) (G_TYPE_INSTANCE_GET_CLASS((obj), GOCONTAINER_TYPE, goContainerClass))

typedef struct goContainer goContainer;
typedef struct goContainerClass goContainerClass;

struct goContainer {
	GtkContainer parent_instance;
	void *gocontainer;
	GPtrArray *children;		/* for forall() */
};

struct goContainerClass {
	GtkContainerClass parent_class;
};

G_DEFINE_TYPE(goContainer, goContainer, GTK_TYPE_CONTAINER)

static void goContainer_init(goContainer *c)
{
	c->children = g_ptr_array_new();
	gtk_widget_set_has_window(GTK_WIDGET(c), FALSE);
}

static void goContainer_dispose(GObject *obj)
{
	g_ptr_array_unref(GOCONTAINER(obj)->children);
	G_OBJECT_CLASS(goContainer_parent_class)->dispose(obj);
}

static void goContainer_finalize(GObject *obj)
{
	G_OBJECT_CLASS(goContainer_parent_class)->finalize(obj);
}

static void goContainer_add(GtkContainer *container, GtkWidget *widget)
{
	gtk_widget_set_parent(widget, GTK_WIDGET(container));
	g_ptr_array_add(GOCONTAINER(container)->children, widget);
}

static void goContainer_remove(GtkContainer *container, GtkWidget *widget)
{
	gtk_widget_unparent(widget);
	g_ptr_array_remove(GOCONTAINER(container)->children, widget);
}

static void goContainer_size_allocate(GtkWidget *widget, GtkAllocation *allocation)
{
	gtk_widget_set_allocation(widget, allocation);
	containerResizing(GOCONTAINER(widget)->gocontainer, allocation);
}

/*
static void goContainer_get_preferred_width(GtkWidget *widget, gint *min, gint *nat)
{
	if (GOCONTAINER(widget)->child != NULL) {
		gtk_widget_get_preferred_width(GOCONTAINER(widget)->child, min, nat);
		return;
	}
	if (min != NULL)
		*min = 0;
	if (nat != NULL)
		*nat = 0;
}

static void goContainer_get_preferred_height(GtkWidget *widget, gint *min, gint *nat)
{
	if (GOCONTAINER(widget)->child != NULL) {
		gtk_widget_get_preferred_height(GOCONTAINER(widget)->child, min, nat);
		return;
	}
	if (min != NULL)
		*min = 0;
	if (nat != NULL)
		*nat = 0;
}
*/

static void goContainer_forall(GtkContainer *container, gboolean includeInternals, GtkCallback callback, gpointer data)
{
	/* TODO is this safe? */
	g_ptr_array_foreach(GOCONTAINER(container)->children, callback, data);
}

static void goContainer_class_init(goContainerClass *class)
{
	G_OBJECT_CLASS(class)->dispose = goContainer_dispose;
	G_OBJECT_CLASS(class)->finalize = goContainer_finalize;
	GTK_WIDGET_CLASS(class)->size_allocate = goContainer_size_allocate;
//	GTK_WIDGET_CLASS(class)->get_preferred_width = goContainer_get_preferred_width;
//	GTK_WIDGET_CLASS(class)->get_preferred_height = goContainer_get_preferred_height;
	GTK_CONTAINER_CLASS(class)->add = goContainer_add;
	GTK_CONTAINER_CLASS(class)->remove = goContainer_remove;
	GTK_CONTAINER_CLASS(class)->forall = goContainer_forall;
}

GtkWidget *newContainer(void *gocontainer)
{
	goContainer *c;

	c = (goContainer *) g_object_new(GOCONTAINER_TYPE, NULL);
	c->gocontainer = gocontainer;
	return GTK_WIDGET(c);
}

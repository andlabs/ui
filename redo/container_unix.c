// 13 august 2014
#include <gtk/gtk.h>

#define CUSTOM_CONTAINER_TYPE (customContainer_get_type())
#define CUSTOM_CONTAINER(obj) (G_TYPE_CHECK_INSTANCE_CAST((obj), CUSTOM_CONTAINER_TYPE, CustomContainer))
#define IS_CUSTOM_CONTAINER(obj) (G_TYPE_CHECK_INSTANCE_TYPE((obj), CUSTOM_CONTAINER_TYPE))
#define CUSTOM_CONTAINER_CLASS(class) (G_TYPE_CHECK_CLASS_CAST((class), CUSTOM_CONTAINER_TYPE, CustomContainerClass))
#define IS_CUSTOM_CONTAINER_CLASS(class) (G_TYPE_CHECK_CLASS_TYPE((class), CUSTOM_CONTAINER_TYPE))
#define CUSTOM_CONTAINER_GET_CLASS(obj) (G_TYPE_INSTANCE_GET_CLASS((obj), CUSTOM_CONTAINER_TYPE, CustomContainerClass))

typedef struct CustomContainer CustomContainer;
typedef struct CustomContainerClass CustomContainerClass;

struct CustomContainer {
	GtkContainer parent_instance;
	GtkWidget *child;
};

struct CustomContainerClass {
	GtkContainerClass parent_class;
};

G_DEFINE_TYPE(CustomContainer, customContainer, GTK_TYPE_CONTAINER)

static void customContainer_init(CustomContainer *c)
{
	gtk_widget_set_has_window(GTK_WIDGET(c), FALSE);
}

static void customContainer_dispose(GObject *obj)
{
	G_OBJECT_CLASS(customContainer_parent_class)->dispose(obj);
}

static void customContainer_finalize(GObject *obj)
{
	G_OBJECT_CLASS(customContainer_parent_class)->finalize(obj);
}

static void customContainer_add(GtkContainer *container, GtkWidget *widget)
{
	gtk_widget_set_parent(widget, GTK_WIDGET(container));
	CUSTOM_CONTAINER(container)->child = widget;
}

static void customContainer_remove(GtkContainer *container, GtkWidget *widget)
{
	gtk_widget_unparent(widget);
	CUSTOM_CONTAINER(container)->child = NULL;
}

static void customContainer_size_allocate(GtkWidget *widget, GtkAllocation *allocation)
{
	gtk_widget_set_allocation(widget, allocation);
	if (CUSTOM_CONTAINER(widget)->child != NULL)
		gtk_widget_size_allocate(CUSTOM_CONTAINER(widget)->child, allocation);
}

static void customContainer_get_preferred_width(GtkWidget *widget, gint *min, gint *nat)
{
	if (CUSTOM_CONTAINER(widget)->child != NULL) {
		gtk_widget_get_preferred_width(CUSTOM_CONTAINER(widget)->child, min, nat);
		return;
	}
	if (min != NULL)
		*min = 0;
	if (nat != NULL)
		*nat = 0;
}

static void customContainer_get_preferred_height(GtkWidget *widget, gint *min, gint *nat)
{
	if (CUSTOM_CONTAINER(widget)->child != NULL) {
		gtk_widget_get_preferred_height(CUSTOM_CONTAINER(widget)->child, min, nat);
		return;
	}
	if (min != NULL)
		*min = 0;
	if (nat != NULL)
		*nat = 0;
}

static void customContainer_forall(GtkContainer *container, gboolean includeInternals, GtkCallback callback, gpointer data)
{
	if (CUSTOM_CONTAINER(container)->child != NULL)
		(*callback)(CUSTOM_CONTAINER(container)->child, data);
}

static void customContainer_class_init(CustomContainerClass *class)
{
	G_OBJECT_CLASS(class)->dispose = customContainer_dispose;
	G_OBJECT_CLASS(class)->finalize = customContainer_finalize;
	GTK_WIDGET_CLASS(class)->size_allocate = customContainer_size_allocate;
//	GTK_WIDGET_CLASS(class)->get_preferred_width = customContainer_get_preferred_width;
//	GTK_WIDGET_CLASS(class)->get_preferred_height = customContainer_get_preferred_height;
	GTK_CONTAINER_CLASS(class)->add = customContainer_add;
	GTK_CONTAINER_CLASS(class)->remove = customContainer_remove;
	GTK_CONTAINER_CLASS(class)->forall = customContainer_forall;
}

int main(void)
{
	gtk_init(NULL, NULL);
	GtkWidget *mainwin = gtk_window_new(GTK_WINDOW_TOPLEVEL);
	GtkWidget *cc = g_object_new(CUSTOM_CONTAINER_TYPE, NULL);
	gtk_container_add(GTK_CONTAINER(cc), gtk_button_new_with_label("Test"));
	gtk_container_add(GTK_CONTAINER(mainwin), cc);
	gtk_widget_show_all(mainwin);
	gtk_main();
	return 0;
}

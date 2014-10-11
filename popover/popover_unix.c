// 11 october 2014
// #qo pkg-config: gtk+-3.0
#include <gtk/gtk.h>

#define GOPOPOVER_TYPE (goPopover_get_type())
#define GOPOPOVER(obj) (G_TYPE_CHECK_INSTANCE_CAST((obj), GOPOPOVER_TYPE, goPopover))
#define IS_GOPOPOVER(obj) (G_TYPE_CHECK_INSTANCE_TYPE((obj), GOPOPOVER_TYPE))
#define GOPOPOVER_CLASS(class) (G_TYPE_CHECK_CLASS_CAST((class), GOPOPOVER_TYPE, goPopoverClass))
#define IS_GOPOPOVER_CLASS(class) (G_TYPE_CHECK_CLASS_TYPE((class), GOPOPOVER_TYPE))
#define GOPOPOVER_GET_CLASS(obj) (G_TYPE_INSTANCE_GET_CLASS((obj), GOPOPOVER_TYPE, goPopoverClass))

typedef struct goPopover goPopover;
typedef struct goPopoverClass goPopoverClass;

struct goPopover {
	GtkBin parent_instance;
	void *gocontainer;
	GdkWindow *gdkwin;
};

struct goPopoverClass {
	GtkBinClass parent_class;
};

G_DEFINE_TYPE(goPopover, goPopover, GTK_TYPE_BIN)

static void goPopover_init(goPopover *p)
{
	gtk_widget_set_has_window(GTK_WIDGET(p), TRUE);
}

static void goPopover_dispose(GObject *obj)
{
	G_OBJECT_CLASS(goPopover_parent_class)->dispose(obj);
}

static void goPopover_finalize(GObject *obj)
{
	G_OBJECT_CLASS(goPopover_parent_class)->finalize(obj);
}

static void goPopover_realize(GtkWidget *widget)
{
	GdkWindowAttr attr;
	goPopover *p = GOPOPOVER(widget);

	attr.x = 0;
	attr.y = 0;
	attr.width = 200;
	attr.height = 200;
	attr.wclass = GDK_INPUT_OUTPUT;
	attr.event_mask = gtk_widget_get_events(GTK_WIDGET(p)) | GDK_POINTER_MOTION_MASK | GDK_BUTTON_MOTION_MASK | GDK_BUTTON_PRESS_MASK | GDK_BUTTON_RELEASE_MASK | GDK_EXPOSURE_MASK | GDK_ENTER_NOTIFY_MASK | GDK_LEAVE_NOTIFY_MASK;
	attr.visual = gtk_widget_get_visual(GTK_WIDGET(p));
	attr.window_type = GDK_WINDOW_CHILD;		// GtkPopover does this; TODO what does GtkWindow(GTK_WINDOW_POPUP) do?
	p->gdkwin = gdk_window_new(gtk_widget_get_parent_window(GTK_WIDGET(p)),
		&attr, GDK_WA_VISUAL);
	gtk_widget_set_window(GTK_WIDGET(p), p->gdkwin);
	gtk_widget_register_window(GTK_WIDGET(p), p->gdkwin);
	gtk_widget_set_realized(GTK_WIDGET(p), TRUE);
}

static void goPopover_map(GtkWidget *widget)
{
	gdk_window_show(GOPOPOVER(widget)->gdkwin);
	GTK_WIDGET_CLASS(goPopover_parent_class)->map(widget);
}

static void goPopover_unmap(GtkWidget *widget)
{
	gdk_window_hide(GOPOPOVER(widget)->gdkwin);
	GTK_WIDGET_CLASS(goPopover_parent_class)->unmap(widget);
}

static gboolean goPopover_draw(GtkWidget *widget, cairo_t *cr)
{
	GtkStyleContext *context;

	context = gtk_widget_get_style_context(widget);
	gtk_render_background(context, cr, 0, 0, 200, 200);
	return TRUE;
}

static void goPopover_class_init(goPopoverClass *class)
{
	G_OBJECT_CLASS(class)->dispose = goPopover_dispose;
	G_OBJECT_CLASS(class)->finalize = goPopover_finalize;
	GTK_WIDGET_CLASS(class)->realize = goPopover_realize;
	GTK_WIDGET_CLASS(class)->map = goPopover_map;
	GTK_WIDGET_CLASS(class)->unmap = goPopover_unmap;
	GTK_WIDGET_CLASS(class)->draw = goPopover_draw;
}

void buttonClicked(GtkWidget *button, gpointer data)
{
	GtkWidget *popover;

	popover = g_object_new(GOPOPOVER_TYPE, NULL);
	gtk_widget_set_parent(popover, gtk_widget_get_parent(button));
	gtk_widget_show(popover);
}

int main(void)
{
	GtkWidget *mainwin;
	GtkWidget *button;

	gtk_init(NULL, NULL);
	mainwin = gtk_window_new(GTK_WINDOW_TOPLEVEL);
	gtk_window_resize(GTK_WINDOW(mainwin), 150, 50);
	g_signal_connect(mainwin, "destroy", gtk_main_quit, NULL);
	button = gtk_button_new_with_label("Click Me");
	g_signal_connect(button, "clicked", G_CALLBACK(buttonClicked), NULL);
	gtk_container_add(GTK_CONTAINER(mainwin), button);
	gtk_widget_show_all(mainwin);
	gtk_main();
	return 0;
}

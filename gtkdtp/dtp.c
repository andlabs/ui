// 9 january 2015
#include "dtp.h"

#define GDTP(x) ((goDateTimePicker *) x)
#define PRIV(x) (GDTP(x)->priv)
#define GDTPC(x) ((goDateTimePickerClass *) x)

struct goDateTimePickerPrivate {
	GtkWidget *openbutton;

	GtkWidget *popup;
	GtkWidget *calendar;
	GtkWidget *spinHours;
	GtkWidget *spinMinutes;
	GtkWidget *spinSeconds;
};

G_DEFINE_TYPE_WITH_CODE(goDateTimePicker, goDateTimePicker, GTK_TYPE_BOX,
	G_ADD_PRIVATE(goDateTimePicker))

enum {
	gtkMargin  = 12,
	gtkXPadding = 12,
	gtkYPadding = 6,
};

static void goDateTimePicker_init(goDateTimePicker *dtp)
{
	goDateTimePickerPrivate *d;
	GtkWidget *arrow;
	GtkWidget *vbox;
	GtkWidget *hbox;

	dtp->priv = goDateTimePicker_get_instance_private(dtp);
	d = dtp->priv;

	// create the actual bar elements
	// TODO the entry field
	// just make a dummy one for testing
	hbox = gtk_entry_new();
	gtk_style_context_add_class(gtk_widget_get_style_context(hbox), GTK_STYLE_CLASS_COMBOBOX_ENTRY);
	gtk_widget_set_hexpand(hbox, TRUE);
	gtk_widget_set_halign(hbox, GTK_ALIGN_FILL);
	gtk_container_add(GTK_CONTAINER(dtp), hbox);
	// the open button
	d->openbutton = gtk_toggle_button_new();
	arrow = gtk_arrow_new(GTK_ARROW_DOWN, GTK_SHADOW_NONE);
	gtk_container_add(GTK_CONTAINER(d->openbutton), arrow);
	// and make them look linked
	// TODO sufficient?
	gtk_style_context_add_class(gtk_widget_get_style_context(GTK_WIDGET(dtp)), "linked");
	// and mark them as visible
	gtk_widget_show_all(d->openbutton);
	// and add them to the bar
	gtk_container_add(GTK_CONTAINER(dtp), d->openbutton);

	// now create the popup that will hold everything
	d->popup = gtk_window_new(GTK_WINDOW_POPUP);
	gtk_window_set_type_hint(GTK_WINDOW(d->popup), GDK_WINDOW_TYPE_HINT_COMBO);
	gtk_window_set_resizable(GTK_WINDOW(d->popup), FALSE);
	vbox = gtk_box_new(GTK_ORIENTATION_VERTICAL, gtkYPadding);
	gtk_container_set_border_width(GTK_CONTAINER(vbox), gtkMargin);
	d->calendar = gtk_calendar_new();
	gtk_container_add(GTK_CONTAINER(vbox), d->calendar);
	gtk_container_add(GTK_CONTAINER(d->popup), vbox);
}

static void goDateTimePicker_dispose(GObject *obj)
{
	goDateTimePickerPrivate *d = PRIV(obj);

	// TODO really with g_clear_object()?
	g_clear_object(&(d->openbutton));
	g_clear_object(&(d->popup));
	// TODO g_object_clear() the children?
	G_OBJECT_CLASS(goDateTimePicker_parent_class)->dispose(obj);
}

static void goDateTimePicker_finalize(GObject *obj)
{
	G_OBJECT_CLASS(goDateTimePicker_parent_class)->finalize(obj);
}

static void goDateTimePicker_class_init(goDateTimePickerClass *class)
{
	G_OBJECT_CLASS(class)->dispose = goDateTimePicker_dispose;
	G_OBJECT_CLASS(class)->finalize = goDateTimePicker_finalize;
}

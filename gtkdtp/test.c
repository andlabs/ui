// 9 january 2015
#include "dtp.h"

// #qo pkg-config: gtk+-3.0

int main(void)
{
	GtkWidget *mainwin;

	gtk_init(NULL, NULL);
	mainwin = gtk_window_new(GTK_WINDOW_TOPLEVEL);
	gtk_container_add(GTK_CONTAINER(mainwin), g_object_new(goDateTimePicker_get_type(), NULL));
	gtk_widget_show_all(mainwin);
	gtk_main();
	return 0;
}

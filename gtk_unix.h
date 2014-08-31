// 16 march 2014

#ifndef __GO_UI_GTK_UNIX_H__
#define __GO_UI_GTK_UNIX_H__

/*
MIN_REQUIRED signals that programs are expected to run on at least GLib 2.32/GTK+ 3.4 and thus deprectation warnings for newer versions (such as gtk_scrolled_window_add_with_viewport() being deprecated in GTK+ 3.8) should be suppressed.
MAX_ALLOWED signals that programs will not use features introduced in newer versions of GLib/GTK+ and that the compiler should warn us if we slip.
Thanks to desrt in irc.gimp.net/#gtk+
*/

// GLib/GObject
#define GLIB_VERSION_MIN_REQUIRED GLIB_VERSION_2_32
#define GLIB_VERSION_MAX_ALLOWED GLIB_VERSION_2_32

// GDK/GTK+
#define GDK_VERSION_MIN_REQUIRED GDK_VERSION_3_4
#define GDK_VERSION_MAX_ALLOWED GDK_VERSION_3_4

// cairo has no such macros (thanks Company in irc.gimp.net/#gtk+)

#include <stdlib.h>
#include <gtk/gtk.h>

// table_unix.c
extern void tableAppendColumn(GtkTreeView *, gint, gchar *, GtkCellRenderer *, gchar *);
typedef struct goTableModel goTableModel;
typedef struct goTableModelClass goTableModelClass;
struct goTableModel {
	GObject parent_instance;
	void *gotable;
};
struct goTableModelClass {
	GObjectClass parent_class;
};
extern goTableModel *newTableModel(void *);
extern void tableUpdate(goTableModel *, gint, gint);

// container_unix.c
extern GtkWidget *newContainer(void *);

#endif

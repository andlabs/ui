/* 16 march 2014 */

/* This header file is a convenience to ensure that the compatibility flags below are included in all Go files that include <gtk/gtk.h> */

/*
MIN_REQUIRED signals that programs are expected to run on at least GTK+ 3 .4 and thus deprectation warnings for newer versions (such as gtk_scrolled_window_add_with_viewport() being deprecated in GTK+ 3.8) should be suppressed.
MAX_ALLOWED signals that programs will not use features introduced in newer versions of GTK+ and that the compiler should warn us if we slip.
Thanks to desrt in irc.gimp.net/#gtk+
*/
#define GDK_VERSION_MIN_REQUIRED GDK_VERSION_3_4
#define GDK_VERSION_MAX_ALLOWED GDK_VERSION_3_4

/* TODO are there equivalent compatibility macros for the other components of GTK+? */

#include <gtk/gtk.h>
/* TODO include <stdlib.h> too */

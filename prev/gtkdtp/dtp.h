// 9 january 2015
#define GLIB_VERSION_MIN_REQUIRED GLIB_VERSION_2_32
#define GLIB_VERSION_MAX_ALLOWED GLIB_VERSION_2_32
#define GDK_VERSION_MIN_REQUIRED GDK_VERSION_3_4
#define GDK_VERSION_MAX_ALLOWED GDK_VERSION_3_4
#include <gtk/gtk.h>

typedef struct goDateTimePicker goDateTimePicker;
typedef struct goDateTimePickerClass goDateTimePickerClass;
typedef struct goDateTimePickerPrivate goDateTimePickerPrivate;

struct goDateTimePicker {
	GtkBox parent_instance;
	goDateTimePickerPrivate *priv;
};

struct goDateTimePickerClass {
	GtkBoxClass parent_class;
};

GType goDateTimePicker_get_type(void);

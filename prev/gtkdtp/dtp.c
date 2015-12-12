// 9 january 2015
#include "dtp.h"

/* notes:
https://git.gnome.org/browse/gtk+/tree/gtk/gtkcombobox.c?h=gtk-3-4
*/

#define GDTP(x) ((goDateTimePicker *) x)
#define PRIV(x) (GDTP(x)->priv)
#define GDTPC(x) ((goDateTimePickerClass *) x)

struct goDateTimePickerPrivate {
	gint year;
	gint month;
	gint day;
};

G_DEFINE_TYPE_WITH_CODE(goDateTimePicker, goDateTimePicker, GTK_TYPE_BOX,
	;)

// TODO figure out how to share these between C and Go
enum {
	gtkMargin  = 12,
	gtkXPadding = 12,
	gtkYPadding = 6,
};

static void goDateTimePicker_init(goDateTimePicker *dtp)
{
	dtp->priv = G_TYPE_INSTANCE_GET_PRIVATE(dtp, goDateTimePicker_get_type(), goDateTimePickerPrivate);
}

static void goDateTimePicker_dispose(GObject *obj)
{
	goDateTimePickerPrivate *d = PRIV(obj);

	// TODO really with g_clear_object()?
	G_OBJECT_CLASS(goDateTimePicker_parent_class)->dispose(obj);
}

static void goDateTimePicker_finalize(GObject *obj)
{
	G_OBJECT_CLASS(goDateTimePicker_parent_class)->finalize(obj);
}

enum {
	pYear = 1,
	pMonth,
	pDay,
	nParams,
};

static GParamSpec *gdtpParams[] = {
	NULL,		// always null
	NULL,		// year
	NULL,		// month
	NULL,		// day
};

static void goDateTimePicker_set_property(GObject *obj, guint prop, const GValue *value, GParamSpec *spec)
{
	goDateTimePickerPrivate *d = PRIV(obj);

	switch (prop) {
	case pYear:
		d->year = g_value_get_int(value);
		break;
	case pMonth:
		d->month = g_value_get_int(value);
		break;
	case pDay:
		d->day = g_value_get_int(value);
		// see note on GtkCalendar comaptibility below
		if (d->day == 0)
			;	// TODO
		break;
	default:
		G_OBJECT_WARN_INVALID_PROPERTY_ID(obj, prop, spec);
		return;
	}
	// TODO refresh everything here
}

static void goDateTimePicker_get_property(GObject *obj, guint prop, GValue *value, GParamSpec *spec)
{
	goDateTimePickerPrivate *d = PRIV(obj);

	switch (prop) {
	case pYear:
		g_value_set_int(value, d->year);
		break;
	case pMonth:
		g_value_set_int(value, d->month);
		break;
	case pDay:
		g_value_set_int(value, d->day);
		break;
	default:
		G_OBJECT_WARN_INVALID_PROPERTY_ID(obj, prop, spec);
		return;
	}
}

static void goDateTimePicker_class_init(goDateTimePickerClass *class)
{
	g_type_class_add_private(class, sizeof (goDateTimePickerPrivate));

	G_OBJECT_CLASS(class)->dispose = goDateTimePicker_dispose;
	G_OBJECT_CLASS(class)->finalize = goDateTimePicker_finalize;
	G_OBJECT_CLASS(class)->set_property = goDateTimePicker_set_property;
	G_OBJECT_CLASS(class)->get_property = goDateTimePicker_get_property;

	// types and values are to be compatible with the 3.4 GtkCalendar parameters
	gdtpParams[pYear] = g_param_spec_int("year",
		"current year",
		"Current year",
		0,
		G_MAXINT >> 9,
		0,
		G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS);
	gdtpParams[pMonth] = g_param_spec_uint("month",
		"current month",
		"Current month",
		0,
		11,
		0,
		G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS);
	// because of the requirement to be compatible with GtkCalendar, we have to follow its rules about dates
	// values are 1..31 with 0 meaning no date selected
	// we will not allow no date to be selected, so we will set the default to 1 instead of 0
	// TODO is this an issue for binding?
	gdtpParams[pDay] = g_param_spec_uint("day",
		"current day",
		"Current day",
		0,
		31,
		1,
		G_PARAM_READWRITE | G_PARAM_STATIC_STRINGS);
	g_object_class_install_properties(G_OBJECT_CLASS(class), nParams, gdtpParams);
}

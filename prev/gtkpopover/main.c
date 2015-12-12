// 10 october 2014
// #qo pkg-config: gtk+-3.0
#define GLIB_VERSION_MIN_REQUIRED GLIB_VERSION_2_32
#define GLIB_VERSION_MAX_ALLOWED GLIB_VERSION_2_32
#define GDK_VERSION_MIN_REQUIRED GDK_VERSION_3_4
#define GDK_VERSION_MAX_ALLOWED GDK_VERSION_3_4
#include <gtk/gtk.h>

typedef gint LONG;

typedef struct POINT POINT;

struct POINT {
	LONG x;
	LONG y;
};

struct popover {
	void *gopopover;

	// a nice consequence of this design is that it allows four arrowheads to jut out at once; in practice only one will ever be used, but hey — simple implementation!
	LONG arrowLeft;
	LONG arrowRight;
	LONG arrowTop;
	LONG arrowBottom;
};

struct popover _p = { NULL, -1, -1, 20, -1 };
struct popover *p = &_p;

#define ARROWHEIGHT 8
#define ARROWWIDTH 8		/* should be the same for smooth lines */

void makePopoverPath(cairo_t *cr, LONG width, LONG height)
{
	POINT pt[20];
	int n;
	LONG xmax, ymax;

	cairo_new_path(cr);
	n = 0;

	// figure out the xmax and ymax of the box
	xmax = width;
	if (p->arrowRight >= 0)
		xmax -= ARROWWIDTH;
	ymax = height;
	if (p->arrowBottom >= 0)
		ymax -= ARROWHEIGHT;

	// the first point is either at (0,0), (0,arrowHeight), (arrowWidth,0), or (arrowWidth,arrowHeight)
	pt[n].x = 0;
	if (p->arrowLeft >= 0)
		pt[n].x = ARROWWIDTH;
	pt[n].y = 0;
	if (p->arrowTop >= 0)
		pt[n].y = ARROWHEIGHT;
	n++;

	// the left side
	pt[n].x = pt[n - 1].x;
	if (p->arrowLeft >= 0) {
		pt[n].y = pt[n - 1].y + p->arrowLeft;
		n++;
		pt[n].x = pt[n - 1].x - ARROWWIDTH;
		pt[n].y = pt[n - 1].y + ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x + ARROWWIDTH;
		pt[n].y = pt[n - 1].y + ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x;
	}
	pt[n].y = ymax;
	n++;

	// the bottom side
	pt[n].y = pt[n - 1].y;
	if (p->arrowBottom >= 0) {
		pt[n].x = pt[n - 1].x + p->arrowBottom;
		n++;
		pt[n].x = pt[n - 1].x + ARROWWIDTH;
		pt[n].y = pt[n - 1].y + ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x + ARROWWIDTH;
		pt[n].y = pt[n - 1].y - ARROWHEIGHT;
		n++;
		pt[n].y = pt[n - 1].y;
	}
	pt[n].x = xmax;
	n++;

	// the right side
	pt[n].x = pt[n - 1].x;
	if (p->arrowRight >= 0) {
		pt[n].y = pt[0].y + p->arrowRight + (ARROWHEIGHT * 2);
		n++;
		pt[n].x = pt[n - 1].x + ARROWWIDTH;
		pt[n].y = pt[n - 1].y - ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x - ARROWWIDTH;
		pt[n].y = pt[n - 1].y - ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x;
	}
	pt[n].y = pt[0].y;
	n++;

	// the top side
	pt[n].y = pt[n - 1].y;
	if (p->arrowTop >= 0) {
		pt[n].x = pt[0].x + p->arrowTop + (ARROWWIDTH * 2);
		n++;
		pt[n].x = pt[n - 1].x - ARROWWIDTH;
		pt[n].y = pt[n - 1].y - ARROWHEIGHT;
		n++;
		pt[n].x = pt[n - 1].x - ARROWWIDTH;
		pt[n].y = pt[n - 1].y + ARROWHEIGHT;
		n++;
		pt[n].y = pt[n - 1].y;
	}
	pt[n].x = pt[0].x;
	n++;

	int i;

	// TODO right and bottom edge
	cairo_set_line_width(cr, 1);
	cairo_move_to(cr, pt[0].x + 0.5, pt[0].y + 0.5);
	for (i = 1; i < n; i++)
		cairo_line_to(cr, pt[i].x + 0.5, pt[i].y + 0.5);
}

void drawPopoverFrame(GtkWidget *widget, cairo_t *cr, LONG width, LONG height, int forceAlpha)
{
	GtkStyleContext *context;
	GdkRGBA background, border;

	// TODO see what GtkPopover itself does
	// TODO drop shadow
	context = gtk_widget_get_style_context(widget);
	gtk_style_context_add_class(widget, GTK_STYLE_CLASS_BACKGROUND);
	gtk_style_context_get_background_color(context, GTK_STATE_FLAG_NORMAL, &background);
	gtk_style_context_get_border_color(context, GTK_STATE_FLAG_NORMAL, &border);
	if (forceAlpha) {
		background.alpha = 1;
		border.alpha = 1;
	}
	makePopoverPath(cr, width, height);
	cairo_set_source_rgba(cr, background.red, background.green, background.blue, background.alpha);
	cairo_fill_preserve(cr);
	cairo_set_source_rgba(cr, border.red, border.green, border.blue, border.alpha);
	cairo_stroke(cr);
}

gboolean popoverDraw(GtkWidget *widget, cairo_t *cr, gpointer data)
{
	gint width, height;

	width = gtk_widget_get_allocated_width(widget);
	height = gtk_widget_get_allocated_height(widget);
	drawPopoverFrame(widget, cr, width, height, 0);
	return FALSE;
}

void popoverSetSize(GtkWidget *widget, LONG width, LONG height)
{
	GdkWindow *gdkwin;
	cairo_t *cr;
	cairo_surface_t *cs;
	cairo_region_t *region;

	gtk_window_resize(GTK_WINDOW(widget), width, height);
	gdkwin = gtk_widget_get_window(widget);
	gdk_window_shape_combine_region(gdkwin, NULL, 0, 0);
	// TODO check errors
	cs = cairo_image_surface_create(CAIRO_FORMAT_ARGB32, width, height);
	cr = cairo_create(cs);
	drawPopoverFrame(widget, cr, width, height, 1);
	region = gdk_cairo_region_create_from_surface(cs);
	gdk_window_shape_combine_region(gdkwin, region, 0, 0);
	cairo_destroy(cr);
	cairo_surface_destroy(cs);
}

int main(void)
{
	GtkWidget *w;

	gtk_init(NULL, NULL);
	w = gtk_window_new(GTK_WINDOW_POPUP);
	gtk_window_set_decorated(GTK_WINDOW(w), FALSE);
	gtk_widget_set_app_paintable(w, TRUE);
	g_signal_connect(w, "draw", G_CALLBACK(popoverDraw), NULL);
	gtk_widget_set_has_window(w, TRUE);
	gtk_widget_realize(w);
	popoverSetSize(w, 200, 200);
	gtk_window_move(GTK_WINDOW(w), 50, 50);
	gtk_widget_show_all(w);
	gtk_main();
	return 0;
}

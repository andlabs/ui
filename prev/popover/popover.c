// 9 october 2014
#include <stdlib.h>
#include <stdint.h>
#include "popover.h"

#define ARROWHEIGHT 8
#define ARROWWIDTH 8		/* should be the same for smooth lines on Windows (TODO is there a better/nicer looking way?) */

struct popover {
	void *gopopover;

	// a nice consequence of this design is that it allows four arrowheads to jut out at once; in practice only one will ever be used, but hey — simple implementation!
	intptr_t arrowLeft;
	intptr_t arrowTop;
	intptr_t arrowRight;
	intptr_t arrowBottom;
};

popover *popoverDataNew(void *gopopover)
{
	popover *p;

	p = (popover *) malloc(sizeof (popover));
	if (p != NULL) {
		p->gopopover = gopopover;
		p->arrowLeft = -1;
		p->arrowTop = 20;//TODO-1;
		p->arrowRight = -1;
		p->arrowBottom = -1;
	}
	return p;
}

int popoverMakeFramePoints(popover *p, intptr_t width, intptr_t height, popoverPoint pt[20])
{
	int n;
	intptr_t xmax, ymax;

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

	return n;
}

void popoverWindowSizeToClientSize(popover *p, popoverRect *r)
{
	r->left++;
	r->top++;
	r->right--;
	r->bottom--;
	if (p->arrowLeft >= 0)
		r->left += ARROWWIDTH;
	if (p->arrowRight >= 0)
		r->right -= ARROWWIDTH;
	if (p->arrowTop >= 0)
		r->top += ARROWHEIGHT;
	if (p->arrowBottom >= 0)
		r->bottom -= ARROWHEIGHT;
}

// TODO window edge detection
popoverRect popoverPointAt(popover *p, popoverRect control, intptr_t width, intptr_t height, unsigned int side)
{
	intptr_t x, y;
	popoverRect out;

	// account for border
	width += 2;
	height += 2;
	p->arrowLeft = -1;
	p->arrowRight = -1;
	p->arrowTop = -1;
	p->arrowBottom = -1;
	// TODO right and bottom
	switch (side) {
	case popoverPointLeft:
		width += ARROWWIDTH;
		p->arrowLeft = height / 2 - ARROWHEIGHT;
		x = control.right;
		y = control.top - ((height - (control.bottom - control.top)) / 2);
		break;
	case popoverPointTop:
		height += ARROWHEIGHT;
		p->arrowTop = width / 2 - ARROWWIDTH;
		x = control.left - ((width - (control.right - control.left)) / 2);
		y = control.bottom;
		break;
	}
	out.left = x;
	out.top = y;
	out.right = x + width;
	out.bottom = y + height;
	return out;
}

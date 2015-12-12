// 11 october 2014

typedef struct popover popover;
typedef struct popoverPoint popoverPoint;
typedef struct popoverRect popoverRect;

struct popoverPoint {
	intptr_t x;
	intptr_t y;
};

struct popoverRect {
	intptr_t left;
	intptr_t top;
	intptr_t right;
	intptr_t bottom;
};

// note the order: flipping sides is as easy as side ^ 1
enum {
	popoverPointLeft,
	popoverPointRight,
	popoverPointTop,
	popoverPointBottom,
};

popover *popoverDataNew(void *);
int popoverMakeFramePoints(popover *, intptr_t, intptr_t, popoverPoint[20]);
void popoverWindowSizeToClientSize(popover *, popoverRect *);
popoverRect popoverPointAt(popover *, popoverRect, intptr_t, intptr_t, unsigned int);

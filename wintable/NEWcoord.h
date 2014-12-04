// 4 december 2014

typedef struct rowcol rowcol;

struct rowcol {
	intptr_t row;
	intptr_t column;
};

static rowcol clientCoordToRowColumn(struct table *t, POINT pt)
{
	// TODO
}

// same as client coordinates, but stored in a lParam (like the various mouse messages provide)
// TODO should we use GetMessagePos() instead?
static rowcol lParamToRowColumn(struct table *t, LPARAM lParam)
{
	POINT pt;

	pt.x = GET_X_LPARAM(lParam);
	pt.y = GET_Y_LPARAM(lParam);
	return clientCoordToRowColumn(t, pt);
}

// returns TRUE if the row is visible and thus has client coordinates; FALSE otherwise
static BOOL rowColumnToClientCoord(struct table *t, rowcol rc, struct POINT *pt)
{
	// TODO
}

// TODO idealCoordToRowColumn/rowColumnToIdealCoord?

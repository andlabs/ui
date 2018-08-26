// 26 august 2018
#include "pkgui.h"
#include "_cgo_export.h"

uiInitOptions *pkguiAllocInitOptions(void)
{
	return (uiInitOptions *) pkguiAlloc(sizeof (uiInitOptions));
}

void pkguiFreeInitOptions(uiInitOptions *o)
{
	free(o);
}

void pkguiQueueMain(uintptr_t n)
{
	uiQueueMain(pkguiDoQueueMain, (void *) n);
}

void pkguiOnShouldQuit(void)
{
	uiOnShouldQuit(pkguiDoOnShouldQuit, NULL);
}

void pkguiWindowOnClosing(uiWindow *w)
{
	uiWindowOnClosing(w, pkguiDoWindowOnClosing, NULL);
}

void pkguiButtonOnClicked(uiButton *b)
{
	uiButtonOnClicked(b, pkguiDoButtonOnClicked, NULL);
}

void pkguiCheckboxOnToggled(uiCheckbox *c)
{
	uiCheckboxOnToggled(c, pkguiDoCheckboxOnToggled, NULL);
}

void pkguiColorButtonOnChanged(uiColorButton *c)
{
	uiColorButtonOnChanged(c, pkguiDoColorButtonOnChanged, NULL);
}

pkguiColorDoubles pkguiAllocColorDoubles(void)
{
	pkguiColorDoubles c;

	c.r = (double *) pkguiAlloc(4 * sizeof (double));
	c.g = c.r + 1;
	c.b = c.g + 1;
	c.a = c.b + 1;
	return c;
}

void pkguiFreeColorDoubles(pkguiColorDoubles c)
{
	free(c.r);
}

void pkguiComboboxOnSelected(uiCombobox *c)
{
	uiComboboxOnSelected(c, pkguiDoComboboxOnSelected, NULL);
}

void pkguiDateTimePickerOnChanged(uiDateTimePicker *d)
{
	uiDateTimePickerOnChanged(d, pkguiDoDateTimePickerOnChanged, NULL);
}

struct tm *pkguiAllocTime(void)
{
	return (struct tm *) pkguiAlloc(sizeof (struct tm));
}

void pkguiFreeTime(struct tm *t)
{
	free(t);
}

void pkguiEditableComboboxOnChanged(uiEditableCombobox *c)
{
	uiEditableComboboxOnChanged(c, pkguiDoEditableComboboxOnChanged, NULL);
}

void pkguiEntryOnChanged(uiEntry *e)
{
	uiEntryOnChanged(e, pkguiDoEntryOnChanged, NULL);
}

void pkguiFontButtonOnChanged(uiFontButton *b)
{
	uiFontButtonOnChanged(b, pkguiDoFontButtonOnChanged, NULL);
}

void pkguiMultilineEntryOnChanged(uiMultilineEntry *e)
{
	uiMultilineEntryOnChanged(e, pkguiDoMultilineEntryOnChanged, NULL);
}

void pkguiRadioButtonsOnSelected(uiRadioButtons *r)
{
	uiRadioButtonsOnSelected(r, pkguiDoRadioButtonsOnSelected, NULL);
}

void pkguiSliderOnChanged(uiSlider *s)
{
	uiSliderOnChanged(s, pkguiDoSliderOnChanged, NULL);
}

void pkguiSpinboxOnChanged(uiSpinbox *s)
{
	uiSpinboxOnChanged(s, pkguiDoSpinboxOnChanged, NULL);
}

uiDrawBrush *pkguiAllocBrush(void)
{
	return (uiDrawBrush *) pkguiAlloc(sizeof (uiDrawBrush));
}

void pkguiFreeBrush(uiDrawBrush *b)
{
	free(b);
}

uiDrawBrushGradientStop *pkguiAllocGradientStops(size_t n)
{
	return (uiDrawBrushGradientStop *) pkguiAlloc(n * sizeof (uiDrawBrushGradientStop));
}

void pkguiFreeGradientStops(uiDrawBrushGradientStop *stops)
{
	free(stops);
}

void pkguiSetGradientStop(uiDrawBrushGradientStop *stops, size_t i, double pos, double r, double g, double b, double a)
{
	stops[i].Pos = pos;
	stops[i].R = r;
	stops[i].G = g;
	stops[i].B = b;
	stops[i].A = a;
}

uiDrawStrokeParams *pkguiAllocStrokeParams(void)
{
	return (uiDrawStrokeParams *) pkguiAlloc(sizeof (uiDrawStrokeParams));
}

void pkguiFreeStrokeParams(uiDrawStrokeParams *p)
{
	free(p);
}

double *pkguiAllocDashes(size_t n)
{
	return (double *) pkguiAlloc(n * sizeof (double));
}

void pkguiFreeDashes(double *dashes)
{
	free(dashes);
}

void pkguiSetDash(double *dashes, size_t i, double dash)
{
	dashes[i] = dash;
}

uiDrawMatrix *pkguiAllocMatrix(void)
{
	return (uiDrawMatrix *) pkguiAlloc(sizeof (uiDrawMatrix));
}

void pkguiFreeMatrix(uiDrawMatrix *m)
{
	free(m);
}

uiUnderlineColor *pkguiNewUnderlineColor(void)
{
	return (uiUnderlineColor *) pkguiAlloc(sizeof (uiUnderlineColor));
}

void pkguiFreeUnderlineColor(uiUnderlineColor *c)
{
	free(c);
}

uiFontDescriptor *pkguiNewFontDescriptor(void)
{
	return (uiFontDescriptor *) pkguiAlloc(sizeof (uiFontDescriptor));
}

void pkguiFreeFontDescriptor(uiFontDescriptor *fd)
{
	free(fd);
}

uiDrawTextLayoutParams *pkguiNewDrawTextLayoutParams(void)
{
	return (uiDrawTextLayoutParams *) pkguiAlloc(sizeof (uiDrawTextLayoutParams));
}

void pkguiFreeDrawTextLayoutParams(uiDrawTextLayoutParams *p)
{
	free(p);
}

uiAreaHandler *pkguiAllocAreaHandler(void)
{
	uiAreaHandler *ah;

	ah = (uiAreaHandler *) pkguiAlloc(sizeof (uiAreaHandler));
	ah->Draw = pkguiDoAreaHandlerDraw;
	ah->MouseEvent = pkguiDoAreaHandlerMouseEvent;
	ah->MouseCrossed = pkguiDoAreaHandlerMouseCrossed;
	ah->DragBroken = pkguiDoAreaHandlerDragBroken;
	ah->KeyEvent = pkguiDoAreaHandlerKeyEvent;
	return ah;
}

void pkguiFreeAreaHandler(uiAreaHandler *ah)
{
	free(ah);
}

// cgo can't generate const, so we need this trampoline
static void realDoTableModelSetCellValue(uiTableModelHandler *mh, uiTableModel *m, int row, int column, const uiTableValue *value)
{
	pkguiDoTableModelSetCellValue(mh, m, row, column, (uiTableValue *) value);
}

const uiTableModelHandler pkguiTableModelHandler = {
	.NumColumns = pkguiDoTableModelNumColumns,
	.ColumnType = pkguiDoTableModelColumnType,
	.NumRows = pkguiDoTableModelNumRows,
	.CellValue = pkguiDoTableModelCellValue,
	.SetCellValue = realDoTableModelSetCellValue,
};

uiTableTextColumnOptionalParams *pkguiAllocTableTextColumnOptionalParams(void)
{
	return (uiTableTextColumnOptionalParams *) pkguiAlloc(sizeof (uiTableTextColumnOptionalParams));
}

void pkguiFreeTableTextColumnOptionalParams(uiTableTextColumnOptionalParams *p)
{
	free(p);
}

uiTableParams *pkguiAllocTableParams(void)
{
	return (uiTableParams *) pkguiAlloc(sizeof (uiTableParams));
}

void pkguiFreeTableParams(uiTableParams *p)
{
	free(p);
}

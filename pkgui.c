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

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

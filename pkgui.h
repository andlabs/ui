// 12 august 2018
#ifndef pkguiHFileIncluded
#define pkguiHFileIncluded

#include <stdlib.h>
#include <string.h>
#include <time.h>
#include "ui.h"

// main.go
extern uiInitOptions *pkguiAllocInitOptions(void);
extern void pkguiFreeInitOptions(uiInitOptions *o);
extern void pkguiQueueMain(uintptr_t n);
extern void pkguiOnShouldQuit(void);

// window.go
extern void pkguiWindowOnClosing(uiWindow *w);

// button.go
extern void pkguiButtonOnClicked(uiButton *b);

// checkbox.go
extern void pkguiCheckboxOnToggled(uiCheckbox *c);

// colorbutton.go
extern void pkguiColorButtonOnChanged(uiColorButton *c);
typedef struct pkguiColorDoubles pkguiColorDoubles;
struct pkguiColorDoubles {
	double *r;
	double *g;
	double *b;
	double *a;
};
extern pkguiColorDoubles pkguiAllocColorDoubles(void);
extern void pkguiFreeColorDoubles(pkguiColorDoubles c);

// combobox.go
extern void pkguiComboboxOnSelected(uiCombobox *c);

// datetimepicker.go
extern void pkguiDateTimePickerOnChanged(uiDateTimePicker *d);
extern struct tm *pkguiAllocTime(void);
extern void pkguiFreeTime(struct tm *t);

// editablecombobox.go
extern void pkguiEditableComboboxOnChanged(uiEditableCombobox *c);

// entry.go
extern void pkguiEntryOnChanged(uiEntry *e);

// fontbutton.go
extern void pkguiFontButtonOnChanged(uiFontButton *b);

// multilineentry.go
extern void pkguiMultilineEntryOnChanged(uiMultilineEntry *e);

// radiobuttons.go
extern void pkguiRadioButtonsOnSelected(uiRadioButtons *r);

// slider.go
extern void pkguiSliderOnChanged(uiSlider *s);

// spinbox.go
extern void pkguiSpinboxOnChanged(uiSpinbox *s);

// draw.go
extern uiDrawBrush *pkguiAllocBrush(void);
extern void pkguiFreeBrush(uiDrawBrush *b);
extern uiDrawBrushGradientStop *pkguiAllocGradientStops(size_t n);
extern void pkguiFreeGradientStops(uiDrawBrushGradientStop *stops);
extern void pkguiSetGradientStop(uiDrawBrushGradientStop *stops, size_t i, double pos, double r, double g, double b, double a);
extern uiDrawStrokeParams *pkguiAllocStrokeParams(void);
extern void pkguiFreeStrokeParams(uiDrawStrokeParams *p);
extern double *pkguiAllocDashes(size_t n);
extern void pkguiFreeDashes(double *dashes);
extern void pkguiSetDash(double *dashes, size_t i, double dash);
extern uiDrawMatrix *pkguiAllocMatrix(void);
extern void pkguiFreeMatrix(uiDrawMatrix *m);

// drawtext.go
extern uiUnderlineColor *pkguiNewUnderlineColor(void);
extern void pkguiFreeUnderlineColor(uiUnderlineColor *c);
extern uiFontDescriptor *pkguiNewFontDescriptor(void);
extern void pkguiFreeFontDescriptor(uiFontDescriptor *fd);
extern uiDrawTextLayoutParams *pkguiNewDrawTextLayoutParams(void);
extern void pkguiFreeDrawTextLayoutParams(uiDrawTextLayoutParams *p);

// areahandler.go
extern uiAreaHandler *pkguiAllocAreaHandler(void);
extern void pkguiFreeAreaHandler(uiAreaHandler *ah);

// tablemodel.go
extern const uiTableModelHandler pkguiTableModelHandler;

// table.go
extern uiTableTextColumnOptionalParams *pkguiAllocTableTextColumnOptionalParams(void);
extern void pkguiFreeTableTextColumnOptionalParams(uiTableTextColumnOptionalParams *p);
extern uiTableParams *pkguiAllocTableParams(void);
extern void pkguiFreeTableParams(uiTableParams *p);

#endif

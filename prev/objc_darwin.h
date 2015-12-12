/* 8 july 2014 */

/* cgo will include this file multiple times */
#ifndef __GO_UI_OBJC_DARWIN_H__
#define __GO_UI_OBJC_DARWIN_H__

#define MAC_OS_X_VERSION_MIN_REQUIRED MAC_OS_X_VERSION_10_7
#define MAC_OS_X_VERSION_MAX_ALLOWED MAC_OS_X_VERSION_10_7

#include <stdlib.h>
#include <stdint.h>
#include <objc/message.h>
#include <objc/objc.h>
#include <objc/runtime.h>

/* Objective-C -> Go types for max safety */
struct xsize {
	intptr_t width;
	intptr_t height;
};

struct xrect {
	intptr_t x;
	intptr_t y;
	intptr_t width;
	intptr_t height;
};

struct xalignment {
	struct xrect rect;
	intptr_t baseline;
};

struct xpoint {
	intptr_t x;
	intptr_t y;
};

/* uitask_darwin.m */
extern id getAppDelegate(void);	/* used by the other .m files */
extern void uiinit(char **);
extern void uimsgloop(void);
extern void uistop(void);
extern void beginModal(void);
extern void endModal(void);
extern void issue(void *);

/* window_darwin.m */
extern id newWindow(intptr_t, intptr_t);
extern void windowSetDelegate(id, void *);
extern void windowSetContentView(id, id);
extern const char *windowTitle(id);
extern void windowSetTitle(id, const char *);
extern void windowShow(id);
extern void windowHide(id);
extern void windowClose(id);
extern id windowContentView(id);
extern void windowRedraw(id);

/* basicctrls_darwin.m */
#define textfieldWidth (96)		/* according to Interface Builder */
extern id newButton(void);
extern void buttonSetDelegate(id, void *);
extern const char *buttonText(id);
extern void buttonSetText(id, char *);
extern id newCheckbox(void);
extern void checkboxSetDelegate(id, void *);
extern BOOL checkboxChecked(id);
extern void checkboxSetChecked(id, BOOL);
extern id finishNewTextField(id, BOOL);
extern id newTextField(void);
extern id newPasswordField(void);
extern void textfieldSetDelegate(id, void *);
extern const char *textfieldText(id);
extern void textfieldSetText(id, char *);
extern id textfieldOpenInvalidPopover(id, char *);
extern void textfieldCloseInvalidPopover(id);
extern BOOL textfieldEditable(id);
extern void textfieldSetEditable(id, BOOL);
extern id newLabel(void);
extern id newGroup(id);
extern const char *groupText(id);
extern void groupSetText(id, char *);
extern id newTextbox(void);
extern char *textboxText(id);
extern void textboxSetText(id, char *);
extern id newProgressBar(void);
extern intmax_t progressbarPercent(id);
extern void progressbarSetPercent(id, intmax_t);

/* container_darwin.m */
extern id newContainerView(void *);
extern void moveControl(id, intptr_t, intptr_t, intptr_t, intptr_t);
extern struct xrect containerBounds(id);

/* tab_darwin.m */
extern id newTab(void);
extern void tabAppend(id, char *, id);
extern struct xsize tabPreferredSize(id);

/* table_darwin.m */
enum {
	colTypeText,
	colTypeImage,
	colTypeCheckbox,
};
extern id newTable(void);
extern void tableAppendColumn(id, intptr_t, char *, int, BOOL);
extern void tableUpdate(id);
extern void tableMakeDataSource(id, void *);
extern struct xsize tablePreferredSize(id);
extern intptr_t tableSelected(id);
extern void tableSelect(id, intptr_t);

/* control_darwin.m */
extern void parent(id, id);
extern void controlSetHidden(id, BOOL);
extern void setStandardControlFont(id);
extern void setSmallControlFont(id);
extern struct xsize controlPreferredSize(id);
extern id newScrollView(id, BOOL);
extern struct xalignment alignmentInfo(id, struct xrect);
extern struct xalignment alignmentInfoFrame(id);

/* area_darwin.h */
extern Class getAreaClass(void);
extern id newArea(void *);
extern BOOL drawImage(void *, intptr_t, intptr_t, intptr_t, intptr_t, intptr_t);
extern const uintptr_t cNSShiftKeyMask;
extern const uintptr_t cNSControlKeyMask;
extern const uintptr_t cNSAlternateKeyMask;
extern const uintptr_t cNSCommandKeyMask;
extern uintptr_t modifierFlags(id);
extern struct xpoint getTranslatedEventPoint(id, id);
extern intptr_t buttonNumber(id);
extern intptr_t clickCount(id);
extern uintptr_t pressedMouseButtons(void);
extern uintptr_t keyCode(id);
extern void areaRepaint(id, struct xrect);
extern void areaRepaintAll(id);
extern void areaTextFieldOpen(id, id, intptr_t, intptr_t);
extern void areaSetTextField(id, id);
extern void areaEndTextFieldEditing(id, id);


/* common_darwin.m */
extern void disableAutocorrect(id);

/* image_darwin.m */
extern id toTableImage(void *, intptr_t, intptr_t, intptr_t);

/* dialog_darwin.m */
extern void openFile(id, void *);

/* warningpopover_darwin.m */
extern id newWarningPopover(char *);
extern void warningPopoverShow(id, id);

/* spinbox_darwin.m */
extern id newSpinbox(void *, intmax_t, intmax_t);
extern id spinboxTextField(id);
extern id spinboxStepper(id);
extern intmax_t spinboxValue(id);
extern void spinboxSetValue(id, intmax_t);

#endif

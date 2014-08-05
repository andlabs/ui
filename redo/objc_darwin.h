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

/* uitask_darwin.m */
extern id getAppDelegate(void);	/* used by the other .m files */
extern BOOL uiinit(void);
extern void uimsgloop(void);
extern void uistop(void);
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
extern void parent(id, id);
extern void controlSetHidden(id, BOOL);
extern void setStandardControlFont(id);
extern void setSmallControlFont(id);
extern id newButton(void);
extern void buttonSetDelegate(id, void *);
extern const char *buttonText(id);
extern void buttonSetText(id, char *);
extern id newCheckbox(void);
extern void checkboxSetDelegate(id, void *);
extern BOOL checkboxChecked(id);
extern void checkboxSetChecked(id, BOOL);
extern id newTextField(void);
extern id newPasswordField(void);
extern const char *textFieldText(id);
extern void textFieldSetText(id, char *);
extern id newLabel(void);

/* container_darwin.m */
extern id newContainerView(void *);
extern void moveControl(id, intptr_t, intptr_t, intptr_t, intptr_t);

/* tab_darwin.m */
extern id newTab(void);
extern void tabAppend(id, char *, id);

/* table_darwin.m */
extern id newTable(void);
extern void tableAppendColumn(id, char *);
extern void tableUpdate(id);
extern void tableMakeDataSource(id, void *);

/* control_darwin.m */
extern id newScrollView(id);

/* xsizing_darwin.m */
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
extern struct xsize controlPrefSize(id);
extern struct xsize tabPrefSize(id);
extern struct xsize areaPrefSize(id);
extern struct xalignment alignmentInfo(id, struct xrect);
extern struct xrect frame(id);

/* area_darwin.h */
struct xpoint {
	intptr_t x;
	intptr_t y;
};
extern id newArea(void *);
extern void drawImage(void *, intptr_t, intptr_t, intptr_t, intptr_t, intptr_t);
extern uintptr_t modifierFlags(id);
extern struct xpoint getTranslatedEventPoint(id, id);
extern intptr_t buttonNumber(id);
extern intptr_t clickCount(id);
extern uintptr_t pressedMouseButtons(void);
extern uintptr_t keyCode(id);
extern void display(id);

#endif

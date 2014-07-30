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
extern id newButton(void);
extern void buttonSetDelegate(id, void *);
extern const char *buttonText(id);
extern void buttonSetText(id, char *);
extern id newCheckbox(void);
extern BOOL checkboxChecked(id);
extern void checkboxSetChecked(id, BOOL);
extern id newTextField(void);
extern id newPasswordField(void);
extern const char *textFieldText(id);
extern void textFieldSetText(id, char *);
extern id newLabel(void);

/* sizing_darwin.m */
extern void moveControl(id, intptr_t, intptr_t, intptr_t, intptr_t);

/* containerctrls_darwin.m */
extern id newTab(void *);
extern id tabAppend(id, char *);

/* table_darwin.m */
extern id newTable(void);
extern void tableAppendColumn(id, char *);
extern void tableUpdate(id);
extern void tableMakeDataSource(id, void *);

/* control_darwin.m */
extern id newScrollView(id);

#endif

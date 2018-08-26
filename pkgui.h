// 12 august 2018
#include <stdlib.h>
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

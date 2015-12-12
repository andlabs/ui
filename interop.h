// 11 december 2015

#ifndef __UI_INTEROP_H__
#define __UI_INTEROP_H__

#include <stdint.h>

extern char *interopInit(void);
extern void interopFreeStr(char *);
extern void interopRun(void);
extern void interopQuit(void);
extern void interopQueueMain(uintptr_t);

#endif

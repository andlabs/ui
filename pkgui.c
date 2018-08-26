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

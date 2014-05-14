// 13 may 2014

extern id toListboxItem(id, id);
extern id fromListboxItem(id, id);
extern id newListboxArray(void);
extern void listboxArrayAppend(id, id);
extern void listboxArrayInsertBefore(id, id, uintptr_t);
extern void listboxArrayDelete(id, uintptr_t);
extern id listboxArrayItemAt(id, uintptr_t);
extern void bindListboxArray(id, id, id, id);
extern id boundListboxArray(id, id);
extern id makeListboxTableColumn(id);
extern id listboxTableColumn(id, id);
extern id makeListbox(id, BOOL);
extern id listboxSelectedRowIndexes(id);
extern uintptr_t listboxIndexesCount(id);
extern uintptr_t listboxIndexesFirst(id);
extern uintptr_t listboxIndexesNext(id, uintptr_t);
extern intptr_t listboxLen(id);
extern void listboxDeselectAll(id);

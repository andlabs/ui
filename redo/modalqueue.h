/* 19 august 2014 */

extern void beginModal(void);
extern void endModal(void);
extern int queueIfModal(void *);

/* needed by the above */
extern void doissue(void *);
extern void modalPanic(char *, char *);

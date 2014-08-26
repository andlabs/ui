// 19 august 2014

#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include "modalqueue.h"

static struct {
	int inmodal;
	void **queue;
	size_t len;
	size_t cap;
} mq = { 0, NULL, 0, 0 };

void beginModal(void)
{
	mq.inmodal = 1;
	if (mq.queue == NULL) {
		mq.cap = 128;
		mq.queue = (void **) malloc(mq.cap * sizeof (void *));
		if (mq.queue == NULL)
			modalPanic("error allocating modal queue", strerror(errno));
		mq.len = 0;
	}
}

void endModal(void)
{
	size_t i;

	mq.inmodal = 0;
	for (i = 0; i < mq.len; i++)
		doissue(mq.queue[i]);
	mq.len = 0;
}

int queueIfModal(void *what)
{
	if (!mq.inmodal)
		return 0;
	mq.queue[mq.len] = what;
	mq.len++;
	if (mq.len >= mq.cap) {
		mq.cap *= 2;
		mq.queue = (void **) realloc(mq.queue, mq.cap * sizeof (void *));
		if (mq.queue == NULL)
			modalPanic("error growing modal queue", strerror(errno));
	}
	return 1;
}

#include <phidget21.h>
#include <pthread.h>

typedef struct handlerNode {
	struct handlerNode * next;
	void * result;
} handlerNode;

typedef struct handler {
	pthread_mutex_t mutex;
	pthread_cond_t cond;
	handlerNode * head;
	handlerNode * tail;
} handler;

handler * newHandler();

void handlerAppendResult(handler *h, void *r);
void * handlerAwait(handler *h);
void handlerFree(handler *h);

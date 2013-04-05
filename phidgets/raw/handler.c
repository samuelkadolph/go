#include <pthread.h>
#include <stdlib.h>
#include "handler.h"

handler * newHandler() {
	handler * h = calloc(1, sizeof(handler));

	pthread_mutex_init(&h->mutex, NULL);
	pthread_cond_init(&h->cond, NULL);
	h->head = NULL;
	h->tail = NULL;

	return h;
}

void handlerAppendResult(handler *h, void *r) {
	pthread_mutex_lock(&h->mutex);

	handlerNode * n = calloc(1, sizeof(handlerNode));
	n->result = r;

	if (h->head == NULL) {
		h->head = h->tail = n;
	} else {
		h->tail->next = n;
		h->tail = n;
	}

	pthread_cond_broadcast(&h->cond);
	pthread_mutex_unlock(&h->mutex);
}

void * handlerAwait(handler *h) {
	pthread_mutex_lock(&h->mutex);

	while (h->head == NULL) {
		pthread_cond_wait(&h->cond, &h->mutex);
	}

	handlerNode * n = h->head;
	h->head = n->next;

	pthread_mutex_unlock(&h->mutex);

	void * r = n->result;
	free(n);

	return r;
}

void handlerFree(handler *h) {
	handlerNode * n = NULL;
	for (n = h->head; n != NULL; n = n->next) {
	  free(n->result);
	  free(n);
	}

	pthread_cond_destroy(&h->cond);
	pthread_mutex_destroy(&h->mutex);

	free(h);
}

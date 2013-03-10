#include <phidget21.h>
#include <stdlib.h>
#include "handler.h"

typedef enum {
  phidgetAttach = 1,
  phidgetConnect,
  phidgetDetach,
  phidgetDisconnect
} eventType;

typedef struct onErrorResult {
  int code;
  const char * string;
} onErrorResult;

void onErrorResultFree(onErrorResult *r);

onErrorResult * onErrorAwait(handler *h);
void onEventAwait(handler *h);

int setOnErrorHandler(CPhidgetHandle p, handler *h);
int setOnEventHandler(CPhidgetHandle p, handler *h, eventType t);
void unsetOnErrorHandler(CPhidgetHandle p);
void unsetOnEventHandler(CPhidgetHandle p, eventType t);

#include "phidget.h"

void onErrorResultFree(onErrorResult * r) {
  free(r);
}

onErrorResult * onErrorAwait(handler *h) {
  return (onErrorResult *)handlerAwait(h);
}

void onEventAwait(handler *h) {
  handlerAwait(h);
}

int onErrorHandler(CPhidgetHandle p, void *ptr, int code, const char *string) {
  handler * h = (handler *)ptr;

  onErrorResult * r = calloc(1, sizeof(onErrorResult));
  r->code = code;
  r->string = string;
  handlerAppendResult(h, r);

  return 0;
}

int onEventHandler(CPhidgetHandle p, void *ptr) {
  handler * h = (handler *)ptr;
  handlerAppendResult(h, NULL);
  return 0;
}

int setOnErrorHandler(CPhidgetHandle p, handler *h) {
  return CPhidget_set_OnError_Handler(p, &onErrorHandler, h);
}

int setOnEventHandler(CPhidgetHandle p, handler *h, eventType t) {
  switch (t) {
    case phidgetAttach:
      return CPhidget_set_OnAttach_Handler(p, &onEventHandler, h);
    case phidgetConnect:
      return CPhidget_set_OnServerConnect_Handler(p, &onEventHandler, h);
    case phidgetDetach:
      return CPhidget_set_OnDetach_Handler(p, &onEventHandler, h);
    case phidgetDisconnect:
      return CPhidget_set_OnServerDisconnect_Handler(p, &onEventHandler, h);
  }
}

void unsetOnErrorHandler(CPhidgetHandle p) {
  CPhidget_set_OnError_Handler(p, NULL, NULL);
}

void unsetOnEventHandler(CPhidgetHandle p, eventType t) {
  switch (t) {
    case phidgetAttach:
      CPhidget_set_OnAttach_Handler(p, NULL, NULL);
    case phidgetConnect:
      CPhidget_set_OnServerConnect_Handler(p, NULL, NULL);
    case phidgetDetach:
      CPhidget_set_OnDetach_Handler(p, NULL, NULL);
    case phidgetDisconnect:
      CPhidget_set_OnServerDisconnect_Handler(p, NULL, NULL);
  }
}

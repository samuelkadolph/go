#include <string.h>
#include "ir.h"

void onCodeResultFree(onCodeResult *r) {
  free(r->data);
  free(r);
}

void onLearnResultFree(onLearnResult *r) {
  free(r->data);
  free(r);
}

void onRawDataResultFree(onRawDataResult *r) {
  free(r->data);
  free(r);
}

onCodeResult * onCodeAwait(handler *h) {
  return (onCodeResult *)handlerAwait(h);
}

onLearnResult * onLearnAwait(handler *h) {
  return (onLearnResult *)handlerAwait(h);
}

onRawDataResult * onRawDataAwait(handler *h) {
  return (onRawDataResult *)handlerAwait(h);
}

int onCodeHandler(CPhidgetIRHandle ir, void *ptr, unsigned char *data, int dataLength, int bitCount, int repeat) {
  handler * h = (handler *)ptr;

  onCodeResult * r = calloc(1, sizeof(onCodeResult));
  r->data = malloc(dataLength * sizeof(char));
  r->dataLength = dataLength;
  r->bitCount = bitCount;
  r->repeat = repeat;
  memcpy(r->data, data, dataLength * sizeof(char));
  handlerAppendResult(h, r);

  return 0;
}

int onLearnHandler(CPhidgetIRHandle ir, void *ptr, unsigned char *data, int dataLength, CPhidgetIR_CodeInfoHandle codeInfo) {
  handler * h = (handler *)ptr;

  onLearnResult * r = calloc(1, sizeof(onLearnResult));
  r->data = malloc(dataLength * sizeof(char));
  r->dataLength = dataLength;
  r->codeInfo = codeInfo;
  memcpy(r->data, data, dataLength * sizeof(char));
  handlerAppendResult(h, r);

  return 0;
}

int onRawDataHandler(CPhidgetIRHandle ir, void *ptr, int *data, int dataLength) {
  handler * h = (handler *)ptr;

  onRawDataResult * r = calloc(1, sizeof(onRawDataResult));
  r->data = malloc(dataLength * sizeof(int));
  r->dataLength = dataLength;
  memcpy(r->data, data, dataLength * sizeof(int));
  handlerAppendResult(h, r);

  return 0;
}

int setOnCodeHandler(CPhidgetIRHandle ir, handler *h) {
  return CPhidgetIR_set_OnCode_Handler(ir, &onCodeHandler, h);
}

int setOnLearnHandler(CPhidgetIRHandle ir, handler *h) {
  return CPhidgetIR_set_OnLearn_Handler(ir, &onLearnHandler, h);
}

int setOnRawDataHandler(CPhidgetIRHandle ir, handler *h) {
  return CPhidgetIR_set_OnRawData_Handler(ir, &onRawDataHandler, h);
}

void unsetOnCodeHandler(CPhidgetIRHandle ir) {
  CPhidgetIR_set_OnCode_Handler(ir, NULL, NULL);
}

void unsetOnLearnHandler(CPhidgetIRHandle ir) {
  CPhidgetIR_set_OnLearn_Handler(ir, NULL, NULL);
}

void unsetOnRawDataHandler(CPhidgetIRHandle ir) {
  CPhidgetIR_set_OnRawData_Handler(ir, NULL, NULL);
}

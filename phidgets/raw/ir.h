#include <phidget21.h>
#include <stdlib.h>
#include "handler.h"

typedef struct onCodeResult {
  unsigned char *data;
  int dataLength;
  int bitCount;
  int repeat;
} onCodeResult;

typedef struct onLearnResult {
  unsigned char *data;
  int dataLength;
  CPhidgetIR_CodeInfoHandle codeInfo;
} onLearnResult;

typedef struct onRawDataResult {
  int *data;
  int dataLength;
} onRawDataResult;

void onCodeResultFree(onCodeResult *r);
void onLearnResultFree(onLearnResult *r);
void onRawDataResultFree(onRawDataResult *r);

onCodeResult * onCodeAwait(handler *h);
onLearnResult * onLearnAwait(handler *h);
onRawDataResult * onRawDataAwait(handler *h);

int setOnCodeHandler(CPhidgetIRHandle ir, handler *h);
int setOnLearnHandler(CPhidgetIRHandle ir, handler *h);
int setOnRawDataHandler(CPhidgetIRHandle ir, handler *h);
void unsetOnCodeHandler(CPhidgetIRHandle ir);
void unsetOnLearnHandler(CPhidgetIRHandle ir);
void unsetOnRawDataHandler(CPhidgetIRHandle ir);

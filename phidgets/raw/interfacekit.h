#include <phidget21.h>
#include <stdlib.h>
#include "handler.h"

typedef struct onChangeResult {
	int index;
	int value;
} onChangeResult;

typedef enum {
	inputChanged = 1,
	outputChanged,
	sensorChanged
} onChangeType;

void onChangeResultFree(onChangeResult *r);

onChangeResult * onChangeAwait(handler *h);

int setOnChangeHandler(CPhidgetInterfaceKitHandle ifk, handler *h, onChangeType t);
void unsetOnChangeHandler(CPhidgetInterfaceKitHandle ifk, onChangeType t);

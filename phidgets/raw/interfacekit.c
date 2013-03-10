#include "interfacekit.h"

void onChangeResultFree(onChangeResult *r) {
	free(r);
}

onChangeResult * onChangeAwait(handler *h) {
	return (onChangeResult *)handlerAwait(h);
}

int onChangeHandler(CPhidgetInterfaceKitHandle ifk, void *ptr, int index, int value) {
	handler * h = (handler *)ptr;

	onChangeResult * r = calloc(1, sizeof(onChangeResult));
	r->index = index;
	r->value = value;
	handlerAppendResult(h, r);

	return 0;
}

int setOnChangeHandler(CPhidgetInterfaceKitHandle ifk, handler *h, onChangeType t) {
	switch (t) {
		case inputChanged:
			return CPhidgetInterfaceKit_set_OnInputChange_Handler(ifk, &onChangeHandler, h);
		case outputChanged:
			return CPhidgetInterfaceKit_set_OnOutputChange_Handler(ifk, &onChangeHandler, h);
		case sensorChanged:
			return CPhidgetInterfaceKit_set_OnSensorChange_Handler(ifk, &onChangeHandler, h);
	}
}

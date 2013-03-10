package raw

// #include "handler.h"
import "C"

func createHandler(f func(h *C.handler) C.int) (*C.handler, error) {
	h := C.newHandler()

	if err := result(f(h)); err != nil {
		C.handlerFree(h)
		return nil, err
	}

	return h, nil
}

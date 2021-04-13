package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// HR is a wrapper for a GHR.
type HR struct {
	View
	private int
}

func NewHR(window *View) *HR {
	handle := C.hrNew(unsafe.Pointer(window))
	return (*HR)(handle)
}

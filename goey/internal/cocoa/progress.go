package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Progress is a wrapper for a NSProgressIndicator.
type Progress struct {
	View
	private int
}

func NewProgress(window *View, min, value, max float64) *Progress {
	handle := C.progressNew(unsafe.Pointer(window), C.double(min), C.double(value), C.double(max))
	return (*Progress)(handle)
}

func (w *Progress) Min() float64 {
	value := C.progressMin(unsafe.Pointer(w))
	return float64(value)
}

func (w *Progress) Max() float64 {
	value := C.progressMax(unsafe.Pointer(w))
	return float64(value)
}

func (w *Progress) Value() float64 {
	value := C.progressValue(unsafe.Pointer(w))
	return float64(value)
}

func (w *Progress) Update(min, value, max float64) {
	C.progressUpdate(unsafe.Pointer(w),
		C.double(min), C.double(value), C.double(max))
}

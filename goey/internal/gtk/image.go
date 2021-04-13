package gtk

// #include "thunks.h"
import "C"
import "unsafe"

func ImageImageData(handle uintptr) []byte {
	length := C.size_t(0)

	data := C.imageImageData(unsafe.Pointer(handle), &length)
	return C.GoBytes(unsafe.Pointer(data), C.int(length))
}

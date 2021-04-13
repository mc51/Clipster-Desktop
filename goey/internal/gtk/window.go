package gtk

// #include "thunks.h"
import "C"
import "unsafe"

type Window interface {
	Widget
	OnDeleteEvent() bool
	OnSizeAllocate(width, height int)
}

//export onDeleteEvent
func onDeleteEvent(handle unsafe.Pointer) bool {
	return widgets[uintptr(handle)].(Window).OnDeleteEvent()
}

//export onSizeAllocate
func onSizeAllocate(handle unsafe.Pointer, width, height int) {
	widgets[uintptr(handle)].(Window).OnSizeAllocate(width, height)
}

func WindowScreenshot(handle uintptr) ([]byte, bool, int, int, int) {
	var data unsafe.Pointer
	var dataLen C.size_t
	var hasAlpha C.bool
	var width, height C.int
	var stride C.unsigned

	C.windowScreenshot(unsafe.Pointer(handle), &data, &dataLen, &hasAlpha, &width, &height, &stride)

	return C.GoBytes(data, C.int(dataLen)), bool(hasAlpha), int(width), int(height), int(stride)
}

package gtk

// #include "thunks.h"
import "C"
import "unsafe"

type Widget interface {
	OnDestroy()
}

type WidgetWithFocus interface {
	Widget
	OnFocus()
	OnBlur()
}

var (
	widgets = map[uintptr]Widget{}
)

func RegisterWidget(handle uintptr, widget Widget) {
	widgets[handle] = widget
}

//export onDestroy
func onDestroy(handle unsafe.Pointer) {
	widgets[uintptr(handle)].OnDestroy()
	delete(widgets, uintptr(handle))
}

//export onFocus
func onFocus(handle unsafe.Pointer) {
	widgets[uintptr(handle)].(WidgetWithFocus).OnFocus()
}

//export onBlur
func onBlur(handle unsafe.Pointer) {
	widgets[uintptr(handle)].(WidgetWithFocus).OnBlur()
}

func WidgetNaturalSize(widget uintptr) (int, int) {
	var width, height C.int

	C.widgetNaturalSize(unsafe.Pointer(widget), &width, &height)
	return int(width), int(height)
}

func WindowSize(window uintptr) (int, int) {
	ret := C.windowSize(unsafe.Pointer(window))
	return int(ret.width), int(ret.height)
}

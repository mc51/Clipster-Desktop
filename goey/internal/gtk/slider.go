package gtk

import "C"
import "unsafe"

type Slider interface {
	WidgetWithFocus
	OnChange(value float64)
}

//export onChangeFloat64
func onChangeFloat64(handle unsafe.Pointer, value float64) {
	widgets[uintptr(handle)].(Slider).OnChange(value)
}

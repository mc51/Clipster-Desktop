package gtk

import "C"
import "unsafe"

type Checkbox interface {
	WidgetWithFocus
	OnChange(value bool)
}

//export onChangeBool
func onChangeBool(handle unsafe.Pointer, value bool) {
	widgets[uintptr(handle)].(Checkbox).OnChange(value)
}

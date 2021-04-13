package gtk

import "C"
import "unsafe"

type Combobox interface {
	WidgetWithFocus
	OnChange(value int)
}

//export onChangeInt
func onChangeInt(handle unsafe.Pointer, value int) {
	widgets[uintptr(handle)].(Combobox).OnChange(value)
}

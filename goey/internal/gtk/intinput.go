package gtk

import "C"
import "unsafe"

type IntInput interface {
	WidgetWithFocus
	OnChange(value int64)
	OnEnterKey(value int64)
}

//export onChangeInt64
func onChangeInt64(handle unsafe.Pointer, value int64) {
	widgets[uintptr(handle)].(IntInput).OnChange(value)
}

//export onEnterKeyInt64
func onEnterKeyInt64(handle unsafe.Pointer, value int64) {
	widgets[uintptr(handle)].(IntInput).OnEnterKey(value)
}

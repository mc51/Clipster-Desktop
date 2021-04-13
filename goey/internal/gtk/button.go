package gtk

import "C"
import "unsafe"

type Button interface {
	WidgetWithFocus
	OnClick()
}

//export onClick
func onClick(handle unsafe.Pointer) {
	widgets[uintptr(handle)].(Button).OnClick()
}

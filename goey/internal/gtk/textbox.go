package gtk

import "C"
import "unsafe"

type Textbox interface {
	WidgetWithFocus
	OnChange(value string)
	OnEnterKey(value string)
}

//export onChangeString
func onChangeString(handle unsafe.Pointer, text *C.char) {
	widgets[uintptr(handle)].(Textbox).OnChange(C.GoString(text))
}

//export onEnterKey
func onEnterKey(handle unsafe.Pointer, text *C.char) {
	widgets[uintptr(handle)].(Textbox).OnEnterKey(C.GoString(text))
}

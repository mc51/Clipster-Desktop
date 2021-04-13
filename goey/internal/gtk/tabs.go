package gtk

import "C"
import "unsafe"

type Tabs interface {
	Widget
	OnChange(value int)
}

//export onChangeTab
func onChangeTab(handle unsafe.Pointer, value int) {
	widgets[uintptr(handle)].(Tabs).OnChange(value)
}

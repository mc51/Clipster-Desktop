package gtk

import "C"
import "unsafe"
import "time"

type DateInput interface {
	WidgetWithFocus
	OnChange(time.Time)
}

//export onChangeTime
func onChangeTime(handle unsafe.Pointer, year, month, day uint) {
	value := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Local)
	widgets[uintptr(handle)].(DateInput).OnChange(value)
}

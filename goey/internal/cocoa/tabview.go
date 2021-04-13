package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// TabView is a wrapper for a NSTabView.
type TabView struct {
	View
	private int
}

type tabviewCallback struct {
	onChange func(int)
}

var (
	tabviewCallbacks = make(map[unsafe.Pointer]tabviewCallback)
)

func NewTabView(window *View) *TabView {
	handle := C.tabviewNew(unsafe.Pointer(window))
	return (*TabView)(handle)
}

func (w *TabView) Close() {
	C.viewClose(unsafe.Pointer(w))
	delete(tabviewCallbacks, unsafe.Pointer(w))
}

func (w *TabView) AddItem(text string) {
	ctext := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctext))
	}()

	C.tabviewAddItem(unsafe.Pointer(w), ctext)
}

func (w *TabView) ContentView(index int) *View {
	view := C.tabviewContentView(unsafe.Pointer(w), C.int(index))
	return (*View)(view)
}

func (w *TabView) ContentInsets() (int, int) {
	size := C.tabviewContentInsets(unsafe.Pointer(w))
	return int(size.width), int(size.height)
}

func (w *TabView) ItemAtIndex(index int) string {
	ctext := C.tabviewItemAtIndex(unsafe.Pointer(w), C.int(index))
	return C.GoString(ctext)
}

func (w *TabView) NumberOfItems() int {
	rc := C.tabviewNumberOfItems(unsafe.Pointer(w))
	return int(rc)
}

func (w *TabView) RemoveItemAtIndex(index int) {
	C.tabviewRemoveItemAtIndex(unsafe.Pointer(w), C.int(index))
}

func (w *TabView) SelectItem(index int) {
	C.tabviewSelectItem(unsafe.Pointer(w), C.int(index))
}

func (w *TabView) SetItemAtIndex(index int, text string) {
	ctext := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctext))
	}()

	C.tabviewSetItemAtIndex(unsafe.Pointer(w), C.int(index), ctext)
}

func (w *TabView) SetOnChange(cb func(int)) {
	tabviewCallbacks[unsafe.Pointer(w)] = tabviewCallback{
		onChange: cb,
	}
}

//export tabviewDidSelectItem
func tabviewDidSelectItem(handle unsafe.Pointer, index int) {
	if cb := tabviewCallbacks[handle]; cb.onChange != nil {
		cb.onChange(index)
	}
}

package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// PopUpButton is a wrapper for a NSPopUpButton.
type PopUpButton struct {
	Control
	private int
}

type popupbuttonCallback struct {
	onChange func(int)
	onFocus  func()
	onBlur   func()
}

var (
	popupbuttonCallbacks = make(map[unsafe.Pointer]popupbuttonCallback)
)

func NewPopUpButton(window *View) *PopUpButton {
	handle := C.popupbuttonNew(unsafe.Pointer(window))
	return (*PopUpButton)(handle)
}

func (w *PopUpButton) Close() {
	C.viewClose(unsafe.Pointer(w))
	delete(popupbuttonCallbacks, unsafe.Pointer(w))
}

func (w *PopUpButton) Callbacks() (func(int), func(), func()) {
	cb := popupbuttonCallbacks[unsafe.Pointer(w)]
	return cb.onChange, cb.onFocus, cb.onBlur
}

func (w *PopUpButton) SetCallbacks(onchange func(int), onfocus func(), onblur func()) {
	popupbuttonCallbacks[unsafe.Pointer(w)] = popupbuttonCallback{
		onChange: onchange,
		onFocus:  onfocus,
		onBlur:   onblur,
	}
}

func (w *PopUpButton) AddItem(text string) {
	ctext := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctext))
	}()

	C.popupbuttonAddItem(unsafe.Pointer(w), ctext)
}

func (w *PopUpButton) ItemAtIndex(index int) string {
	rc := C.popupbuttonItemAtIndex(unsafe.Pointer(w), C.int(index))
	return C.GoString(rc)
}

func (w *PopUpButton) NumberOfItems() int {
	rc := C.popupbuttonNumberOfItems(unsafe.Pointer(w))
	return int(rc)
}

func (w *PopUpButton) RemoveAllItems() {
	C.popupbuttonRemoveAllItems(unsafe.Pointer(w))
}

func (w *PopUpButton) SetValue(value int, unset bool) {
	C.popupbuttonSetValue(unsafe.Pointer(w), C.int(value), toBool(unset))
}

func (w *PopUpButton) Value() int {
	rc := C.popupbuttonValue(unsafe.Pointer(w))
	return int(rc)
}

//export popupbuttonOnChange
func popupbuttonOnChange(handle unsafe.Pointer, value int) {
	if cb := popupbuttonCallbacks[handle]; cb.onChange != nil {
		cb.onChange(value)
	}
}

//export popupbuttonOnFocus
func popupbuttonOnFocus(handle unsafe.Pointer) {
	if cb := popupbuttonCallbacks[handle]; cb.onFocus != nil {
		cb.onFocus()
	}
}

//export popupbuttonOnBlur
func popupbuttonOnBlur(handle unsafe.Pointer) {
	if cb := popupbuttonCallbacks[handle]; cb.onBlur != nil {
		cb.onBlur()
	}
}

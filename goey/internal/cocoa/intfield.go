package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// IntField is a wrapper for a NSTextField and NSStepper.
type IntField struct {
	Control
	private int
}

type intfieldCallback struct {
	onChange func(int64)
	onFocus  func()
	onBlur   func()
}

var (
	intfieldCallbacks = make(map[unsafe.Pointer]intfieldCallback)
)

func NewIntField(window *View, value, min, max int64) *IntField {
	handle := C.intfieldNew(unsafe.Pointer(window), C.int64_t(value), C.int64_t(min), C.int64_t(max))
	return (*IntField)(handle)
}

func (w *IntField) Close() {
	C.intfieldClose(unsafe.Pointer(w))
	delete(intfieldCallbacks, unsafe.Pointer(w))
}

func (w *IntField) Callbacks() (func(int64), func(), func()) {
	cb := intfieldCallbacks[unsafe.Pointer(w)]
	return cb.onChange, cb.onFocus, cb.onBlur
}

func (w *IntField) SetCallbacks(onchange func(int64), onfocus func(), onblur func()) {
	intfieldCallbacks[unsafe.Pointer(w)] = intfieldCallback{
		onChange: onchange,
		onFocus:  onfocus,
		onBlur:   onblur,
	}
}

func (w *IntField) IsEditable() bool {
	return C.intfieldIsEditable(unsafe.Pointer(w)) != 0
}

func (w *IntField) Max() int64 {
	return int64(C.intfieldMax(unsafe.Pointer(w)))
}

func (w *IntField) Min() int64 {
	return int64(C.intfieldMin(unsafe.Pointer(w)))
}

func (w *IntField) Placeholder() string {
	ctext := C.intfieldPlaceholder(unsafe.Pointer(w))
	return C.GoString(ctext)
}

func (w *IntField) SetEditable(value bool) {
	C.intfieldSetEditable(unsafe.Pointer(w), toBool(value))
}

func (c *IntField) SetFrame(x, y, dx, dy int) {
	C.intfieldSetFrame(unsafe.Pointer(c), C.int(x), C.int(y), C.int(dx), C.int(dy))
}

func (w *IntField) SetValue(value, min, max int64) {
	C.intfieldSetValue(unsafe.Pointer(w), C.int64_t(value), C.int64_t(min), C.int64_t(max))
}

func (w *IntField) SetPlaceholder(text string) {
	ctitle := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.intfieldSetPlaceholder(unsafe.Pointer(w), ctitle)
}

func (w *IntField) Value() int64 {
	value := C.intfieldValue(unsafe.Pointer(w))
	return int64(value)
}

//export intfieldOnChange
func intfieldOnChange(handle unsafe.Pointer, value int64) {
	if cb := intfieldCallbacks[handle]; cb.onChange != nil {
		cb.onChange(value)
	}
}

//export intfieldOnFocus
func intfieldOnFocus(handle unsafe.Pointer) {
	if cb := intfieldCallbacks[handle]; cb.onFocus != nil {
		cb.onFocus()
	}
}

//export intfieldOnBlur
func intfieldOnBlur(handle unsafe.Pointer) {
	if cb := intfieldCallbacks[handle]; cb.onBlur != nil {
		cb.onBlur()
	}
}

package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// TextField is a wrapper for a NSTextField.
type TextField struct {
	Control
	private int
}

type textfieldCallback struct {
	onChange   func(string)
	onFocus    func()
	onBlur     func()
	onEnterKey func(string)
}

var (
	textfieldCallbacks = make(map[unsafe.Pointer]textfieldCallback)
)

func NewTextField(window *View, title string, password bool) *TextField {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.textfieldNew(unsafe.Pointer(window), ctitle, toBool(password))
	return (*TextField)(handle)
}

func (w *TextField) Close() {
	C.viewClose(unsafe.Pointer(w))
	delete(textfieldCallbacks, unsafe.Pointer(w))
}

func (w *TextField) Callbacks() (func(string), func(), func(), func(string)) {
	cb := textfieldCallbacks[unsafe.Pointer(w)]
	return cb.onChange, cb.onFocus, cb.onBlur, cb.onEnterKey
}

func (w *TextField) SetCallbacks(onchange func(string), onfocus func(), onblur func(), onenterkey func(string)) {
	textfieldCallbacks[unsafe.Pointer(w)] = textfieldCallback{
		onChange:   onchange,
		onFocus:    onfocus,
		onBlur:     onblur,
		onEnterKey: onenterkey,
	}
}

func (w *TextField) IsEditable() bool {
	return C.textfieldIsEditable(unsafe.Pointer(w)) != 0
}

func (w *TextField) IsPassword() bool {
	return C.textfieldIsPassword(unsafe.Pointer(w)) != 0
}

func (w *TextField) Placeholder() string {
	ctext := C.textfieldPlaceholder(unsafe.Pointer(w))
	return C.GoString(ctext)
}

func (w *TextField) SetEditable(value bool) {
	C.textfieldSetEditable(unsafe.Pointer(w), toBool(value))
}

func (w *TextField) SetValue(text string) {
	ctitle := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.textfieldSetValue(unsafe.Pointer(w), ctitle)
}

func (w *TextField) SetPlaceholder(text string) {
	ctitle := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.textfieldSetPlaceholder(unsafe.Pointer(w), ctitle)
}

func (w *TextField) Value() string {
	ctext := C.textfieldValue(unsafe.Pointer(w))
	return C.GoString(ctext)
}

//export textfieldOnChange
func textfieldOnChange(handle unsafe.Pointer, text *C.char) {
	if cb := textfieldCallbacks[handle]; cb.onChange != nil {
		cb.onChange(C.GoString(text))
	}
}

//export textfieldOnFocus
func textfieldOnFocus(handle unsafe.Pointer) {
	if cb := textfieldCallbacks[handle]; cb.onFocus != nil {
		cb.onFocus()
	}
}

//export textfieldOnBlur
func textfieldOnBlur(handle unsafe.Pointer) {
	if cb := textfieldCallbacks[handle]; cb.onBlur != nil {
		cb.onBlur()
	}
}

//export textfieldOnEnterKey
func textfieldOnEnterKey(handle unsafe.Pointer, text *C.char) {
	if cb := textfieldCallbacks[handle]; cb.onEnterKey != nil {
		cb.onEnterKey(C.GoString(text))
	}
}

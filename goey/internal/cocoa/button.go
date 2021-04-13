package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Button is a wrapper for a NSButton.
type Button struct {
	Control
	private int
}

type buttonCallback struct {
	onClick  func()
	onChange func(bool)
	onFocus  func()
	onBlur   func()
}

var (
	buttonCallbacks = make(map[unsafe.Pointer]buttonCallback)
)

func NewButton(window *View, title string) *Button {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.buttonNew(unsafe.Pointer(window), ctitle)
	return (*Button)(handle)
}

func NewCheckButton(window *View, title string, value bool) *Button {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.buttonNewCheck(unsafe.Pointer(window), ctitle, toBool(value))
	return (*Button)(handle)
}

func (w *Button) Close() {
	C.viewClose(unsafe.Pointer(w))
	delete(buttonCallbacks, unsafe.Pointer(w))
}

func (w *Button) PerformClick() {
	C.buttonPerformClick(unsafe.Pointer(w))
}

func (w *Button) Callbacks() (func(), func(bool), func(), func()) {
	cb := buttonCallbacks[unsafe.Pointer(w)]
	return cb.onClick, cb.onChange, cb.onFocus, cb.onBlur
}

func (w *Button) SetCallbacks(onclick func(), onchange func(bool), onfocus func(), onblur func()) {
	buttonCallbacks[unsafe.Pointer(w)] = buttonCallback{
		onClick:  onclick,
		onChange: onchange,
		onFocus:  onfocus,
		onBlur:   onblur,
	}
}

func (w *Button) IsDefault() bool {
	rc := C.buttonIsDefault(unsafe.Pointer(w))
	return rc != 0
}

func (w *Button) SetDefault(value bool) {
	C.buttonSetDefault(unsafe.Pointer(w), toBool(value))
}

func (w *Button) State() bool {
	rc := C.buttonState(unsafe.Pointer(w))
	return int(rc) != 0
}

func (w *Button) SetState(value bool) {
	C.buttonSetState(unsafe.Pointer(w), toBool(value))
}

func (w *Button) Title() string {
	cstring := C.buttonTitle(unsafe.Pointer(w))
	return C.GoString(cstring)
}

func (w *Button) SetTitle(title string) {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.buttonSetTitle(unsafe.Pointer(w), ctitle)
}

//export buttonOnClick
func buttonOnClick(handle unsafe.Pointer) {
	if cb := buttonCallbacks[handle]; cb.onClick != nil {
		cb.onClick()
	}
}

//export buttonOnChange
func buttonOnChange(handle unsafe.Pointer, value bool) {
	if cb := buttonCallbacks[handle]; cb.onChange != nil {
		cb.onChange(value)
	}
}

//export buttonOnFocus
func buttonOnFocus(handle unsafe.Pointer) {
	if cb := buttonCallbacks[handle]; cb.onFocus != nil {
		cb.onFocus()
	}
}

//export buttonOnBlur
func buttonOnBlur(handle unsafe.Pointer) {
	if cb := buttonCallbacks[handle]; cb.onBlur != nil {
		cb.onBlur()
	}
}

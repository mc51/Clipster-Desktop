package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Slider is a wrapper for a NSSlider.
type Slider struct {
	Control
	private int
}

type sliderCallback struct {
	onChange func(float64)
	onFocus  func()
	onBlur   func()
}

var (
	sliderCallbacks = make(map[unsafe.Pointer]sliderCallback)
)

func NewSlider(window *View, min, value, max float64) *Slider {
	handle := C.sliderNew(unsafe.Pointer(window), C.double(min), C.double(value), C.double(max))
	return (*Slider)(handle)
}

func (w *Slider) Callbacks() (func(float64), func(), func()) {
	cb := sliderCallbacks[unsafe.Pointer(w)]
	return cb.onChange, cb.onFocus, cb.onBlur
}

func (w *Slider) SetCallbacks(onchange func(float64), onfocus func(), onblur func()) {
	sliderCallbacks[unsafe.Pointer(w)] = sliderCallback{
		onChange: onchange,
		onFocus:  onfocus,
		onBlur:   onblur,
	}
}

func (w *Slider) Min() float64 {
	value := C.sliderMin(unsafe.Pointer(w))
	return float64(value)
}

func (w *Slider) Max() float64 {
	value := C.sliderMax(unsafe.Pointer(w))
	return float64(value)
}

func (w *Slider) Value() float64 {
	value := C.sliderValue(unsafe.Pointer(w))
	return float64(value)
}

func (w *Slider) Update(min, value, max float64) {
	C.sliderUpdate(unsafe.Pointer(w),
		C.double(min), C.double(value), C.double(max))
}

//export sliderOnChange
func sliderOnChange(handle unsafe.Pointer, value float64) {
	if cb := sliderCallbacks[handle]; cb.onChange != nil {
		cb.onChange(value)
	}
}

//export sliderOnFocus
func sliderOnFocus(handle unsafe.Pointer) {
	if cb := sliderCallbacks[handle]; cb.onFocus != nil {
		cb.onFocus()
	}
}

//export sliderOnBlur
func sliderOnBlur(handle unsafe.Pointer) {
	if cb := sliderCallbacks[handle]; cb.onBlur != nil {
		cb.onBlur()
	}
}
